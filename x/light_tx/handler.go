package light_tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/light_tx/keeper"
	"github.com/cosmos/cosmos-sdk/x/light_tx/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgSendLightTx:
			msgServer := keeper.NewMsgServerImpl(k)
			_, err := msgServer.SendLightTx(sdk.WrapSDKContext(ctx), msg)
			if err != nil {
				return nil, err
			}
			return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized light_tx message type: %T", msg)
		}
	}
}
