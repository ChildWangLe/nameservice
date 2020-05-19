package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/neo.wang/nameservice/x/nameservice/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryResolve = "resolve"
	QueryWhois   = "whois"
	QueryNames   = "names"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResolve:
			return handleQueryResolve(ctx, path[1:], req, keeper)
		case QueryWhois:
			return handleQueryWhois(ctx, path[1:], req, keeper)
		case QueryNames:
			return handleQueryNames(ctx, req, keeper)
		}
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown NameService query endpoint")
	}
}

func handleQueryResolve(ctx sdk.Context,
	path []string,
	req abci.RequestQuery,
	keeper Keeper) ([]byte, error) {
	value := keeper.ResolveName(ctx, path[0])
	if 0 >= len(value) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "could not resolve name")
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, types.QueryResResolve{Value: value})
	if nil != err {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	fmt.Println(string(res))
	return res, nil
}

func handleQueryWhois(ctx sdk.Context,
	path []string,
	req abci.RequestQuery,
	keeper Keeper) ([]byte, error) {
	whois := keeper.GetWhois(ctx, path[0])
	res, err := codec.MarshalJSONIndent(keeper.cdc, whois)
	if nil != err {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func handleQueryNames(ctx sdk.Context,
	req abci.RequestQuery,
	keeper Keeper) ([]byte, error) {
	var nameList types.QueryResNames
	for it := keeper.GetNamesIterator(ctx); it.Valid(); it.Next() {
		nameList = append(nameList, string(it.Key()))
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, nameList)
	if nil != err {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
