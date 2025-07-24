package keeper

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type Keeper struct {
	// 필요한 의존성 (예: storeKey, codec 등)
	storeKey storetypes.StoreKey
}

func NewKeeper(storeKey storetypes.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
	}
}
