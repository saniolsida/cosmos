package tx

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/cosmos/cosmos-sdk/fullnode_bridge/types"
)

func BroadcastLightTx(msg types.LightTxMessage) (string, error) {
	var args []string

	if msg.Original != nil {
		// SolarData 전송
		args = []string{
			"tx", "lighttx", "send-light-tx",
			msg.Original.DeviceID,
			msg.Original.Timestamp,
			fmt.Sprintf("%.2f", msg.Original.TotalEnergy),
			msg.Hash,
			msg.Signature,
			msg.Pubkey,
		}
	} else if msg.REC != nil {
		// RECMeta 전송
		args = []string{
			"tx", "lighttx", "send-light-tx",
			msg.REC.FacilityID,
			msg.REC.FacilityName,
			msg.REC.Location,
			msg.REC.TechnologyType,
			msg.REC.CapacityMW,
			msg.REC.RegistrationDate,
			msg.REC.CertifiedId,
			msg.REC.IssueData,
			msg.REC.GenerationStartDate,
			msg.REC.GenerationEndDate,
			msg.REC.MeasuredVolumeMWh,
			msg.REC.RetiredDate,
			msg.REC.RetirementPurpose,
			msg.REC.Status,
			msg.REC.Timestamp,
			msg.Hash,
			msg.Signature,
			msg.Pubkey,
		}
	} else {
		return "", fmt.Errorf("no valid data to send (both Original and REC are nil)")
	}

	// 공통 플래그 추가
	args = append(args,
		"--from", "alice",
		"--home", "/root/cosmos/cosmos-sdk/private/.simapp",
		"--chain-id", "learning-chain-1",
		"--keyring-backend", "test",
		"--broadcast-mode", "block",
		"--node", "http://localhost:26657",
		"--yes",
		"--output", "json",
	)

	cmd := exec.Command("/root/cosmos/cosmos-sdk/build/simd", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("simd error: %v\noutput: %s", err, string(output))
	}

	var resp struct {
		TxHash string `json:"txhash"`
	}
	if err := json.Unmarshal(output, &resp); err != nil {
		return "", fmt.Errorf("failed to parse tx response: %v\noutput: %s", err, string(output))
	}

	return resp.TxHash, nil
}

// SendStakeToAddress.go
func SendStakeToAddress(toAddr string) (string, error) {
	// 예시 CLI 호출
	cmd := exec.Command("/root/cosmos/cosmos-sdk/build/simd", "tx", "bank", "send",
		"alice", toAddr, "1stake",
		"--fees", "0.01stake",
		"--chain-id", "learning-chain-1",
		"--home", "/root/cosmos/cosmos-sdk/private/.simapp",
		"--yes", "--keyring-backend", "test", "--broadcast-mode", "sync")

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func QueryBalance(address string) (string, error) {
	// simd CLI를 통한 잔고 조회
	cmd := exec.Command("/root/cosmos/cosmos-sdk/build/simd", "query", "bank", "balances", address,
		"--node", "tcp://localhost:26657", // RPC 노드 주소 (필요 시 수정 가능)
		"--home", "/root/cosmos/cosmos-sdk/private/.simapp",
		"--output", "json")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("잔고 조회 실패: %v\n출력: %s", err, string(out))
	}
	return string(out), nil
}

func SendRewardTx(toAddr string, power float64) (string, error) {
	// 발전량이 0 이하이면 트랜잭션 안 보냄
	if power <= 0 {
		return "", fmt.Errorf("[Kafka: reward] 보상할 발전량이 없습니다")
	}

	// 소수점 버림
	amount := int64(power)
	amountStr := strconv.FormatInt(amount, 10)

	fmt.Printf("Kafka: [reward] 보상 트랜잭션 준비 중: 주소=%s, 발전량=%.2f → 지급액=%dstake\n", toAddr, power, amount)

	// 트랜잭션 실행 명령
	cmd := exec.Command("/root/cosmos/cosmos-sdk/build/simd", "tx", "reward", "reward-solar-power",
		toAddr, amountStr,
		"--from", "alice",
		"--chain-id", "learning-chain-1",
		"--home", "/root/cosmos/cosmos-sdk/private/.simapp",
		"--gas", "auto",
		"--yes",
		"--keyring-backend", "test",
		"--broadcast-mode", "sync")

	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		fmt.Println("[Kafka: reward] 트랜잭션 전송 실패:", err)
		fmt.Println("[Kafka: reward] 출력 내용:", output)
		return output, err
	}

	return output, nil
}
