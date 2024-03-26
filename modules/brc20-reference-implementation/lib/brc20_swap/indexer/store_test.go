package indexer

import (
	"fmt"
	"testing"

	"brc20query/lib/brc20_swap/decimal"
	"brc20query/lib/brc20_swap/model"
)

var (
	pg_host     = "localhost"
	pg_port     = 5432
	pg_user     = "postgres"
	pg_password = "postgres"
	pg_dbname   = "swapdev"
	psqlInfo    string
)

func init() {
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg_host, pg_port, pg_user, pg_password, pg_dbname)
}

func TestSavaDataToDb(t *testing.T) {
	createKey1 := model.NFTCreateIdxKey{
		Height:     820000,
		IdxInBlock: 0,
	}
	createKey2 := model.NFTCreateIdxKey{
		Height:     830000,
		IdxInBlock: 0,
	}

	g := &BRC20ModuleIndexer{
		AllHistory:                    []*model.BRC20History{},
		UserAllHistory:                map[string]*model.BRC20UserHistory{},
		InscriptionsTickerInfoMap:     map[string]*model.BRC20TokenInfo{},
		UserTokensBalanceData:         map[string]map[string]*model.BRC20TokenBalance{},
		TokenUsersBalanceData:         map[string]map[string]*model.BRC20TokenBalance{},
		InscriptionsValidBRC20DataMap: map[string]*model.InscriptionBRC20InfoResp{},
		InscriptionsTransferRemoveMap: map[string]struct{}{},
		InscriptionsValidTransferMap: map[string]*model.InscriptionBRC20TickInfo{
			createKey1.String(): &model.InscriptionBRC20TickInfo{
				Tick: "ordi", TxId: "txid", Vout: 1, Satoshi: 1, Offset: 0, PkScript: "pkscript",
				InscriptionNumber: 1,
				Amount:            decimal.MustNewDecimalFromString("10000000000000000000", 18),
				Meta:              &model.InscriptionBRC20Data{InscriptionId: "ordi"},
			},
		},
		InscriptionsInvalidTransferMap: map[string]*model.InscriptionBRC20TickInfo{
			createKey2.String(): &model.InscriptionBRC20TickInfo{
				Tick: "ordi", TxId: "txid", Vout: 1, Satoshi: 1, Offset: 0, PkScript: "pkscript",
				InscriptionNumber: 1,
				Amount:            decimal.MustNewDecimalFromString("10000000000000000000", 18),
				Meta:              &model.InscriptionBRC20Data{InscriptionId: "ordi"},
			},
		},
		ModulesInfoMap: map[string]*model.BRC20ModuleSwapInfo{
			"module1": &model.BRC20ModuleSwapInfo{
				ID:   "0x123456789",
				Name: "brc20-swap",
			},
		},
		UsersModuleWithTokenMap:                     map[string]string{},
		UsersModuleWithLpTokenMap:                   map[string]string{},
		InscriptionsApproveRemoveMap:                map[string]struct{}{},
		InscriptionsValidApproveMap:                 map[string]*model.InscriptionBRC20SwapInfo{},
		InscriptionsInvalidApproveMap:               map[string]*model.InscriptionBRC20SwapInfo{},
		InscriptionsCondApproveRemoveMap:            map[string]struct{}{},
		InscriptionsValidConditionalApproveMap:      map[string]*model.InscriptionBRC20SwapConditionalApproveInfo{},
		InscriptionsInvalidConditionalApproveMap:    map[string]*model.InscriptionBRC20SwapConditionalApproveInfo{},
		InscriptionsCommitRemoveMap:                 map[string]struct{}{},
		InscriptionsValidCommitMap:                  map[string]*model.InscriptionBRC20Data{},
		InscriptionsInvalidCommitMap:                map[string]*model.InscriptionBRC20Data{},
		InscriptionsValidCommitMapById:              map[string]*model.InscriptionBRC20Data{},
		InscriptionsWithdrawRemoveMap:               map[string]struct{}{},
		InscriptionsValidWithdrawMap:                map[string]*model.InscriptionBRC20SwapInfo{},
		InscriptionsInvalidWithdrawMap:              map[string]*model.InscriptionBRC20SwapInfo{},
		ThisTxId:                                    "",
		TxStaticTransferStatesForConditionalApprove: []*model.TransferStateForConditionalApprove{},
	}
	g.SaveDataToDB(psqlInfo, 0)
}

func TestLoadDataFromDb(t *testing.T) {
	g := &BRC20ModuleIndexer{}
	g.LoadDataFromDB(psqlInfo, 0)
}
