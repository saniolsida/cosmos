package mapping

import (
	"context"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/mapping/client/cli"
	"github.com/cosmos/cosmos-sdk/x/mapping/client/rest"
	"github.com/cosmos/cosmos-sdk/x/mapping/keeper"
	"github.com/cosmos/cosmos-sdk/x/mapping/types"
)

// AppModuleBasic implements the basic methods for the module
type AppModuleBasic struct{}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mapping",
		Short: "Mapping module transactions",
	}

	cmd.AddCommand(
		cli.CmdRegisterDevice(),
	)

	return cmd
}

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.ModuleName)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(
		context.Background(),
		mux,
		types.NewQueryClient(clientCtx),
	)
	if err != nil {
		panic(err)
	}
}

func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(clientCtx, rtr)
}

// Name returns the module name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers legacy amino codec
func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {}

// RegisterLegacyAminoCodec is unused but required
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers interface types
func (AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns default genesis state as raw bytes
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (AppModuleBasic) ValidateGenesis(
	cdc codec.JSONCodec,
	_ client.TxEncodingConfig, // 사용하지 않으면 '_' 로 무시 가능
	bz json.RawMessage,
) error {
	return nil
}

// AppModule implements the full module
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule
func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

// Name returns the module name
func (am AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants does nothing
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns an empty route (for legacy)
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, nil)
}

// QuerierRoute returns module's query route name
func (am AppModule) QuerierRoute() string { return types.QuerierRoute }

// LegacyQuerierHandler returns nil (not used in new version)
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers msg service handlers
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// InitGenesis initializes state from genesis file
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {

	return []abci.ValidatorUpdate{}
}

// ConsensusVersion returns module version
func (AppModule) ConsensusVersion() uint64 {
	return 1
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	// BeginBlock 로직 없으면 비워도 됩니다.
}

func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
