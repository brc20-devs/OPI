package model

import (
	"fmt"

	"brc20query/lib/brc20_swap/utils"
)

type BRC20HistoryBase struct {
	Type  uint8 // inscribe-deploy/inscribe-mint/inscribe-transfer/transfer/send/receive
	Valid bool

	TxId   string
	Idx    uint32
	Vout   uint32
	Offset uint64

	PkScriptFrom string
	PkScriptTo   string
	Satoshi      uint64
	Fee          int64

	Height    uint32
	TxIdx     uint32
	BlockTime uint32
}

// history
type BRC20History struct {
	BRC20HistoryBase

	Inscription *InscriptionBRC20TickInfoResp

	// param
	Tick   string
	Amount string

	// state
	OverallBalance      string
	TransferableBalance string
	AvailableBalance    string
}

func NewBRC20History(ticker string, historyType uint8, isValid bool, isTransfer bool,
	from *InscriptionBRC20TickInfo, bal *BRC20TokenBalance, to *InscriptionBRC20Data) *BRC20History {
	history := &BRC20History{
		BRC20HistoryBase: BRC20HistoryBase{
			Type:      historyType,
			Valid:     isValid,
			Height:    to.Height,
			TxIdx:     to.TxIdx,
			BlockTime: to.BlockTime,
			Fee:       to.Fee,
		},
		Inscription: &InscriptionBRC20TickInfoResp{
			Height:            from.Height,
			Data:              from.Data,
			InscriptionNumber: from.InscriptionNumber,
			InscriptionId:     fmt.Sprintf("%si%d", utils.HashString([]byte(from.TxId)), from.Idx),
			Satoshi:           from.Satoshi,
		},
		Tick:   ticker,
		Amount: from.Amount.String(),
	}
	if isTransfer {
		history.TxId = to.TxId
		history.Vout = to.Vout
		history.Offset = to.Offset
		history.Idx = to.Idx
		history.PkScriptFrom = from.PkScript
		history.PkScriptTo = to.PkScript
		history.Satoshi = to.Satoshi
		if history.Satoshi == 0 {
			history.PkScriptTo = history.PkScriptFrom
		}

	} else {
		history.TxId = from.TxId
		history.Vout = from.Vout
		history.Offset = from.Offset
		history.Idx = from.Idx
		history.PkScriptTo = from.PkScript
		history.Satoshi = from.Satoshi
	}

	if bal != nil {
		history.OverallBalance = bal.AvailableBalance.Add(bal.TransferableBalance).String()
		history.TransferableBalance = bal.TransferableBalance.String()
		history.AvailableBalance = bal.AvailableBalance.String()
	} else {
		history.OverallBalance = "0"
		history.TransferableBalance = "0"
		history.AvailableBalance = "0"
	}
	return history
}
