package indexer

import (
	"log"
	"time"

	"github.com/unisat-wallet/libbrc20-indexer/decimal"
	"github.com/unisat-wallet/libbrc20-indexer/loader"
	"github.com/unisat-wallet/libbrc20-indexer/model"
)

// func (g *BRC20ModuleIndexer) SaveDataToDB(height int) {

func (g *BRC20ModuleIndexer) PurgeHistoricalData() {
	// purge history
	g.AllHistory = make([]uint32, 0) // fixme
	g.InscriptionsTransferRemoveMap = make(map[string]uint32, 0)
	g.InscriptionsApproveRemoveMap = make(map[string]uint32, 0)
	g.InscriptionsCondApproveRemoveMap = make(map[string]uint32, 0)
	g.InscriptionsCommitRemoveMap = make(map[string]uint32, 0)
}

func (g *BRC20ModuleIndexer) SaveDataToDB(height uint32) {
	tx := loader.MustBegin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			if err := tx.Commit(); err != nil {
				log.Fatal("SaveDataToDB tx commit failed: ", err)
			}
		}
	}()

	// ticker info
	loader.SaveDataToDBTickerInfoMap(tx, height, g.InscriptionsTickerInfoMap)
	loader.SaveDataToDBTickerBalanceMap(tx, height, g.TokenUsersBalanceData)
	// loader.SaveDataToDBTickerHistoryMap(tx, height, g.AllHistory)  // fixme

	loader.SaveDataToDBTransferStateMap(tx, height, g.InscriptionsTransferRemoveMap)
	loader.SaveDataToDBValidTransferMap(tx, height, g.InscriptionsValidTransferMap)

	// module info
	loader.SaveDataToDBModuleInfoMap(tx, height, g.ModulesInfoMap)
	loader.SaveDataToDBModuleHistoryMap(tx, height, g.ModulesInfoMap)
	loader.SaveDataToDBModuleCommitChainMap(tx, height, g.ModulesInfoMap)
	loader.SaveDataToDBModuleUserBalanceMap(tx, height, g.ModulesInfoMap)
	loader.SaveDataToDBModulePoolLpBalanceMap(tx, height, g.ModulesInfoMap)
	loader.SaveDataToDBModuleUserLpBalanceMap(tx, height, g.ModulesInfoMap)

	loader.SaveDataToDBSwapApproveStateMap(tx, height, g.InscriptionsApproveRemoveMap)
	loader.SaveDataToDBSwapApproveMap(tx, height, g.InscriptionsValidApproveMap)

	loader.SaveDataToDBSwapCondApproveStateMap(tx, height, g.InscriptionsCondApproveRemoveMap)
	loader.SaveDataToDBSwapCondApproveMap(tx, height, g.InscriptionsValidConditionalApproveMap)

	loader.SaveDataToDBSwapCommitStateMap(tx, height, g.InscriptionsCommitRemoveMap)
	loader.SaveDataToDBSwapCommitMap(tx, height, g.InscriptionsValidCommitMap)

	loader.SaveDataToDBSwapWithdrawStateMap(tx, height, g.InscriptionsWithdrawRemoveMap)
	loader.SaveDataToDBSwapWithdrawMap(tx, height, g.InscriptionsValidWithdrawMap)
}

