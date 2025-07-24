package types

import (
	_ "github.com/gogo/protobuf/proto" // ensure proto interface included
)

func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

func (gs *GenesisState) Validate() error {
	// 유효성 검사 로직, 없으면 nil 반환
	return nil
}
