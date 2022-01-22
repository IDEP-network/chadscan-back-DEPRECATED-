package node

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/shopspring/decimal"
	"github.com/everstake/cosmoscan-api/log"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	precision = 6

	DepositPeriodProposalStatus = "PROPOSAL_STATUS_DEPOSIT_PERIOD"
	VotingPeriodProposalStatus  = "PROPOSAL_STATUS_VOTING_PERIOD"
	PassedProposalStatus        = "PROPOSAL_STATUS_PASSED"
	RejectedProposalStatus      = "PROPOSAL_STATUS_REJECTED"
	FailedProposalStatus        = "PROPOSAL_STATUS_FAILED"
)

var PrecisionDiv = decimal.New(1, precision)

type (
	API struct {
		cfg    config.Config
		client *http.Client
	}
	CommunityPool struct {
		Pool []struct {
			Denom  string          `json:"denom"`
			Amount decimal.Decimal `json:"amount"`
		} `json:"pool"`
	}
	Validators struct {
		Validators []Validator `json:"validators"`
	}
	Validator struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Tokens          uint64          `json:"tokens,string"`
		DelegatorShares decimal.Decimal `json:"delegator_shares"`
		Description     struct {
			Moniker  string `json:"moniker"`
			Identity string `json:"identity"`
			Website  string `json:"website"`
			Details  string `json:"details"`
		} `json:"description"`
		UnbondingHeight uint64    `json:"unbonding_height,string"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          decimal.Decimal `json:"rate"`
				MaxRate       decimal.Decimal `json:"max_rate"`
				MaxChangeRate decimal.Decimal `json:"max_change_rate"`
			} `json:"commission_rates"`
		} `json:"commission"`
		MaxChangeRate decimal.Decimal `json:"max_change_rate"`
	}
	Inflation struct {
		Inflation decimal.Decimal `json:"inflation"`
	}
	AmountResult struct {
		Balances []struct {
			Denom  string          `json:"denom"`
			Amount decimal.Decimal `json:"amount"`
		} `json:"balances"`
	}
	StakingPool struct {
		Pool struct {
			NotBondedTokens decimal.Decimal `json:"not_bonded_tokens"`
			BondedTokens    decimal.Decimal `json:"bonded_tokens"`
		} `json:"pool"`
	}
	Supply struct {
		Amount struct {
			Denom  string          `json:"denom"`
			Amount decimal.Decimal `json:"amount"`
		} `json:"amount"`
	}
	StakeResult struct {
		DelegationResponses []struct {
			Delegation struct {
				DelegatorAddress string          `json:"delegator_address"`
				ValidatorAddress string          `json:"validator_address"`
				Shares           decimal.Decimal `json:"shares"`
			} `json:"delegation"`
		} `json:"delegation_responses"`
	}
	UnbondingResult struct {
		UnbondingResponses []struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Entries          []struct {
				Balance decimal.Decimal `json:"balance"`
			} `json:"entries"`
		} `json:"unbonding_responses"`
	}
	ProposalsResult struct {
		Proposals []struct {
			Content struct {
				Type        string `json:"@type"`
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"content"`
			ProposalID       uint64 `json:"proposal_id,string"`
			Status           string `json:"status"`
			FinalTallyResult struct {
				Yes        int64 `json:"yes,string"`
				Abstain    int64 `json:"abstain,string"`
				No         int64 `json:"no,string"`
				NoWithVeto int64 `json:"no_with_veto,string"`
			} `json:"final_tally_result"`
			SubmitTime     time.Time `json:"submit_time"`
			DepositEndTime time.Time `json:"deposit_end_time"`
			TotalDeposit   []struct {
				Amount decimal.Decimal `json:"amount"`
			} `json:"total_deposit"`
			VotingStartTime time.Time `json:"voting_start_time"`
			VotingEndTime   time.Time `json:"voting_end_time"`
		} `json:"proposals"`
	}
	ProposalProposer struct {
		Proposal struct {
			ProposalID uint64 `json:"proposal_id,string"`
			Proposer   string `json:"proposer"`
		} `json:"proposal"`
	}
	DelegatorValidatorStakeResult struct {
		DelegationResponse struct {
			Delegation struct {
				DelegatorAddress string          `json:"delegator_address"`
				ValidatorAddress string          `json:"validator_address"`
				Shares           decimal.Decimal `json:"shares"`
			} `json:"delegation"`
			Balance struct {
				Denom  string          `json:"denom"`
				Amount decimal.Decimal `json:"amount"`
			} `json:"balance"`
		} `json:"delegation_response"`
	}
	ProposalVotersResult struct {
		Result []struct {
			ProposalID uint64 `json:"proposal_id,string"`
			Voter      string `json:"voter"`
			Option     string `json:"option"`
		} `json:"result"`
	}
	ProposalTallyResult struct {
		Tally struct {
			Yes        int64 `json:"yes,string"`
			Abstain    int64 `json:"abstain,string"`
			No         int64 `json:"no,string"`
			NoWithVeto int64 `json:"no_with_veto,string"`
		} `json:"tally"`
	}
	Balance struct {
		Denom string `json:"denom"`
		Amount string `json:"amount"`
	}
	WalletAddressResult struct {
		Balances []Balance `json:"balances"`
		Pagination struct {
			NextKey string `json:"next_key"`
			Total string `json:"total"`
		} `json:"pagination"`
	}
	Signature struct {
		BlockIDFlag int64 `json:"block_id_flag"`
		Signature  string `json:"signature"`
		Timestamp  string `json:"timestamp"`
		ValidatorAddress  string `json:"validator_address"`
	}
	BlockHeightResult struct {
		Block struct {
			Data struct {
				Txs []string `json:"txs"`
				TxsNum int
			} `json:"data"`
			Evidence struct {
				Evidence []string `json:"evidence"`
			} `json:"evidence"`
			Header struct {
				AppHash  string `json:"app_hash"`
				ChainID  string `json:"chain_id"`
				ConsensusHash  string `json:"consensus_hash"`
				DataHash  string `json:"data_hash"`
				EvidenceHash  string `json:"evidence_hash"`
				Height  string `json:"height"`
				LastBlockID struct {
					Hash  string `json:"hash"`
					Parts struct {
						Hash  string `json:"hash"`
						Total int64 `json:"total"`
					} `json:"parts"`
				} `json:"last_block_id"`
				LastCommitHash  string `json:"last_commit_hash"`
				LastResultsHash  string `json:"last_results_hash"`
				NextValidatorsHash  string `json:"next_validators_hash"`
				ProposerAddress  string `json:"proposer_address"`
				Time  string `json:"time"`
				ValidatorsHash  string `json:"validators_hash"`
				Version struct {
					Block  string `json:"block"`
				}`json:"version"`
			} `json:"header"`
			LastCommit struct {
				BlockID struct {
					Hash string `json:"hash"`
					Parts struct {
						Hash string `json:"hash"`
						Total int64 `json:"total"`
					} `json:"block_id"`
				} `json:"block_id"`
				Height  string `json:"height"`
				Round int64 `json:"round"`
				Signatures []Signature `json:"signatures"`
			} `json:"last_commit"`

		} `json:"block"`
		BlockID  struct {
			Hash  string `json:"hash"`
			Parts struct {
				Hash string `json:"hash"`
				Total int64 `json:"total"`
			} `json:"parts"`
		} `json:"block_id"`
	}
	ValidatorMiscResult struct {
		Height string `json:"height"`
		Result struct {
			Commission struct {
				CommissionRates struct {
					MaxChangeRate string `json:"max_change_rate"` 
					MaxRate string `json:"max_rate"`
					Rate string `json:"rate"`
				} `json:"commission_rates"`
				UpdateTime string `json:"update_time"`
			} `json:"commission"`
			ConsensusPubKey struct {
				Type string `json:"type"`
				Value string `json:"value"`
			} `json:"consensus_pubkey"`
			DelegatorShares string `json:"delegator_shares"`
			Description struct {
				Moniker string `json:"moniker"`
			} `json:"description"`
			MinSelfDelegation string `json:"min_self_delegation"`
			OperatorAddress string `json:"operator_address"`
			Status int64 `json:"status"`
			Tokens uint64 `json:"tokens,string"`
			UnboundingTime string `json:"unbounding_time"`
		} `json:"result"`
	}
	ValidatorCommunityPoolResult struct {
		Height string `json:"height"`
		Results []struct {
			Denom string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"result"`
	}
	SelfBondReward struct {
		Denom string `json:"denom"`
		Amount string `json:"amonut"`
	}
	SelfDelegateCommission struct {
		Denom string `json:"denom"`
                Amount string `json:"amonut"`
	}
	ValidatorSelfDelegateResult struct {
		Height string `json:"height"`
		Result struct {
			OperatorAddress string `json:"operator_address"`
			SelfBondRewards []SelfBondReward `json:"self_bond_rewards"`
			ValCommission struct {
				SelfDelegateCommissions []SelfDelegateCommission `json:"commission"`
			} `json:"val_commission"`
		} `json:"result"`
	}
	ValidatorPowerChangeUndelegateResult struct {
		Unbonding_Responses []struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Entries          []struct {
				Balance string `json:"balance"`
				CompletionTime string `json:"completion_time"`
				CreationHeight string `json:"creation_height"`
				InitialBalance string `json:"initial_balance"`
			} `json:"entries"`
		} `json:"unbonding_responses"`
	}
	ValidatorSlashingResult struct {
		Slashes []struct {
			ValidatorPeriod string `json:"validator_period"`
			Fraction string `json:"fraction"`
		} `json:"slashes"`
		Pagination struct {
			NextKey string `json:"next_key"`
			Total   string `json:"total"`
		} `json:"pagination"`
	}
	ValidatorProposerPriorityResult struct {
		Height string `json:"height"`
		Result struct {
			BlockHeight string `json:"block_height"`
			Validators []struct {
                                Address string `json:"address"`
                                PubKey struct {
					Type string `json:"type"`
					Value string `json:"value"`
				} `json:"pub_key"`
				VotingPower string `json:"voting_power"`
				ProposerPriority string `json:"proposer_priority"`
			} `json:"validators"`
		} `json:"result"`
	}
	ValidatorAggInfoResult struct {
		Name string `json:"name"`
		Delegated float64 `json:"delegated"`
		DelegatedPercent float64 `json:"delegated_percent"`
		CommissionPercent float64 `json:"commission_percent"`
		Status int64 `json:"status"`
		UptimePercent float64 `json:"uptime_percent"`
	}

)

func NewAPI(cfg config.Config) *API {
	return &API{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (api API) request(endpoint string, data interface{}) error {
	url := fmt.Sprintf("%s/%s", api.cfg.Parser.Node, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %s", err.Error())
	}
	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %s", err.Error())
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = json.Unmarshal(d, data)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return nil
}

func (api API) GetCommunityPoolAmount() (amount decimal.Decimal, err error) {
	var cp CommunityPool
	err = api.request("cosmos/distribution/v1beta1/community_pool", &cp)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	for _, p := range cp.Pool {
		//if p.Denom == "uxprt" {
		if p.Denom == "idep" {
			amount = amount.Add(p.Amount)
		}
	}
	return amount.Div(PrecisionDiv), nil
}

func (api API) GetValidators() (items []Validator, err error) {
	var validators Validators
	log.Info("node.api.GetValidators() called!")
	err = api.request("cosmos/staking/v1beta1/validators?pagination.limit=10000", &validators)
	if err != nil {
		return nil, fmt.Errorf("request: %s", err.Error())
	}
	return validators.Validators, nil
}

func (api API) GetInflation() (amount decimal.Decimal, err error) {
	var inflation Inflation
	err = api.request("cosmos/mint/v1beta1/inflation", &inflation)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return inflation.Inflation.Mul(decimal.New(100, 0)), nil
}

func (api API) GetTotalSupply() (amount decimal.Decimal, err error) {
	var s Supply
	//err = api.request("/cosmos/bank/v1beta1/supply/uxprt", &s)
	err = api.request("/cosmos/bank/v1beta1/supply/idep", &s)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return s.Amount.Amount.Div(PrecisionDiv), nil
}

func (api API) GetStakingPool() (sp StakingPool, err error) {
	err = api.request("cosmos/staking/v1beta1/pool", &sp)
	if err != nil {
		return sp, fmt.Errorf("request: %s", err.Error())
	}
	sp.Pool.BondedTokens = sp.Pool.BondedTokens.Div(PrecisionDiv)
	sp.Pool.NotBondedTokens = sp.Pool.NotBondedTokens.Div(PrecisionDiv)
	return sp, nil
}

func (api API) GetBalance(address string) (amount decimal.Decimal, err error) {
	var result AmountResult
	err = api.request(fmt.Sprintf("cosmos/bank/v1beta1/balances/%s", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	for _, b := range result.Balances {
		//if b.Denom == "uxprt" {
		if b.Denom == "idep" {
			amount = amount.Add(b.Amount)
		}
	}
	return amount.Div(PrecisionDiv), nil
}

func (api API) GetStake(address string) (amount decimal.Decimal, err error) {
	var result StakeResult
	err = api.request(fmt.Sprintf("cosmos/staking/v1beta1/delegations/%s?pagination.limit=10000", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	shares := decimal.Zero
	for _, r := range result.DelegationResponses {
		shares = shares.Add(r.Delegation.Shares)
	}
	return shares.Div(PrecisionDiv), nil
}

func (api API) GetUnbonding(address string) (amount decimal.Decimal, err error) {
	var result UnbondingResult
	err = api.request(fmt.Sprintf("cosmos/staking/v1beta1/delegators/%s/unbonding_delegations?pagination.limit=10000", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	for _, r := range result.UnbondingResponses {
		for _, entry := range r.Entries {
			amount = amount.Add(entry.Balance)
		}
	}
	amount = amount.Div(PrecisionDiv)
	return amount, nil
}

func (api API) GetProposals() (proposals ProposalsResult, err error) {
	err = api.request("cosmos/gov/v1beta1/proposals?pagination.limit=10000", &proposals)
	if err != nil {
		return proposals, fmt.Errorf("request: %s", err.Error())
	}
	return proposals, nil
}

func (api API) GetDelegatorValidatorStake(delegator string, validator string) (amount decimal.Decimal, err error) {
	var result DelegatorValidatorStakeResult
	err = api.request(fmt.Sprintf("cosmos/staking/v1beta1/validators/%s/delegations/%s", validator, delegator), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return result.DelegationResponse.Delegation.Shares.Div(PrecisionDiv), nil
}

func (api API) ProposalTallyResult(id uint64) (result ProposalTallyResult, err error) {
	err = api.request(fmt.Sprintf("/cosmos/gov/v1beta1/proposals/%d/tally", id), &result)
	if err != nil {
		return result, fmt.Errorf("request: %s", err.Error())
	}
	return result, nil
}

func (api API) GetBlockHeight(blockHeight string) (result BlockHeightResult, err error) {

	err = api.request(fmt.Sprintf("blocks/%s", blockHeight), &result)
        if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }
        return result, nil
}

func (api API) GetWalletAddress(walletAddr string) (result WalletAddressResult, err error) {

	err = api.request(fmt.Sprintf("cosmos/bank/v1beta1/balances/%s", walletAddr), &result)
        if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }
        return result, nil
}

func (api API) GetValidatorMisc(validatorAddr string) (result ValidatorMiscResult, err error) {

        err = api.request(fmt.Sprintf("staking/validators/%s", validatorAddr), &result)
        if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }
        return result, nil
}

func (api API) GetValidatorDelegations(validatorAddr string) (result DelegatorValidatorStakeResult, err error) {

        err = api.request(fmt.Sprintf("staking/validators/%s/delegations", validatorAddr), &result)
        if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }
        return result, nil
}

func (api API) GetValidatorCommunityPool() (result ValidatorCommunityPoolResult, err error) {

	err = api.request("distribution/community_pool", &result)
	if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }
        return result, nil

}

func (api API) GetValidatorSelfDelegate(validatorAddr string) (result ValidatorSelfDelegateResult, err error) {

        err = api.request(fmt.Sprintf("distribution/validators/%s", validatorAddr), &result)
	if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }
        return result, nil

}

func (api API) GetValidatorPowerChangeUndelegate(validatorAddr string) (result ValidatorPowerChangeUndelegateResult, err error) {

	err = api.request(fmt.Sprintf("cosmos/staking/v1beta1/validators/%s/unbonding_delegations", validatorAddr), &result)
        if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }
        return result, nil

}

func (api API) GetValidatorSlashing(validatorAddr string) (result ValidatorSlashingResult, err error) {

        err = api.request(fmt.Sprintf("cosmos/distribution/v1beta1/validators/%s/slashes", validatorAddr), &result)
        if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }
        return result, nil

}

func (api API) GetValidatorProposerPriority(validatorAddr string) (result ValidatorProposerPriorityResult, err error) {

	var validatorMiscResult  ValidatorMiscResult
	var validatorProposerPriorityResult ValidatorProposerPriorityResult

        err = api.request(fmt.Sprintf("staking/validators/%s", validatorAddr), &validatorMiscResult)
        if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }

	err = api.request("validatorsets/latest", &result)
        if err != nil {
                return result, fmt.Errorf("request: %s", err.Error())
        }

	//return only the entry that is having consensus key address = validatorMiscResult.result.consensus_pubkey.value

	//log.Info("result = %s", result)
	log.Info("block height = %s", result.Result.BlockHeight)

	consensus_pkey := validatorMiscResult.Result.ConsensusPubKey.Value
	log.Info("consensus_pkey = %s", consensus_pkey)

        for _, v := range result.Result.Validators {
                if v.PubKey.Value == consensus_pkey {
			log.Info("Found matched consensus_pkey!")
			validatorProposerPriorityResult.Result.Validators = append(validatorProposerPriorityResult.Result.Validators, v)
			validatorProposerPriorityResult.Result.BlockHeight = result.Result.BlockHeight
			validatorProposerPriorityResult.Height = result.Height
			break
		}
        }

	return validatorProposerPriorityResult, nil

}
