package loader

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"brc20query/lib/brc20_swap/model"
)

var SwapDB *sql.DB

const (
	pg_host     = "10.16.11.95"
	pg_port     = 5432
	pg_user     = "postgres"
	pg_password = "postgres"
	pg_dbname   = "postgres"
)

func Init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg_host, pg_port, pg_user, pg_password, pg_dbname)
	var err error
	SwapDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Connect PG Failed: ", err)
	}

	SwapDB.SetMaxOpenConns(2000)
	SwapDB.SetMaxIdleConns(1000)
}

// brc20_ticker_info
func SaveDataToDBTickerInfoMap(height int,
	inscriptionsTickerInfoMap map[string]*model.BRC20TokenInfo,
) {
	stmtTickerInfo, err := SwapDB.Prepare(`
INSERT INTO brc20_ticker_info(block_height, tick, max_supply, decimals, limit_per_mint, remaining_supply, pkscript_deployer)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	for ticker, info := range inscriptionsTickerInfoMap {
		// save ticker info
		res, err := stmtTickerInfo.Exec(height, info.Ticker,
			info.Deploy.Max.String(),
			info.Deploy.Decimal,
			info.Deploy.Limit.String(),
			info.Deploy.Max.Sub(info.Deploy.TotalMinted).String(),
			info.Deploy.PkScript,
		)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}
}

func SaveDataToDBTickerBalanceMap(height int,
	tokenUsersBalanceData map[string]map[string]*model.BRC20TokenBalance,
) {
	stmtUserBalance, err := SwapDB.Prepare(`
INSERT INTO brc20_user_balance(block_height, tick, pkscript, available_balance, transferable_balance)
VALUES ($1, $2, $3, $4, $5)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for ticker, info := range inscriptionsTickerInfoMap {
		// holders
		for holder, balanceData := range tokenUsersBalanceData[ticker] {
			// save balance db
			res, err := stmtUserBalance.Exec(height, info.Ticker,
				balanceData.PkScript,
				balanceData.AvailableBalance.String(),
				balanceData.TransferableBalance.String(),
			)
			if err != nil {
				log.Fatal("PG Statements Exec Wrong: ", err)
			}
			id, err := res.RowsAffected()
			if err != nil {
				log.Fatal("PG Affecte Wrong: ", err)
			}
			fmt.Println(id)
		}
	}
}

func SaveDataToDBTickerHistoryMap(height int,
	inscriptionsTickerInfoMap map[string]*model.BRC20TokenInfo,
) {
	stmtBRC20History, err := SwapDB.Prepare(`
INSERT INTO brc20_history(block_height, tick,
    history_type,
    valid,
    txid,
    idx,
    vout,
    output_value,
    output_offset,
    pkscript_from,
    pkscript_to,
    fee,
    txidx,
    block_time,
    inscription_number,
    inscription_id,
    inscription_content,
	 amount,
	 available_balance,
	 transferable_balance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for ticker, info := range inscriptionsTickerInfoMap {
		nValid := 0
		for _, h := range info.History {
			if h.Valid {
				nValid++
			}
		}

		// history
		for _, h := range info.History {
			if !h.Valid {
				continue
			}

			{
				res, err := stmtBRC20History.Exec(height, info.Ticker,
					h.Type, h.Valid,
					h.TxId, h.Idx, h.Vout, h.Satoshi, h.Offset,
					h.PkScriptFrom, h.PkScriptTo,
					h.Fee,
					h.TxIdx, h.BlockTime,
					h.Inscription.InscriptionNumber, h.Inscription.InscriptionId,
					"", // content
					h.Amount.String(), h.AvailableBalance.String(), h.TransferableBalance.String(),
				)
				if err != nil {
					log.Fatal("PG Statements Exec Wrong: ", err)
				}
				id, err := res.RowsAffected()
				if err != nil {
					log.Fatal("PG Affecte Wrong: ", err)
				}
				fmt.Println(id)
			}
		}
	}
}

func SaveDataToDBTransferStateMap(height int,
	inscriptionsTransferRemoveMap map[string]struct{},
) {
	stmtTransferState, err := SwapDB.Prepare(`
INSERT INTO brc20_transfer_state(block_height, create_key, moved)
VALUES ($1, $2, $3)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for createKey := range inscriptionsTransferRemoveMap {
		res, err := stmtTransferState.Exec(height, createKey, true)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}
}

func SaveDataToDBValidTransferMap(height int,
	inscriptionsValidTransferMap map[string]*model.InscriptionBRC20TickInfo,
) {
	stmtValidTransfer, err := SwapDB.Prepare(`
INSERT INTO brc20_valid_transfer(block_height, tick, pkscript, amount,
    inscription_number, inscription_id, txid, vout, output_value, output_offset)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for _, transferInfo := range inscriptionsValidTransferMap {
		res, err := stmtValidTransfer.Exec(height, transferInfo.Tick,
			transferInfo.PkScript,
			transferInfo.Amount.String(),
			transferInfo.InscriptionNumber, transferInfo.Meta.GetInscriptionId(),
			transferInfo.TxId, transferInfo.Vout, transferInfo.Satoshi, transferInfo.Offset,
		)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}

}

func SaveDataToDBTransferStateMap(height int,
	inscriptionsTransferRemoveMap map[string]struct{},
) {
	stmtTransferState, err := SwapDB.Prepare(`
INSERT INTO brc20_transfer_state(block_height, create_key, moved)
VALUES ($1, $2, $3)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for createKey := range inscriptionsTransferRemoveMap {
		res, err := stmtTransferState.Exec(height, createKey, true)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}
}

func SaveDataToDBValidTransferMap(height int,
	inscriptionsValidTransferMap map[string]*model.InscriptionBRC20TickInfo,
) {
	stmtValidTransfer, err := SwapDB.Prepare(`
INSERT INTO brc20_valid_transfer(block_height, tick, pkscript, amount,
    inscription_number, inscription_id, txid, vout, output_value, output_offset)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for _, transferInfo := range inscriptionsValidTransferMap {
		res, err := stmtValidTransfer.Exec(height, transferInfo.Tick,
			transferInfo.PkScript,
			transferInfo.Amount.String(),
			transferInfo.InscriptionNumber, transferInfo.Meta.GetInscriptionId(),
			transferInfo.TxId, transferInfo.Vout, transferInfo.Satoshi, transferInfo.Offset,
		)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}

}

func SaveDataToDBTickerInfoMap(height int,
	inscriptionsTransferRemoveMap map[string]struct{},
) {
	stmtApproveState, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_approve_state(block_height, create_key, moved)
VALUES ($1, $2, $3)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	for createKey := range inscriptionsTransferRemoveMap {
		res, err := stmtTransferState.Exec(height, createKey, true)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}
}

func SaveDataToDBTickerInfoMap(height int,
	inscriptionsValidTransferMap map[string]*model.InscriptionBRC20TickInfo,
) {
	stmtValidApprove, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_valid_approve(block_height, module_id, tick, pkscript, amount,
    inscription_number, inscription_id, txid, vout, output_value, output_offset)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	for _, approveInfo := range balanceData.ValidApproveMap {
		res, err := stmtValidApprove.Exec(height, moduleId, ticker,
			approveInfo.PkScript,
			approveInfo.Amount.String(),
			approveInfo.InscriptionNumber, approveInfo.Meta.GetInscriptionId(),
			approveInfo.TxId, approveInfo.Vout, approveInfo.Satoshi, approveInfo.Offset,
		)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}
}

// cond approve
func SaveDataToDBTickerInfoMap(height int,
	inscriptionsTransferRemoveMap map[string]struct{},
) {

	stmtCondApproveState, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_cond_approve_state(block_height, create_key, moved)
VALUES ($1, $2, $3)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	for createKey := range inscriptionsTransferRemoveMap {
		res, err := stmtTransferState.Exec(height, createKey, true)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}
}

func SaveDataToDBTickerInfoMap(height int,
	inscriptionsTransferRemoveMap map[string]struct{},
) {

	stmtValidCondApprove, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_valid_cond_approve(block_height, module_id, tick, pkscript, amount,
    inscription_number, inscription_id, txid, vout, output_value, output_offset)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	for _, condApproveInfo := range balanceData.ValidCondApproveMap {
		res, err := stmtValidApprove.Exec(height, moduleId, ticker,
			condApproveInfo.PkScript,
			condApproveInfo.Amount.String(),
			condApproveInfo.InscriptionNumber, condApproveInfo.Meta.GetInscriptionId(),
			condApproveInfo.TxId, condApproveInfo.Vout, condApproveInfo.Satoshi, condApproveInfo.Offset,
		)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}
}

func SaveDataToDBTickerInfoMap(height int,
	inscriptionsTransferRemoveMap map[string]struct{},
) {

	stmtValidWithdraw, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_valid_withdraw(block_height, module_id, tick, pkscript, amount,
    inscription_number, inscription_id, txid, vout, output_value, output_offset)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for _, withdrawInfo := range balanceData.ValidWithdrawMap {
		res, err := stmtValidApprove.Exec(height, moduleId, ticker,
			withdrawInfo.PkScript,
			withdrawInfo.Amount.String(),
			withdrawInfo.InscriptionNumber, withdrawInfo.Meta.GetInscriptionId(),
			withdrawInfo.TxId, withdrawInfo.Vout, withdrawInfo.Satoshi, withdrawInfo.Offset,
		)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)
	}

}

func SaveDataToDBModuleInfoMap(fname string,
	modulesInfoMap map[string]*model.BRC20ModuleSwapInfo) {

	stmtSwapInfo, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_info(block_height, module_id,
	 name,
    pkscript_deployer,
    pkscript_sequencer,
    pkscript_gas_to,
    pkscript_lp_fee,
	 gas_tick,
    fee_rate_swap
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for moduleId, info := range modulesInfoMap {
		// save swap info db
		res, err := stmtSwapInfo.Exec(height, moduleId,
			info.Name,
			info.DeployerPkScript,
			info.SequencerPkScript,
			info.GasToPkScript,
			info.LpFeePkScript,
			info.GasTick,
			info.FeeRateSwap,
		)
		if err != nil {
			log.Fatal("PG Statements Exec Wrong: ", err)
		}
		id, err := res.RowsAffected()
		if err != nil {
			log.Fatal("PG Affecte Wrong: ", err)
		}
		fmt.Println(id)

	}
}

func SaveDataToDBModuleHistoryMap(fname string,
	modulesInfoMap map[string]*model.BRC20ModuleSwapInfo) {

	stmtSwapHistory, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_history(block_height, module_id,
    history_type,
    valid,
    txid,
    idx,
    vout,
    output_value,
    output_offset,
    pkscript_from,
    pkscript_to,
    fee,
    txidx,
    block_time,
    inscription_number,
    inscription_id,
    inscription_content
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for moduleId, info := range modulesInfoMap {

		nValid := 0
		// history
		for _, h := range info.History {
			if h.Valid {
				nValid++
			}
			if !h.Valid {
				continue
			}

			{
				res, err := stmtSwapHistory.Exec(height, moduleId,
					h.Type, h.Valid,
					h.TxId, h.Idx, h.Vout, h.Satoshi, h.Offset,
					h.PkScriptFrom, h.PkScriptTo,
					h.Fee,
					h.TxIdx, h.BlockTime,
					h.Inscription.InscriptionNumber, h.Inscription.InscriptionId,
					"", // content
				)
				if err != nil {
					log.Fatal("PG Statements Exec Wrong: ", err)
				}
				id, err := res.RowsAffected()
				if err != nil {
					log.Fatal("PG Affecte Wrong: ", err)
				}
				fmt.Println(id)
			}

		}

	}
}

func SaveDataToDBModuleCommitInfoMap(fname string,
	modulesInfoMap map[string]*model.BRC20ModuleSwapInfo) {

	stmtSwapCommitState, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_commit_state(block_height, module_id, commit_id, valid, connected)
VALUES ($1, $2, $3, $4, $5)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	stmtSwapCommitInfo, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_history(block_height, module_id,
    history_type,
    valid,
    txid,
    idx,
    vout,
    output_value,
    output_offset,
    pkscript_from,
    pkscript_to,
    fee,
    txidx,
    block_time,
    inscription_number,
    inscription_id,
    inscription_content
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for moduleId, info := range modulesInfoMap {
		// commit state
		commitState := make(map[string]*[2]bool)
		for commitId := range info.CommitInvalidMap {
			if state, ok := commitState[commitId]; !ok {
				commitState[commitId] = &[2]bool{false, false}
			} else {
				state[0] = false
			}
		}

		for commitId := range info.CommitIdMap {
			if state, ok := commitState[commitId]; !ok {
				commitState[commitId] = &[2]bool{true, false}
			} else {
				state[0] = true
			}
		}
		for commitId := range info.CommitIdChainMap {
			if state, ok := commitState[commitId]; !ok {
				commitState[commitId] = &[2]bool{true, true}
			} else {
				state[1] = true
			}
		}

		// save commit state db
		for commitId, state := range info.commitState {
			res, err := stmtSwapCommitState.Exec(height, moduleId,
				commitId,
				state[0], // valid
				state[1], // connected
			)
			if err != nil {
				log.Fatal("PG Statements Exec Wrong: ", err)
			}
			id, err := res.RowsAffected()
			if err != nil {
				log.Fatal("PG Affecte Wrong: ", err)
			}
			fmt.Println(id)
		}

	}
}

func SaveDataToDBModuleUserBalanceMap(fname string,
	modulesInfoMap map[string]*model.BRC20ModuleSwapInfo) {

	stmtUserBalance, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_user_balance(block_height, module_id, tick,
    pkscript, swap_balance, available_balance, approveable_balance, cond_approveable_balance, withdrawable_balance)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}

	for moduleId, info := range modulesInfoMap {
		for ticker, holdersMap := range info.TokenUsersBalanceDataMap {
			// holders
			for holder, balanceData := range holdersMap {
				// save balance db
				res, err := stmtUserBalance.Exec(height, moduleId, ticker,
					balanceData.PkScript,
					balanceData.SwapAccountBalance.String(),
					balanceData.AvailableBalance.String(),
					balanceData.ApproveableBalance.String(),
					balanceData.CondApproveableBalance.String(),
					balanceData.WithdrawableBalance.String(),
				)
				if err != nil {
					log.Fatal("PG Statements Exec Wrong: ", err)
				}
				id, err := res.RowsAffected()
				if err != nil {
					log.Fatal("PG Affecte Wrong: ", err)
				}
				fmt.Println(id)
			}
		}
	}
}

func SaveDataToDBModulePoolLpBalanceMap(fname string,
	modulesInfoMap map[string]*model.BRC20ModuleSwapInfo) {

	stmtPoolBalance, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_pool_balance(block_height, module_id, tick0, tick0_balance, tick1, tick1_balance, lp_balance)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	for moduleId, info := range modulesInfoMap {
		for ticker, swap := range info.SwapPoolTotalBalanceDataMap {
			// save swap balance db
			res, err := stmtPoolBalance.Exec(height, moduleId,
				swap.Tick[0],
				swap.TickBalance[0],
				swap.Tick[1],
				swap.TickBalance[1],
				swap.LpBalance.String(),
			)
			if err != nil {
				log.Fatal("PG Statements Exec Wrong: ", err)
			}
			id, err := res.RowsAffected()
			if err != nil {
				log.Fatal("PG Affecte Wrong: ", err)
			}
			fmt.Println(id)
		}
	}
}

func SaveDataToDBModuleUserLpBalanceMap(fname string,
	modulesInfoMap map[string]*model.BRC20ModuleSwapInfo) {

	stmtLpBalance, err := SwapDB.Prepare(`
INSERT INTO brc20_swap_user_lp_balance(block_height, module_id, pool, pkscript, lp_balance)
VALUES ($1, $2, $3, $4, $5)
`)
	if err != nil {
		log.Fatal("PG Statements Wrong: ", err)
	}
	for moduleId, info := range modulesInfoMap {
		for ticker, holdersMap := range info.LPTokenUsersBalanceMap {
			// holders
			for holder, balanceData := range holdersMap {
				// save balance db
				res, err := stmtLpBalance.Exec(height, moduleId, ticker,
					holder,
					balanceData.String(),
				)
				if err != nil {
					log.Fatal("PG Statements Exec Wrong: ", err)
				}
				id, err := res.RowsAffected()
				if err != nil {
					log.Fatal("PG Affecte Wrong: ", err)
				}
				fmt.Println(id)
			}
		}
	}

}

func SaveDataToDBModuleTickInfoMap(moduleId string, condStateBalanceDataMap map[string]*model.BRC20ModuleConditionalApproveStateBalance,
	inscriptionsTickerInfoMap, userTokensBalanceData map[string]map[string]*model.BRC20ModuleTokenBalance) {

	// condStateBalanceDataMap
	for ticker, stateBalance := range condStateBalanceDataMap {
		fmt.Printf("  module deposit/withdraw state: %s deposit: %s, match: %s, new: %s, cancel: %s, wait: %s\n",
			ticker,
			stateBalance.BalanceDeposite.String(),
			stateBalance.BalanceApprove.String(),
			stateBalance.BalanceNewApprove.String(),
			stateBalance.BalanceCancelApprove.String(),

			stateBalance.BalanceNewApprove.Sub(
				stateBalance.BalanceApprove).Sub(
				stateBalance.BalanceCancelApprove).String(),
		)
	}

	fmt.Printf("\n")
}
