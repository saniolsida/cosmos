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
			fmt.Sprintf("%.2f", msg.Original.Power),
			fmt.Sprintf("%.2f", msg.Original.Voltage),
			fmt.Sprintf("%.2f", msg.Original.PowerOutput),
			msg.Hash,
			msg.Signature,
			msg.Pubkey,
		}
	} else if msg.REC != nil {
		// RECMeta 전송
		args = []string{
			"tx", "lighttx", "send-rec-tx",
			msg.REC.FacilityId,
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
	cmd := exec.Command("simd", "tx", "bank", "send",
		"alice", "cosmos1n4q3249rl4rh7vqwve4rxdxpa5yeg5dkf32k7f", "1stake",
		"--fees", "200stake",
		"--chain-id", "learning-chain-1",
		"--home", "/root/cosmos/cosmos-sdk/private/.simapp",
		"--yes", "--keyring-backend", "test", "--broadcast-mode", "sync")

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func SendRewardTx(toAddr string, power float64) (string, error) {
	// 발전량이 0 이하이면 트랜잭션 안 보냄
	if power <= 0 {
		return "", fmt.Errorf("보상할 발전량이 없습니다")
	}

	// 소수점 버림
	amount := int64(power)
	amountStr := strconv.FormatInt(amount, 10)

	// 트랜잭션 실행 명령
	cmd := exec.Command("/root/cosmos/cosmos-sdk/build/simd", "tx", "reward", "reward-solar-power",
		toAddr, amountStr,
		"--from", "alice",
		"--chain-id", "learning-chain-1",
		"--home", "/root/cosmos/cosmos-sdk/private/.simapp",
		"--fees", "0.01stake",
		"--gas", "auto",
		"--yes",
		"--keyring-backend", "test",
		"--broadcast-mode", "sync")

	out, err := cmd.CombinedOutput()
	return string(out), err
}
