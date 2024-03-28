package model

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/utils"
)

// nft create point on create
type NFTCreateIdxKey struct {
	Height     uint32 // Height of NFT show in block onCreate
	IdxInBlock uint64 // Index of NFT show in block onCreate
}

func (p *NFTCreateIdxKey) String() string {
	var key [12]byte
	binary.LittleEndian.PutUint32(key[0:4], p.Height)
	binary.LittleEndian.PutUint64(key[4:12], p.IdxInBlock)
	return string(key[:])
}

// event raw data
type InscriptionBRC20Data struct {
	IsTransfer bool
	TxId       string `json:"-"`
	Idx        uint32 `json:"-"`
	Vout       uint32 `json:"-"`
	Offset     uint64 `json:"-"`

	Satoshi  uint64 `json:"-"`
	PkScript string `json:"-"`
	Fee      int64  `json:"-"`

	InscriptionNumber int64
	ContentBody       []byte
	CreateIdxKey      string

	Height    uint32 // Height of NFT show in block onCreate
	TxIdx     uint32
	BlockTime uint32
	Sequence  uint16

	// for cache
	InscriptionId string
}

func (data *InscriptionBRC20Data) GetInscriptionId() string {
	if data.InscriptionId == "" {
		data.InscriptionId = fmt.Sprintf("%si%d", utils.HashString([]byte(data.TxId)), data.Idx)
	}
	return data.InscriptionId
}

type InscriptionBRC20InfoResp struct {
	Operation    string `json:"op,omitempty"`
	BRC20Tick    string `json:"tick,omitempty"`
	BRC20Max     string `json:"max,omitempty"`
	BRC20Limit   string `json:"lim,omitempty"`
	BRC20Amount  string `json:"amt,omitempty"`
	BRC20Decimal string `json:"decimal,omitempty"`
	BRC20Minted  string `json:"minted,omitempty"`
}

// decode protocal
type InscriptionBRC20ProtocalContent struct {
	Proto     string `json:"p,omitempty"`
	Operation string `json:"op,omitempty"`
}

func (body *InscriptionBRC20ProtocalContent) Unmarshal(contentBody []byte) (err error) {
	var bodyMap map[string]interface{} = make(map[string]interface{}, 8)
	if err := json.Unmarshal(contentBody, &bodyMap); err != nil {
		return err
	}
	if v, ok := bodyMap["p"].(string); ok {
		body.Proto = v
	}
	if v, ok := bodyMap["op"].(string); ok {
		body.Operation = v
	}
	return nil
}

// decode mint/transfer
type InscriptionBRC20MintTransferContent struct {
	Proto       string `json:"p,omitempty"`
	Operation   string `json:"op,omitempty"`
	BRC20Tick   string `json:"tick,omitempty"`
	BRC20Amount string `json:"amt,omitempty"`
}

func (body *InscriptionBRC20MintTransferContent) Unmarshal(contentBody []byte) (err error) {
	var bodyMap map[string]interface{} = make(map[string]interface{}, 8)
	if err := json.Unmarshal(contentBody, &bodyMap); err != nil {
		return err
	}
	if v, ok := bodyMap["p"].(string); ok {
		body.Proto = v
	}
	if v, ok := bodyMap["op"].(string); ok {
		body.Operation = v
	}
	if v, ok := bodyMap["tick"].(string); ok {
		body.BRC20Tick = v
	}
	if v, ok := bodyMap["amt"].(string); ok {
		body.BRC20Amount = v
	}
	return nil
}

// decode deploy data
type InscriptionBRC20DeployContent struct {
	Proto        string `json:"p,omitempty"`
	Operation    string `json:"op,omitempty"`
	BRC20Tick    string `json:"tick,omitempty"`
	BRC20Max     string `json:"max,omitempty"`
	BRC20Limit   string `json:"lim,omitempty"`
	BRC20Decimal string `json:"dec,omitempty"`
}

