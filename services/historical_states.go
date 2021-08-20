package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
	"sort"
	"time"
)

func (s ServiceFacade) KeepHistoricalState() {
	for {
		states, err := s.dao.GetHistoricalStates(filters.HistoricalState{Limit: 1})
		if err != nil {
			log.Error("KeepHistoricalState: dao.GetHistoricalStates: %s", err.Error())
			<-time.After(time.Second * 10)
			continue
		}
		tn := time.Now()
		if len(states) != 0 {
			lastState := states[0]
			if tn.Sub(lastState.CreatedAt.Time) < time.Hour {
				point := lastState.CreatedAt.Time.Add(time.Hour)
				<-time.After(point.Sub(tn))
			}
		}
		state, err := s.makeState()
		if err != nil {
			log.Error("KeepHistoricalState: makeState: %s", err.Error())
			<-time.After(time.Minute * 10)
			continue
		}
		for {
			if err = s.dao.CreateHistoricalStates([]dmodels.HistoricalState{state}); err == nil {
				break
			}
			log.Error("KeepHistoricalState: dao.CreateHistoricalStates: %s", err.Error())
			<-time.After(time.Second * 10)
		}
		<-time.After(time.Minute)
	}
}

func (s ServiceFacade) Test() (state dmodels.HistoricalState, err error) {
	return s.makeState()
}

func (s ServiceFacade) makeState() (state dmodels.HistoricalState, err error) {
	state.InflationRate, err = s.node.GetInflation()
	if err != nil {
		return state, fmt.Errorf("node.GetInflation: %s", err.Error())
	}
	state.InflationRate = state.InflationRate.Truncate(2)
	state.CommunityPool, err = s.node.GetCommunityPoolAmount()
	if err != nil {
		return state, fmt.Errorf("node.GetCommunityPoolAmount: %s", err.Error())
	}
	state.CommunityPool = state.CommunityPool.Truncate(2)
	totalSupply, err := s.node.GetTotalSupply()
	if err != nil {
		return state, fmt.Errorf("node.GetTotalSupply: %s", err.Error())
	}
	stakingPool, err := s.node.GetStakingPool()
	if err != nil {
		return state, fmt.Errorf("node.GetStakingPool: %s", err.Error())
	}
	if !totalSupply.IsZero() {
		state.StakedRatio = stakingPool.Pool.BondedTokens.Div(totalSupply).Mul(decimal.New(100, 0)).Truncate(2)
	}
	validators, err := s.node.GetValidators()
	if err != nil {
		return state, fmt.Errorf("node.GetValidators: %s", err.Error())
	}
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].DelegatorShares.GreaterThan(validators[j].DelegatorShares)
	})
	if len(validators) >= 20 {
		top20Stake := decimal.Zero
		for i := 0; i < 20; i++ {
			top20Stake = top20Stake.Add(validators[i].DelegatorShares)
		}
		top20Stake = top20Stake.Div(node.PrecisionDiv)
		if !stakingPool.Pool.BondedTokens.IsZero() {
			state.Top20Weight = top20Stake.Div(stakingPool.Pool.BondedTokens).Mul(decimal.New(100, 0)).Truncate(2)
		}
	}

	state.CirculatingSupply = totalSupply.Truncate(2)

	state.Price, state.TradingVolume, err = s.cm.GetMarketData()
	if err != nil {
		return state, fmt.Errorf("cmc.GetMarketData: %s", err.Error())
	}
	state.MarketCap = state.CirculatingSupply.Mul(state.Price).Truncate(2)

	if state.Price.IsZero() {
		return state, fmt.Errorf("cmc not found currency")
	}
	state.CreatedAt = dmodels.NewTime(time.Now())
	return state, nil
}

func (s *ServiceFacade) GetHistoricalState() (state smodels.HistoricalState, err error) {
	models, err := s.dao.GetHistoricalStates(filters.HistoricalState{Limit: 1})
	if err != nil {
		return state, fmt.Errorf("dao.GetHistoricalStates: %s", err.Error())
	}
	if len(models) == 0 {
		return state, fmt.Errorf("not found any states")
	}
	state.Current = models[0]
	state.PriceAgg, err = s.dao.GetAggHistoricalStatesByField(filters.Agg{
		By:   filters.AggByHour,
		From: dmodels.NewTime(time.Now().Add(-time.Hour * 24)),
	}, "his_price")
	if err != nil {
		return state, fmt.Errorf("dao.GetAggHistoricalStatesByField: %s", err.Error())
	}
	state.MarketCapAgg, err = s.dao.GetAggHistoricalStatesByField(filters.Agg{
		By:   filters.AggByHour,
		From: dmodels.NewTime(time.Now().Add(-time.Hour * 24)),
	}, "his_market_cap")
	if err != nil {
		return state, fmt.Errorf("dao.GetAggHistoricalStatesByField: %s", err.Error())
	}
	state.StakedRatioAgg, err = s.dao.GetAggHistoricalStatesByField(filters.Agg{
		By:   filters.AggByDay,
		From: dmodels.NewTime(time.Now().Add(-time.Hour * 24 * 30)),
	}, "his_staked_ratio")
	if err != nil {
		return state, fmt.Errorf("dao.GetAggHistoricalStatesByField: %s", err.Error())
	}
	return state, nil
}
