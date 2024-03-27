package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func DecodeTokensFromSwapPair(tickPair string) (token0, token1 string, err error) {
	if len(tickPair) != 9 || tickPair[4] != '/' {
		return "", "", errors.New("func: removeLiq tickPair invalid")
	}
	token0 = tickPair[:4]
	token1 = tickPair[5:]

	return token0, token1, nil
}

// single sha256 hash
func GetSha256(data []byte) (hash []byte) {
	sha := sha256.New()
	sha.Write(data[:])
	hash = sha.Sum(nil)
	return
}

func GetHash256(data []byte) (hash []byte) {
	sha := sha256.New()
	sha.Write(data[:])
	tmp := sha.Sum(nil)
	sha.Reset()
	sha.Write(tmp)
	hash = sha.Sum(nil)
	return
}

func HashString(data []byte) (res string) {
	length := 32
	var reverseData [32]byte

	// need reverse
	for i := 0; i < length; i++ {
		reverseData[i] = data[length-i-1]
	}
	return hex.EncodeToString(reverseData[:])
}

func ReverseBytes(data []byte) (result []byte) {
	for _, b := range data {
		result = append([]byte{b}, result...)
	}
	return result
}

// GetAddressFromScript Use btcsuite to get address
func GetAddressFromScript(script []byte, params *chaincfg.Params) (string, error) {
	scriptClass, addresses, _, err := txscript.ExtractPkScriptAddrs(script, params)
	if err != nil {
		return "", fmt.Errorf("failed to get address: %v", err)
	}

	if len(addresses) == 0 {
		return "", fmt.Errorf("noaddress")
	}

	if scriptClass == txscript.NonStandardTy {
		return "", fmt.Errorf("non-standard")
	}

	return addresses[0].EncodeAddress(), nil
}

func GetModuleFromScript(script []byte) (module string, ok bool) {
	if len(script) < 34 || len(script) > 38 {
		return "", false
	}
	if script[0] != 0x6a {
		return "", false
	}
	if int(script[1])+2 != len(script) {
		return "", false
	}

	var idx uint32
	if script[1] <= 32 {
		idx = uint32(0)
	} else if script[1] <= 33 {
		idx = uint32(script[34])
	} else if script[1] <= 34 {
		idx = uint32(binary.LittleEndian.Uint16(script[34:36]))
	} else if script[1] <= 35 {
		idx = uint32(script[34]) | uint32(script[35])<<8 | uint32(script[36])<<16
	} else if script[1] <= 36 {
		idx = binary.LittleEndian.Uint32(script[34:38])
	}

	module = fmt.Sprintf("%si%d", HashString(script[2:34]), idx)
	return module, true
}

func GetInnerSwapPoolNameByToken(token0, token1 string) (poolPair string) {
	token0 = strings.ToLower(token0)
	token1 = strings.ToLower(token1)

	if token0 > token1 {
		poolPair = fmt.Sprintf("%s/%s", token1, token0)
	} else {
		poolPair = fmt.Sprintf("%s/%s", token0, token1)
	}
	return poolPair
}
