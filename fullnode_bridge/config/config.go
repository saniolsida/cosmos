package config

var (
	// Kafka 브로커 IP 및 포트
	KafkaBrokers = []string{"192.168.0.132:10000"}

	// Kafka Consumer
	TopicLightTx         = "light-vote-topic"     // 태양광 발전량 메세지 토픽
	TopicAccountCreate   = "create-address-topic" // 회원가입 메세지 토픽
	TopicVoteMember      = "user-count-topic"     // 전체 유권자 수 갱신 토픽
	TopicDeviceToAddress = "device-address-topic" // 오라클 -> 풀노드: 주소값 반환 토픽

	TopicAccountResult = "result-account-topic" // 회원가입 후 잔액 메세지 토픽

	//Kafka Producer
	TopicDeviceToAddressRequest = "device-address-request-topic" // 풀노드 -> 오라클: 주소값 요청 토픽
	TopicRequestMemberCount     = "request-user-count-topic"
	//Group
	TopicLightTxGroup         = "light-tx-group"
	TopicDeviceToAddressGroup = "address-device-group"
	TopicAccountGroup         = "account-create-group"
)
