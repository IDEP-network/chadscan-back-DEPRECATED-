package services

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	//"../cosmos-sdk/types"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/helpers"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/persistenceOne/persistenceCore/application"
	//"persistenceCore/application"
	"github.com/shopspring/decimal"
	"sort"
	"strconv"
	"time"
)

const validatorsMapCacheKey = "validators_map"
const validatorsCacheKey = "validators"
const mostJailedValidators = "mostJailedValidators"

func (s *ServiceFacade) UpdateValidatorsMap() {
	mp, err := s.makeValidatorMap()
	if err != nil {
		log.Error("UpdateValidatorsMap: makeValidatorMap: %s", err.Error())
		return
	}
	s.dao.CacheSet(validatorsMapCacheKey, mp, time.Minute*30)
}

func (s *ServiceFacade) GetValidatorMap() (map[string]node.Validator, error) {
	data, found := s.dao.CacheGet(validatorsMapCacheKey)
	if found {
		return data.(map[string]node.Validator), nil
	}
	mp, err := s.makeValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("makeValidatorMap: %s", err.Error())
	}
	return mp, nil
}

func (s *ServiceFacade) makeValidatorMap() (map[string]node.Validator, error) {
	mp := make(map[string]node.Validator)
	validators, err := s.node.GetValidators()
	if err != nil {
		return nil, fmt.Errorf("node.GetValidators: %s", err.Error())
	}
	for _, validator := range validators {
		mp[validator.OperatorAddress] = validator
	}
	return mp, nil
}

func (s *ServiceFacade) GetStakingPie() (pie smodels.Pie, err error) {
	stakingPool, err := s.node.GetStakingPool()
	if err != nil {
		return pie, fmt.Errorf("node.GetStakingPool: %s", err.Error())
	}
	pie.Total = stakingPool.Pool.BondedTokens
	validatorsMap, err := s.GetValidatorMap()
	if err != nil {
		return pie, fmt.Errorf("s.GetValidatorMap: %s", err.Error())
	}
	var validators []node.Validator
	for _, v := range validatorsMap {
		validators = append(validators, v)
	}
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].DelegatorShares.GreaterThan(validators[j].DelegatorShares)
	})
	if len(validators) < 20 {
		return pie, fmt.Errorf("not enought validators")
	}
	parts := make([]smodels.PiePart, 20)
	for i := 0; i < 20; i++ {
		parts[i] = smodels.PiePart{
			Label: validators[i].OperatorAddress,
			Title: validators[i].Description.Moniker,
			Value: validators[i].DelegatorShares.Div(node.PrecisionDiv),
		}
	}
	pie.Parts = parts
	return pie, nil
}

func (s *ServiceFacade) GetValidators() (validators []smodels.Validator, err error) {
	data, found := s.dao.CacheGet(validatorsCacheKey)
	if found {
		return data.([]smodels.Validator), nil
	}
	return nil, fmt.Errorf("not found in cache")
}

func (s *ServiceFacade) UpdateValidators() {
	validators, err := s.makeValidators()
	if err != nil {
		log.Error("UpdateValidators: makeValidators: %s", err.Error())
		return
	}
	s.dao.CacheSet(validatorsCacheKey, validators, time.Hour)
}

