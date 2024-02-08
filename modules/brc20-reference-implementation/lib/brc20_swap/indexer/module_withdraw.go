package indexer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"brc20query/lib/brc20_swap/model"
	"brc20query/lib/brc20_swap/utils"
)

func (g *BRC20ModuleIndexer) ProcessInscribeWithdraw(data *model.InscriptionBRC20Data) error {
	var body model.InscriptionBRC20ModuleWithdrawContent
	if err := json.Unmarshal(data.ContentBody, &body); err != nil {
		log.Printf("parse module withdraw json failed. txid: %s",
			hex.EncodeToString(utils.ReverseBytes([]byte(data.TxId))),
		)
		return err
	}

	// lower case only
	if body.Module != strings.ToLower(body.Module) {
		return errors.New("module id invalid")
	}

	if _, ok := g.ModulesInfoMap[body.Module]; !ok { // invalid module
		return errors.New("module invalid")
	}

	// black module
	return nil
}
