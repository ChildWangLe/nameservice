package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func resolveNameHandler(cliCtx context.CLIContext,
	storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		paramType := vars[restName]
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/resolve/%s", storeName, paramType))
		if nil != err {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func whoisHandler(cliCtx context.CLIContext,
	storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		paramType := vars[restName]
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/whois/%s", storeName, paramType))
		if nil != err {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func namesHandler(cliCtx context.CLIContext,
	storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/names", storeName))
		if nil != err {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
