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
