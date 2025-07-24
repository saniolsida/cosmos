// x/reward/handler.go

package reward

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/reward/keeper"
	"github.com/cosmos/cosmos-sdk/x/reward/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgRewardSolarPower:
			// 보상 처리 로직
			err := k.RewardSolarPower(ctx, msg.Address, msg.Amount)
			if err != nil {
				return nil, err
			}
			return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized reward message type: %T", msg)
		}
	}
}
