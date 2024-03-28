package indexer

import (
	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/loader"
	"brc20query/lib/brc20_swap/model"
	"log"
)

// func (g *BRC20ModuleIndexer) SaveDataToDB(height int) {

func (g *BRC20ModuleIndexer) SaveDataToDB(dbConnInfo string, height uint32) {
	loader.Init(dbConnInfo)
	defer loader.SwapDB.Close()

	// ticker info
	loader.SaveDataToDBTickerInfoMap(height, g.InscriptionsTickerInfoMap)
	loader.SaveDataToDBTickerBalanceMap(height, g.TokenUsersBalanceData)
	loader.SaveDataToDBTickerHistoryMap(height, g.AllHistory)

	loader.SaveDataToDBTransferStateMap(height, g.InscriptionsTransferRemoveMap)
	loader.SaveDataToDBValidTransferMap(height, g.InscriptionsValidTransferMap)

	// module info
	loader.SaveDataToDBModuleInfoMap(height, g.ModulesInfoMap)
	loader.SaveDataToDBModuleHistoryMap(height, g.ModulesInfoMap)
	loader.SaveDataToDBModuleCommitChainMap(height, g.ModulesInfoMap)
	loader.SaveDataToDBModuleUserBalanceMap(height, g.ModulesInfoMap)
	loader.SaveDataToDBModulePoolLpBalanceMap(height, g.ModulesInfoMap)
	loader.SaveDataToDBModuleUserLpBalanceMap(height, g.ModulesInfoMap)

	loader.SaveDataToDBSwapApproveStateMap(height, g.InscriptionsApproveRemoveMap)
	loader.SaveDataToDBSwapApproveMap(height, g.InscriptionsValidApproveMap)

	loader.SaveDataToDBSwapCondApproveStateMap(height, g.InscriptionsCondApproveRemoveMap)
	loader.SaveDataToDBSwapCondApproveMap(height, g.InscriptionsValidConditionalApproveMap)

	loader.SaveDataToDBSwapCommitStateMap(height, g.InscriptionsCommitRemoveMap)
	loader.SaveDataToDBSwapCommitMap(height, g.InscriptionsValidCommitMap)

	loader.SaveDataToDBSwapWithdrawStateMap(height, g.InscriptionsWithdrawRemoveMap)
	loader.SaveDataToDBSwapWithdrawMap(height, g.InscriptionsValidWithdrawMap)

}

func (g *BRC20ModuleIndexer) LoadDataFromDB(dbConnInfo string, height int) {
	loader.Init(dbConnInfo)
	defer loader.SwapDB.Close()

	var err error
	if g.InscriptionsTickerInfoMap, err = loader.LoadFromDbTickerInfoMap(); err != nil {
		log.Fatal("LoadFromDBTickerInfoMap failed: ", err)
	}

	if g.UserTokensBalanceData, err = loader.LoadFromDbUserTokensBalanceData(nil, nil); err != nil {
		log.Fatal("LoadFromDBUserTokensBalanceData failed: ", err)
	}
	g.TokenUsersBalanceData = loader.UserTokensBalanceMap2TokenUsersBalanceMap(g.UserTokensBalanceData)

	if g.InscriptionsTransferRemoveMap, err = loader.LoadFromDBTransferStateMap(); err != nil {
		log.Fatal("LoadFromDBTransferStateMap failed: ", err)
	}

	if g.InscriptionsValidTransferMap, err = loader.LoadFromDBValidTransferMap(); err != nil {
		log.Fatal("LoadFromDBvalidTransferMap failed: ", err)
	}

	if g.ModulesInfoMap, err = loader.LoadFromDBModuleInfoMap(); err != nil {
		log.Fatal("LoadFromDBModuleInfoMap failed: ", err)
	}

	for mid, info := range g.ModulesInfoMap {
		loadFromDBSwapModuleInfo(mid, info)
	}

	if g.InscriptionsApproveRemoveMap, err = loader.LoadFromDBSwapApproveStateMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapApproveStateMap failed: ", err)
	}

	if g.InscriptionsValidApproveMap, err = loader.LoadFromDBSwapApproveMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapApproveMap failed: ", err)
	}

	if g.InscriptionsCondApproveRemoveMap, err = loader.LoadFromDBSwapCondApproveStateMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCondApproveStateMap failed: ", err)
	}

	if g.InscriptionsValidConditionalApproveMap, err = loader.LoadFromDBSwapCondApproveMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCondApproveMap failed: ", err)
	}

	if g.InscriptionsCommitRemoveMap, err = loader.LoadFromDBSwapCommitStateMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCommitStateMap failed: ", err)
	}
	if g.InscriptionsValidCommitMap, err = loader.LoadFromDBSwapCommitMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCommitMap failed: ", err)
	}

	if g.InscriptionsWithdrawRemoveMap, err = loader.LoadFromDBSwapWithdrawStateMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapWithdrawStateMap failed: ", err)
	}
	if g.InscriptionsValidWithdrawMap, err = loader.LoadFromDBSwapWithdrawMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapWithdrawMap failed: ", err)
	}
}

func loadFromDBSwapModuleInfo(mid string, info *model.BRC20ModuleSwapInfo) {
	if hm, err := loader.LoadFromDBModuleHistoryMap(mid); err != nil {
		log.Fatal("LoadFromDBModuleHistoryMap failed: ", err)
	} else {
		for _, history := range hm {
			info.History = history
		}
	}

	if ccs, err := loader.LoadModuleCommitChain(mid, nil); err != nil {
		log.Fatal("LoadModuleCommitChain failed: ", err)
	} else {
		for _, cc := range ccs {
			if cc.Valid && cc.Connected {
				info.CommitIdChainMap[cc.CommitID] = struct{}{}
			} else if cc.Valid && !cc.Connected {
				info.CommitIdMap[cc.CommitID] = struct{}{}
			} else {
				info.CommitInvalidMap[cc.CommitID] = struct{}{}
			}
		}
	}

	// [tick][address]balanceData
	if tabm, err := loader.LoadFromDBModuleUserBalanceMap(mid, nil, nil); err != nil {
		log.Fatal("LoadFromDBModuleUserBalanceMap failed: ", err)
	} else {
		info.TokenUsersBalanceDataMap = tabm
		for tick, abs := range tabm {
			for addr, balance := range abs {
				if _, ok := info.UsersTokenBalanceDataMap[addr]; !ok {
					info.UsersTokenBalanceDataMap[addr] = make(map[string]*model.BRC20ModuleTokenBalance)
				}
				// [address][tick]balanceData
				info.UsersTokenBalanceDataMap[addr][tick] = balance
			}
		}
	}

	if poolBalanceMap, err := loader.LoadFromDBModulePoolLpBalanceMap(mid, nil); err != nil {
		log.Fatal("LoadFromDBModulePoolLpBalanceMap failed: ", err)
	} else {
		info.SwapPoolTotalBalanceDataMap = poolBalanceMap
	}

	// [pool][address]balance
	if userLpBalanceMap, err := loader.LoadFromDBModuleUserLpBalanceMap(mid, nil, nil); err != nil {
		log.Fatal("LoadFromDBModuleUserLpBalanceMap failed: ", err)
	} else {
		info.LPTokenUsersBalanceMap = userLpBalanceMap

		for pool, abs := range userLpBalanceMap {
			for addr, balance := range abs {
				if _, ok := info.UsersLPTokenBalanceMap[addr]; !ok {
					info.UsersLPTokenBalanceMap[addr] = make(map[string]*decimal.Decimal)
				}
				// [address][pool]balance
				info.UsersLPTokenBalanceMap[addr][pool] = balance
			}
		}
	}

}
