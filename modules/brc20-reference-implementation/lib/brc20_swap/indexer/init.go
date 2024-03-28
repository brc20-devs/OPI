package indexer

import (
	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/model"
	"log"
	"strings"
)

type BRC20ModuleIndexer struct {
	// brc20 base
	AllHistory     []*model.BRC20History
	UserAllHistory map[string]*model.BRC20UserHistory

	InscriptionsTickerInfoMap     map[string]*model.BRC20TokenInfo
	UserTokensBalanceData         map[string]map[string]*model.BRC20TokenBalance // [address][ticker]balance
	TokenUsersBalanceData         map[string]map[string]*model.BRC20TokenBalance // [ticker][address]balance
	InscriptionsValidBRC20DataMap map[string]*model.InscriptionBRC20InfoResp

	// inner valid transfer
	InscriptionsTransferRemoveMap map[string]struct{} // todo
	InscriptionsValidTransferMap  map[string]*model.InscriptionBRC20TickInfo
	// inner invalid transfer
	InscriptionsInvalidTransferMap map[string]*model.InscriptionBRC20TickInfo

	// module
	// all modules info
	ModulesInfoMap map[string]*model.BRC20ModuleSwapInfo

	// module of users [address]moduleid
	UsersModuleWithTokenMap map[string]string

	// module lp of users [address]moduleid
	UsersModuleWithLpTokenMap map[string]string

	// runtime for approve
	InscriptionsApproveRemoveMap  map[string]struct{}                        // todo
	InscriptionsValidApproveMap   map[string]*model.InscriptionBRC20SwapInfo // inner valid approve [create_key]
	InscriptionsInvalidApproveMap map[string]*model.InscriptionBRC20SwapInfo

	// runtime for conditional approve
	InscriptionsCondApproveRemoveMap         map[string]struct{} // todo
	InscriptionsValidConditionalApproveMap   map[string]*model.InscriptionBRC20SwapConditionalApproveInfo
	InscriptionsInvalidConditionalApproveMap map[string]*model.InscriptionBRC20SwapConditionalApproveInfo

	// runtime for commit
	InscriptionsCommitRemoveMap  map[string]struct{}                    // todo
	InscriptionsValidCommitMap   map[string]*model.InscriptionBRC20Data // inner valid commit by key
	InscriptionsInvalidCommitMap map[string]*model.InscriptionBRC20Data

	InscriptionsValidCommitMapById map[string]*model.InscriptionBRC20Data // inner valid commit by id

	// runtime for withdraw
	InscriptionsWithdrawRemoveMap  map[string]struct{}                        // todo
	InscriptionsValidWithdrawMap   map[string]*model.InscriptionBRC20SwapInfo // inner valid withdraw by key
	InscriptionsInvalidWithdrawMap map[string]*model.InscriptionBRC20SwapInfo

	// for gen approve event
	ThisTxId                                    string
	TxStaticTransferStatesForConditionalApprove []*model.TransferStateForConditionalApprove

	// save flag
	Durty bool
}

func (g *BRC20ModuleIndexer) GetBRC20HistoryByUser(pkScript string) (userHistory *model.BRC20UserHistory) {
	if history, ok := g.UserAllHistory[pkScript]; !ok {
		userHistory = &model.BRC20UserHistory{}
		g.UserAllHistory[pkScript] = userHistory
	} else {
		userHistory = history
	}
	return userHistory
}

func (g *BRC20ModuleIndexer) GetBRC20HistoryByUserForAPI(pkScript string) (userHistory *model.BRC20UserHistory) {
	if history, ok := g.UserAllHistory[pkScript]; !ok {
		userHistory = &model.BRC20UserHistory{}
	} else {
		userHistory = history
	}
	return userHistory
}

func (g *BRC20ModuleIndexer) initBRC20() {
	// all history
	g.AllHistory = make([]*model.BRC20History, 0)

	// user history
	g.UserAllHistory = make(map[string]*model.BRC20UserHistory, 0)

	// all ticker info
	g.InscriptionsTickerInfoMap = make(map[string]*model.BRC20TokenInfo, 0)

	// ticker of users
	g.UserTokensBalanceData = make(map[string]map[string]*model.BRC20TokenBalance, 0)

	// ticker holders
	g.TokenUsersBalanceData = make(map[string]map[string]*model.BRC20TokenBalance, 0)

	// valid brc20 inscriptions
	g.InscriptionsValidBRC20DataMap = make(map[string]*model.InscriptionBRC20InfoResp, 0)

	// inner valid transfer
	g.InscriptionsValidTransferMap = make(map[string]*model.InscriptionBRC20TickInfo, 0)
	// inner invalid transfer
	g.InscriptionsInvalidTransferMap = make(map[string]*model.InscriptionBRC20TickInfo, 0)
}

