package consumer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/fullnode_bridge/tx"
	"github.com/cosmos/cosmos-sdk/fullnode_bridge/types"

	"github.com/cosmos/cosmos-sdk/fullnode_bridge/config"

	"github.com/IBM/sarama"

	"crypto/sha256"

	"github.com/btcsuite/btcutil/bech32"
	"golang.org/x/crypto/ripemd160"
)

type lightTxHandler struct{}

type SignatureEntry struct {
	TxMsg     types.LightTxMessage
	Address   string
	Timestamp time.Time
}

var (
	VoteMap   = make(map[string][]SignatureEntry) // hash -> 서명자 목록
	DeviceID  = make(map[string]string)           // hash -> device_id
	VoteMutex sync.Mutex
)

var (
	VoteMemberCount int // 데이터베이스 멤버 수 기록 변수
)

func (h *lightTxHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *lightTxHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *lightTxHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Println("[Solar data][Raw Message]:", string(msg.Value)) // 👉 수신된 원본 메시지 출력

		var txMsg types.LightTxMessage
		if err := json.Unmarshal(msg.Value, &txMsg); err != nil {
			fmt.Println("[Solar data] 메시지 파싱 실패:", err)
			continue
		}

		for len(txMsg.Pubkey)%4 != 0 {
			txMsg.Pubkey += "="
		}
		pubkeyBytes, err := base64.StdEncoding.DecodeString(txMsg.Pubkey)
		if err != nil {
			fmt.Println("[Solar data] 퍼블릭키 디코딩 실패:", err)
			continue
		}

		if len(pubkeyBytes) != 33 {
			fmt.Println("❌ 잘못된 퍼블릭키 길이:", len(pubkeyBytes))
			continue
		}

		address, err := PubKeyToAddress(pubkeyBytes)
		if err != nil {
			fmt.Println("주소 생성 실패:", err)
			continue
		}
		VoteMutex.Lock()
		VoteMap[txMsg.Hash] = append(VoteMap[txMsg.Hash], SignatureEntry{
			TxMsg:     txMsg,
			Address:   address,
			Timestamp: time.Now(),
		})
		// 해시별로 device_id 또는 facility_id 저장
		if txMsg.Original != nil && txMsg.Original.DeviceID != "" {
			DeviceID[txMsg.Hash] = txMsg.Original.DeviceID
		} else if txMsg.REC != nil && txMsg.REC.FacilityId != "" {
			DeviceID[txMsg.Hash] = txMsg.REC.FacilityId
		} else {
			fmt.Println("[Solar data] DeviceID와 FacilityId 모두 존재하지 않음, 저장 안 함:", txMsg.Hash)
		}
		VoteMutex.Unlock()

		session.MarkMessage(msg, "")
	}
	return nil
}

func PubKeyToAddress(pubKeyBytes []byte) (string, error) {
	// 1. SHA-256
	sha := sha256.Sum256(pubKeyBytes)

	// 2. RIPEMD-160
	ripemd := ripemd160.New()
	_, err := ripemd.Write(sha[:])
	if err != nil {
		return "", err
	}
	pubKeyHash := ripemd.Sum(nil) // 20바이트

	// 3. Bech32 인코딩
	converted, err := bech32.ConvertBits(pubKeyHash, 8, 5, true)
	if err != nil {
		return "", err
	}
	address, err := bech32.Encode("cosmos", converted)
	if err != nil {
		return "", err
	}

	return address, nil
}

func StartVoteEvaluator() {
	fmt.Println("[Solar data] StartVoteEvaluator 시작됨")

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			now := time.Now()
			fmt.Println("[Solar data] 투표 수집 시작:", now.Format(time.RFC3339))

			VoteMutex.Lock()
			for hash, entries := range VoteMap {
				if len(entries) == 0 {
					fmt.Printf("[Solar data] Tx: [%s] entries 없음. 건너뜀\n", hash)
					continue
				}

				elapsed := now.Sub(entries[0].Timestamp)
				fmt.Printf("[Solar data]  [%s] entry 수: %d, 경과시간: %.1f초\n", hash, len(entries), elapsed.Seconds())

				if elapsed < 10*time.Second {
					fmt.Printf("[Solar data] (%.1f초 경과). 투표 검증 중\n", elapsed.Seconds())
					continue
				}

				// 주소 중복 제거
				unique := map[string]bool{}
				for _, e := range entries {
					unique[e.Address] = true
				}
				var uniqueList []string
				for k := range unique {
					uniqueList = append(uniqueList, k)
				}

				if len(unique) >= 1 {
					txMsg := entries[0].TxMsg
					fmt.Println("[Solar data] 트랜잭션 전송 시도 중...")

					txHash, err := tx.BroadcastLightTx(txMsg)
					if err != nil {
						fmt.Println("[Solar data] 트랜잭션 전송 실패:", err)
					} else {
						fmt.Printf("[Solar data] 트랜잭션 전송 성공: %s\n", txHash)
						fmt.Printf("[Solar data] → 서명자 주소 목록: %v\n", uniqueList)

						// 디바이스 아이디에 매칭되는 주소 get
						var userAddress = "cosmos1234"
						tx.SendRewardTx(userAddress, txMsg.Original.Power)
					}

					delete(VoteMap, hash)
					fmt.Printf("[Solar data] [%s] voteMap에서 제거됨\n", hash)
				} else {
					fmt.Printf("[Solar data] 고유 주소 없음. 트랜잭션 전송 안 함\n")
				}
			}
			VoteMutex.Unlock()
		}
	}()
}

