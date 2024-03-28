package indexer

import (
	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/loader"
	"brc20query/lib/brc20_swap/model"
	"brc20query/logger"
	"log"
	"time"

	"go.uber.org/zap"
)

// func (g *BRC20ModuleIndexer) SaveDataToDB(height int) {

func (g *BRC20ModuleIndexer) PurgeHistoricalData() {
	// purge history
	g.AllHistory = make([]*model.BRC20History, 0)
	g.InscriptionsTransferRemoveMap = make(map[string]uint32, 0)
	g.InscriptionsApproveRemoveMap = make(map[string]uint32, 0)
	g.InscriptionsCondApproveRemoveMap = make(map[string]uint32, 0)
	g.InscriptionsCommitRemoveMap = make(map[string]uint32, 0)
}

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

	var (
		err error
		st  time.Time
	)

	st = time.Now()
	if g.InscriptionsTickerInfoMap, err = loader.LoadFromDbTickerInfoMap(); err != nil {
		log.Fatal("LoadFromDBTickerInfoMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBTickerInfoMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsTickerInfoMap)),
	)

	st = time.Now()
	if g.UserTokensBalanceData, err = loader.LoadFromDbUserTokensBalanceData(nil, nil); err != nil {
		log.Fatal("LoadFromDBUserTokensBalanceData failed: ", err)
	}
	g.TokenUsersBalanceData = loader.UserTokensBalanceMap2TokenUsersBalanceMap(g.UserTokensBalanceData)
	logger.Log.Debug("LoadFromDBUserTokensBalanceData",
		zap.String("duration", time.Since(st).String()),
		zap.Int("ticks", len(g.TokenUsersBalanceData)),
		zap.Int("addresses", len(g.UserTokensBalanceData)),
	)

	st = time.Now()
	if g.InscriptionsTransferRemoveMap, err = loader.LoadFromDBTransferStateMap(); err != nil {
		log.Fatal("LoadFromDBTransferStateMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBTransferStateMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsTransferRemoveMap)),
	)

	st = time.Now()
	if g.InscriptionsValidTransferMap, err = loader.LoadFromDBValidTransferMap(); err != nil {
		log.Fatal("LoadFromDBvalidTransferMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBvalidTransferMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsValidTransferMap)),
	)

	st = time.Now()
	if g.ModulesInfoMap, err = loader.LoadFromDBModuleInfoMap(); err != nil {
		log.Fatal("LoadFromDBModuleInfoMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBModuleInfoMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.ModulesInfoMap)),
	)

	st = time.Now()
	if g.InscriptionsApproveRemoveMap, err = loader.LoadFromDBSwapApproveStateMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapApproveStateMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBSwapApproveStateMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsApproveRemoveMap)),
	)

	st = time.Now()
	if g.InscriptionsValidApproveMap, err = loader.LoadFromDBSwapApproveMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapApproveMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBSwapApproveMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsValidApproveMap)),
	)

	st = time.Now()
	if g.InscriptionsCondApproveRemoveMap, err = loader.LoadFromDBSwapCondApproveStateMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCondApproveStateMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBSwapCondApproveStateMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsCondApproveRemoveMap)),
	)

	st = time.Now()
	if g.InscriptionsValidConditionalApproveMap, err = loader.LoadFromDBSwapCondApproveMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCondApproveMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBSwapCondApproveMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsValidConditionalApproveMap)),
	)

	st = time.Now()
	if g.InscriptionsCommitRemoveMap, err = loader.LoadFromDBSwapCommitStateMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCommitStateMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBSwapCommitStateMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsCommitRemoveMap)),
	)

	st = time.Now()
	if g.InscriptionsValidCommitMap, err = loader.LoadFromDBSwapCommitMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapCommitMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBSwapCommitMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsValidCommitMap)),
	)

	st = time.Now()
	if g.InscriptionsWithdrawRemoveMap, err = loader.LoadFromDBSwapWithdrawStateMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapWithdrawStateMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBSwapWithdrawStateMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsWithdrawRemoveMap)),
	)

	st = time.Now()
	if g.InscriptionsValidWithdrawMap, err = loader.LoadFromDBSwapWithdrawMap(nil); err != nil {
		log.Fatal("LoadFromDBSwapWithdrawMap failed: ", err)
	}
	logger.Log.Debug("LoadFromDBSwapWithdrawMap",
		zap.String("duration", time.Since(st).String()),
		zap.Int("count", len(g.InscriptionsValidWithdrawMap)),
	)

	for mid, info := range g.ModulesInfoMap {
		logger.Log.Debug("loadFromDBSwapModuleInfo", zap.String("moduleId", mid))
		loadFromDBSwapModuleInfo(mid, info)
	}
}

func loadFromDBSwapModuleInfo(mid string, info *model.BRC20ModuleSwapInfo) {
	var st = time.Now()
	if hm, err := loader.LoadFromDBModuleHistoryMap(mid); err != nil {
		log.Fatal("LoadFromDBModuleHistoryMap failed: ", err)
	} else {
		logger.Log.Debug("LoadFromDBModuleHistoryMap",
			zap.String("duration", time.Since(st).String()), zap.Int("count", len(hm)),
		)
		for _, history := range hm {
			info.History = history
		}
	}

	st = time.Now()
	if ccs, err := loader.LoadModuleCommitChain(mid, nil); err != nil {
		log.Fatal("LoadModuleCommitChain failed: ", err)
	} else {
		logger.Log.Debug("LoadModuleCommitChain",
			zap.String("duration", time.Since(st).String()), zap.Int("count", len(ccs)))
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

		logger.Log.Debug("LoadFromDBModuleUserBalanceMap",
			zap.String("duration", time.Since(st).String()),
			zap.Int("ticks", len(tabm)),
			zap.Int("addresses", len(info.UsersTokenBalanceDataMap)),
		)
	}

	st = time.Now()
	if poolBalanceMap, err := loader.LoadFromDBModulePoolLpBalanceMap(mid, nil); err != nil {
		log.Fatal("LoadFromDBModulePoolLpBalanceMap failed: ", err)
	} else {
		logger.Log.Debug("LoadFromDBModulePoolLpBalanceMap",
			zap.String("duration", time.Since(st).String()), zap.Int("count", len(poolBalanceMap)))
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

		logger.Log.Debug("LoadFromDBModuleUserLpBalanceMap",
			zap.String("duration", time.Since(st).String()),
			zap.Int("pools", len(userLpBalanceMap)),
			zap.Int("addresses", len(info.UsersLPTokenBalanceMap)),
		)
	}

}
