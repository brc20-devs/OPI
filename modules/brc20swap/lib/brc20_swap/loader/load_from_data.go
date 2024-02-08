package loader

import (
	"brc20query/lib/brc20_swap/model"
	"brc20query/lib/brc20_swap/utils"
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetBRC20InputDataLineCount(fname string) (int, error) {
	file, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	max := 128 * 1024 * 1024
	buf := make([]byte, max)
	scanner.Buffer(buf, max)

	count := 0
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func LoadBRC20InputData(fname string, brc20Datas chan *model.InscriptionBRC20Data, endHeight int) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	max := 128 * 1024 * 1024
	buf := make([]byte, max)
	scanner.Buffer(buf, max)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")

		if len(fields) != 13 {
			return fmt.Errorf("invalid data format")
		}

		var data model.InscriptionBRC20Data

		sequence, err := strconv.ParseUint(fields[0], 10, 16)
		if err != nil {
			return err
		}
		data.Sequence = uint16(sequence)
		data.IsTransfer = (data.Sequence > 0)

		txid, err := hex.DecodeString(fields[1])
		if err != nil {
			return err
		}
		data.TxId = string(utils.ReverseBytes([]byte(txid)))

		idx, err := strconv.ParseUint(fields[2], 10, 32)
		if err != nil {
			return err
		}
		data.Idx = uint32(idx)

		vout, err := strconv.ParseUint(fields[3], 10, 32)
		if err != nil {
			return err
		}
		data.Vout = uint32(vout)

		offset, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return err
		}
		data.Offset = uint64(offset)

		satoshi, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			return err
		}
		data.Satoshi = uint64(satoshi)

		pkScript, err := hex.DecodeString(fields[6])
		if err != nil {
			return err
		}
		data.PkScript = string(pkScript)

		inscriptionNumber, err := strconv.ParseInt(fields[7], 10, 64)
		if err != nil {
			return err
		}
		data.InscriptionNumber = int64(inscriptionNumber)

		content, err := hex.DecodeString(fields[8])
		if err != nil {
			return err
		}
		data.ContentBody = content
		data.CreateIdxKey = string(fields[9])

		height, err := strconv.ParseUint(fields[10], 10, 32)
		if err != nil {
			return err
		}
		data.Height = uint32(height)

		txIdx, err := strconv.ParseUint(fields[11], 10, 32)
		if err != nil {
			return err
		}
		data.TxIdx = uint32(txIdx)

		blockTime, err := strconv.ParseUint(fields[12], 10, 32)
		if err != nil {
			return err
		}
		data.BlockTime = uint32(blockTime)

		if int(height) > endHeight {
			break
		}

		brc20Datas <- &data
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
