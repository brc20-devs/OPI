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
func ProcessUpdateLatestBRC20SwapInit(endHeight int) {
	brc20DatasLoad := make(chan *brc20swapModel.InscriptionBRC20Data, 10240)
	brc20DatasDump := make(chan *brc20swapModel.InscriptionBRC20Data, 10240)
	brc20DatasParse := make(chan *brc20swapModel.InscriptionBRC20Data, 10240)

	inputFileName := "./data/log_file.txt"
	log.Printf("loading data...")
	totalDataCount, err := brc20swapLoader.GetBRC20InputDataLineCount(inputFileName)
	if err != nil {
		log.Printf("invalid input, %s", err)
		return
	}

	go func(endHeight int) {
		if err := brc20swapLoader.LoadBRC20InputDataFromOrdLog(inputFileName, brc20DatasLoad, endHeight); err != nil {
			log.Printf("invalid input, %s", err)
		}
		close(brc20DatasLoad)
	}(endHeight)

	go func() {
		for data := range brc20DatasLoad {
			brc20DatasParse <- data
			brc20DatasDump <- data
		}
		close(brc20DatasParse)
		close(brc20DatasDump)
	}()

	go func() {
		brc20swapLoader.DumpBRC20InputData("./data/brc20.input.txt", brc20DatasDump, true)
	}()

	g := &brc20swapIndexer.BRC20ModuleIndexer{}
	g.ProcessUpdateLatestBRC20Init(brc20DatasParse, totalDataCount)

	model.GSwap = g

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_CONF_HOST, DB_CONF_PORT, DB_CONF_USER, DB_CONF_PASSWD, DB_CONF_DATABASE)

	g.SaveDataToDB(psqlInfo, endHeight)

	brc20swapLoader.DumpTickerInfoMap("./data/brc20.output.txt",
		g.InscriptionsTickerInfoMap,
		g.UserTokensBalanceData,
		g.TokenUsersBalanceData,
	)

	brc20swapLoader.DumpModuleInfoMap("./data/brc20-module.output.txt",
		g.ModulesInfoMap,
	)
}