func StartSolarKafkaConsumer() {
	brokers := config.KafkaBrokers
	topic := config.TopicLightTx
	groupID := config.TopicLightTxGroup // 모든 서버에서 동일하게 설정해야 함

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_1_0_0
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		panic(fmt.Sprintf("[Solar data] ConsumerGroup 생성 실패: %v", err))
	}

	go func() {
		for {
			err := consumerGroup.Consume(context.Background(), []string{topic}, &lightTxHandler{})
			if err != nil {
				fmt.Printf("[Solar data] Consume 중 오류 발생: %v\n", err)
			}
		}
	}()

	fmt.Println("[Solar data] Kafka Consumer Group 수신 대기 중...")
	StartVoteEvaluator() // 참여자 수집 + 평가 루틴 시작
}

// 회원가입 알고리즘

type accountHandler struct {
	producer    sarama.SyncProducer
	resultTopic string
}

func (h *accountHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *accountHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *accountHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var authMsg types.AuthMessage
		if err := json.Unmarshal(msg.Value, &authMsg); err != nil {
			fmt.Println("[Account] 메시지 파싱 실패:", err)
			continue
		}

		fmt.Println("[Account] 주소 활성화 요청:", authMsg.ID)

		// 주소로 1 stake 전송 (tx.SendStakeToAddress 함수로 정의)
		output, err := tx.SendStakeToAddress(authMsg.ID)
		if err != nil {
			fmt.Println("[Account] 송금 실패:", err)
			continue
		}

		fmt.Println("[Account] 송금 성공:\n", output)

		producerMsg := &sarama.ProducerMessage{
			Topic: h.resultTopic,
			Value: sarama.StringEncoder(output),
		}

		_, _, err = h.producer.SendMessage(producerMsg)
		if err != nil {
			fmt.Println("[Account] 결과 메시지 전송 실패:", err)
		} else {
			fmt.Println("[Account] 결과 메시지 전송 완료")
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func StartAccountConsumer() {
	brokers := config.KafkaBrokers
	topic := config.TopicAccountCreate
	resultTopic := config.TopicAccountResult
	groupID := config.TopicAccountGroup

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_1_0_0
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Producer 생성
	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		panic(fmt.Sprintf("[Account] Kafka producer 생성 실패: %v", err))
	}

	// ConsumerGroup 생성
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		panic(fmt.Sprintf("[Account] Kafka ConsumerGroup 생성 실패: %v", err))
	}

	handler := &accountHandler{
		producer:    producer,
		resultTopic: resultTopic,
	}

	go func() {
		for {
			err := consumerGroup.Consume(context.Background(), []string{topic}, handler)
			if err != nil {
				fmt.Printf("[Account] Consume 오류: %v\n", err)
			}
		}
	}()

	fmt.Println("[Account] Kafka Consumer Group 수신 대기 중...")
}

func StartVoteMemberConsumer() {
	brokers := config.KafkaBrokers
	topic := config.TopicVoteMember
	partition := int32(0)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_1_0_0

	consumer, err := sarama.NewConsumer(brokers, saramaConfig)
	if err != nil {
		panic(fmt.Sprintf("[Kafka: VoteMember] Consumer 생성 실패: %v", err))
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		panic(fmt.Sprintf("[Kafka: VoteMember] 파티션 구독 실패: %v", err))
	}

	go func() {
		fmt.Println("[Kafka: VoteMember] Kafka Partition Consumer 수신 대기 중...")
		for msg := range partitionConsumer.Messages() {
			fmt.Printf("[Kafka: VoteMember] 수신 메시지: %s\n", string(msg.Value))

			// 여기에 메시지 파싱 및 전역 변수 갱신 로직 추가
		}
	}()
}
