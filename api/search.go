package api


import (
	"github.com/everstake/cosmoscan-api/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetWalletAddress(w http.ResponseWriter, r *http.Request) {

	log.Info("api.GetWalletAddress() entered")

	address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetWalletAddress(address)
        if err != nil {
                log.Error("API GetWalletAddress: svc.GetWalletAddress: %s", err.Error())
                jsonError(w)
                return
        }

	//jsonBadRequest(w, "GetWalletAddress is incomplete ")

	jsonData(w, resp)
}

func (api *API) GetTransactionHash(w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetTransactionHash() entered")

        hash, ok := mux.Vars(r)["hash"]
        if !ok || hash == "" {
                jsonBadRequest(w, "invalid hash")
                return
        }

	log.Info("transaction hash = %s", hash)

        resp, err := api.svc.GetTransactionHash(hash)
        if err != nil {
                log.Error("API GetTransactionHash: svc.GetTransactionHash: %s", err.Error())
                jsonError(w)
                return
        }

        //jsonBadRequest(w, "GetTransactionHash is incomplete ")
        jsonData(w, resp)

}


func (api *API) GetBlockHash(w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetBlockHash() entered")

        hash, ok := mux.Vars(r)["hash"]
        if !ok || hash == "" {
                jsonBadRequest(w, "invalid hash")
                return
        }

	log.Info("block hash = %s", hash)

        resp, err := api.svc.GetBlockHash(hash)
        if err != nil {
                log.Error("API GetBlockHash: svc.GetBlockHash: %s", err.Error())
                jsonError(w)
                return
        }

        //jsonBadRequest(w, "GetBlockHash is incomplete ")
        jsonData(w, resp)

}


func (api *API) GetBlockHeight(w http.ResponseWriter, r *http.Request) {

	log.Info("GetBlockHeight() entered")

	height, ok := mux.Vars(r)["height"]
	if !ok || height == "" {
		jsonBadRequest(w, "invalid height")
		return
	}

	resp, err := api.svc.GetBlockHeight(height)
	if err != nil {
		log.Error("API GetBlockHeight: svc.GetBlockHeight: %s", err.Error())
		jsonError(w)
		return
	}

	log.Info("block height = %s", height)
	//jsonBadRequest(w, "GetBlockHeight is incomplete ")

	jsonData(w, resp)

}

