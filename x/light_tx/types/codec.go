package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	// 1. MsgSendLightTx를 등록
	msg := &MsgSendLightTx{}

	// 2. 인터페이스 등록
	reg.RegisterInterface(
		sdk.MsgInterfaceProtoName,
		(*sdk.Msg)(nil),
	)

	reg.RegisterImplementations(
		(*sdk.Msg)(nil),
		msg,
	)

	// 4. MsgServiceDesc 등록
	msgservice.RegisterMsgServiceDesc(reg, &_Msg_serviceDesc)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// 아직 특별히 등록할 타입이 없으면 그냥 빈 함수로 둬도 됨
}

var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
