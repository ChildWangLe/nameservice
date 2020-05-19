package nameservice

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	WhoisRecords []Whois `json:"whois_records"`
}

func NewGenesisState(WhoisRecords []Whois) GenesisState {
	return GenesisState{
		WhoisRecords: WhoisRecords,
	}
}

func ValidateGenesis(gs GenesisState) error {
	for _, e := range gs.WhoisRecords {
		if nil == e.Owner {
			return fmt.Errorf(`invalid "Whois" record: value%s, error: missing owner`,
				e.Value)
		}
		if 0 >= len(e.Value) {
			return fmt.Errorf(`invalid "Whois" record: owner: %s, error: missing value`,
				e.Owner)
		}
		if nil == e.Price {
			return fmt.Errorf(`invalid "Whois" record: value: %s, error: missing price`,
				e.Value)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		WhoisRecords: []Whois{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, gs GenesisState) {
	for _, e := range gs.WhoisRecords {
		keeper.SetWhois(ctx, e.Value, e)
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	var arr []Whois
	for it := keeper.GetNamesIterator(ctx); it.Valid(); it.Next() {
		name := string(it.Key())
		whois := keeper.GetWhois(ctx, name)
		arr = append(arr, whois)
	}
	return GenesisState{WhoisRecords: arr}
}
