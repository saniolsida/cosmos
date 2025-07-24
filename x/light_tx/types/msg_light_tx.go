package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSendLightTx{}

// 생성자 함수 (SolarData용)
func NewMsgSendLightTxWithSolar(
	creator string,
	original *SolarData,
	hash, signature, pubkey string,
) *MsgSendLightTx {
	return &MsgSendLightTx{
		Creator:   creator,
		Payload:   &MsgSendLightTx_Original{Original: original},
		Hash:      hash,
		Signature: signature,
		Pubkey:    pubkey,
	}
}

// 생성자 함수 (RECMeta용)
func NewMsgSendLightTxWithREC(
	creator string,
	rec *RECMeta,
	hash, signature, pubkey string,
) *MsgSendLightTx {
	return &MsgSendLightTx{
		Creator:   creator,
		Payload:   &MsgSendLightTx_Rec{Rec: rec},
		Hash:      hash,
		Signature: signature,
		Pubkey:    pubkey,
	}
}

// 라우팅 정보
func (msg *MsgSendLightTx) Route() string {
	return ModuleName
}

func (msg *MsgSendLightTx) Type() string {
	return "SendLightTx"
}

// 서명자 추출
func (msg *MsgSendLightTx) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(fmt.Sprintf("invalid creator address: %s", err))
	}
	return []sdk.AccAddress{creator}
}

// 서명 바이트
func (msg *MsgSendLightTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// 기본 유효성 검사
func (msg *MsgSendLightTx) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	switch payload := msg.Payload.(type) {
	case *MsgSendLightTx_Original:
		if payload.Original == nil {
			return fmt.Errorf("original solar data is nil")
		}
	case *MsgSendLightTx_Rec:
		if payload.Rec == nil {
			return fmt.Errorf("REC metadata is nil")
		}
	default:
		return fmt.Errorf("either SolarData or RECMeta must be provided")
	}

	if msg.Signature == "" {
		return fmt.Errorf("signature is required")
	}
	if msg.Pubkey == "" {
		return fmt.Errorf("pubkey is required")
	}
	return nil
}