func (g *BRC20ModuleIndexer) initModule() {
	// all modules info
	g.ModulesInfoMap = make(map[string]*model.BRC20ModuleSwapInfo, 0)

	// module of users [address]moduleid
	g.UsersModuleWithTokenMap = make(map[string]string, 0)

	// swap
	// module of users [address]moduleid
	g.UsersModuleWithLpTokenMap = make(map[string]string, 0)

	// runtime for approve
	g.InscriptionsValidApproveMap = make(map[string]*model.InscriptionBRC20SwapInfo, 0)
	g.InscriptionsInvalidApproveMap = make(map[string]*model.InscriptionBRC20SwapInfo, 0)

	// runtime for conditional approve
	g.InscriptionsValidConditionalApproveMap = make(map[string]*model.InscriptionBRC20SwapConditionalApproveInfo, 0)
	g.InscriptionsInvalidConditionalApproveMap = make(map[string]*model.InscriptionBRC20SwapConditionalApproveInfo, 0)

	// runtime for commit
	g.InscriptionsValidCommitMap = make(map[string]*model.InscriptionBRC20Data, 0) // inner valid commit
	g.InscriptionsInvalidCommitMap = make(map[string]*model.InscriptionBRC20Data, 0)

	g.InscriptionsValidCommitMapById = make(map[string]*model.InscriptionBRC20Data, 0) // inner valid commit
}

func (g *BRC20ModuleIndexer) GenerateApproveEventsByTransfer(inscription, tick, from, to string, amt *decimal.Decimal) (events []*model.ConditionalApproveEvent) {
	transStateStatic := &model.TransferStateForConditionalApprove{
		Tick:          tick,
		From:          from,
		To:            to,
		Balance:       decimal.NewDecimalCopy(amt), // maybe no need copy
		InscriptionId: inscription,
		Max:           amt.String(),
	}
	// First, globally save the transfer status.
	g.TxStaticTransferStatesForConditionalApprove = append(g.TxStaticTransferStatesForConditionalApprove, transStateStatic)

	// Then process each module one by one.
	for _, moduleInfo := range g.ModulesInfoMap {
		if g.ThisTxId != moduleInfo.ThisTxId {
			// For the first time processing the transfer event within the module, you need to clear the status first.
			moduleInfo.TransferStatesForConditionalApprove = nil
			moduleInfo.ApproveStatesForConditionalApprove = nil
			moduleInfo.ThisTxId = g.ThisTxId
		}

		// Skip processing the transfer directly when there is no approve status.
		if len(moduleInfo.ApproveStatesForConditionalApprove) == 0 {
			continue
		}

		transState := &model.TransferStateForConditionalApprove{
			Tick:          tick,
			From:          from,
			To:            to,
			Balance:       decimal.NewDecimalCopy(amt), // maybe no need copy
			InscriptionId: inscription,
			Max:           amt.String(),
		}

		innerEvents := moduleInfo.GenerateApproveEventsByTransfer(transState)
		events = append(events, innerEvents...)
	}
	return events
}

func (g *BRC20ModuleIndexer) GenerateApproveEventsByApprove(owner string, balance *decimal.Decimal,
	data *model.InscriptionBRC20Data, approveInfo *model.InscriptionBRC20SwapConditionalApproveInfo) (events []*model.ConditionalApproveEvent) {
	if moduleInfo, ok := g.ModulesInfoMap[approveInfo.Module]; ok {
		log.Printf("generate approve event. module: %s", moduleInfo.ID)

		if g.ThisTxId != moduleInfo.ThisTxId {
			// First appearance, clear status
			moduleInfo.TransferStatesForConditionalApprove = nil
			moduleInfo.ApproveStatesForConditionalApprove = nil
			moduleInfo.ThisTxId = g.ThisTxId
			log.Printf("generate approve event. init")
		}

		// First appearance of approve, copy all global transfer events.
		if len(moduleInfo.ApproveStatesForConditionalApprove) == 0 {
			moduleInfo.TransferStatesForConditionalApprove = nil
			for _, s := range g.TxStaticTransferStatesForConditionalApprove {
				moduleInfo.TransferStatesForConditionalApprove = append(moduleInfo.TransferStatesForConditionalApprove, s)
			}
			log.Printf("generate approve event. copy transfer")
		}

		log.Printf("generate approve event. balance: %s", balance.String())
		innerEvents := moduleInfo.GenerateApproveEventsByApprove(owner, balance, data, approveInfo)
		events = append(events, innerEvents...)
	}
	return events
}

