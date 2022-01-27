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

func (api *API) GetTransactionDetail(w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetTransactionDetail() entered")

        hash, ok := mux.Vars(r)["hash"]
        if !ok || hash == "" {
                jsonBadRequest(w, "invalid hash")
                return
        }

        resp, err := api.svc.GetTransactionDetail(hash)
        if err != nil {
                log.Error("API GetTransactionDetail: svc.GetTransactionDetail: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)

}
