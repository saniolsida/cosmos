package keeper

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/mapping/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	logger   log.Logger
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, logger log.Logger) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: key,
		logger:   logger,
	}
}

func (k Keeper) SetDeviceMapping(ctx sdk.Context, deviceID string, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := []byte("device:" + deviceID) // ✅ 여기에 key 정의
	store.Set(key, address.Bytes())

	k.Logger(ctx).Info("📝 SetDeviceMapping",
		"store_key_str", string(key),
		"storeKey_ptr", fmt.Sprintf("%p", k.storeKey),
	)
	k.Logger(ctx).Info("✅ SetDeviceMapping called",
		"device_id", deviceID,
		"address", address.String(),
		"store_key", string(key),
	)
}

func (k Keeper) GetDeviceAddress(ctx sdk.Context, deviceID string) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	key := []byte("device:" + deviceID)
	bz := store.Get(key)

	k.Logger(ctx).Info("GetDeviceAddress called",
		"height", ctx.BlockHeight(),
		"is_check_tx", ctx.IsCheckTx(),
	)

	k.Logger(ctx).Info("📦 storeKey info",
		"device_id", deviceID,
		"key", string(key),
		"storeKey_ptr", fmt.Sprintf("%p", k.storeKey),
	)
	if bz == nil {
		k.Logger(ctx).Info("❌ GetDeviceAddress: 디바이스를 찾을 수 없음",
			"device_id", deviceID,
			"store_key", string(key),
		)
		return nil, false
	}

	address := sdk.AccAddress(bz)

	k.Logger(ctx).Info("✅ GetDeviceAddress: 디바이스 조회 성공",
		"device_id", deviceID,
		"store_key", string(key),
		"address", address.String(),
	)

	return address, true
}

// ✅ keeper/keeper.go
func (k Keeper) GetAddressFromDeviceID(
	ctx context.Context,
	req *types.QueryGetAddressFromDeviceIDRequest,
) (*types.QueryGetAddressFromDeviceIDResponse, error) {

	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	address, found := k.GetDeviceAddress(sdkCtx, req.DeviceId)

	k.Logger(sdkCtx).Info("🧪 GetAddressFromDeviceID",
		"is_check_tx", sdkCtx.IsCheckTx(),
		"context_type", fmt.Sprintf("%T", ctx),
		"caller_stack", string(debug.Stack()),
	)

	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "device_id %s not found", req.DeviceId)
	}

	return &types.QueryGetAddressFromDeviceIDResponse{
		Address: address.String(),
	}, nil
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger()
}