func (s *ServiceFacade) makeValidators() (validators []smodels.Validator, err error) {
	nodeValidators, err := s.node.GetValidators()
	if err != nil {
		return nil, fmt.Errorf("node.GetValidators: %s", err.Error())
	}
	stakingPool, err := s.node.GetStakingPool()
	if err != nil {
		return nil, fmt.Errorf("node.GetStakingPool: %s", err.Error())
	}
	for _, v := range nodeValidators {
		consAddress, err := helpers.GetHexAddressFromBase64PK(v.ConsensusPubkey.Key)
		if err != nil {
			return nil, fmt.Errorf("helpers.GetHexAddressFromBase64PK: %s", err.Error())
		}
		blockProposed, err := s.dao.GetProposedBlocksTotal(filters.BlocksProposed{Proposers: []string{consAddress}})
		if err != nil {
			return nil, fmt.Errorf("dao.GetProposedBlocksTotal: %s", err.Error())
		}

		addressBytes, err := types.GetFromBech32(v.OperatorAddress, application.Bech32PrefixValAddr)
		if err != nil {
			return nil, fmt.Errorf("types.GetFromBech32: %s", err.Error())
		}
		address, err := types.AccAddressFromHex(hex.EncodeToString(addressBytes))
		if err != nil {
			return nil, fmt.Errorf("types.AccAddressFromHex: %s", err.Error())
		}
		totalVotes, err := s.dao.GetTotalVotesByAddress(address.String())
		if err != nil {
			return nil, fmt.Errorf("dao.GetTotalVotesByAddress: %s", err.Error())
		}

		delegatorsTotal, err := s.dao.GetDelegatorsTotal(filters.Delegators{Validators: []string{v.OperatorAddress}})
		if err != nil {
			return nil, fmt.Errorf("dao.GetDelegatorsTotal: %s", err.Error())
		}

		power24Change, err := s.dao.GetVotingPower(filters.VotingPower{
			TimeRange: filters.TimeRange{
				From: dmodels.NewTime(time.Now().Add(-time.Hour * 24)),
				To:   dmodels.NewTime(time.Now()),
			},
			Validators: []string{v.OperatorAddress},
		})

		selfStake, err := s.node.GetDelegatorValidatorStake(address.String(), v.OperatorAddress)
		if err != nil {
			return nil, fmt.Errorf("node.GetDelegatorValidatorStake: %s", err.Error())
		}

		power := v.DelegatorShares.Div(node.PrecisionDiv)
		percentPower := decimal.Zero
		if !stakingPool.Pool.BondedTokens.IsZero() {
			percentPower = power.Div(stakingPool.Pool.BondedTokens).Mul(decimal.NewFromInt(100)).Truncate(2)
		}

		validators = append(validators, smodels.Validator{
			Title:           v.Description.Moniker,
			Power:           power,
			PercentPower:    percentPower,
			SelfStake:       selfStake,
			Fee:             v.Commission.CommissionRates.Rate,
			BlocksProposed:  blockProposed,
			Delegators:      delegatorsTotal,
			Power24Change:   power24Change,
			GovernanceVotes: totalVotes,
			Website:         v.Description.Website,
			OperatorAddress: v.OperatorAddress,
			AccAddress:      address.String(),
			ConsAddress:     consAddress,
		})
	}

	sort.Slice(validators, func(i, j int) bool {
		return validators[i].Power.Cmp(validators[j].Power) == 1
	})

	return validators, nil
}

func (s *ServiceFacade) GetTopProposedBlocksValidators() (items []dmodels.ValidatorValue, err error) {
	data, found := s.dao.CacheGet(topProposedBlocksValidatorsKey)
	if found {
		return data.([]dmodels.ValidatorValue), nil
	}
	items, err = s.dao.GetTopProposedBlocksValidators()
	if err != nil {
		return nil, fmt.Errorf("dao.GetTopProposedBlocksValidators: %s", err.Error())
	}
	validators, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidators: %s", err.Error())
	}
	mp := make(map[string]string)
	for _, validator := range validators {
		address, err := helpers.GetHexAddressFromBase64PK(validator.ConsensusPubkey.Key)
		if err != nil {
			return nil, fmt.Errorf("helpers.GetHexAddressFromBase64PK: %s", err.Error())
		}
		mp[address] = validator.Description.Moniker
	}
	for i, item := range items {
		title, found := mp[item.Validator]
		if found {
			items[i] = dmodels.ValidatorValue{
				Validator: title,
				Value:     item.Value,
			}
		}
	}
	s.dao.CacheSet(topProposedBlocksValidatorsKey, items, time.Minute*60)
	return items, nil
}

func (s *ServiceFacade) GetMostJailedValidators() (items []dmodels.ValidatorValue, err error) {
	data, found := s.dao.CacheGet(mostJailedValidators)
	if found {
		return data.([]dmodels.ValidatorValue), nil
	}
	items, err = s.dao.GetMostJailedValidators()
	if err != nil {
		return nil, fmt.Errorf("dao.GetMostJailedValidators: %s", err.Error())
	}
	validators, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidators: %s", err.Error())
	}
	mp := make(map[string]string)
	for _, validator := range validators {
		mp[validator.OperatorAddress] = validator.Description.Moniker
	}
	for i, item := range items {
		title, found := mp[item.Validator]
		if found {
			items[i] = dmodels.ValidatorValue{
				Validator: title,
				Value:     item.Value,
			}
		}
	}
	s.dao.CacheSet(mostJailedValidators, items, time.Minute*60)
	return items, nil
}