func (body *InscriptionBRC20DeployContent) Unmarshal(contentBody []byte) (err error) {
	var bodyMap map[string]interface{} = make(map[string]interface{}, 8)
	if err := json.Unmarshal(contentBody, &bodyMap); err != nil {
		return err
	}
	if v, ok := bodyMap["p"].(string); ok {
		body.Proto = v
	}
	if v, ok := bodyMap["op"].(string); ok {
		body.Operation = v
	}
	if v, ok := bodyMap["tick"].(string); ok {
		body.BRC20Tick = v
	}
	if v, ok := bodyMap["max"].(string); ok {
		body.BRC20Max = v
	}
	if _, ok := bodyMap["lim"]; !ok {
		body.BRC20Limit = body.BRC20Max
	} else {
		if v, ok := bodyMap["lim"].(string); ok {
			body.BRC20Limit = v
		}
	}

	if _, ok := bodyMap["dec"]; !ok {
		body.BRC20Decimal = decimal.MAX_PRECISION_STRING
	} else {
		if v, ok := bodyMap["dec"].(string); ok {
			body.BRC20Decimal = v
		}
	}

	return nil
}

// all ticker (state and history)
type BRC20TokenInfo struct {
	UpdateHeight uint32

	Ticker string
	Deploy *InscriptionBRC20TickInfo

	// empty
	History                 []*BRC20History
	HistoryMint             []*BRC20History
	HistoryInscribeTransfer []*BRC20History
	HistoryTransfer         []*BRC20History
}

type InscriptionBRC20TransferInfo struct {
	Tick   string
	Amount *decimal.Decimal
	Data   *InscriptionBRC20Data
}

// inscription info, with mint state
type InscriptionBRC20TickInfo struct {
	Data   *InscriptionBRC20InfoResp `json:"data"`
	Tick   string
	Amount *decimal.Decimal `json:"-"`
	//ContentBody []byte           `json:"content"`
	Meta *InscriptionBRC20Data

	Max   *decimal.Decimal `json:"-"`
	Limit *decimal.Decimal `json:"-"`

	TotalMinted        *decimal.Decimal `json:"-"`
	ConfirmedMinted    *decimal.Decimal `json:"-"`
	ConfirmedMinted1h  *decimal.Decimal `json:"-"`
	ConfirmedMinted24h *decimal.Decimal `json:"-"`

	MintTimes uint32 `json:"-"`
	Decimal   uint8  `json:"-"`

	TxId   string `json:"-"`
	Idx    uint32 `json:"-"`
	Vout   uint32 `json:"-"`
	Offset uint64 `json:"-"`

	Satoshi  uint64 `json:"-"`
	PkScript string `json:"-"`

	InscriptionNumber int64  `json:"inscriptionNumber"`
	CreateIdxKey      string `json:"-"`
	Height            uint32 `json:"-"`
	TxIdx             uint32 `json:"-"`
	BlockTime         uint32 `json:"-"`

	CompleteHeight    uint32 `json:"-"`
	CompleteBlockTime uint32 `json:"-"`

	InscriptionNumberStart int64 `json:"-"`
	InscriptionNumberEnd   int64 `json:"-"`
}

func (in *InscriptionBRC20TickInfo) DeepCopy() (copy *InscriptionBRC20TickInfo) {
	copy = &InscriptionBRC20TickInfo{
		Data:    in.Data,
		Decimal: in.Decimal,

		TxId:   in.TxId,
		Idx:    in.Idx,
		Vout:   in.Vout,
		Offset: in.Offset,

		Satoshi:  in.Satoshi,
		PkScript: in.PkScript,

		InscriptionNumber: in.InscriptionNumber,
		CreateIdxKey:      in.CreateIdxKey,
		Height:            in.Height,
		TxIdx:             in.TxIdx,
		BlockTime:         in.BlockTime,

		// runtime value
		Max:                decimal.NewDecimalCopy(in.Max),
		Limit:              decimal.NewDecimalCopy(in.Limit),
		TotalMinted:        decimal.NewDecimalCopy(in.TotalMinted),
		ConfirmedMinted:    decimal.NewDecimalCopy(in.ConfirmedMinted),
		ConfirmedMinted1h:  decimal.NewDecimalCopy(in.ConfirmedMinted1h),
		ConfirmedMinted24h: decimal.NewDecimalCopy(in.ConfirmedMinted24h),
		Amount:             decimal.NewDecimalCopy(in.Amount),

		MintTimes: in.MintTimes,

		CompleteHeight:    in.CompleteHeight,
		CompleteBlockTime: in.CompleteBlockTime,

		InscriptionNumberStart: in.InscriptionNumberStart,
		InscriptionNumberEnd:   in.InscriptionNumberEnd,
	}
	return copy
}

