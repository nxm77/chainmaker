/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
)

var AssetCode = "asset1"

func TestInsertIDADetail(t *testing.T) {
	idaDetails := []*db.IDAAssetDetail{
		{
			AssetCode: AssetCode,
		},
	}
	err := InsertIDADetail(ChainID, idaDetails)
	if err != nil {
		t.Errorf("InsertIDADetail failed: %s", err.Error())
	}
}

func TestGetIDAAssetList(t *testing.T) {
	offset := 0
	limit := 10
	_, err := GetIDAAssetList(offset, limit, ChainID, ContractAddr1, AssetCode)
	if err != nil {
		t.Errorf("GetIDAAssetList failed: %s", err.Error())
	}
}

func TestGetIDAAssetCount(t *testing.T) {
	_, err := GetIDAAssetCount(ChainID, ContractAddr1, AssetCode)
	if err != nil {
		t.Errorf("GetIDAAssetCount failed: %s", err.Error())
	}
}

func TestGetIDAAssetDetailByCode(t *testing.T) {
	_, err := GetIDAAssetDetailByCode(ChainID, ContractAddr1, AssetCode)
	if err != nil {
		t.Errorf("GetIDAAssetDetailByCode failed: %s", err.Error())
	}
}

func TestUpdateIDADetailByCode(t *testing.T) {
	idaDetail := &db.IDAAssetDetail{
		AssetCode: AssetCode,
	}
	UpdateIDADetailByCode(ChainID, idaDetail)
}

func TestGetIDAAssetDetailMapByCodes(t *testing.T) {
	assetCodes := []string{"asset1", "asset2"}
	_, err := GetIDAAssetDetailMapByCodes(ChainID, assetCodes)
	if err != nil {
		t.Errorf("GetIDAAssetDetailMapByCodes failed: %s", err.Error())
	}
}

func TestGetIDAAssetListWithEmptyParams(t *testing.T) {
	offset := 0
	limit := 10
	_, err := GetIDAAssetList(offset, limit, ChainID, ContractAddr1, AssetCode)
	if err != nil {
		t.Errorf("GetIDAAssetList failed: %s", err.Error())
	}
}

func TestGetIDAAssetCountWithEmptyParams(t *testing.T) {
	_, err := GetIDAAssetCount(ChainID, ContractAddr1, AssetCode)
	if err != nil {
		t.Errorf("GetIDAAssetCount failed: %s", err.Error())
	}
}

func TestGetIDAAssetDetailByCodeWithEmptyParams(t *testing.T) {
	_, err := GetIDAAssetDetailByCode(ChainID, ContractAddr1, "")
	if err == nil {
		t.Errorf("GetIDAAssetDetailByCode failed")
	}
}

func TestUpdateIDADetailByCodeWithEmptyParams(t *testing.T) {
	idaDetail := &db.IDAAssetDetail{
		AssetCode: "",
	}
	err := UpdateIDADetailByCode(ChainID, idaDetail)
	if err != nil {
		t.Errorf("UpdateIDADetailByCode failed: %s", err.Error())
	}
}

func TestGetIDAAssetDetailMapByCodesWithEmptyParams(t *testing.T) {
	assetCodes := []string{}
	_, err := GetIDAAssetDetailMapByCodes(ChainID, assetCodes)
	if err == nil {
		t.Errorf("GetIDAAssetDetailMapByCodes failed")
	}
}
