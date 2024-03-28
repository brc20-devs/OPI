package brc20

import (
	brc20swapIndexer "brc20query/lib/brc20_swap/indexer"
	brc20swapLoader "brc20query/lib/brc20_swap/loader"
	brc20swapModel "brc20query/lib/brc20_swap/model"
	"brc20query/model"
	"fmt"
	"log"
	"os"
)

var (
	DB_CONF_USER     = os.Getenv("DB_USER")
	DB_CONF_HOST     = os.Getenv("DB_HOST")
	DB_CONF_PORT     = os.Getenv("DB_PORT")
	DB_CONF_DATABASE = os.Getenv("DB_DATABASE")
	DB_CONF_PASSWD   = os.Getenv("DB_PASSWD")
)

// ProcessUpdateLatestBRC20SwapInit
func ProcessUpdateLatestBRC20SwapInit(startHeight, endHeight int) {
	brc20DatasLoad := make(chan *brc20swapModel.InscriptionBRC20Data, 10240)
	brc20DatasDump := make(chan *brc20swapModel.InscriptionBRC20Data, 10240)
	brc20DatasParse := make(chan *brc20swapModel.InscriptionBRC20Data, 10240)

	inputFileName := "./data/log_file.txt"
	log.Printf("loading data...")
	if _, err := brc20swapLoader.GetBRC20InputDataLineCount(inputFileName); err != nil {
		log.Printf("invalid input, %s", err)
		return
	}

	go func(endHeight int) {
		if err := brc20swapLoader.LoadBRC20InputDataFromOrdLog(inputFileName, brc20DatasLoad, startHeight, endHeight); err != nil {
			log.Printf("invalid input, %s", err)
		}
		close(brc20DatasLoad)
	}(endHeight)

	go func() {
		for data := range brc20DatasLoad {
			brc20DatasParse <- data
			brc20DatasDump <- data
		}

		// finish
		brc20DatasParse <- &brc20swapModel.InscriptionBRC20Data{}

		close(brc20DatasParse)
		close(brc20DatasDump)
	}()

	go func() {
		brc20swapLoader.DumpBRC20InputData("./data/brc20.input.txt", brc20DatasDump, true)
	}()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_CONF_HOST, DB_CONF_PORT, DB_CONF_USER, DB_CONF_PASSWD, DB_CONF_DATABASE)

	g := &brc20swapIndexer.BRC20ModuleIndexer{}
	g.Init()

	log.Printf("loading database...")
	g.LoadDataFromDB(psqlInfo, startHeight)
	log.Printf("load database ok")

	brc20DatasPerHeight := []*brc20swapModel.InscriptionBRC20Data{}
	lastHeight := uint32(startHeight)
	for data := range brc20DatasParse {
		if len(brc20DatasPerHeight) > 0 && lastHeight != data.Height {
			g.ProcessUpdateLatestBRC20Loop(brc20DatasPerHeight, len(brc20DatasPerHeight))
			if g.Durty {
				log.Printf("height: %d, saving database...", lastHeight)
				g.SaveDataToDB(psqlInfo, lastHeight)
				log.Printf("save database ok")
			}

			brc20DatasPerHeight = []*brc20swapModel.InscriptionBRC20Data{}
		}
		lastHeight = data.Height
		brc20DatasPerHeight = append(brc20DatasPerHeight, data)
	}

	model.GSwap = g

	log.Printf("dumping output...")
	brc20swapLoader.DumpTickerInfoMap("./data/brc20.output.txt",
		g.InscriptionsTickerInfoMap,
		g.UserTokensBalanceData,
		g.TokenUsersBalanceData,
	)

	brc20swapLoader.DumpModuleInfoMap("./data/brc20-module.output.txt",
		g.ModulesInfoMap,
	)
	log.Printf("dump output ok")
}

// SELECT relname, n_live_tup AS row_count FROM pg_stat_user_tables ORDER BY relname DESC;
