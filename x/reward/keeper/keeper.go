package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/reward/types"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	bankKeeper    types.BankKeeper // 인터페이스로 정의 필요
	AccountKeeper types.AccountKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		bankKeeper:    bankKeeper,
		AccountKeeper: accountKeeper,
	}
}

func (k Keeper) RewardSolarPower(ctx sdk.Context, to string, amount string) error {

	if k.bankKeeper == nil {
		fmt.Println("❌ bankKeeper is nil!!!")
		panic("bankKeeper is nil")
	}

	amt, ok := sdk.NewIntFromString(amount)
	if !ok {
		return fmt.Errorf("잘못된 amount 형식: %s", amount)
	}

	coins := sdk.NewCoins(sdk.NewCoin("stake", amt))

	toAddr, err := sdk.AccAddressFromBech32(to)
	if err != nil {
		return err
	}

	//코인 발행
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coins)
}
