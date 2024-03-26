package indexer

import (
	"fmt"
	"testing"
)

var (
	pg_host     = "localhost"
	pg_port     = 5432
	pg_user     = "postgres"
	pg_password = "postgres"
	//pg_dbname   = "brc20_swap_test"
	pg_dbname = "swaptest"
)

func TestLoadDataFromDb(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg_host, pg_port, pg_user, pg_password, pg_dbname)
	g := &BRC20ModuleIndexer{}
	g.LoadDataFromDB(psqlInfo, 0)
}
