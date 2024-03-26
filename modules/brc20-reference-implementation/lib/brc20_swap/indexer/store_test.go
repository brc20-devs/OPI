package indexer

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

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
	g := &BRC20ModuleIndexer{
		InscriptionsValidTransferMap: map[string]*model.InscriptionBRC20TickInfo{
			fmt.Sprintf("create_key:%s", uuid.NewString()): &model.InscriptionBRC20TickInfo{
				Tick:              "ordi",
				PkScript:          "pkscript",
				Amount:            decimal.MustNewDecimalFromString("10000000000000000000", 18),
				InscriptionNumber: 1,
				Meta: &model.InscriptionBRC20Data{
					InscriptionId: "ordi",
				},
				TxId:    "txid",
				Vout:    1,
				Satoshi: 1,
				Offset:  0,
			},
		},
	}
	g.SaveDataToDB(psqlInfo, 0)
}

func TestLoadDataFromDb(t *testing.T) {
	g := &BRC20ModuleIndexer{}
	g.LoadDataFromDB(psqlInfo, 0)
}
