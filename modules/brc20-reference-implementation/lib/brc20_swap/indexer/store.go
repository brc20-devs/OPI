package indexer

import (
	"brc20query/lib/brc20_swap/loader"
	"log"
)

// func (g *BRC20ModuleIndexer) SaveDataToDB(height int) {

func (g *BRC20ModuleIndexer) SaveDataToDB(dbConnInfo string, height int) {
	loader.Init(dbConnInfo)
	defer loader.SwapDB.Close()

	// ticker info
	loader.SaveDataToDBTickerInfoMap(height, g.InscriptionsTickerInfoMap)
	loader.SaveDataToDBTickerBalanceMap(height, g.TokenUsersBalanceData)
	loader.SaveDataToDBTickerHistoryMap(height, g.InscriptionsTickerInfoMap)

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
	// TODO: g.InscriptionsTickerInfoMap.History？

	if g.UserTokensBalanceData, err = loader.LoadFromDbUserTokensBalanceData(nil, nil); err != nil {
		log.Fatal("LoadFromDBUserTokensBalanceData failed: ", err)
	}
	g.TokenUsersBalanceData = loader.UserTokensBalanceMap2TokenUsersBalanceMap(g.UserTokensBalanceData)

	if g.InscriptionsTransferRemoveMap, err = loader.LoadFromDBTransferStateMap(); err != nil {
		log.Fatal("LoadFromDBTransferStateMap failed: ", err)
	}

	if g.InscriptionsInvalidTransferMap, err = loader.LoadFromDBValidTransferMap(); err != nil {
		log.Fatal("LoadFromDBvalidTransferMap failed: ", err)
	}
}
