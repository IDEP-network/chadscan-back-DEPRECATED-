package services

import (
	"fmt"

	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
)

const topProposedBlocksValidatorsKey = "topProposedBlocksValidatorsKey"
const rewardPerBlock = 4.0

func (s *ServiceFacade) GetAggBlocksCount(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggBlocksCount(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggBlocksCount: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetBlocks(filter filters.Blocks) (resp smodels.PaginatableResponse, err error) {
	dBlocks, err := s.dao.GetBlocks(filter)
	if err != nil {
		return resp, fmt.Errorf("dao.GetBlocks: %s", err.Error())
	}
	total, err := s.dao.GetBlocksCount(filter)
	if err != nil {
		return resp, fmt.Errorf("dao.GetBlocksCount: %s", err.Error())
	}
	validators, err := s.makeValidatorMap()
	if err != nil {
		return resp, fmt.Errorf("s.makeValidatorMap: %s", err.Error())
	}
	var blocks []dmodels.Block
	for _, b := range dBlocks {
		var proposer string
		validator, ok := validators[b.Proposer]
		if ok {
			proposer = validator.Description.Moniker
		}
		blocks = append(blocks, dmodels.Block{
			ID:          				b.ID,
			Hash:            b.Hash,
			Proposer:        proposer,
			CreatedAt:       b.CreatedAt,
		})
	}
	return smodels.PaginatableResponse{
		Items: blocks,
		Total: total,
	}, nil
}

func (s *ServiceFacade) GetAggBlocksDelay(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggBlocksDelay(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggBlocksDelay: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAggUniqBlockValidators(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggUniqBlockValidators(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggUniqBlockValidators: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetValidatorBlocksStat(validatorAddress string) (stat smodels.ValidatorBlocksStat, err error) {
	validator, err := s.GetValidator(validatorAddress)
	if err != nil {
		return stat, fmt.Errorf("GetValidator: %s", err.Error())
	}
	stat.Proposed, err = s.dao.GetProposedBlocksTotal(filters.BlocksProposed{
		Proposers: []string{validator.ConsAddress},
	})
	if err != nil {
		return stat, fmt.Errorf("dao.GetProposedBlocksTotal: %s", err.Error())
	}
	stat.MissedValidations, err = s.dao.GetMissedBlocksCount(filters.MissedBlocks{
		Validators: []string{validator.ConsAddress},
	})
	if err != nil {
		return stat, fmt.Errorf("dao.GetMissedBlocksCount: %s", err.Error())
	}
	stat.Revenue = decimal.NewFromFloat(rewardPerBlock).Mul(decimal.NewFromInt(int64(stat.Proposed)))
	return stat, nil
}

func (s *ServiceFacade) GetBlockTransactionCounts(blockHeight string) (transactioncounts int, err error) {

	log.Info("service.GetBlockTransactionCounts() entered")

	transactioncounts = 0

	block, err := s.node.GetBlockHeight(blockHeight)
	if err != nil {
		return transactioncounts, fmt.Errorf("node.GetBlockTransactionCounts - GetBlockHeight: %s", err.Error())
	}

	transactioncounts = len(block.Block.Data.Txs)

	return transactioncounts, nil
}
