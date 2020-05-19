package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	restName = "name"
)

func RegisterRoutes(cliCtx context.CLIContext,
	r *mux.Router,
	storeName string) {
	for _, e := range []struct {
		Path   string
		Method string
		Func   func(context.CLIContext, string) http.HandlerFunc
	}{
		{fmt.Sprintf("/%s/names", storeName), http.MethodGet, namesHandler},
		{fmt.Sprintf("/%s/names", storeName), http.MethodPost, buyNameHandler},
		{fmt.Sprintf("/%s/names", storeName), http.MethodPut, setNameHandler},
		{fmt.Sprintf("/%s/names/{%s}", storeName, restName), http.MethodGet, resolveNameHandler},
		{fmt.Sprintf("/%s/names/{%s}/whois", storeName, restName), http.MethodGet, whoisHandler},
		{fmt.Sprintf("/%s/names", storeName), http.MethodDelete, deleteNameHandler},
		{fmt.Sprintf("/%s/names/{%s}/can_auction", storeName, restName), http.MethodPost, setCanAuctionHandler},
		{fmt.Sprintf("/%s/names/{%s}/offset", storeName, restName), http.MethodPost, offsetHandler},
	} {
		r.HandleFunc(e.Path, e.Func(cliCtx, restName)).Methods(e.Method)
	}
}
