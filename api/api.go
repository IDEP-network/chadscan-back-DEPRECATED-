package api

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type API struct {
	dao          dao.DAO
	cfg          config.Config
	svc          services.Services
	router       *mux.Router
	queryDecoder *schema.Decoder
}

type errResponse struct {
	Error string `json:"error"`
	Msg   string `json:"msg"`
}

func NewAPI(cfg config.Config, svc services.Services, dao dao.DAO) *API {
	sd := schema.NewDecoder()
	sd.IgnoreUnknownKeys(true)
	sd.RegisterConverter(dmodels.Time{}, func(s string) reflect.Value {
		timestamp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}
		}
		t := dmodels.NewTime(time.Unix(timestamp, 0))
		return reflect.ValueOf(t)
	})
	return &API{
		cfg:          cfg,
		dao:          dao,
		svc:          svc,
		queryDecoder: sd,
	}
}

func (api *API) Title() string {
	return "API"
}

func (api *API) Run() error {
	api.router = mux.NewRouter()
	api.loadRoutes()

	http.Handle("/", api.router)
	log.Info("Listen API server on %s port", api.cfg.API.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", api.cfg.API.Port), nil)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) Stop() error {
	return nil
}

func (api *API) loadRoutes() {

	api.router = mux.NewRouter()

	api.router.
		PathPrefix("/static").
		Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./resources/static"))))

	wrapper := negroni.New()

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins:   api.cfg.API.AllowedHosts,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-User-Env", "Sec-Fetch-Mode"},
	}))

	// public
	HandleActions(api.router, wrapper, "", []*Route{
		{Path: "/", Method: http.MethodGet, Func: api.Index},
		{Path: "/health", Method: http.MethodGet, Func: api.Health},
		{Path: "/api", Method: http.MethodGet, Func: api.GetSwaggerAPI},

		{Path: "/meta", Method: http.MethodGet, Func: api.GetMetaData},
		{Path: "/historical-state", Method: http.MethodGet, Func: api.GetHistoricalState},
		{Path: "/transactions/fee/agg", Method: http.MethodGet, Func: api.GetAggTransactionsFee},
		{Path: "/transfers/volume/agg", Method: http.MethodGet, Func: api.GetAggTransfersVolume},
		{Path: "/operations/count/agg", Method: http.MethodGet, Func: api.GetAggOperationsCount},
		{Path: "/block/hash/{hash}", Method: http.MethodGet, Func: api.GetBlockHash},
		{Path: "/block/height/{height}", Method: http.MethodGet, Func: api.GetBlockHeight},
		{Path: "/blocks/count/agg", Method: http.MethodGet, Func: api.GetAggBlocksCount},
		{Path: "/blocks/delay/agg", Method: http.MethodGet, Func: api.GetAggBlocksDelay},
		{Path: "/blocks/validators/uniq/agg", Method: http.MethodGet, Func: api.GetAggUniqBlockValidators},
		{Path: "/blocks/operations/agg", Method: http.MethodGet, Func: api.GetAvgOperationsPerBlock},
		{Path: "/delegations/volume/agg", Method: http.MethodGet, Func: api.GetAggDelegationsVolume},
		{Path: "/undelegations/volume/agg", Method: http.MethodGet, Func: api.GetAggUndelegationsVolume},
		{Path: "/unbonding/volume/agg", Method: http.MethodGet, Func: api.GetAggUnbondingVolume},
		{Path: "/bonded-ratio/agg", Method: http.MethodGet, Func: api.GetAggBondedRatio},
		{Path: "/network/stats", Method: http.MethodGet, Func: api.GetNetworkStats},
		{Path: "/staking/pie", Method: http.MethodGet, Func: api.GetStakingPie},
		{Path: "/proposals", Method: http.MethodGet, Func: api.GetProposals},
		{Path: "/proposals/votes", Method: http.MethodGet, Func: api.GetProposalVotes},
		{Path: "/proposals/deposits", Method: http.MethodGet, Func: api.GetProposalDeposits},
		{Path: "/proposals/chart", Method: http.MethodGet, Func: api.GetProposalChartData},
		{Path: "/transaction/hash/{hash}", Method: http.MethodGet, Func: api.GetTransactionHash},
		{Path: "/validators", Method: http.MethodGet, Func: api.GetValidators},
		{Path: "/validators/33power/agg", Method: http.MethodGet, Func: api.GetAggValidators33Power},
		{Path: "/validators/top/proposed", Method: http.MethodGet, Func: api.GetTopProposedBlocksValidators},
		{Path: "/validators/top/jailed", Method: http.MethodGet, Func: api.GetMostJailedValidators},
		{Path: "/validators/fee/ranges", Method: http.MethodGet, Func: api.GetFeeRanges},
		{Path: "/validators/delegators/total", Method: http.MethodGet, Func: api.GetValidatorsDelegatorsTotal},
		{Path: "/accounts/whale/agg", Method: http.MethodGet, Func: api.GetAggWhaleAccounts},
		{Path: "/validator/{address}/balance", Method: http.MethodGet, Func: api.GetValidatorBalance},
		{Path: "/validator/{address}/delegations/agg", Method: http.MethodGet, Func: api.GetValidatorDelegationsAgg},
		{Path: "/validator/{address}/delegators/agg", Method: http.MethodGet, Func: api.GetValidatorDelegatorsAgg},
		{Path: "/validator/{address}/blocks/stats", Method: http.MethodGet, Func: api.GetValidatorBlocksStat},
		{Path: "/validator/{address}", Method: http.MethodGet, Func: api.GetValidator},
		{Path: "/validator/{address}/delegators", Method: http.MethodGet, Func: api.GetValidatorDelegators},
		{Path: "/wallet/address/{address}", Method: http.MethodGet, Func: api.GetWalletAddress},
		{Path: "/validator/misc/{address}", Method: http.MethodGet, Func: api.GetValidatorMisc},
		{Path: "/validator/delegations/{address}", Method: http.MethodGet, Func: api.GetValidatorDelegations},
		{Path: "/validator/governance/{address}", Method: http.MethodGet, Func: api.GetValidatorGovernance},
		{Path: "/validator/transfer/{address}", Method: http.MethodGet, Func: api.GetValidatorTransfer},
		{Path: "/validator/staking/{address}", Method: http.MethodGet, Func: api.GetValidatorStaking},
		{Path: "/validator/detail/communitypool", Method: http.MethodGet, Func: api.GetValidatorCommunityPool},
		{Path: "/validator/selfdelegate/{address}", Method: http.MethodGet, Func: api.GetValidatorSelfDelegate},
		{Path: "/validator/powerchange/delegate/{address}", Method: http.MethodGet, Func: api.GetValidatorPowerChangeDelegate},
		{Path: "/validator/powerchange/undelegate/{address}", Method: http.MethodGet, Func: api.GetValidatorPowerChangeUndelegate},
		{Path: "/validator/distribution/{address}", Method: http.MethodGet, Func: api.GetValidatorDistribution},
		{Path: "/validator/last100blocks/{address}", Method: http.MethodGet, Func: api.GetValidatorLast100Blocks},
		{Path: "/validator/missedblocks/{address}", Method: http.MethodGet, Func: api.GetValidatorMissedBlocks},
		{Path: "/validator/slashing/{address}", Method: http.MethodGet, Func: api.GetValidatorSlashing},
		{Path: "/validator/proposerpriority/{address}", Method: http.MethodGet, Func: api.GetValidatorProposerPriority},
		{Path: "/latest/transactions/{count}", Method: http.MethodGet, Func: api.GetLatestTransactions},
	})

}

func jsonData(writer http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(bytes)
}

func jsonError(writer http.ResponseWriter) {
	writer.WriteHeader(500)
	bytes, err := json.Marshal(errResponse{
		Error: "service_error",
	})
	if err != nil {
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(bytes)
}

func jsonBadRequest(writer http.ResponseWriter, msg string) {
	bytes, err := json.Marshal(errResponse{
		Error: "bad_request",
		Msg:   msg,
	})
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(400)
	writer.Write(bytes)
}

func (api *API) GetSwaggerAPI(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("./resources/templates/swagger.html")
	if err != nil {
		log.Error("GetSwaggerAPI: ReadFile", zap.Error(err))
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Error("GetSwaggerAPI: Write", zap.Error(err))
		return
	}
}
