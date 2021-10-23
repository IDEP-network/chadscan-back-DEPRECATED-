package clickhouse

import (
	//"fmt"
	"github.com/Masterminds/squirrel"
	//"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	//"github.com/everstake/cosmoscan-api/smodels"
	//"github.com/shopspring/decimal"
	"github.com/everstake/cosmoscan-api/log"
)

func (db DB) GetTransactionHash(txnHash string) (txn []dmodels.Transaction, err error) {

	log.Info("dao.clickhouse.GetTransactionHash() entered")

	//q := squirrel.Select("*").From(dmodels.TransactionsTable).
	//	Where(squirrel.Eq{"trn_hash": txnHash})

	//the following is working:
	//q := squirrel.Select("*").From(dmodels.TransactionsTable).Limit(2)

	//q := squirrel.Select("*").From(dmodels.TransactionsTable).Where(squirrel.Eq{"Hash": txnHash})
	q := squirrel.Select("*").From(dmodels.TransactionsTable).Where(squirrel.Eq{"trn_hash": txnHash})

	//log.Info("q=%s", q)
        err = db.Find(&txn, q)

	return txn, err
}

func (db DB) GetBlockHash(blockHash string) (block []dmodels.Block, err error) {

	log.Info("dao.clickhouse.GetBlockHash() entered")

	q := squirrel.Select("*").From(dmodels.BlocksTable).Where(squirrel.Eq{"blk_hash": blockHash})
	err = db.Find(&block, q)

	return block, err
}
