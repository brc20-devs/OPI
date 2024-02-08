package brc20_swap

import (
	"brc20query/lib/brc20_swap/constant"
	"brc20query/lib/brc20_swap/event"
	"brc20query/lib/brc20_swap/indexer"
	"brc20query/lib/brc20_swap/loader"
	"brc20query/lib/brc20_swap/model"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

var LOAD_TESTNET = false
var LOAD_EVENTS = true

// go test -v -run TestBRC20Swap\$ .
func TestBRC20Swap(t *testing.T) {

	if LOAD_TESTNET {
		constant.GlobalNetParams = &chaincfg.TestNet3Params
		constant.TICKS_ENABLED = "sats ordi trac oshi btcs oxbt texo cncl meme honk zbit vmpx pepe mxrc   doge eyee"

		constant.TICKS_ENABLED = "sats ordi trac oshi btcs oxbt texo cncl meme honk zbit vmpx pepe mxrc   doge eyee test 💰 your domo"

		constant.TICKS_ENABLED = ""
		indexer.MODULE_SWAP_SOURCE_INSCRIPTION_ID = "eabfbf7cba3509134582c2216709527ddde716d3be96beababc16c8f28d5fd31i0"
	} else {
		// mainnet
		constant.GlobalNetParams = &chaincfg.MainNetParams
		constant.TICKS_ENABLED = "sats ordi trac oshi btcs oxbt texo cncl meme honk zbit vmpx pepe mxrc   doge eyee"
		indexer.MODULE_SWAP_SOURCE_INSCRIPTION_ID = "93ce120ff87364c261a534fea4c39196a615f449412fb3547a185d92306a39b8i0"
	}

	// brc20Datas, err := loader.LoadBRC20InputJsonData("./data/brc20swap.input.conf")
	// if err != nil {
	// 	t.Logf("load json failed: %s", err)
	// }
	// loader.DumpBRC20InputData("./data/brc20swap.input.txt", brc20Datas, false)

	// if err := indexer.InitResultDataFromFile("./data/brc20swap.results.json"); err != nil {
	// 	t.Logf("load json failed: %s", err)
	// }

	brc20Datas := make([]*model.InscriptionBRC20Data, 0)

	var err error
	if LOAD_EVENTS {
		constant.DEBUG = true

		t.Logf("start loading event")
		if datas, err := event.InitTickDataFromFile("./data/brc20swap.ticks.json"); err != nil {
			t.Logf("load tick json failed: %s", err)
			return
		} else {
			brc20Datas = append(brc20Datas, datas...)
		}
		if datas, err := event.GenerateBRC20InputDataFromEvents("./data/brc20swap.events.json"); err != nil {
			t.Logf("load event json failed: %s", err)
			return
		} else {
			brc20Datas = append(brc20Datas, datas...)
		}
		loader.DumpBRC20InputData("./data/brc20swap.events.input.txt", brc20Datas, false)

	} else {
		t.Logf("start loading data")

		if brc20Datas, err = loader.LoadBRC20InputData("./data/brc20swap.input.txt"); err != nil {
			t.Logf("load json failed: %s", err)
		}

	}
	t.Logf("start init")
	brc20DatasSplit := len(brc20Datas) - 0

	g := &indexer.BRC20ModuleIndexer{}
	g.ProcessUpdateLatestBRC20Init(brc20Datas[0:brc20DatasSplit])

	// next half
	t.Logf("start deep copy")
	// newg := g.DeepCopy()
	newg := g

	t.Logf("start process")
	newg.ProcessUpdateLatestBRC20Loop(brc20Datas[brc20DatasSplit:])

	t.Logf("dump swap")
	loader.DumpTickerInfoMap("./data/brc20swap.output.txt",
		newg.InscriptionsTickerInfoMap,
		newg.UserTokensBalanceData,
		newg.TokenUsersBalanceData,
	)

	t.Logf("dump module")
	loader.DumpModuleInfoMap("./data/brc20module.output.txt",
		newg.ModulesInfoMap,
	)

}
