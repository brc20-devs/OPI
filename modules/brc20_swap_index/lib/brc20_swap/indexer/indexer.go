package indexer

import (
	"bytes"
	"log"

	"brc20query/lib/brc20_swap/constant"
	"brc20query/lib/brc20_swap/model"
)

func isJson(contentBody []byte) bool {
	if len(contentBody) < 40 {
		return false
	}

	content := bytes.TrimSpace(contentBody)
	if !bytes.HasPrefix(content, []byte("{")) {
		return false
	}
	if !bytes.HasSuffix(content, []byte("}")) {
		return false
	}

	return true
}

// ProcessUpdateLatestBRC20Loop
func (g *BRC20ModuleIndexer) ProcessUpdateLatestBRC20Loop(brc20Datas []*model.InscriptionBRC20Data, totalDataCount int) {
	log.Printf("process swap update. total %d", totalDataCount)

	g.Durty = false
	for idx, data := range brc20Datas {
		percent := idx * 100 / totalDataCount
		// is sending transfer
		if data.IsTransfer {
			// module conditional approve
			if condApproveInfo, isInvalid := g.GetConditionalApproveInfoByKey(data.CreateIdxKey); condApproveInfo != nil {
				if err := g.ProcessConditionalApprove(data, condApproveInfo, isInvalid); err != nil {
					log.Printf("process conditional approve move failed: %s", err)
				} else {
					g.Durty = true
				}
				continue
			}

			// not first move
			if data.Sequence != 1 {
				continue
			}

			// transfer
			if transferInfo, isInvalid := g.GetTransferInfoByKey(data.CreateIdxKey); transferInfo != nil {
				g.InscriptionsTransferRemoveMap[data.CreateIdxKey] = data.Height
				g.Durty = true

				if err := g.ProcessTransfer(data, transferInfo, isInvalid); err != nil {
					log.Printf("process transfer move failed: %s", err)
				} else {
					g.Durty = true
				}
				continue
			}

			// module approve
			if approveInfo, isInvalid := g.GetApproveInfoByKey(data.CreateIdxKey); approveInfo != nil {
				g.InscriptionsApproveRemoveMap[data.CreateIdxKey] = data.Height
				g.Durty = true

				if err := g.ProcessApprove(data, approveInfo, isInvalid); err != nil {
					log.Printf("process approve move failed: %s", err)
				} else {
					g.Durty = true
				}
				continue
			}

			// module commit
			if commitFrom, isInvalid := g.GetCommitInfoByKey(data.CreateIdxKey); commitFrom != nil {
				g.InscriptionsCommitRemoveMap[data.CreateIdxKey] = data.Height
				g.Durty = true

				if err := g.ProcessCommit(commitFrom, data, isInvalid); err != nil {
					log.Printf("process commit move failed: %s", err)
				} else {
					g.Durty = true
				}
				continue
			}

			continue
		}

		// inscribe as fee
		if data.Satoshi == 0 {
			continue
		}

		if ok := isJson(data.ContentBody); !ok {
			// log.Println("not json")
			continue
		}

		// protocal, lower case only
		body := new(model.InscriptionBRC20ProtocalContent)
		if err := body.Unmarshal(data.ContentBody); err != nil {
			// log.Println("Unmarshal failed", err, string(data.ContentBody))
			continue
		}

		// is inscribe deploy/mint/transfer
		if body.Proto != constant.BRC20_P &&
			body.Proto != constant.BRC20_P_MODULE &&
			body.Proto != constant.BRC20_P_SWAP {
			// log.Println("not proto")
			continue
		}

		var process func(*model.InscriptionBRC20Data) error
		if body.Proto == constant.BRC20_P && body.Operation == constant.BRC20_OP_DEPLOY {
			process = g.ProcessDeploy
		} else if body.Proto == constant.BRC20_P && body.Operation == constant.BRC20_OP_MINT {
			process = g.ProcessMint
		} else if body.Proto == constant.BRC20_P && body.Operation == constant.BRC20_OP_TRANSFER {
			process = g.ProcessInscribeTransfer
		} else if body.Proto == constant.BRC20_P_MODULE && body.Operation == constant.BRC20_OP_MODULE_DEPLOY {
			process = g.ProcessCreateModule
		} else if body.Proto == constant.BRC20_P_MODULE && body.Operation == constant.BRC20_OP_MODULE_WITHDRAW {
			process = g.ProcessInscribeWithdraw
		} else if body.Proto == constant.BRC20_P_SWAP && body.Operation == constant.BRC20_OP_SWAP_APPROVE {
			process = g.ProcessInscribeApprove
		} else if body.Proto == constant.BRC20_P_SWAP && body.Operation == constant.BRC20_OP_SWAP_CONDITIONAL_APPROVE {
			process = g.ProcessInscribeConditionalApprove
		} else if body.Proto == constant.BRC20_P_SWAP && body.Operation == constant.BRC20_OP_SWAP_COMMIT {
			process = g.ProcessInscribeCommit
		} else {
			continue
		}

		if err := process(data); err != nil {
			if body.Operation == constant.BRC20_OP_MINT {
				if constant.DEBUG {
					log.Printf("(%d %%) process failed: %s", percent, err)
				}
			} else {
				log.Printf("(%d %%) process failed: %s", percent, err)
			}
		} else {
			g.Durty = true
		}
	}

	// for _, holdersBalanceMap := range g.TokenUsersBalanceData {
	// 	for key, balance := range holdersBalanceMap {
	// 		if balance.AvailableBalance.Sign() == 0 && balance.TransferableBalance.Sign() == 0 {
	// 			delete(holdersBalanceMap, key)
	// 		}
	// 	}
	// }
	if !g.Durty {
		return
	}

	log.Printf("process swap finish. ticker: %d, users: %d, tokens: %d, validInscription: %d, validTransfer: %d, invalidTransfer: %d",
		len(g.InscriptionsTickerInfoMap),
		len(g.UserTokensBalanceData),
		len(g.TokenUsersBalanceData),

		len(g.InscriptionsValidBRC20DataMap),

		len(g.InscriptionsValidTransferMap),
		len(g.InscriptionsInvalidTransferMap),
	)

	nswap := 0
	for _, m := range g.ModulesInfoMap {
		nswap += len(m.SwapPoolTotalBalanceDataMap)
	}

	nuser := 0
	for _, m := range g.ModulesInfoMap {
		nuser += len(m.UsersTokenBalanceDataMap)
	}

	log.Printf("process swap finish. module: %d, swap: %d, users: %d, validApprove: %d, invalidApprove: %d, validCommit: %d, invalidCommit: %d",
		len(g.ModulesInfoMap),
		nswap,
		nuser,

		len(g.InscriptionsValidApproveMap),
		len(g.InscriptionsInvalidApproveMap),

		len(g.InscriptionsValidCommitMap),
		len(g.InscriptionsInvalidCommitMap),
	)
}

// ProcessUpdateLatestBRC20Init
func (g *BRC20ModuleIndexer) ProcessUpdateLatestBRC20Init(brc20Datas []*model.InscriptionBRC20Data, totalDataCount int) {
	log.Printf("process swap init. total %d", len(brc20Datas))

	g.initBRC20()
	g.initModule()

	g.ProcessUpdateLatestBRC20Loop(brc20Datas, totalDataCount)
}