func (s *ServiceFacade) GetFeeRanges() (items []smodels.FeeRange, err error) {
	point := int64(10)
	min := decimal.Zero
	max := decimal.Zero
	validatorsMap, err := s.GetValidatorMap()
	for _, validator := range validatorsMap {
		if min.IsZero() && max.IsZero() {
			min = validator.Commission.CommissionRates.Rate
			max = validator.Commission.CommissionRates.Rate
			continue
		}
		if validator.Commission.CommissionRates.Rate.LessThan(min) {
			min = validator.Commission.CommissionRates.Rate
		}
		if validator.Commission.CommissionRates.Rate.GreaterThan(max) {
			max = validator.Commission.CommissionRates.Rate
		}
	}
	step := max.Sub(min).Div(decimal.NewFromInt(point))
	for i := int64(1); i <= point; i++ {
		var validators []smodels.FeeRangeValidator
		from := step.Mul(decimal.NewFromInt(i)).Sub(step)
		to := step.Mul(decimal.NewFromInt(i))

		for _, validator := range validatorsMap {
			rate := validator.Commission.CommissionRates.Rate
			if rate.GreaterThan(from) && rate.LessThanOrEqual(to) {
				validators = append(validators, smodels.FeeRangeValidator{
					Validator: validator.Description.Moniker,
					Fee:       rate,
				})
			}
		}
		items = append(items, smodels.FeeRange{
			From:       from,
			To:         to,
			Validators: validators,
		})
	}
	return items, nil
}

func (s *ServiceFacade) GetValidatorsDelegatorsTotal() (values []dmodels.ValidatorValue, err error) {
	validatorsMap, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidatorMap: %s", err.Error())
	}
	values, err = s.dao.GetValidatorsDelegatorsTotal()
	if err != nil {
		return nil, fmt.Errorf("dao.GetValidatorsDelegatorsTotal: %s", err.Error())
	}
	for i, v := range values {
		validator, found := validatorsMap[v.Validator]
		if found {
			values[i].Validator = validator.Description.Moniker
		}
	}
	return values, nil
}

func (s *ServiceFacade) GetValidator(address string) (validator smodels.Validator, err error) {
	data, found := s.dao.CacheGet(validatorsCacheKey)
	if !found {
		return validator, fmt.Errorf("not found in cache")
	}
	validators, ok := data.([]smodels.Validator)
	if !ok {
		return validator, fmt.Errorf("can`t cast to current type")
	}
	for _, v := range validators {
		if v.OperatorAddress == address {
			return v, nil
		}
	}
	return validator, fmt.Errorf("not found validator with address: %s", address)
}

func (s *ServiceFacade) GetValidatorBalance(valAddress string) (balance smodels.Balance, err error) {
	validator, err := s.GetValidator(valAddress)
	if err != nil {
		return balance, fmt.Errorf("GetValidator: %s", err.Error())
	}
	balance.SelfDelegated = validator.SelfStake
	balance.OtherDelegated = validator.Power.Sub(validator.SelfStake)
	addressBytes, err := types.GetFromBech32(valAddress, application.Bech32PrefixValAddr)
	if err != nil {
		return balance, fmt.Errorf("types.GetFromBech32: %s", err.Error())
	}
	address, err := types.AccAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		return balance, fmt.Errorf("types.AccAddressFromHex: %s", err.Error())
	}
	balance.Available, err = s.node.GetBalance(address.String())
	if err != nil {
		return balance, fmt.Errorf("node.GetBalance: %s", err.Error())
	}
	return balance, nil
}

func (s *ServiceFacade) GetValidatorMisc(valAddress string) (result node.ValidatorMiscResult, err error) {

	log.Info("service.GetValidatorMisc() entered")

	result, err = s.node.GetValidatorMisc(valAddress)
	if err != nil {
		return node.ValidatorMiscResult{}, fmt.Errorf("node.GetValidatorMiscResult: %s", err.Error())
	}

	return result, nil
}

