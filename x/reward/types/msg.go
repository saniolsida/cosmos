// x/reward/types/msgs.go

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgRewardSolarPower{}

func NewMsgRewardSolarPower(creator, address, power string) *MsgRewardSolarPower {
	return &MsgRewardSolarPower{
		Creator: creator,
		Address: address,
		Amount:  power,
	}
}

func (msg *MsgRewardSolarPower) Route() string {
	return ModuleName
}

func (msg *MsgRewardSolarPower) Type() string {
	return "RewardSolarPower"
}

func (msg *MsgRewardSolarPower) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err) // 안전하지 않지만 예시에서는 사용
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRewardSolarPower) GetSignBytes() []byte {
	// LegacyAmino JSON 직렬화
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRewardSolarPower) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return err
	}
	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return err
	}
	return nil
}
