package model

import (
	"fmt"

	"brc20query/lib/brc20_swap/utils"
)

// history
type BRC20ModuleHistory struct {
	BRC20HistoryBase
	Inscription InscriptionBRC20SwapInfoResp
	Data        any
	// no state
}

func NewBRC20ModuleHistory(isTransfer bool, historyType uint8, from, to *InscriptionBRC20Data, data any, isValid bool) *BRC20ModuleHistory {
	history := &BRC20ModuleHistory{
		BRC20HistoryBase: BRC20HistoryBase{
			Type:  historyType,
			Valid: isValid,
		},
		Inscription: InscriptionBRC20SwapInfoResp{
			Height:            from.Height,
			ContentBody:       from.ContentBody, // to.Content is empty on transfer
			InscriptionNumber: from.InscriptionNumber,
			InscriptionId:     fmt.Sprintf("%si%d", utils.HashString([]byte(from.TxId)), from.Idx),
		},
	}
	if isTransfer {
		history.TxId = to.TxId
		history.Vout = to.Vout
		history.Offset = to.Offset
		history.Idx = to.Idx
		history.PkScriptFrom = from.PkScript
		history.PkScriptTo = to.PkScript
		history.Satoshi = to.Satoshi

		history.Height = to.Height
		history.TxIdx = to.TxIdx
		history.BlockTime = to.BlockTime

	} else {
		history.TxId = from.TxId
		history.Vout = from.Vout
		history.Offset = from.Offset
		history.Idx = from.Idx
		history.PkScriptTo = from.PkScript
		history.Satoshi = from.Satoshi

		history.Height = from.Height
		history.TxIdx = from.TxIdx
		history.BlockTime = from.BlockTime
	}
	history.Data = data
	return history
}

// history
type BRC20SwapHistoryApproveData struct {
	Tick   string `json:"tick"`
	Amount string `json:"amount"` // current amt
}

// history
type BRC20SwapHistoryCondApproveData struct {
	Tick                  string `json:"tick"`
	Amount                string `json:"amount"`      // current amt
	Balance               string `json:"balance"`     // current balance
	TransferInscriptionId string `json:"transfer"`    // transfer inscription id
	TransferMax           string `json:"transferMax"` // transfer inscription id
}
