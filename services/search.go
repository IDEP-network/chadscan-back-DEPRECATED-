package services

import (
	//"encoding/hex"
	"fmt"
	//"github.com/cosmos/cosmos-sdk/types"
	//"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	//"github.com/everstake/cosmoscan-api/services/helpers"
	"github.com/everstake/cosmoscan-api/services/node"
	//"github.com/everstake/cosmoscan-api/smodels"
	//"github.com/persistenceOne/persistenceCore/application"
	//"github.com/shopspring/decimal"
	//"sort"
	//"time"
)

func (s *ServiceFacade) GetWalletAddress (walletAddr string) (result node.WalletAddressResult, err error) {

	log.Info("services.GetWalletAddress() entered")

	result, err  = s.node.GetWalletAddress(walletAddr)
        if err != nil {
                return node.WalletAddressResult{}, fmt.Errorf("node.GetWalletAddressResult: %s", err.Error())
        }

	//return nil, fmt.Errorf("services.GetWalletAddress() incomplete", err.Error())
	return result, nil
}

func (s *ServiceFacade) GetTransactionHash (txnHash string) (txn []dmodels.Transaction, err error) {

        log.Info("services.GetTransactionHash() entered")

	txn, err = s.dao.GetTransactionHash(txnHash)
	if err != nil {
		return nil, fmt.Errorf("dao.GetTransactionHash: %s", err.Error())
	}

        //return nil, fmt.Errorf("GetTransactionHash() incomplete", err.Error())
	return txn, nil
}


func (s *ServiceFacade) GetBlockHash(blockHash string) (block []dmodels.Block, err error) {

        log.Info("services.GetBlockHash() entered")

	block, err = s.dao.GetBlockHash(blockHash)
        if err != nil {
                return nil, fmt.Errorf("dao.GetBlockHash: %s", err.Error())
        }

        //return nil, fmt.Errorf("GetBlockHash() incomplete", err.Error())
	return block, nil
}


func (s *ServiceFacade) GetBlockHeight(blockHeight string) (result node.BlockHeightResult, err error) {

	log.Info("service.GetBlockHeight() entered")

	result, err = s.node.GetBlockHeight(blockHeight)
	if err != nil {
		return node.BlockHeightResult{}, fmt.Errorf("node.GetBlockHeight: %s", err.Error())
	}

	//return nil, fmt.Errorf("GetBlockHeight() incomplete", err.Error())

	return result, nil
}
