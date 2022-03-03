package api

import (
	"net/http"

	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/gorilla/mux"
)

func (api *API) GetAggBlocksCount(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggBlocksCount)
}

func (api *API) GetBlocks(w http.ResponseWriter, r *http.Request) {
	var filter filters.Blocks
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	if filter.Limit == 0 {
		filter.Limit = 100
	}
	if filter.Limit  > 1000 {
		filter.Limit = 1000
	}
	resp, err := api.svc.GetBlocks(filter)
	if err != nil {
		log.Error("API GetBlocks: svc.GetBlocks: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
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