func (copyDup *BRC20ModuleIndexer) deepCopyBRC20Data(base *BRC20ModuleIndexer) {
	// history
	// for _, h := range base.AllHistory {
	// 	copyDup.AllHistory = append(copyDup.AllHistory, h)
	// }

	for k, v := range base.InscriptionsTickerInfoMap {
		tinfo := &model.BRC20TokenInfo{
			Ticker: v.Ticker,
			Deploy: v.Deploy.DeepCopy(),
		}

		// history
		tinfo.History = make([]*model.BRC20History, len(v.History))
		copy(tinfo.History, v.History)

		tinfo.HistoryMint = make([]*model.BRC20History, len(v.HistoryMint))
		copy(tinfo.HistoryMint, v.HistoryMint)

		tinfo.HistoryInscribeTransfer = make([]*model.BRC20History, len(v.HistoryInscribeTransfer))
		copy(tinfo.HistoryInscribeTransfer, v.HistoryInscribeTransfer)

		tinfo.HistoryTransfer = make([]*model.BRC20History, len(v.HistoryTransfer))
		copy(tinfo.HistoryTransfer, v.HistoryTransfer)

		// set info
		copyDup.InscriptionsTickerInfoMap[k] = tinfo
	}

	for u, userTokens := range base.UserTokensBalanceData {
		userTokensCopy := make(map[string]*model.BRC20TokenBalance, 0)
		copyDup.UserTokensBalanceData[u] = userTokensCopy
		for uniqueLowerTicker, v := range userTokens {
			tb := v.DeepCopy()
			userTokensCopy[uniqueLowerTicker] = tb

			tokenUsers, ok := copyDup.TokenUsersBalanceData[uniqueLowerTicker]
			if !ok {
				tokenUsers = make(map[string]*model.BRC20TokenBalance, 0)
				copyDup.TokenUsersBalanceData[uniqueLowerTicker] = tokenUsers
			}
			tokenUsers[u] = tb
		}
	}

	for k, v := range base.InscriptionsValidBRC20DataMap {
		copyDup.InscriptionsValidBRC20DataMap[k] = v
	}

	// transferInfo
	for k, v := range base.InscriptionsValidTransferMap {
		copyDup.InscriptionsValidTransferMap[k] = v
	}

	log.Printf("deepCopyBRC20Data finish. total: %d", len(base.InscriptionsTickerInfoMap))
}

func (copyDup *BRC20ModuleIndexer) cherryPickBRC20Data(base *BRC20ModuleIndexer, pickUsersPkScript, pickTokensTick map[string]bool) {

	for lowerTick := range pickTokensTick {
		v, ok := base.InscriptionsTickerInfoMap[lowerTick]
		if !ok {
			continue
		}

		tinfo := &model.BRC20TokenInfo{
			Ticker: v.Ticker,
			Deploy: v.Deploy.DeepCopy(),
		}
		copyDup.InscriptionsTickerInfoMap[lowerTick] = tinfo
	}
	for u := range pickUsersPkScript {
		userTokens, ok := base.UserTokensBalanceData[u]
		if !ok {
			continue
		}
		userTokensCopy := make(map[string]*model.BRC20TokenBalance, 0)
		for lowerTick := range pickTokensTick {
			balance, ok := userTokens[lowerTick]
			if !ok {
				continue
			}
			userTokensCopy[lowerTick] = balance.DeepCopy()
		}
		copyDup.UserTokensBalanceData[u] = userTokensCopy
	}

	for u, userTokens := range copyDup.UserTokensBalanceData {
		for uniqueLowerTicker, balance := range userTokens {
			tokenUsers, ok := copyDup.TokenUsersBalanceData[uniqueLowerTicker]
			if !ok {
				tokenUsers = make(map[string]*model.BRC20TokenBalance, 0)
				copyDup.TokenUsersBalanceData[uniqueLowerTicker] = tokenUsers
			}
			tokenUsers[u] = balance
		}
	}

	log.Printf("cherryPickBRC20Data finish. total: %d", len(copyDup.InscriptionsTickerInfoMap))
}

