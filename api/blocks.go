package api

import (
	"net/http"

	"github.com/everstake/cosmoscan-api/log"
	"github.com/gorilla/mux"
)

func (api *API) GetAggBlocksCount(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggBlocksCount)
}

func (api *API) GetBlocks(w http.ResponseWriter, r *http.Request) {
	api.blockHandler(w, r, api.svc.GetBlocks)
}

func (api *API) GetAggBlocksDelay(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggBlocksDelay)
}

func (api *API) GetAggUniqBlockValidators(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggUniqBlockValidators)
}

func (api *API) GetBlockTransactionCounts(w http.ResponseWriter, r *http.Request) {

	log.Info("GetBlockTransactionCounts() entered")

	height, ok := mux.Vars(r)["height"]
	if !ok || height == "" {
		jsonBadRequest(w, "invalid height")
		return
	}

	resp, err := api.svc.GetBlockTransactionCounts(height)
	if err != nil {
		log.Error("API GetBlockTransactionCounts: svc.GetBlockTransactionCounts: %s", err.Error())
		jsonError(w)
		return
	}

	log.Info("block height = %s", height)

	jsonData(w, resp)

}
