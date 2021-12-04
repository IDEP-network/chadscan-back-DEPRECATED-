package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/everstake/cosmoscan-api/log"
)

func (api *API) GetAggTransactionsFee(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggTransactionsFee)
}

func (api *API) GetAggOperationsCount(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggOperationsCount)
}

func (api *API) GetAvgOperationsPerBlock(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAvgOperationsPerBlock)
}

func (api *API) GetLatestTransactions(w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetLatestTransactions() entered")

        count, ok := mux.Vars(r)["count"]
        if !ok || count == "" {
                jsonBadRequest(w, "invalid count")
                return
        }

        resp, err := api.svc.GetLatestTransactions(count)
        if err != nil {
                log.Error("API GetLatestTransactions: svc.GetLatestTransactions: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)

}