func (s *ServiceFacade) GetValidatorDelegations(valAddress string) (result node.DelegatorValidatorStakeResult, err error) {

        log.Info("service.GetValidatorDelegations() entered")

        result, err = s.node.GetValidatorDelegations(valAddress)
        if err != nil {
                return node.DelegatorValidatorStakeResult{}, fmt.Errorf("node.GetValidatorDelegationsResult: %s", err.Error())
        }

        return result, nil
}

func (s *ServiceFacade) GetValidatorGovernance(selfDelegateAddr string) (votes []dmodels.ProposalVote, err error) {

	log.Info("service.GetValidatorGovernance() entered")

	votes, err = s.dao.GetValidatorGovernance(selfDelegateAddr)
	if err != nil {
		return nil, fmt.Errorf("dao.GetValidatorGovernance: %s", err.Error())
	}

	return votes, nil

}

func (s *ServiceFacade) GetValidatorTransfer(selfDelegateAddr string)(transfers []dmodels.Transfer, err error) {

	log.Info("service.GetValidatorTransfer() entered")

        transfers, err = s.dao.GetValidatorTransfer(selfDelegateAddr)
        if err != nil {
                return nil, fmt.Errorf("dao.GetValidatorTransfer: %s", err.Error())
        }

	return transfers, nil
}

func (s *ServiceFacade) GetValidatorCommunityPool()(result node.ValidatorCommunityPoolResult, err error) {

	log.Info("service.GetValidatorCommunityPool() entered")

	result, err = s.node.GetValidatorCommunityPool()
	        if err != nil {
                return node.ValidatorCommunityPoolResult{}, fmt.Errorf("node.GetValidatorCommunityPoolResult: %s", err.Error())
        }

	return result, nil
}

func (s *ServiceFacade) GetValidatorSelfDelegate(validatorAddr string)(result node.ValidatorSelfDelegateResult, err error) {

        log.Info("service.GetValidatorSelfDelegate() entered")

        result, err = s.node.GetValidatorSelfDelegate(validatorAddr)
                if err != nil {
                return node.ValidatorSelfDelegateResult{}, fmt.Errorf("node.GetValidatorSelfDelegate: %s", err.Error())
        }

        return result, nil
}

func (s *ServiceFacade) GetValidatorPowerChangeDelegate(validatorAddr string)(delegations []dmodels.Delegation, err error) {

        log.Info("service.GetValidatorPowerChangeDelegate() entered")

        delegations, err = s.dao.GetValidatorPowerChangeDelegate(validatorAddr)
                if err != nil {
                return nil, fmt.Errorf("dao.GetValidatorPowerChangeDelegate: %s", err.Error())
        }

        return delegations, nil
}

func (s *ServiceFacade) GetValidatorPowerChangeUndelegate(validatorAddr string)(result node.ValidatorPowerChangeUndelegateResult, err error) {

        log.Info("service.GetValidatorPowerChangeUndelegate() entered")

        result, err = s.node.GetValidatorPowerChangeUndelegate(validatorAddr)
                if err != nil {
                return node.ValidatorPowerChangeUndelegateResult{}, fmt.Errorf("node.GetValidatorPowerChangeUndelegate: %s", err.Error())
        }

        return result, nil
}

func (s *ServiceFacade) GetValidatorDistribution(validatorAddr string)(delegator_rewards []dmodels.DelegatorReward, err error) {

        log.Info("service.GetValidatorDistribution() entered")

        delegator_rewards, err = s.dao.GetValidatorDistribution(validatorAddr)
                if err != nil {
                return nil, fmt.Errorf("dao.GetValidatorDistribution: %s", err.Error())
        }

        return delegator_rewards, nil
}

func (s *ServiceFacade) GetValidatorStaking(validatorAddr string)(validator_stakings []dmodels.ValidatorStaking, err error) {

        log.Info("service.GetValidatorStaking() entered")

        validator_stakings, err = s.dao.GetValidatorStaking(validatorAddr)
                if err != nil {
                return nil, fmt.Errorf("dao.GetValidatorStaking: %s", err.Error())
        }

        return validator_stakings, nil
}

