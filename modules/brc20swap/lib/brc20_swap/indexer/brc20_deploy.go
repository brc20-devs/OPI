package indexer

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"brc20query/lib/brc20_swap/constant"
	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/model"
)

func (g *BRC20ModuleIndexer) ProcessDeploy(data *model.InscriptionBRC20Data) error {
	body := new(model.InscriptionBRC20DeployContent)
	if err := body.Unmarshal(data.ContentBody); err != nil {
		return nil
	}

	// check tick/amt
	if len(body.BRC20Tick) != 4 {
		return nil
		// return errors.New("deploy, tick length not 4")
	}
	uniqueLowerTicker := strings.ToLower(body.BRC20Tick)

	// tick enable, fixme
	if constant.TICKS_ENABLED != "" {
		if strings.Contains(uniqueLowerTicker, " ") {
			return nil
		}
		if !strings.Contains(constant.TICKS_ENABLED, uniqueLowerTicker) {
			return nil
		}
	}

	if _, ok := g.InscriptionsTickerInfoMap[uniqueLowerTicker]; ok { // dup ticker
		return nil
		// return errors.New("deploy, but tick exist")
	}
	if body.BRC20Max == "" { // without max
		log.Printf("deploy, but max missing. ticker: %s",
			uniqueLowerTicker,
		)
		return errors.New("deploy, but max missing")
	}

	tinfo := model.NewInscriptionBRC20TickInfo(body.BRC20Tick, body.Operation, data)
	tinfo.Data.BRC20Max = body.BRC20Max
	tinfo.Data.BRC20Limit = body.BRC20Limit
	tinfo.Data.BRC20Decimal = body.BRC20Decimal
	tinfo.Data.BRC20Minted = "0"
	tinfo.InscriptionNumberStart = data.InscriptionNumber

	// dec
	if dec, err := strconv.ParseUint(tinfo.Data.BRC20Decimal, 10, 64); err != nil || dec > 18 {
		// dec invalid
		log.Printf("deploy, but dec invalid. ticker: %s, dec: %s",
			uniqueLowerTicker,
			tinfo.Data.BRC20Decimal,
		)
		return errors.New("deploy, but dec invalid")
	} else {
		tinfo.Decimal = uint8(dec)
	}

	// max
	if max, err := decimal.NewDecimalFromString(body.BRC20Max, int(tinfo.Decimal)); err != nil {
		// max invalid
		log.Printf("deploy, but max invalid. ticker: %s, max: '%s'",
			uniqueLowerTicker,
			body.BRC20Max,
		)
		return errors.New("deploy, but max invalid")
	} else {
		if max.Sign() <= 0 || max.IsOverflowUint64() {
			return errors.New("deploy, but max invalid (range)")
		}
		tinfo.Max = max
	}

	// lim
	if lim, err := decimal.NewDecimalFromString(tinfo.Data.BRC20Limit, int(tinfo.Decimal)); err != nil {
		// limit invalid
		log.Printf("deploy, but limit invalid. ticker: %s, limit: '%s'",
			uniqueLowerTicker,
			tinfo.Data.BRC20Limit,
		)
		return errors.New("deploy, but lim invalid")
	} else {
		if lim.Sign() <= 0 || lim.IsOverflowUint64() {
			return errors.New("deploy, but lim invalid (range)")
		}
		tinfo.Limit = lim
	}

	tokenInfo := &model.BRC20TokenInfo{Ticker: body.BRC20Tick, Deploy: tinfo}
	g.InscriptionsTickerInfoMap[uniqueLowerTicker] = tokenInfo

	tokenBalance := &model.BRC20TokenBalance{Ticker: body.BRC20Tick, PkScript: data.PkScript}

	history := model.NewBRC20History(constant.BRC20_HISTORY_TYPE_N_INSCRIBE_DEPLOY, true, false, tinfo, nil, data)
	tokenBalance.History = append(tokenBalance.History, history)
	tokenInfo.History = append(tokenInfo.History, history)

	// user history
	userHistory := g.GetBRC20HistoryByUser(string(data.PkScript))
	userHistory.History = append(userHistory.History, history)
	// all history
	g.AllHistory = append(g.AllHistory, history)

	// init user tokens
	var userTokens map[string]*model.BRC20TokenBalance
	if tokens, ok := g.UserTokensBalanceData[string(data.PkScript)]; !ok {
		userTokens = make(map[string]*model.BRC20TokenBalance, 0)
		g.UserTokensBalanceData[string(data.PkScript)] = userTokens
	} else {
		userTokens = tokens
	}
	userTokens[uniqueLowerTicker] = tokenBalance

	// init token users
	tokenUsers := make(map[string]*model.BRC20TokenBalance, 0)
	tokenUsers[string(data.PkScript)] = tokenBalance
	g.TokenUsersBalanceData[uniqueLowerTicker] = tokenUsers

	// g.InscriptionsValidBRC20DataMap[data.CreateIdxKey] = tinfo.Data
	return nil
}
