package reward

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/reward/client/cli"
	"github.com/cosmos/cosmos-sdk/x/reward/keeper"
	"github.com/cosmos/cosmos-sdk/x/reward/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

// -----------------------
// AppModuleBasic
// -----------------------

var _ module.AppModuleBasic = AppModuleBasic{}

type AppModuleBasic struct{}

var ModuleAccountPermissions = map[string][]string{
	types.ModuleName: {authtypes.Minter},
}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	// return cdc.MustMarshalJSON(types.DefaultGenesis())
	return nil
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router)              {}
func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

// -----------------------
// AppModule
// -----------------------

type AppModule struct {
	AppModuleBasic
	keeper        keeper.Keeper
	accountKeeper authkeeper.AccountKeeper
}

func NewAppModule(k keeper.Keeper, ak authkeeper.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		accountKeeper:  ak,
	}
}

func (AppModule) ConsensusVersion() uint64 {
	return 1
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	// gRPC 서비스를 등록하지 않는 경우, 이 메서드는 비워두어도 됩니다.
}

func (am AppModule) Name() string {
	return types.ModuleName
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

func (am AppModule) QuerierRoute() string                                  { return "" }
func (am AppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier { return nil }

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	logger := ctx.Logger()

	// 모듈 계정이 존재하지 않으면 생성
	account := am.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	if account == nil {
		moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)
		am.accountKeeper.SetModuleAccount(ctx, moduleAcc)
		logger.Info("✅ Reward 모듈 계정 생성 완료", "module", types.ModuleName)
	} else {
		logger.Info("ℹ️ Reward 모듈 계정이 이미 존재함", "module", types.ModuleName)
	}

	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the current state as types.GenesisState
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
