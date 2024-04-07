package indexer

import (
	"encoding/hex"
	"errors"
	"log"

	"brc20query/lib/brc20_swap/constant"
	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/model"
	"brc20query/lib/brc20_swap/utils"
)

func (g *BRC20ModuleIndexer) ProcessCommitFunctionGasFee(moduleInfo *model.BRC20ModuleSwapInfo, userPkScript string, gasAmt *decimal.Decimal) error {

	tokenBalance := moduleInfo.GetUserTokenBalance(moduleInfo.GasTick, userPkScript)
	// fixme: Must use the confirmed amount
	if tokenBalance.SwapAccountBalance.Cmp(gasAmt) < 0 {
		address, err := utils.GetAddressFromScript([]byte(userPkScript), constant.GlobalNetParams)
		if err != nil {
			address = hex.EncodeToString([]byte(userPkScript))
		}

		log.Printf("gas[%s] user[%s], balance %s", moduleInfo.GasTick, address, tokenBalance)
		return errors.New("gas fee: token balance insufficient")
	}

	gasToBalance := moduleInfo.GetUserTokenBalance(moduleInfo.GasTick, moduleInfo.GasToPkScript)

	// User Real-time gas Balance Update
	tokenBalance.SwapAccountBalance = tokenBalance.SwapAccountBalance.Sub(gasAmt)
	gasToBalance.SwapAccountBalance = gasToBalance.SwapAccountBalance.Add(gasAmt)

	tokenBalance.UpdateHeight = g.CurrentHeight
	gasToBalance.UpdateHeight = g.CurrentHeight

	log.Printf("gas fee[%s]: %s user: %s, gasTo: %s", moduleInfo.GasTick, gasAmt, tokenBalance.SwapAccountBalance, gasToBalance.SwapAccountBalance)
	return nil
}
