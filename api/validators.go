package api

import (
	"github.com/everstake/cosmoscan-api/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetTopProposedBlocksValidators(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetTopProposedBlocksValidators()
	if err != nil {
		log.Error("API GetTopProposedBlocksValidators: svc.GetTopProposedBlocksValidators: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

func (api *API) GetMostJailedValidators(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetMostJailedValidators()
	if err != nil {
		log.Error("API GetMostJailedValidators: svc.GetMostJailedValidators: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

func (api *API) GetFeeRanges(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetFeeRanges()
	if err != nil {
		log.Error("API GetFeeRanges: svc.GetFeeRanges: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

func (api *API) GetValidators(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetValidators()
	if err != nil {
		log.Error("API GetValidators: svc.GetValidators: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

func (api *API) GetValidatorsDelegatorsTotal(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetValidatorsDelegatorsTotal()
	if err != nil {
		log.Error("API GetValidatorsDelegatorsTotal: svc.GetValidatorsDelegatorsTotal: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetValidator(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	resp, err := api.svc.GetValidator(address)
	if err != nil {
		log.Error("API GetValidator: svc.GetValidator: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetValidatorBalance(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	resp, err := api.svc.GetValidatorBalance(address)
	if err != nil {
		log.Error("API GetValidatorBalance: svc.GetValidatorBalance: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetValidatorBlocksStat(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	resp, err := api.svc.GetValidatorBlocksStat(address)
	if err != nil {
		log.Error("API GetValidatorBlocksStat: svc.GetValidatorBlocksStat: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetValidatorMisc (w http.ResponseWriter, r *http.Request) {
	log.Info("api.GetValidatorMisc() entered!")

	address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

	resp, err := api.svc.GetValidatorMisc(address)
        if err != nil {
                log.Error("API GetValidatorMisc: svc.GetValidatorMisc: %s", err.Error())
                jsonError(w)
                return
        }

	//jsonBadRequest(w, "GetValidatorMisc is incomplete ")

	jsonData(w, resp)
}

func (api *API) GetValidatorDelegations (w http.ResponseWriter, r *http.Request) {
        log.Info("api.GetValidatoDelegations() entered!")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

	resp, err := api.svc.GetValidatorDelegations(address)
        if err != nil {
                log.Error("API GetValidatorDelegations: svc.GetValidatorDelegations: %s", err.Error())
                jsonError(w)
                return
        }

        //jsonBadRequest(w, "GetValidatoDelegations is incomplete ")

	jsonData(w, resp)
}

func (api *API) GetValidatorGovernance (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorGovernance() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

	resp, err := api.svc.GetValidatorGovernance(address)
        if err != nil {
                log.Error("API GetValidatorGovernance: svc.GetValidatorGovernance: %s", err.Error())
                jsonError(w)
                return
        }

	//jsonBadRequest(w, "GetValidatorGovernance is incomplete ")

        jsonData(w, resp)

}

func (api *API) GetValidatorTransfer (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorTransfer() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

	resp, err := api.svc.GetValidatorTransfer(address)
        if err != nil {
                log.Error("API GetValidatorTransfer: svc.GetValidatorTransfer: %s", err.Error())
                jsonError(w)
                return
        }

        //jsonBadRequest(w, "GetValidatorTransfer is incomplete ")

        jsonData(w, resp)

}

func (api *API) GetValidatorCommunityPool (w http.ResponseWriter, r *http.Request) {

	log.Info("api.GetValidatorCommunityPool() entered")

	resp, err := api.svc.GetValidatorCommunityPool()
        if err != nil {
                log.Error("API GetValidatorCommunityPool: svc.GetValidatorCommunityPool: %s", err.Error())
                jsonError(w)
                return
        }

	//jsonBadRequest(w, "GetValidatorCommunityPool is incomplete ")

	jsonData(w, resp)
}

func (api *API) GetValidatorSelfDelegate (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorSelfDelegate() entered")

	address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorSelfDelegate(address)
        if err != nil {
                log.Error("API GetValidatorSelfDelegate: svc.GetValidatorSelfDelegate: %s", err.Error())
                jsonError(w)
                return
        }

        //jsonBadRequest(w, "GetValidatorSelfDelegate is incomplete ")

        jsonData(w, resp)
}

func (api *API) GetValidatorPowerChangeDelegate (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorPowerChangeDelegate() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorPowerChangeDelegate(address)
        if err != nil {
                log.Error("API GetValidatorPowerChangeDelegate: svc.GetValidatorPowerChangeDelegate: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorPowerChangeUndelegate (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorPowerChangeUndelegate() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorPowerChangeUndelegate(address)
        if err != nil {
                log.Error("API GetValidatorPowerChangeUndelegate: svc.GetValidatorPowerChangeUndelegate: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorDistribution (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorDistribution() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorDistribution(address)
        if err != nil {
                log.Error("API GetValidatorDistribution: svc.GetValidatorDistribution: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorStaking (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorStaking() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorStaking(address)
        if err != nil {
                log.Error("API GetValidatorStaking: svc.GetValidatorStaking: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorLast100Blocks (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorLast100Blocks() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorLast100Blocks(address)
        if err != nil {
                log.Error("API GetValidatorLast100Blocks: svc.GetValidatorLast100Blocks: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorMissedBlocks (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorMissedBlocks() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorMissedBlocks(address)
        if err != nil {
                log.Error("API GetValidatorMissedBlocks: svc.GetValidatorMissedBlocks: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorSlashing (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorSlashing() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorSlashing(address)
        if err != nil {
                log.Error("API GetValidatorSlashing: svc.GetValidatorSlashing: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorProposerPriority (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorProposerPriority() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorProposerPriority(address)
        if err != nil {
                log.Error("API GetValidatorProposerPriority: svc.GetValidatorProposerPrioirty: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorUptimePercent (w http.ResponseWriter, r *http.Request) {

        log.Info("api.GetValidatorUptimePercent() entered")

        address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

        resp, err := api.svc.GetValidatorUptimePercent(address)
        if err != nil {
                log.Error("API GetValidatorUptimePercent: svc.GetValidatorUptimePercent: %s", err.Error())
                jsonError(w)
                return
        }

        jsonData(w, resp)
}

func (api *API) GetValidatorAggInfo (w http.ResponseWriter, r *http.Request) {

	log.Info("api.GetValidatorAggInfo() entered")

	address, ok := mux.Vars(r)["address"]
        if !ok || address == "" {
                jsonBadRequest(w, "invalid address")
                return
        }

	resp, err := api.svc.GetValidatorAggInfo(address)
        if err != nil {
                log.Error("API GetValidatorAggInfo: svc.GetValidatorAggInfo: %s", err.Error())
                jsonError(w)
                return
        }

	jsonData(w, resp)
}
