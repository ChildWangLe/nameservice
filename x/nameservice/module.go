package nameservice

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gorilla/mux"
	"github.com/neo.wang/nameservice/x/nameservice/client/cli"
	"github.com/neo.wang/nameservice/x/nameservice/client/rest"
	"github.com/neo.wang/nameservice/x/nameservice/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModuleBasic struct{}

func (object AppModuleBasic) Name() string {
	return ModuleName
}

func (object AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (object AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (object AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var gs GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &gs)
	if nil != err {
		return err
	}
	return ValidateGenesis(gs)
}

func (object AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, StoreKey)
}

func (object AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

func (object AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

type AppModule struct {
	AppModuleBasic
	keeper     Keeper
	bankKeeper bank.Keeper
}

func NewAppModule(keeper Keeper, bankKeeper bank.Keeper) AppModule {
	return AppModule{
		keeper:     keeper,
		bankKeeper: bankKeeper,
	}
}

func (object AppModule) Name() string {
	return ModuleName
}

func (object AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (object AppModule) Route() string {
	return RouteKey
}

func (object AppModule) NewHandler() sdk.Handler {
	return NewHandler(object.keeper)
}

func (object AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (object AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(object.keeper)
}

func (object AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	object.keeper.CountdownAllCanAuction(ctx, req.Header.Height, func(m map[string]types.Whois) {
		for name, whois := range m {
			whois.Owner = whois.AuctionUser
			whois.CanAuction = false
			whois.Price = whois.AuctionPrice
			whois.AuctionDealBlockHeight = 0
			whois.AuctionUser = sdk.AccAddress{}
			whois.AuctionPrice = sdk.Coins{}
			object.keeper.SetWhois(ctx, name, whois)
		}
	})
}

func (object AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (object AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var gs GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &gs)
	InitGenesis(ctx, object.keeper, gs)
	return []abci.ValidatorUpdate{}
}

func (object AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, object.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}
