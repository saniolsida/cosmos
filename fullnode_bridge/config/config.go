package config

var (
	// Kafka 브로커 IP 및 포트
	KafkaBrokers = []string{"192.168.0.14:10001"}

	// Kafka 토픽 이름들
	TopicLightTx       = "light-vote-topic"
	TopicAccountCreate = "account-create-topic"
	TopicAccountResult = "account-result-topic"

	TopicLightTxGroup = "light-tx-group"
	TopicAccountGroup = "account-create-group"

	TopicVoteMember = "vote-member-topic"
)
