package loader

import (
	"brc20query/logger"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/unisat-wallet/libbrc20-indexer/model"
	"go.uber.org/zap"
)

func LoadBRC20InputDataFromDB(ctx context.Context, brc20Datas chan *model.InscriptionBRC20Data, startHeight int, endHeight int) error {
	logger.Log.Info("LoadBRC20InputDataFromDB", zap.Int("startHeight", startHeight), zap.Int("endHeight", endHeight))

	row := SwapDB.QueryRow(`select block_height from block_hashes order by block_height desc limit 1;`)
	var metaMaxHeight int
	if err := row.Scan(&metaMaxHeight); err != nil {
		return err
	}
	if endHeight <= 0 || metaMaxHeight < endHeight {
		endHeight = metaMaxHeight
	}

	for height := startHeight; height < endHeight; height++ {
		batchLimit := 10240
		st := time.Now()
		offset := 0
		count := 0
		blkDatas := make([]*model.InscriptionBRC20Data, 0)
		for ; ; offset += batchLimit {
			if datas, err := loadBRC20InputDataFromDBOnBatch(height, batchLimit, offset); err != nil {
				return err
			} else if len(datas) == 0 {
				break
			} else {
				blkDatas = append(blkDatas, datas...)
				count += len(datas)
				if len(datas) < batchLimit {
					break
				}
			}
		}
		logger.Log.Debug("LoadBRC20InputDataFromDB",
			zap.Int("height", height),
			zap.Int("count", count),
			zap.String("duration", time.Since(st).String()))

		for _, data := range blkDatas {
			brc20Datas <- data
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
	return nil
}

func loadBRC20InputDataFromDBOnBatch(height int, queryLimit int, queryOffset int) (datas []*model.InscriptionBRC20Data, err error) {
	sql := fmt.Sprintf(`
SELECT ts.block_height, ts.inscription_id, ts.txcnt, ts.old_satpoint, ts.new_satpoint,
	ts.new_pkscript, n2id.inscription_number, c.content, c.text_content, h.block_time
FROM ord_transfers AS ts
INNER JOIN ord_number_to_id AS n2id ON ts.inscription_id = n2id.inscription_id
INNER JOIN ord_content AS c ON ts.inscription_id = c.inscription_id
INNER JOIN block_hashes AS h ON ts.block_height = h.block_height 
WHERE ts.block_height = %d AND n2id.cursed_for_brc20 = false
ORDER BY ts.id LIMIT %d OFFSET %d
`, height, queryLimit, queryOffset)

	rows, err := SwapDB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			block_height       uint32
			inscription_id     string
			txcnt              uint32
			old_satpoint       string
			new_satpoint       string
			new_pkscript       string
			inscription_number int64
			content            []byte
			text_content       []byte
			block_time         uint32
		)

		if err := rows.Scan(
			&block_height, &inscription_id, &txcnt, &old_satpoint, &new_satpoint,
			&new_pkscript, &inscription_number, &content, &text_content, &block_time); err != nil {
			return datas, err
		}

		var (
			is_transfer bool   = false
			txid        string = ""
			vout        uint64 = 0
			offset      uint64 = 0
			contentBody []byte = nil
			err         error  = nil
		)

		is_transfer = false
		if txcnt > 0 {
			is_transfer = true
		}

		contentBody = content
		if len(content) == 0 {
			contentBody = text_content
		}

		{
			parts := strings.Split(new_satpoint, ":")
			txid = parts[0]

			vout, err = strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				logger.Log.Debug("loadBRC20InputDataFromDBOnBatch", zap.String("inscription_id", inscription_id))
				return datas, err
			}

			offset, err = strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				logger.Log.Debug("loadBRC20InputDataFromDBOnBatch", zap.String("inscription_id", inscription_id))
				return datas, err
			}
		}

		data := model.InscriptionBRC20Data{
			IsTransfer:        is_transfer,
			TxId:              txid,
			Idx:               uint32(inscription_number),
			Vout:              uint32(vout),
			Offset:            offset,
			Satoshi:           546,
			PkScript:          new_pkscript,
			Fee:               0,
			InscriptionNumber: inscription_number,
			ContentBody:       contentBody,
			CreateIdxKey:      inscription_id,
			Height:            block_height,
			TxIdx:             0,
			BlockTime:         block_time,
			Sequence:          uint16(txcnt),
			InscriptionId:     inscription_id,
		}

		datas = append(datas, &data)
	}

	return datas, nil
}
