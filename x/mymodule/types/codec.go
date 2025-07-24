package types

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Msg 타입이 있다면 여기에 등록하세요. 예시:
	// cdc.RegisterConcrete(MsgDoSomething{}, "mymodule/MsgDoSomething", nil)
}
