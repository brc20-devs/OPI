package brc20

import (
	"brc20query/logger"
	"brc20query/model"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	brc20swapIndexer "github.com/unisat-wallet/libbrc20-indexer/indexer"
	brc20swapLoader "github.com/unisat-wallet/libbrc20-indexer/loader"
	brc20swapModel "github.com/unisat-wallet/libbrc20-indexer/model"
	"go.uber.org/zap"
)

var (
	DB_CONF_USER     = os.Getenv("DB_USER")
	DB_CONF_HOST     = os.Getenv("DB_HOST")
	DB_CONF_PORT     = os.Getenv("DB_PORT")
	DB_CONF_DATABASE = os.Getenv("DB_DATABASE")
	DB_CONF_PASSWD   = os.Getenv("DB_PASSWD")
)

// ProcessUpdateLatestBRC20SwapInit
func ProcessUpdateLatestBRC20SwapInit(ctx context.Context, startHeight, endHeight int) {
	brc20DatasLoad := make(chan *brc20swapModel.InscriptionBRC20Data, 10240)
	brc20DatasDump := make(chan interface{}, 10240)
	brc20DatasParse := make(chan *brc20swapModel.InscriptionBRC20Data, 10240)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_CONF_HOST, DB_CONF_PORT, DB_CONF_USER, DB_CONF_PASSWD, DB_CONF_DATABASE)
	brc20swapLoader.Init(psqlInfo)

	dbHeight, err := brc20swapLoader.GetBrc20LatestHeightFromDB()
	if err != nil {
		log.Panicf("get db height error: %v", err)
	}
	if dbHeight > startHeight {
		startHeight = int(dbHeight) + 1
	}

	go func() {
		if err := brc20swapLoader.LoadBRC20InputDataFromDB(ctx, brc20DatasLoad, startHeight, endHeight); err != nil {
			log.Panicf("load input data from db error: %v", err)
		}
		close(brc20DatasLoad)
	}()

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

	g := &brc20swapIndexer.BRC20ModuleIndexer{}
	g.Init()

	logger.Log.Info("loading from database...")
	st := time.Now()
	g.LoadDataFromDB(startHeight)
	logger.Log.Info("load from database done", zap.String("elapse", time.Since(st).String()))

	brc20DatasPerHeight := make([]*brc20swapModel.InscriptionBRC20Data, 0, 1024)
	lastHeight := uint32(startHeight)
	blockSt := time.Now()
	for data := range brc20DatasParse {
		if len(brc20DatasPerHeight) > 0 && lastHeight != data.Height {
			brc20DatasPerHeightChan := make(chan interface{}, 10240)
			go func() {
				for _, data := range brc20DatasPerHeight {
					brc20DatasPerHeightChan <- data
				}
				close(brc20DatasPerHeightChan)
			}()

			{
				st := time.Now()
				g.ProcessUpdateLatestBRC20Loop(brc20DatasPerHeightChan, nil)
				logger.Log.Debug("process brc20 data",
					zap.Uint32("height", lastHeight),
					zap.String("elapse", time.Since(st).String()))
			}

			select {
			case <-ctx.Done():
				return
			default:
			}

			if g.Durty {
				st := time.Now()
				g.SaveDataToDB(lastHeight)
				logger.Log.Debug("save to database",
					zap.Uint32("height", lastHeight),
					zap.String("elapse", time.Since(st).String()))
				g.PurgeHistoricalData()
			}

			brc20DatasPerHeight = make([]*brc20swapModel.InscriptionBRC20Data, 0, 1024)
			logger.Log.Debug("process block",
				zap.Uint32("height", lastHeight),
				zap.String("elapse", time.Since(blockSt).String()))
			blockSt = time.Now()
		}
		lastHeight = data.Height
		brc20DatasPerHeight = append(brc20DatasPerHeight, data)
	}

	for _, holdersBalanceMap := range g.TokenUsersBalanceData {
		for key, balance := range holdersBalanceMap {
			if balance.AvailableBalance.Sign() == 0 && balance.TransferableBalance.Sign() == 0 {
				delete(holdersBalanceMap, key)
			}
		}
	}

	model.GSwap = g

	log.Printf("dumping output...")
	brc20swapLoader.DumpTickerInfoMap("./data/brc20.output.txt",
		g.HistoryData,
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
