package loader

import (
	"fmt"
	"log"
	"strings"

	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/model"
)

func buildSQLWhereInStr(colVals map[string][]string) (conds []string, args []any) {
	conds = make([]string, 0)
	args = make([]any, 0)
	argIndex := 1

	for col, vals := range colVals {
		if len(vals) == 0 {
			continue
		}

		phs := make([]string, 0, len(vals))
		for _, val := range vals {
			phs = append(phs, fmt.Sprintf("$%d", argIndex))
			args = append(args, val)
			argIndex += 1
		}
		conds = append(conds, fmt.Sprintf("%s IN (%s)", col, strings.Join(phs, ",")))
	}

	return conds, args
}

func LoadFromDbTickerInfoMap() (map[string]*model.BRC20TokenInfo, error) {
	rows, err := SwapDB.Query(`
SELECT t1.block_height, t1.tick, t1.max_supply, t1.decimals, t1.limit_per_mint, t1.remaining_supply, t1.pkscript_deployer
FROM brc20_ticker_info t1 
INNER JOIN (
	SELECT MAX(block_height) as block_height, tick FROM brc20_ticker_info GROUP BY tick
) t2 ON t1.block_height = t2.block_height AND t1.tick = t2.tick;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		height    int
		tick      string
		max       string
		decimals  uint8
		limit     string
		remaining string
		pkscript  string
	)
	tickerInfoMap := make(map[string]*model.BRC20TokenInfo)

	for rows.Next() {
		if err := rows.Scan(&height, &tick, &max, &decimals, &limit, &remaining, &pkscript); err != nil {
			return nil, err
		}

		nremaining := decimal.MustNewDecimalFromString(remaining, int(decimals))
		nmax := decimal.MustNewDecimalFromString(max, int(decimals))
		minted := nmax.Sub(nremaining)

		tickerInfoMap[tick] = &model.BRC20TokenInfo{
			Ticker: tick,
			Deploy: &model.InscriptionBRC20TickInfo{
				Max:         nmax,
				Decimal:     decimals,
				Limit:       decimal.MustNewDecimalFromString(limit, int(decimals)),
				TotalMinted: minted,
			},
		}
	}

	return tickerInfoMap, nil
}

func LoadFromDbUserTokensBalanceData(pkscripts, ticks []string) (map[string]map[string]*model.BRC20TokenBalance, error) {
	inConds, inCondArgs := buildSQLWhereInStr(map[string][]string{
		"pkscript": pkscripts,
		"tick":     ticks,
	})
	condSql := ""
	if len(inConds) > 0 {
		condSql = "WHERE " + strings.Join(inConds, " AND ")
	}

	sql := fmt.Sprintf(`
SELECT t1.tick, t1.pkscript, t1.block_height, t1.available_balance, t1.transferable_balance
FROM brc20_user_balance t1
INNER JOIN (
	SELECT tick, pkscript, MAX(block_height) AS max_block_height
	FROM brc20_user_balance %s GROUP BY tick, pkscript
) t2 ON t1.tick = t2.tick AND t1.pkscript = t2.pkscript AND t1.block_height = t2.max_block_height;
`, condSql)
	args := inCondArgs

	log.Printf("sql: %s", sql)
	log.Printf("args: %v", args)

	rows, err := SwapDB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		tick         string
		pkscript     string
		height       int
		available    string
		transferable string
	)
	userTokensBalanceMap := make(map[string]map[string]*model.BRC20TokenBalance)

	for rows.Next() {
		if err := rows.Scan(&tick, &pkscript, &height, &available, &transferable); err != nil {
			return nil, err
		}

		balance := &model.BRC20TokenBalance{
			Ticker:              tick,
			PkScript:            pkscript,
			AvailableBalance:    decimal.MustNewDecimalFromString(available, 0),
			TransferableBalance: decimal.MustNewDecimalFromString(transferable, 0),
		}

		if _, ok := userTokensBalanceMap[tick]; !ok {
			userTokensBalanceMap[tick] = make(map[string]*model.BRC20TokenBalance)
		}

		userTokensBalanceMap[pkscript][tick] = balance
	}

	return userTokensBalanceMap, nil
}

func UserTokensBalanceMap2TokenUsersBalanceMap(userTokensMap map[string]map[string]*model.BRC20TokenBalance) map[string]map[string]*model.BRC20TokenBalance {
	tokenUsersMap := make(map[string]map[string]*model.BRC20TokenBalance)

	for pkscript, userTokensBalance := range userTokensMap {
		for tick, balance := range userTokensBalance {
			if _, ok := tokenUsersMap[tick]; !ok {
				tokenUsersMap[tick] = make(map[string]*model.BRC20TokenBalance)
			}

			tokenUsersMap[tick][pkscript] = balance
		}
	}
	return tokenUsersMap
}

func LoadFromDBTransferStateMap() (res map[string]struct{}, err error) {
	rows, err := SwapDB.Query(`
SELECT t1.block_height, t1.create_key FROM brc20_transfer_state  t1 
INNER JOIN (
	SELECT MAX(block_height) as block_height, create_key FROM brc20_transfer_state GROUP BY create_key
) t2 ON t1.block_height = t2.block_height AND t1.create_key = t2.create_key
WHERE t1.moved = true;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		height     int
		create_key string
	)
	for rows.Next() {
		if err := rows.Scan(&height, &create_key); err != nil {
			return nil, err
		}
		res[create_key] = struct{}{}
	}

	return res, nil
}

func LoadFromDBValidTransferMap() (res map[string]*model.InscriptionBRC20TickInfo, err error) {
	rows, err := SwapDB.Query(`
SELECT t1.block_height, t1.create_key, t1.tick, t1.pkscript, t1.amount, 
	   t1.inscription_number, t1.inscription_id, 
	   t1.txid, t1.vout, t1.output_value, t1.output_offset
FROM brc20_valid_transfer t1
INNER JOIN (
	SELECT MAX(block_height) as block_height, create_key FROM brc20_valid_transfer GROUP BY create_key
) t2 ON t1.block_height = t2.block_height AND t1.create_key = t2.create_key;
`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res = make(map[string]*model.InscriptionBRC20TickInfo)
	for rows.Next() {
		meta := model.InscriptionBRC20Data{
			// TODO: other meta attrs
		}
		t := model.InscriptionBRC20TickInfo{
			Meta: &meta,
		}
		if err := rows.Scan(&t.Height, &t.CreateIdxKey, &t.Tick, &t.PkScript, &t.Amount,
			&t.InscriptionNumber, &meta.InscriptionId,
			&t.TxId, &t.Vout, &t.Satoshi, &t.Offset,
		); err != nil {
			return nil, err
		}
		res[t.CreateIdxKey] = &t
		log.Println("amount", t.Amount.String())
	}
	return res, nil
}
