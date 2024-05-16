package loader

import (
	"brc20query/lib/brc20_swap/model"
	"brc20query/logger"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func LoadBRC20InputDataFromDB(brc20Datas chan *model.InscriptionBRC20Data, startHeight int, endHeight int) error {
	logger.Log.Info("LoadBRC20InputDataFromDB", zap.Int("startHeight", startHeight), zap.Int("endHeight", endHeight))

	batchLimit := 1024
	for offset := 0; ; offset += batchLimit {
		if lastHeight, err := loadBRC20InputDataFromDBOnBatch(
			brc20Datas, startHeight, endHeight, batchLimit, offset); err != nil {
			return err
		} else if lastHeight > int32(endHeight) || lastHeight == -1 {
			return nil
		} else if offset%10240 == 0 {
			logger.Log.Debug("LoadBRC20InputDataFromDB", zap.Int32("height", lastHeight), zap.Int("count", offset))
		}
	}
}

func loadBRC20InputDataFromDBOnBatch(brc20Datas chan *model.InscriptionBRC20Data,
	startHeight, endHeight int,
	queryLimit int, queryOffset int) (lastHeight int32, err error) {
	sql := fmt.Sprintf(`
SELECT ts.block_height, ts.inscription_id, ts.txcnt, ts.old_satpoint, ts.new_satpoint, 
	ts.new_pkscript, n2id.inscription_number, c.content, c.text_content, h.block_time
FROM ord_transfers AS ts 
LEFT JOIN ord_number_to_id AS n2id ON ts.inscription_id = n2id.inscription_id
LEFT JOIN ord_content AS c ON ts.inscription_id = c.inscription_id
LEFT JOIN block_hashes AS h ON ts.block_height = h.block_height 
WHERE ts.block_height >= %d AND ts.block_height < %d AND n2id.cursed_for_brc20 = false
ORDER BY ts.id LIMIT %d OFFSET %d
`, startHeight, endHeight, queryLimit, queryOffset)

	lastHeight = -1 // -1: no new block or error
	rows, err := SwapDB.Query(sql)
	if err != nil {
		return lastHeight, err
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
			return lastHeight, err
		}

		lastHeight = int32(block_height)
		if int(block_height) > endHeight {
			break
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
				return lastHeight, errors.WithMessagef(err, "inscription_id: %s", inscription_id)
			}

			offset, err = strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				return lastHeight, errors.WithMessagef(err, "inscription_id: %s", inscription_id)
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

		brc20Datas <- &data
	}

	return lastHeight, nil
}