func NewInscriptionBRC20TickInfo(tick, operation string, data *InscriptionBRC20Data) *InscriptionBRC20TickInfo {
	info := &InscriptionBRC20TickInfo{
		Data: &InscriptionBRC20InfoResp{
			BRC20Tick: tick,
			Operation: operation,
		},
		Decimal: 18,

		TxId:   data.TxId,
		Idx:    data.Idx,
		Vout:   data.Vout,
		Offset: data.Offset,

		Satoshi:  data.Satoshi,
		PkScript: data.PkScript,

		InscriptionNumber: data.InscriptionNumber,
		CreateIdxKey:      data.CreateIdxKey,
		Height:            data.Height,
		TxIdx:             data.TxIdx,
		BlockTime:         data.BlockTime,
	}
	return info
}

// all history for user
type BRC20UserHistory struct {
	History []*BRC20History
}

// state of address for each tick, (balance and history)
type BRC20TokenBalance struct {
	UpdateHeight uint32

	Ticker               string
	PkScript             string
	AvailableBalance     *decimal.Decimal
	AvailableBalanceSafe *decimal.Decimal
	TransferableBalance  *decimal.Decimal
	ValidTransferMap     map[string]*InscriptionBRC20TickInfo

	History                 []*BRC20History
	HistoryMint             []*BRC20History
	HistoryInscribeTransfer []*BRC20History
	HistorySend             []*BRC20History
	HistoryReceive          []*BRC20History
}

func (bal *BRC20TokenBalance) OverallBalance() *decimal.Decimal {
	return bal.AvailableBalance.Add(bal.TransferableBalance)
}

func (in *BRC20TokenBalance) DeepCopy() (tb *BRC20TokenBalance) {
	tb = &BRC20TokenBalance{
		Ticker:               in.Ticker,
		PkScript:             in.PkScript,
		AvailableBalanceSafe: decimal.NewDecimalCopy(in.AvailableBalanceSafe),
		AvailableBalance:     decimal.NewDecimalCopy(in.AvailableBalance),
		TransferableBalance:  decimal.NewDecimalCopy(in.TransferableBalance),
	}

	tb.ValidTransferMap = make(map[string]*InscriptionBRC20TickInfo, len(in.ValidTransferMap))
	for k, v := range in.ValidTransferMap {
		tb.ValidTransferMap[k] = v.DeepCopy()
	}

	tb.History = make([]*BRC20History, len(in.History))
	copy(tb.History, in.History)

	tb.HistoryMint = make([]*BRC20History, len(in.HistoryMint))
	copy(tb.HistoryMint, in.HistoryMint)

	tb.HistoryInscribeTransfer = make([]*BRC20History, len(in.HistoryInscribeTransfer))
	copy(tb.HistoryInscribeTransfer, in.HistoryInscribeTransfer)

	tb.HistorySend = make([]*BRC20History, len(in.HistorySend))
	copy(tb.HistorySend, in.HistorySend)

	tb.HistoryReceive = make([]*BRC20History, len(in.HistoryReceive))
	copy(tb.HistoryReceive, in.HistoryReceive)
	return tb
}

// history inscription info
type InscriptionBRC20TickInfoResp struct {
	Height            uint32                    `json:"-"`
	Data              *InscriptionBRC20InfoResp `json:"data"`
	InscriptionNumber int64                     `json:"inscriptionNumber"`
	InscriptionId     string                    `json:"inscriptionId"`
	Satoshi           uint64                    `json:"satoshi"`
	Confirmations     int                       `json:"confirmations"`
}
