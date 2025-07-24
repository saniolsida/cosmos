package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mapping/types"
)

type msgServer struct {
	Keeper
	types.UnimplementedMsgServer
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) RegisterDevice(goCtx context.Context, msg *types.MsgRegisterDevice) (*types.MsgRegisterDeviceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	m.Keeper.Logger(ctx).Info("🧪 RegisterDevice called", "is_check_tx", ctx.IsCheckTx(), "height", ctx.BlockHeight())

	// Bech32 문자열 주소를 AccAddress로 변환
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	// 매핑 저장
	m.Keeper.SetDeviceMapping(ctx, msg.DeviceId, address)

	// 로그
	m.Logger(ctx).Info("✅ 디바이스 등록 성공",
		"device_id", msg.DeviceId,
		"address", address.String(),
	)

	return &types.MsgRegisterDeviceResponse{}, nil
}
