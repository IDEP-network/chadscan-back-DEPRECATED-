package clickhouse

import (
        "fmt"
        "github.com/Masterminds/squirrel"
        //"github.com/everstake/cosmoscan-api/dao/filters"
        "github.com/everstake/cosmoscan-api/dmodels"
        //"github.com/everstake/cosmoscan-api/smodels"
        //"github.com/shopspring/decimal"
        "github.com/everstake/cosmoscan-api/log"
)

func (db DB) GetValidatorGovernance(selfDelegateAddr string)(votes []dmodels.ProposalVote, err error) {

	log.Info("dao.clickhouse.GetValidatorGovernance() entered")

	q := squirrel.Select("*").From(dmodels.ProposalVotesTable).Where(squirrel.Eq{"prv_voter": selfDelegateAddr})
	err = db.Find(&votes, q)

	return votes, err
}

func (db DB) GetValidatorTransfer(selfDelegateAddr string)(transfers []dmodels.Transfer, err error) {

        log.Info("dao.clickhouse.GetValidatorTransfer() entered")
	padding := "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" //padding 21 zeroes
	param := selfDelegateAddr + padding

	log.Info("param = %s", param)

        q := squirrel.Select("*").From(dmodels.TransfersTable).Where(squirrel.Or{
		squirrel.Eq{"trf_to": param},
		squirrel.Eq{"trf_from": param},
		})

        err = db.Find(&transfers, q)

        return transfers, err
}

func (db DB) GetValidatorPowerChangeDelegate(validatorAddr string)(delegations []dmodels.Delegation, err error) {

	log.Info("dao.clickhouse.GetValidatorPowerChangeDelegate() entered")
	padding := "\x00\x00\x00\x00\x00\x00\x00" //padding 7 zeroes
	param := validatorAddr + padding

	log.Info("param = %s", param)

	q := squirrel.Select("*").From(dmodels.DelegationsTable).Where(squirrel.Eq{"dlg_validator": validatorAddr + padding})

	err = db.Find(&delegations, q)

	return delegations, err

}

func (db DB) GetValidatorDistribution(validatorAddr string)(delegator_rewards []dmodels.DelegatorReward, err error) {

        log.Info("dao.clickhouse.GetValidatorDistribution() entered")
        padding := "\x00\x00\x00\x00\x00\x00\x00" //padding 7 zeroes
        param := validatorAddr + padding

        log.Info("param = %s", param)

        q := squirrel.Select("*").From(dmodels.DelegatorRewardsTable).Where(squirrel.Eq{"der_validator": validatorAddr + padding})

        err = db.Find(&delegator_rewards, q)

        return delegator_rewards, err

}

func (db DB) GetValidatorStaking(validatorAddr string)(validator_stakings []dmodels.ValidatorStaking, err error) {

	log.Info("dao.clickhouse.GetValidatorStaking() entered")
        padding := "\x00\x00\x00\x00\x00\x00\x00" //padding 7 zeroes
        param := validatorAddr + padding

        log.Info("param = %s", param)

	q := squirrel.Select("dlg_id, dlg_tx_hash, dlg_delegator, dlg_validator, dlg_amount, dlg_created_at, 'delegate' as dlg_type").From(dmodels.DelegationsTable).Where(squirrel.Eq{"dlg_validator": validatorAddr + padding})

        err = db.Find(&validator_stakings, q)

	return validator_stakings, err

}

func (db DB) GetValidatorLast100Blocks(consensusAddr string)(blocks []dmodels.Block, err error) {

        log.Info("dao.clickhouse.GetValidatorLast100Blocks() entered")

        q := squirrel.Select("*").From(dmodels.BlocksTable).Where(squirrel.Eq{"blk_proposer": consensusAddr}).OrderBy("blk_created_at desc").Limit(100)

        err = db.Find(&blocks, q)

        return blocks, err

}

func (db DB) GetValidatorMissedBlocks(consensusAddr string)(missed_blocks []dmodels.MissedBlock, err error) {

	log.Info("dao.clickhouse.GetValidatorMissedBlocks() entered")

        q := squirrel.Select("*").From(dmodels.MissedBlocks).Where(squirrel.Eq{"mib_validator": consensusAddr})

        err = db.Find(&missed_blocks, q)

        return missed_blocks, err

}

func (db DB) GetValidatorUptimePercent(consensusAddr string) (uptimepercent float64, err error) {

        log.Info("dao.clickhouse.GetValidatorUptimePercent() entered")

	var latestBlock []dmodels.Block	//latest block
	var prevBlock []dmodels.Block	//the block that is NUM_BLOCKS_TO_CHECK  blocks prior to the latest block
	var numMissedBlocks uint64
	var NUM_BLOCKS_TO_CHECK uint64

	NUM_BLOCKS_TO_CHECK = 10000

	q := squirrel.Select("*").From(dmodels.BlocksTable).OrderBy("blk_id desc").Limit(1)
	err = db.Find(&latestBlock, q)

	if err != nil {
		return 0, fmt.Errorf("dao.clickhouse.GetValidatorUptimePercent - 0: %s", err.Error())
	}

	var prevBlockId uint64
	if latestBlock[0].ID < NUM_BLOCKS_TO_CHECK {
		prevBlockId = 1
	} else {
		prevBlockId = latestBlock[0].ID
	}

	q = squirrel.Select("*").From(dmodels.BlocksTable).Where(squirrel.Eq{"blk_id": prevBlockId})
        err = db.Find(&prevBlock, q)

	if err != nil {
                return 0, fmt.Errorf("dao.clickhouse.GetValidatorUptimePercent - 1: %s", err.Error())
	}

	q = squirrel.Select("count (*)").From(dmodels.MissedBlocks).Where( squirrel.And{
		squirrel.Eq{"mib_validator": consensusAddr},
		squirrel.GtOrEq{"mib_created_at": prevBlock[0].CreatedAt},
		squirrel.LtOrEq{"mib_created_at": latestBlock[0].CreatedAt},
		})

        err = db.FindFirst(&numMissedBlocks, q)

	if err != nil {
                return 0, fmt.Errorf("dao.clickhouse.GetValidatorUptimePercent - 2: %s", err.Error())
	}

	uptimepercent = 100*float64(NUM_BLOCKS_TO_CHECK - numMissedBlocks)/float64(NUM_BLOCKS_TO_CHECK)

	return uptimepercent, nil

}

