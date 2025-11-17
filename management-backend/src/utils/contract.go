/*
Package utils comment
Copyright (C) BABEC. All rights reserved.
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	loggers "management_backend/src/logger"
	"strings"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/common/v2/evmutils"
	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"
	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	commonutils "chainmaker.org/chainmaker/utils/v2"
	"github.com/ethereum/go-ethereum/accounts/abi"

	"management_backend/src/global"
)

// GetEvmMethodsByAbi get evm methods by abi
func GetEvmMethodsByAbi(content []byte) (string, int, error) {
	myAbi, err := abi.JSON(strings.NewReader(string(content)))
	if err != nil {
		return "", -1, err
	}
	var methods = make([]*global.Method, 0)

	// 0：正常方法 1：构造函数
	functionType := global.FUNCTION
	if len(myAbi.Constructor.Inputs) > 0 {
		functionType = global.CONSTRUCTOR
	}
	for methodName, methodVale := range myAbi.Methods {
		method := &global.Method{
			MethodFunc: "invoke",
		}
		method.MethodName = methodName

		var methodKeyStr string
		inputs := methodVale.Inputs
		for _, input := range inputs {
			methodKeyStr = methodKeyStr + input.Name + ","
		}

		methodKeyStr = strings.TrimRight(methodKeyStr, ",")
		method.MethodKey = methodKeyStr
		methods = append(methods, method)
	}

	methodJson, err := json.Marshal(methods)
	if err != nil {
		return "", -1, err
	}

	methodStr := string(methodJson)
	if methodStr == global.NULL {
		methodStr = ""
	}
	return methodStr, functionType, nil
}

// GetConstructorKeyValuePair get constructor key value pair
func GetConstructorKeyValuePair(chainMode, hashType string, crtBytes []byte,
	content []byte, kvParams []*pbcommon.KeyValuePair) ([]*pbcommon.KeyValuePair, error) {
	var addrInt *evmutils.Int
	var err error
	if chainMode == global.PUBLIC {
		publicKey, publicKeyErr := asym.PublicKeyFromPEM([]byte(crtBytes))
		if publicKeyErr != nil {
			return nil, publicKeyErr
		}
		addrInt, err = commonutils.PkToAddrInt(publicKey, pbconfig.AddrType_ETHEREUM, crypto.HashAlgoMap[hashType])
		if err != nil {
			return nil, err
		}
	} else {
		_, _, client1AddrSki, skiErr :=
			MakeAddrAndSkiFromCrtBytes(crtBytes)
		if skiErr != nil {
			return nil, skiErr
		}
		addrInt, err = evmutils.MakeAddressFromHex(client1AddrSki)
		if err != nil {
			return nil, err
		}
	}
	addr := evmutils.BigToAddress(addrInt)
	myAbi, err := abi.JSON(strings.NewReader(string(content)))
	if err != nil {
		return nil, err
	}

	var paramMapList []Param
	if len(kvParams) > 0 {
		for _, param := range kvParams {
			for _, input := range myAbi.Constructor.Inputs {
				if param.Key == input.Name {
					paramMap := make(Param)
					paramMap[input.Type.String()] = string(param.Value)
					paramMapList = append(paramMapList, paramMap)
				}
			}
		}
	} else {
		paramMap := make(Param)
		paramMap["address"] = addr
		paramMapList = append(paramMapList, paramMap)
	}

	args, err := ConvertParam(paramMapList)
	if err != nil {
		return nil, err
	}

	loggers.WebLogger.Infof("evm args: %v", args)

	dataByte, err := myAbi.Pack("", args...)
	if err != nil {
		return nil, err
	}

	data := hex.EncodeToString(dataByte)
	pairs := []*pbcommon.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(data),
		},
	}

	return pairs, nil
}

// MakeAddrAndSkiFromCrtBytes MakeAddrAndSkiFromCrtBytes
func MakeAddrAndSkiFromCrtBytes(crtBytes []byte) (string, string, string, error) {
	blockCrt, _ := pem.Decode(crtBytes)
	crt, err := bcx509.ParseCertificate(blockCrt.Bytes)
	if err != nil {
		return "", "", "", err
	}

	ski := hex.EncodeToString(crt.SubjectKeyId)
	addrInt, err := evmutils.MakeAddressFromHex(ski)
	if err != nil {
		return "", "", "", err
	}

	loggers.WebLogger.Info(fmt.Sprintf("get address and ski: 0x%s", addrInt.AsStringKey()))

	return addrInt.String(), fmt.Sprintf("0x%x", addrInt.AsStringKey()), ski, nil
}