func (g *BRC20ModuleIndexer) LoadDataFromDB(height int) {
	var (
		err error
		st  time.Time
	)

	st = time.Now()
	if g.InscriptionsTickerInfoMap, err = loader.LoadFromDbTickerInfoMap(); err != nil {
		log.Fatal("LoadFromDBTickerInfoMap failed: ", err)
	}
	log.Printf("LoadFromDBTickerInfoMap, duration: %s, count: %d",
		time.Since(st).String(),
		len(g.InscriptionsTickerInfoMap),
	)

	st = time.Now()
	if g.UserTokensBalanceData, err = loader.LoadFromDbUserTokensBalanceData(g.InscriptionsTickerInfoMap, nil, nil); err != nil {
		log.Fatal("LoadFromDBUserTokensBalanceData failed: ", err)
	}
	g.TokenUsersBalanceData = loader.UserTokensBalanceMap2TokenUsersBalanceMap(g.InscriptionsTickerInfoMap, g.UserTokensBalanceData)
	log.Printf("LoadFromDBUserTokensBalanceData, duration: %s, ticks: %d, addresses: %d",
		time.Since(st).String(),
		len(g.TokenUsersBalanceData),
		len(g.UserTokensBalanceData),
	)

	// st = time.Now()
	// if g.InscriptionsTransferRemoveMap, err = loader.LoadFromDBTransferStateMap(); err != nil {
	// 	log.Fatal("LoadFromDBTransferStateMap failed: ", err)
	// }
	// log.Printf("LoadFromDBTransferStateMap, duration: %s, count: %d",
	// 	time.Since(st).String(),
	// 	len(g.InscriptionsTransferRemoveMap),
	// )

	st = time.Now()
	if g.InscriptionsValidTransferMap, err = loader.LoadFromDBValidTransferMap(); err != nil {
		log.Fatal("LoadFromDBvalidTransferMap failed: ", err)
	}
	log.Printf("LoadFromDBvalidTransferMap, duration: %s, count: %d",
		time.Since(st).String(),
		len(g.InscriptionsValidTransferMap),
	)

	st = time.Now()
	if g.ModulesInfoMap, err = loader.LoadFromDBModuleInfoMap(); err != nil {
		log.Fatal("LoadFromDBModuleInfoMap failed: ", err)
	}
	log.Printf("LoadFromDBModuleInfoMap, duration: %s, count: %d",
		time.Since(st).String(),
		len(g.ModulesInfoMap),
	)

	// st = time.Now()
	// if g.InscriptionsApproveRemoveMap, err = loader.LoadFromDBSwapApproveStateMap(nil); err != nil {
	// 	log.Fatal("LoadFromDBSwapApproveStateMap failed: ", err)
	// }
	// log.Printf("LoadFromDBSwapApproveStateMap",
	// 	zap.String("duration", time.Since(st).String()),
	// 	zap.Int("count", len(g.InscriptionsApproveRemoveMap)),
	// )

	st = time.Now()
	if g.InscriptionsValidApproveMap, err = loader.LoadFromDBSwapApproveMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapApproveMap failed: ", err)
	}
	log.Printf("LoadFromDBSwapApproveMap, duration: %s, count: %d",
		time.Since(st).String(),
		len(g.InscriptionsValidApproveMap),
	)

	// st = time.Now()
	// if g.InscriptionsCondApproveRemoveMap, err = loader.LoadFromDBSwapCondApproveStateMap(nil); err != nil {
	// 	log.Fatal("LoadFromDBSwapCondApproveStateMap failed: ", err)
	// }
	// log.Printf("LoadFromDBSwapCondApproveStateMap, duration: %s, count: %d",
	// 	time.Since(st).String(),
	// 	len(g.InscriptionsCondApproveRemoveMap),
	// )

	st = time.Now()
	if g.InscriptionsValidConditionalApproveMap, err = loader.LoadFromDBSwapCondApproveMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCondApproveMap failed: ", err)
	}
	log.Printf("LoadFromDBSwapCondApproveMap, duration: %s, count: %d",
		time.Since(st).String(),
		len(g.InscriptionsValidConditionalApproveMap),
	)

	// st = time.Now()
	// if g.InscriptionsCommitRemoveMap, err = loader.LoadFromDBSwapCommitStateMap(nil); err != nil {
	// 	log.Fatal("LoadFromDBSwapCommitStateMap failed: ", err)
	// }
	// log.Printf("LoadFromDBSwapCommitStateMap, duration: %s, count: %d",
	// 	time.Since(st).String(),
	// 	len(g.InscriptionsCommitRemoveMap),
	// )

	st = time.Now()
	if g.InscriptionsValidCommitMap, err = loader.LoadFromDBSwapCommitMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCommitMap failed: ", err)
	}
	log.Printf("LoadFromDBSwapCommitMap, duration: %s, count: %d",
		time.Since(st).String(),
		len(g.InscriptionsValidCommitMap),
	)

	// st = time.Now()
	// if g.InscriptionsWithdrawRemoveMap, err = loader.LoadFromDBSwapWithdrawStateMap(nil); err != nil {
	// 	log.Fatal("LoadFromDBSwapWithdrawStateMap failed: ", err)
	// }
	// log.Printf("LoadFromDBSwapWithdrawStateMap, duration: %s, count: %d",
	//    time.Since(st).String(),
	// 	len(g.InscriptionsWithdrawRemoveMap),
	// )

	st = time.Now()
	if g.InscriptionsValidWithdrawMap, err = loader.LoadFromDBSwapWithdrawMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapWithdrawMap failed: ", err)
	}
	log.Printf("LoadFromDBSwapWithdrawMap, duration: %s, count: %d",
		time.Since(st).String(),
		len(g.InscriptionsValidWithdrawMap),
	)

	for mid, info := range g.ModulesInfoMap {
		log.Printf("loadFromDBSwapModuleInfo, moduleId: %s", mid)
		loadFromDBSwapModuleInfo(mid, info)
	}
}