func (s *ServiceFacade) GetValidatorLast100Blocks(consensusAddr string)(blocks []dmodels.Block, err error) {

        log.Info("service.GetValidatorLast100Blocks() entered")

        blocks, err = s.dao.GetValidatorLast100Blocks(consensusAddr)
                if err != nil {
                return nil, fmt.Errorf("dao.GetValidatorLast100Blocks: %s", err.Error())
        }

        return blocks, nil
}

func (s *ServiceFacade) GetValidatorMissedBlocks(consensusAddr string)(missed_blocks []dmodels.MissedBlock, err error) {

        log.Info("service.GetValidatorMissedBlocks() entered")

        missed_blocks, err = s.dao.GetValidatorMissedBlocks(consensusAddr)
                if err != nil {
                return nil, fmt.Errorf("dao.GetValidatorMissedBlocks: %s", err.Error())
        }

        return missed_blocks, nil
}

func (s *ServiceFacade) GetValidatorSlashing(validatorAddr string)(result node.ValidatorSlashingResult, err error) {

        log.Info("service.GetValidatorSlashing() entered")

        result, err = s.node.GetValidatorSlashing(validatorAddr)
                if err != nil {
		return node.ValidatorSlashingResult{}, fmt.Errorf("node.GetValidatorSlashing: %s", err.Error())
        }

        return result, nil
}

func (s *ServiceFacade) GetValidatorProposerPriority(validatorAddr string)(result node.ValidatorProposerPriorityResult, err error) {

        log.Info("service.GetValidatorProposerPriority() entered")

        result, err = s.node.GetValidatorProposerPriority(validatorAddr)
                if err != nil {
                return node.ValidatorProposerPriorityResult{}, fmt.Errorf("node.GetValidatorProposerPriority: %s", err.Error())
        }

        return result, nil
}

func (s *ServiceFacade) GetValidatorUptimePercent(consensusAddr string)(uptimepercent float64, err error) {

	log.Info("service.GetValidatorUptimePercent() entered")

	uptimepercent, err = s.dao.GetValidatorUptimePercent(consensusAddr)
	if err != nil {
		return 0, fmt.Errorf("dao.GetValidatorUptimePercent: %s", err.Error())
	}

        return uptimepercent, nil
}

func (s *ServiceFacade) GetValidatorAggInfo(validatorAddr string)(result node.ValidatorAggInfoResult, err error) {

	log.Info("service.GetValidatorAggInfo() entered")

	validatorResult, err := s.GetValidator(validatorAddr)
	if err != nil {
                return node.ValidatorAggInfoResult{}, fmt.Errorf("GetValidatorAggInfo() -> GetValidator: %s", err.Error())
        }

	stakingPoolResult, err := s.node.GetStakingPool()
	if err != nil {
		log.Error("GetValidatorAggInfo() -> s.node.GetStakingPool: %s", err.Error())
		return
	}

	miscResult, err := s.GetValidatorMisc(validatorAddr)
	if err != nil {
                return node.ValidatorAggInfoResult{}, fmt.Errorf("GetValidatorAggInfo() -> s.GetValidatorMisc: %s", err.Error())
        }

	cons_address := validatorResult.ConsAddress
	uptimepercentResult, err := s.GetValidatorUptimePercent(cons_address)
	if err != nil {
		return node.ValidatorAggInfoResult{}, fmt.Errorf("GetValidatorAggInfo() -> s.GetValidatorUptimePercent: %s", err.Error())
	}

	delegated := float64(miscResult.Result.Tokens)/100000000

	bonded_tokens, _ := stakingPoolResult.Pool.BondedTokens.Float64()

	log.Info("GetValidatorAggInfo() -> bonded_tokens(0): %g", bonded_tokens)

	bonded_tokens = bonded_tokens/100
	commission, _ := strconv.ParseFloat(miscResult.Result.Commission.CommissionRates.Rate, 64)

	log.Info("GetValidatorAggInfo() -> bonded_tokens(1): %g", bonded_tokens)

	result.Name = validatorResult.Title
        result.Delegated = delegated
	result.DelegatedPercent = 100*delegated/bonded_tokens
	result.CommissionPercent = 100*commission
	result.Status = miscResult.Result.Status
	result.UptimePercent = uptimepercentResult

	return result, nil
}

