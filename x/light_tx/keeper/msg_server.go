package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors" // 👈 이 줄 추가
	lighttype "github.com/cosmos/cosmos-sdk/x/light_tx/types"
)

type msgServer struct {
	Keeper
	lighttype.UnimplementedMsgServer
}

func NewMsgServerImpl(k Keeper) lighttype.MsgServer {
	return &msgServer{Keeper: k}
}

func (k msgServer) SendLightTx(goCtx context.Context, msg *lighttype.MsgSendLightTx) (*lighttype.MsgSendLightTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 분기 처리: oneof payload 중 어떤 타입이 들어왔는지 확인
	switch payload := msg.Payload.(type) {
	case *lighttype.MsgSendLightTx_Original:
		data := payload.Original
		k.Logger(ctx).Info("📩 Received LightTx (SolarData)",
			"creator", msg.Creator,
			"device_id", data.DeviceId,
			"timestamp", data.Timestamp,
			"power", data.Power,
			"voltage", data.Voltage,
			"power_output", data.PowerOutput,
			"hash", msg.Hash,
			"signature", msg.Signature,
			"pubkey", msg.Pubkey,
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent("light_tx_solar",
				sdk.NewAttribute("creator", msg.Creator),
				sdk.NewAttribute("device_id", data.DeviceId),
				sdk.NewAttribute("hash", msg.Hash),
				sdk.NewAttribute("signature", msg.Signature),
			),
		)

	case *lighttype.MsgSendLightTx_Rec:
		data := payload.Rec
		k.Logger(ctx).Info("📩 Received LightTx (RECMeta)",
			"creator", msg.Creator,
			"facility_id", data.FacilityId,
			"facility_name", data.FacilityName,
			"location", data.Location,
			"technology_type", data.TechnologyType,
			"capacity_mw", data.CapacityMw,
			"registration_date", data.RegistrationDate,
			"certified_id", data.CertifiedId,
			"issue_date", data.IssueData,
			"generation_start", data.GenerationStartDate,
			"generation_end", data.GenerationEndDate,
			"measured_volume", data.MeasuredVolume_MWh,
			"retired_date", data.RetiredDate,
			"retirement_purpose", data.RetirementPurpose,
			"status", data.Status,
			"timestamp", data.Timestamp,
			"hash", msg.Hash,
			"signature", msg.Signature,
			"pubkey", msg.Pubkey,
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent("light_tx_rec",
				sdk.NewAttribute("creator", msg.Creator),
				sdk.NewAttribute("facility_id", data.FacilityId),
				sdk.NewAttribute("hash", msg.Hash),
				sdk.NewAttribute("signature", msg.Signature),
			),
		)

	default:
		k.Logger(ctx).Error("❌ Unknown payload type in MsgSendLightTx")
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown payload type")
	}

	return &lighttype.MsgSendLightTxResponse{Result: "OK"}, nil
}
