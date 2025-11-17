/*
Package contract_invoke comment
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package contract_invoke

import (
	"encoding/hex"
	"strings"

	"management_backend/src/ctrl/ca"
	"management_backend/src/db"
	"management_backend/src/global"
	"management_backend/src/utils"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"
	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	commonutils "chainmaker.org/chainmaker/utils/v2"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// GetEvmKv get evm kv
func GetEvmKv(abiKey, methodName, chainMode, hashType string,
	parameterParams []*ParameterParams, crtBytes []byte) ([]*pbcommon.KeyValuePair, string, error) {
	id, userId, hash, err := ca.ResolveUploadKey(abiKey)
	if err != nil {
		return nil, "", err
	}
	upload, err := db.GetUploadByIdAndUserIdAndHash(id, userId, hash)
	if err != nil {
		return nil, "", err
	}

	myAbi, err := abi.JSON(strings.NewReader(string(upload.Content)))
	if err != nil {
		return nil, "", err
	}
	method := myAbi.Methods[methodName]

	var paramMapList []utils.Param

	for _, input := range method.Inputs {
		for _, parameterParam := range parameterParams {
			if parameterParam.Key == input.Name {
				paramMap := make(utils.Param)
				paramMap[input.Type.String()] = parameterParam.Value
				paramMapList = append(paramMapList, paramMap)
			}
		}
	}

	if len(method.Inputs) > 0 && len(paramMapList) < 1 {
		var addr string
		if chainMode == global.PUBLIC {
			publicKey, publicKeyErr := asym.PublicKeyFromPEM([]byte(crtBytes))
			if publicKeyErr != nil {
				return nil, "", publicKeyErr
			}
			addr, err = commonutils.PkToAddrStr(publicKey, pbconfig.AddrType_ETHEREUM, crypto.HashAlgoMap[hashType])
		} else {
			_, addr, _, err =
				utils.MakeAddrAndSkiFromCrtBytes(crtBytes)
		}

		if err != nil {
			return nil, "", err
		}
		paramMap := make(utils.Param)
		paramMap["address"] = addr
		paramMapList = append(paramMapList, paramMap)
	}

	paramBytes, err := GetPaddedParam(&method, paramMapList)
	if err != nil {
		return nil, "", err
	}
	inputData := append(method.ID, paramBytes...)

	inputDataHexStr := hex.EncodeToString(inputData)

	kvs := []*pbcommon.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(inputDataHexStr),
		},
	}

	return kvs, inputDataHexStr[0:8], nil
}

// GetPaddedParam get padded param
func GetPaddedParam(method *abi.Method, param []utils.Param) ([]byte, error) {
	values, err := utils.ConvertParam(param)
	if err != nil {
		return nil, err
	}
	// convert params to bytes
	return method.Inputs.PackValues(values)
}
