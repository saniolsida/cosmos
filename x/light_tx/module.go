package light_tx

import (
	"encoding/json" // ✅ json.RawMessage

	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types" // ✅ sdk.Context, sdk.InvariantRegistry
	"github.com/cosmos/cosmos-sdk/types/module"
	cli "github.com/cosmos/cosmos-sdk/x/light_tx/client/cli"
	"github.com/cosmos/cosmos-sdk/x/light_tx/keeper"
	lighttype "github.com/cosmos/cosmos-sdk/x/light_tx/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return lighttype.ModuleName
}

func (AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	lighttype.RegisterInterfaces(reg)
}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

// 여기서 MsgServer 등록
func (am AppModule) RegisterServices(cfg module.Configurator) {
	lighttype.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (AppModule) ConsensusVersion() uint64 {
	return 1
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(lighttype.DefaultGenesis())
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {

	return []abci.ValidatorUpdate{}
}

func (AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (AppModule) QuerierRoute() string {
	return ""
}

func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(lighttype.RouterKey, NewHandler(am.keeper))
}

func (AppModule) ValidateGenesis(
	cdc codec.JSONCodec,
	txCfg client.TxEncodingConfig,
	bz json.RawMessage,
) error {
	var data lighttype.GenesisState
	return cdc.UnmarshalJSON(bz, &data)
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	lighttype.RegisterLegacyAminoCodec(cdc)
}
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	// 아직 쿼리 기능이 없다면 빈 명령어 반환
	return &cobra.Command{
		Use:   lighttype.ModuleName,
		Short: "Querying commands for the light_tx module",
	}
}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        lighttype.ModuleName,
		Short:                      "Transactions for the light_tx module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(cli.CmdSendLightTx())
	return cmd
}

func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	// return cdc.MustMarshalJSON(lighttype.DefaultGenesis())
	return nil
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// 여기서는 아무 것도 안 해도 됩니다. 기본 동작만 원할 경우 빈 함수로 구현
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}
