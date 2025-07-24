package mymodule

import (
	"context"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	abci "github.com/tendermint/tendermint/abci/types"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/mymodule/keeper"
	"github.com/cosmos/cosmos-sdk/x/mymodule/types"
	mymodulev1beta1 "github.com/cosmos/cosmos-sdk/x/mymodule/types/cosmos/mymodule/v1beta1"
)

// AppModuleBasic implements the AppModuleBasic interface
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (AppModuleBasic) ValidateGenesis(
	cdc codec.JSONCodec,
	txEncodingCfg client.TxEncodingConfig,
	bz json.RawMessage,
) error {
	// 검증 로직이 없으면 nil만 반환해도 됩니다
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	// 여기에 REST 라우트 등록 (필요 없으면 비워둬도 됨)
}

// 반드시 AppModuleBasic 구조체에 추가해야 함!
func (AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (am AppModule) Route() sdk.Route {
	// 현재 메시지 기반 라우팅을 사용하지 않으면 빈 Route 리턴
	return sdk.NewRoute(types.ModuleName, nil)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := mymodulev1beta1.RegisterQueryHandlerClient(
		context.Background(),
		mux,
		mymodulev1beta1.NewQueryClient(clientCtx),
	); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) GetTxCmd() *cobra.Command    { return nil }
func (AppModuleBasic) GetQueryCmd() *cobra.Command { return nil }

type Keeper = keeper.Keeper

// AppModule implements the AppModule interface
type AppModule struct {
	AppModuleBasic
	Keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		Keeper:         k,
	}
}

// LegacyQuerierHandler implements the AppModule interface but is unused (legacy support)
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) Name() string { return types.ModuleName }

func (am AppModule) RegisterServices(cfg module.Configurator) {
	mymodulev1beta1.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.Keeper))
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// 반드시 v0.45에서 요구하는 ConsensusVersion 추가
func (am AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) QuerierRoute() string {
	return ""
}

func (am AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}
