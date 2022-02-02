package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
	"github.com/everstake/cosmoscan-api/log"
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

func (s *ServiceFacade) GetBlocks(filter filters.Blocks) (block []dmodels.Block, err error) {
	blocks, err := s.dao.GetBlocks(filter)
	if err != nil {
		return blocks, fmt.Errorf("dao.GetBlocks: %s", err.Error())
	}
	return blocks, nil
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