func (copyDup *BRC20ModuleIndexer) deepCopyModuleData(base *BRC20ModuleIndexer) {

	for module, info := range base.ModulesInfoMap {
		copyDup.ModulesInfoMap[module] = info.DeepCopy()
	}

	// module of users
	for k, v := range base.UsersModuleWithTokenMap {
		copyDup.UsersModuleWithTokenMap[k] = v
	}

	// module lp of users
	for k, v := range base.UsersModuleWithLpTokenMap {
		copyDup.UsersModuleWithLpTokenMap[k] = v
	}

	// approveInfo
	for k, v := range base.InscriptionsValidApproveMap {
		copyDup.InscriptionsValidApproveMap[k] = v
	}
	for k, v := range base.InscriptionsInvalidApproveMap {
		copyDup.InscriptionsInvalidApproveMap[k] = v
	}

	// conditional approveInfo
	for k, v := range base.InscriptionsValidConditionalApproveMap {
		copyDup.InscriptionsValidConditionalApproveMap[k] = v.DeepCopy()
	}
	for k, v := range base.InscriptionsInvalidConditionalApproveMap {
		copyDup.InscriptionsInvalidConditionalApproveMap[k] = v.DeepCopy()
	}

	// commitInfo
	for k, v := range base.InscriptionsValidCommitMap {
		copyDup.InscriptionsValidCommitMap[k] = v
	}
	for k, v := range base.InscriptionsInvalidCommitMap {
		copyDup.InscriptionsInvalidCommitMap[k] = v
	}

	for k, v := range base.InscriptionsValidCommitMapById {
		copyDup.InscriptionsValidCommitMapById[k] = v
	}

	// runtime state
	copyDup.ThisTxId = base.ThisTxId
	for _, v := range base.TxStaticTransferStatesForConditionalApprove {
		copyDup.TxStaticTransferStatesForConditionalApprove = append(copyDup.TxStaticTransferStatesForConditionalApprove, v.DeepCopy())
	}

	log.Printf("deepCopyModuleData finish. total: %d", len(base.ModulesInfoMap))
}

func (copyDup *BRC20ModuleIndexer) cherryPickModuleData(base *BRC20ModuleIndexer, module string, pickUsersPkScript, pickTokensTick, pickPoolsPair map[string]bool) {

	info, ok := base.ModulesInfoMap[module]
	if ok {
		copyDup.ModulesInfoMap[module] = info.CherryPick(pickUsersPkScript, pickTokensTick, pickPoolsPair)
	}

	// Data required for verification
	for k, v := range base.InscriptionsValidCommitMapById {
		copyDup.InscriptionsValidCommitMapById[k] = v
	}
	log.Printf("cherryPickModuleData finish. total: %d", len(base.ModulesInfoMap))
}

func (g *BRC20ModuleIndexer) Init() {
	g.initBRC20()
	g.initModule()
}

func (base *BRC20ModuleIndexer) DeepCopy() (copyDup *BRC20ModuleIndexer) {
	copyDup = &BRC20ModuleIndexer{}
	copyDup.initBRC20()
	copyDup.initModule()

	copyDup.deepCopyBRC20Data(base)
	copyDup.deepCopyModuleData(base)
	return copyDup
}

func (base *BRC20ModuleIndexer) CherryPick(module string, pickUsersPkScript, pickTokensTick, pickPoolsPair map[string]bool) (copyDup *BRC20ModuleIndexer) {
	copyDup = &BRC20ModuleIndexer{}
	copyDup.initBRC20()
	copyDup.initModule()

	moduleInfo, ok := base.ModulesInfoMap[module]
	if ok {
		lowerTick := strings.ToLower(moduleInfo.GasTick)
		pickTokensTick[lowerTick] = true
	}
	copyDup.cherryPickBRC20Data(base, pickUsersPkScript, pickTokensTick)
	copyDup.cherryPickModuleData(base, module, pickUsersPkScript, pickTokensTick, pickPoolsPair)
	return copyDup
}
