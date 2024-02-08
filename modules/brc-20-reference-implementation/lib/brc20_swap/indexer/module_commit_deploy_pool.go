package indexer

import (
	"errors"
	"log"

	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/model"
)

func (g *BRC20ModuleIndexer) ProcessCommitFunctionDeployPool(moduleInfo *model.BRC20ModuleSwapInfo, f *model.SwapFunctionData) error {
	token0, token1 := f.Params[0], f.Params[1]

	poolPair := GetLowerPairNameByToken(token0, token1)
	if _, ok := moduleInfo.SwapPoolTotalBalanceDataMap[poolPair]; ok {
		return errors.New("deploy: twice")
	}
	poolPairReverse := GetLowerPairNameByToken(token1, token0)
	if _, ok := moduleInfo.SwapPoolTotalBalanceDataMap[poolPairReverse]; ok {
		return errors.New("deploy: twice")
	}

	poolPair = GetLowerInnerPairNameByToken(token0, token1)

	// lp token balance of address in module [pool][address]balance
	moduleInfo.LPTokenUsersBalanceMap[poolPair] = make(map[string]*decimal.Decimal, 0)

	// swap total balance
	// total balance of pool in module [pool]balanceData
	moduleInfo.SwapPoolTotalBalanceDataMap[poolPair] = &model.BRC20ModulePoolTotalBalance{
		Tick:    [2]string{token0, token1},
		History: make([]*model.BRC20ModuleHistory, 0), // fixme:
		// balance
	}
	log.Printf("[%s] pool deploy pool [%s]", moduleInfo.ID, poolPair)
	return nil
}
