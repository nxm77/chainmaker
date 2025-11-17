/*
Package utils comment
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"bytes"
	"crypto/rand"
	"errors"
	loggers "management_backend/src/logger"
	random "math/rand"
	"reflect"
	"strconv"

	"chainmaker.org/chainmaker/common/v2/evmutils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	// RandomRange random
	RandomRange = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

// Param param
type Param map[string]interface{}

// Base64Encode base64Encode
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode base64Decode
func Base64Decode(data string) []byte {
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return decodeBytes
}

// RandomString random string
func RandomString(len int) string {
	var container string
	b := bytes.NewBufferString(RandomRange)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(RandomRange[randomInt.Int64()])
	}
	return container
}

// CurrentMillSeconds current mill seconds
func CurrentMillSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// CurrentSeconds current seconds
func CurrentSeconds() int64 {
	return time.Now().UnixNano() / 1e9
}

// PathExists path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ConvertToPercent convert to percent
func ConvertToPercent(percent float64) string {
	doubleDecimal := fmt.Sprintf("%.2f", 100*percent)
	// Strip trailing zeroes
	for doubleDecimal[len(doubleDecimal)-1] == '0' {
		doubleDecimal = doubleDecimal[:len(doubleDecimal)-1]
	}
	// Strip the decimal point if it's trailing.
	if doubleDecimal[len(doubleDecimal)-1] == '.' {
		doubleDecimal = doubleDecimal[:len(doubleDecimal)-1]
	}
	return fmt.Sprintf("%s%%", doubleDecimal)
}

// GetHostFromAddress get host from address
func GetHostFromAddress(addr string) string {
	s := strings.Split(addr, ":")
	return s[0]
}

var (
	globalChainMaxHeightMp sync.Map //chainid - blockheight(string-> int64)
)

// GetMaxHeight 获取链最大高度
func GetMaxHeight(chainId string) int64 {
	height, _ := globalChainMaxHeightMp.Load(chainId)
	return height.(int64)
}

// SetMaxHeight 对链最大高度赋值
func SetMaxHeight(chainId string, height int64) {
	globalChainMaxHeightMp.Store(chainId, height)
}

// nolint
func RandomSleepTime() int {
	return random.Intn(20) + 50

}

// ConvertParam Convert Param
func ConvertParam(param []Param) ([]interface{}, error) {
	values := make([]interface{}, 0)

	for _, p := range param {
		if len(p) != 1 {
			return nil, fmt.Errorf("invalid param %+v", p)
		}
		for k, v := range p {
			if k == "uint" {
				k = "uint256"
			} else if strings.HasPrefix(k, "uint[") {
				k = strings.Replace(k, "uint[", "uint256[", 1)
			}
			ty, err := abi.NewType(k, "", nil)
			if err != nil {
				return nil, fmt.Errorf("invalid param %+v: %+v", p, err)
			}
			if ty.T == abi.SliceTy || ty.T == abi.ArrayTy {
				if ty.Elem.T == abi.AddressTy {
					v, err = ConvertTyToAddress(v)
					if err != nil {
						return nil, err
					}
				}
				if (ty.Elem.T == abi.IntTy || ty.Elem.T == abi.UintTy) && reflect.TypeOf(v).Elem().Kind() == reflect.Interface {
					v, err = GetTmpData(ty, v)
					if err != nil {
						return nil, err
					}
				}
			}

			if ty.T == abi.AddressTy {
				if v, err = ConvertToAddress(v); err != nil {
					return nil, err
				}
			}

			if (ty.T == abi.IntTy || ty.T == abi.UintTy) && reflect.TypeOf(v).Kind() == reflect.String {
				v = ConvertToInt(ty, v)
			}
			values = append(values, v)
		}
	}
	return values, nil
}

// ConvertTyToAddress Convert TyToAddress
func ConvertTyToAddress(v interface{}) (interface{}, error) {
	tmp, ok := v.([]interface{})
	if !ok {
		loggers.WebLogger.Error("v not is []interface")
		return nil, errors.New("v not is []interface")
	}
	v = make([]common.Address, 0)
	for i := range tmp {
		addr, err := ConvertToAddress(tmp[i])
		if err != nil {
			return nil, err
		}
		v = append(v.([]common.Address), addr)
	}
	return v, nil
}

// GetTmpData Get Tmp Data
func GetTmpData(ty abi.Type, v interface{}) (res interface{}, err error) {
	if ty.Elem.Size > 64 {
		tmp := make([]*big.Int, 0)
		for _, i := range v.([]interface{}) {
			if s, ok := i.(string); ok {
				value, _ := new(big.Int).SetString(s, 10)
				tmp = append(tmp, value)
			} else {
				return nil, fmt.Errorf("abi: cannot use %T as type string as argument", i)
			}
		}
		res = tmp
	} else {
		tmpI := make([]interface{}, 0)
		for _, i := range v.([]interface{}) {
			if s, ok := i.(string); ok {
				value, valueErr := strconv.ParseUint(s, 10, ty.Elem.Size)
				if valueErr != nil {
					return nil, valueErr
				}
				tmpI = append(tmpI, value)
			} else {
				return nil, fmt.Errorf("abi: cannot use %T as type string as argument", i)
			}
		}
		switch ty.Elem.Size {
		case 8:
			tmp := make([]uint8, len(tmpI))
			for i, sv := range tmpI {
				tmp[i] = uint8(sv.(uint64))
			}
			res = tmp
		case 16:
			tmp := make([]uint16, len(tmpI))
			for i, sv := range tmpI {
				tmp[i] = uint16(sv.(uint64))
			}
			res = tmp
		case 32:
			tmp := make([]uint32, len(tmpI))
			for i, sv := range tmpI {
				tmp[i] = uint32(sv.(uint64))
			}
			res = tmp
		case 64:
			tmp := make([]uint64, len(tmpI))
			ok := false
			for i, sv := range tmpI {
				tmp[i], ok = sv.(uint64)
				if !ok {
					loggers.WebLogger.Error("v not is uint64")
				}
			}
			res = tmp
		}
	}
	return res, nil
}

// ConvertToAddress Convert To Address
func ConvertToAddress(v interface{}) (common.Address, error) {
	switch data := v.(type) {
	case string:
		if !common.IsHexAddress(data) {
			return common.Address{}, fmt.Errorf("invalid address %s", data)
		}
		return common.HexToAddress(data), nil
	case evmutils.Address:
		if a, ok := v.(evmutils.Address); ok {
			return common.BytesToAddress(a[:]), nil
		}
	}
	return common.Address{}, fmt.Errorf("invalid address %v", v)
}

// ConvertToInt Convert To Int
func ConvertToInt(ty abi.Type, v interface{}) interface{} {
	if ty.T == abi.IntTy && ty.Size <= 64 {
		tmp, _ := strconv.ParseInt(v.(string), 10, ty.Size)
		switch ty.Size {
		case 8:
			v = int8(tmp)
		case 16:
			v = int16(tmp)
		case 32:
			v = int32(tmp)
		case 64:
			v = int64(tmp)
		}
	} else if ty.T == abi.UintTy && ty.Size <= 64 {
		tmp, _ := strconv.ParseUint(v.(string), 10, ty.Size)
		switch ty.Size {
		case 8:
			v = uint8(tmp)
		case 16:
			v = uint16(tmp)
		case 32:
			v = uint32(tmp)
		case 64:
			v = uint64(tmp)
		}
	} else {
		v, _ = new(big.Int).SetString(v.(string), 10)
	}
	return v
}
