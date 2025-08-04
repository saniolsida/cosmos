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
func (m *msgServer) RewardSolarPower(goCtx context.Context, msg *types.MsgRewardSolarPower) (*types.MsgRewardSolarPowerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fmt.Println("msgServer.RewardSolarPower called")

	// 여기서 내부 로직 호출
	err := m.Keeper.RewardSolarPower(ctx, msg.Address, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgRewardSolarPowerResponse{}, nil
}
