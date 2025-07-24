package types

import (
	// "github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// 인터페이스 등록
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Msg 타입이 있다면 여기에 등록하세요. 예시:
	// registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDoSomething{})
}
