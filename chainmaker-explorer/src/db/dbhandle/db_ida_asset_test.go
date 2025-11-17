/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"testing"

	"chainmaker_web/src/db"

	"github.com/google/uuid"
)

func TestInsertIDAAttachments(t *testing.T) {
	newUUID := uuid.New().String()
	idaAttachments := []*db.IDAAssetAttachment{
		{ID: newUUID, AssetCode: "asset1", Url: "Url"},
	}
	err := InsertIDAAttachments(ChainID, idaAttachments)
	if err != nil {
		t.Errorf("InsertIDAAttachments failed: %s", err)
	}
}

func TestDeleteIDAAttachments(t *testing.T) {
	assetCodes := []string{"asset1", "asset2"}
	err := DeleteIDAAttachments(ChainID, assetCodes)
	if err != nil {
		t.Errorf("DeleteIDAAttachments failed: %s", err)
	}
}

func TestInsertIDAAssetApi(t *testing.T) {
	idaAssetApis := []*db.IDAApiAsset{
		{
			ID:        uuid.New().String(),
			AssetCode: "asset1",
			Url:       "Url",
		},
	}
	err := InsertIDAAssetApi(ChainID, idaAssetApis)
	if err != nil {
		t.Errorf("InsertIDAAssetApi failed: %s", err)
	}
}

func TestDeleteIDAApis(t *testing.T) {
	assetCodes := []string{"asset1", "asset2"}
	err := DeleteIDAApis(ChainID, assetCodes)
	if err != nil {
		t.Errorf("DeleteIDAApis failed: %s", err)
	}
}

func TestInsertIDAAssetData(t *testing.T) {
	idaDatas := []*db.IDADataAsset{
		{
			ID:        uuid.New().String(),
			AssetCode: "asset1",
			FieldName: "FieldName",
		},
	}
	err := InsertIDAAssetData(ChainID, idaDatas)
	if err != nil {
		t.Errorf("InsertIDAAssetData failed: %s", err)
	}
}

func TestDeleteIDADatas(t *testing.T) {
	assetCodes := []string{"asset1", "asset2"}
	err := DeleteIDADatas(ChainID, assetCodes)
	if err != nil {
		t.Errorf("DeleteIDADatas failed: %s", err)
	}
}

func TestGetIDAAssetAttachmentByCode(t *testing.T) {
	assetCode := "asset5"
	_, err := GetIDAAssetAttachmentByCode(ChainID, assetCode)
	if err != nil {
		t.Errorf("GetIDAAssetAttachmentByCode failed: %s", err)
	}
}

func TestGetIDAAssetDataByCode(t *testing.T) {
	assetCode := "asset2"
	_, err := GetIDAAssetDataByCode(ChainID, assetCode)
	if err != nil {
		t.Errorf("GetIDAAssetDataByCode failed: %s", err)
	}
}

func TestGetIDAAssetApiByCode(t *testing.T) {
	assetCode := "asset3"
	_, err := GetIDAAssetApiByCode(ChainID, assetCode)
	if err != nil {
		t.Errorf("GetIDAAssetApiByCode failed: %s", err)
	}
}