func loadFromDBSwapModuleInfo(mid string, info *model.BRC20ModuleSwapInfo) {
	var st = time.Now()
	if hm, err := loader.LoadFromDBModuleHistoryMap(mid); err != nil {
		log.Fatal("LoadFromDBModuleHistoryMap failed: ", err)
	} else {
		log.Printf("LoadFromDBModuleHistoryMap, duration: %s, count: %d",
			time.Since(st).String(),
			len(hm),
		)
		for _, history := range hm {
			info.History = history
		}
	}

	st = time.Now()
	if ccs, err := loader.LoadModuleCommitChain(mid, nil); err != nil {
		log.Fatal("LoadModuleCommitChain failed: ", err)
	} else {
		log.Printf("LoadModuleCommitChain, duration: %s, count: %d",
			time.Since(st).String(),
			len(ccs))
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
	st = time.Now()
	if tabm, err := loader.LoadFromDBModuleUserBalanceMap(mid, nil, nil); err != nil {
		log.Fatal("LoadFromDBModuleUserBalanceMap failed: ", err)
	} else {
		info.TokenUsersBalanceDataMap = tabm
		info.UsersTokenBalanceDataMap = make(map[string]map[string]*model.BRC20ModuleTokenBalance)
		for tick, abs := range tabm {
			for addr, balance := range abs {
				if _, ok := info.UsersTokenBalanceDataMap[addr]; !ok {
					info.UsersTokenBalanceDataMap[addr] = make(map[string]*model.BRC20ModuleTokenBalance)
				}
				// [address][tick]balanceData
				info.UsersTokenBalanceDataMap[addr][tick] = balance
			}
		}

		log.Printf("LoadFromDBModuleUserBalanceMap, duration: %s, count: %d, addresses: %d",
			time.Since(st).String(),
			len(tabm),
			len(info.UsersTokenBalanceDataMap),
		)
	}

	st = time.Now()
	if poolBalanceMap, err := loader.LoadFromDBModulePoolLpBalanceMap(mid, nil); err != nil {
		log.Fatal("LoadFromDBModulePoolLpBalanceMap failed: ", err)
	} else {
		log.Printf("LoadFromDBModulePoolLpBalanceMap, duration: %s, count: %d",
			time.Since(st).String(),
			len(poolBalanceMap))
		info.SwapPoolTotalBalanceDataMap = poolBalanceMap
	}

	// [pool][address]balance
	st = time.Now()
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

		log.Printf("LoadFromDBModuleUserLpBalanceMap, duration: %s, count: %d, addresses: %d",
			time.Since(st).String(),
			len(userLpBalanceMap),
			len(info.UsersLPTokenBalanceMap),
		)
	}

}
