package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgRegisterDevice) Route() string {
	return RouterKey // ex: "mapping"
}

// Type returns the message type
func (msg *MsgRegisterDevice) Type() string {
	return "RegisterDevice"
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRegisterDevice) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator address cannot be empty")
	}
	if msg.DeviceId == "" {
		return fmt.Errorf("device ID cannot be empty")
	}
	if msg.Address == "" {
		return fmt.Errorf("mapped address cannot be empty")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgRegisterDevice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err) // 에러가 발생하면 앱 종료
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgRegisterDevice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
