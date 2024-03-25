package indexer

import (
	"brc20query/lib/brc20_swap/loader"
)

func (g *BRC20ModuleIndexer) SaveDataToDB(height int) {

	loader.SaveDataToDBTickerInfoMap(height,
		g.InscriptionsTickerInfoMap,
		g.UserTokensBalanceData,
		g.TokenUsersBalanceData,
	)

	loader.SaveDataToDBModuleInfoMap(height, g.ModulesInfoMap)
}

func (g *BRC20ModuleIndexer) LoadDataFromDB(height int) {

	loader.LoadDataToDBTickerInfoMap(height,
		g.InscriptionsTickerInfoMap,
		g.UserTokensBalanceData,
		g.TokenUsersBalanceData,
	)

	loader.LoadDataToDBModuleInfoMap(height, g.ModulesInfoMap)
}
