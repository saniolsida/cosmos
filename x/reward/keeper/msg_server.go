package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/reward/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// func (m *msgServer) RewardSolarPower(goCtx context.Context, msg *types.MsgRewardSolarPower) (*types.MsgRewardSolarPowerResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(goCtx)
// 	fmt.Println("✅ RewardSolarPower called in msgServer!")

// 	// 숫자 문자열을 Int로 파싱
// 	amountInt, ok := sdk.NewIntFromString(msg.Amount)
// 	if !ok {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount must be integer string")
// 	}

// 	// "stake" 단위를 붙여 Coins로 생성
// 	coins := sdk.NewCoins(sdk.NewCoin("stake", amountInt))

// 	// 주소 변환
// 	toAddr, err := sdk.AccAddressFromBech32(msg.Address)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 모듈 계정 → 사용자에게 코인 전송
// 	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coins)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &types.MsgRewardSolarPowerResponse{}, nil
// }

func (m *msgServer) RewardSolarPower(goCtx context.Context, msg *types.MsgRewardSolarPower) (*types.MsgRewardSolarPowerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fmt.Println("✅ msgServer.RewardSolarPower called")

	// 여기서 내부 로직 호출
	err := m.Keeper.RewardSolarPower(ctx, msg.Address, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgRewardSolarPowerResponse{}, nil
}
