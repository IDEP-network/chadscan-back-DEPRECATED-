package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/services/node"
)

func (s *ServiceFacade) GetAggTransactionsFee(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggTransactionsFee(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggTransactionsFee: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAggOperationsCount(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggOperationsCount(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggOperationsCount: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAvgOperationsPerBlock(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAvgOperationsPerBlock(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAvgOperationsPerBlock: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetLatestTransactions (count string) (latest_transactions []dmodels.Transaction, err error) {

        log.Info("services.GetLatestTransactions() entered")

	latest_transactions, err = s.dao.GetLatestTransactions(count)
	if err != nil {
		return nil, fmt.Errorf("dao.GetLatestTransactions: %s", err.Error())
	}

	return latest_transactions, nil
}

func (s *ServiceFacade) GetTransactionDetail (hash string) (result node.TransactionHashResult, err error) {

        log.Info("services.GetTransactionDetail() entered")

        result, err = s.node.GetTransactionDetail(hash)
        if err != nil {
                log.Error("GetTransactionDetail() -> s.node.GetTransactionDetail: %s", err.Error())
                return
        }

        return result, nil
}
