/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"testing"
	"time"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"github.com/google/uuid"
	"github.com/test-go/testify/assert"
)

func TestDealInsertIDAAssetsData(t *testing.T) {
	idaInfo := getIdaAssetTest("ida_asset.json")
	// Test case 1: Insert new IDA assets data
	idaContractMap1 := map[string]*db.IDAContract{
		"contractAddr1": {
			ContractName:      "Contract1",
			ContractAddr:      "contractAddr1",
			TotalNormalAssets: 0,
			TotalAssets:       0,
			BlockHeight:       0,
		},
	}

	idaEventResult1 := &model.IDAEventData{
		IDACreatedMap: map[string][]*db.IDACreatedInfo{
			"contractAddr1": {
				{
					IDAInfo:   idaInfo,
					EventTime: 12345678,
				},
			},
		},
	}
	expectedResult1 := &db.IDAAssetsDataDB{
		IDAAssetDetail: []*db.IDAAssetDetail{
			{
				ID:                "newUUID",
				AssetCode:         "asset1",
				ContractName:      "Contract1",
				ContractAddr:      "contractAddr1",
				AssetName:         "Asset1",
				AssetEnName:       "Asset1",
				Category:          1,
				ImmediatelySupply: true,
				DataScale:         "100条",
				IndustryTitle:     "",
				Summary:           "Summary",
				Creator:           "Creator",
				Holder:            "Holder",
				TxID:              "txID",
				UserCategories:    "政府用户, 企业用户, 个人用户",
				UpdateCycleType:   common.UpdateCycleTypeStatic,
				UpdateTimeSpan:    "1 day",
				CreatedTime:       idaEventResult1.IDACreatedMap["contractAddr1"][0].EventTime,
				UpdatedTime:       idaEventResult1.IDACreatedMap["contractAddr1"][0].EventTime,
				SupplyTime:        time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		IDAAssetAttachment: []*db.IDAAssetAttachment{
			{
				ID:          "newUUID",
				AssetCode:   "asset1",
				Url:         "url",
				ContextType: 1,
			},
			{
				ID:          "newUUID",
				AssetCode:   "asset1",
				Url:         "url",
				ContextType: 1,
			},
		},
		IDAAssetData: []*db.IDADataAsset{
			{
				ID:          "newUUID",
				AssetCode:   "asset1",
				FieldName:   "column1",
				FieldType:   "string",
				FieldLength: 0,
			},
		},
		IDAAssetApi: []*db.IDAApiAsset{
			{
				ID:           "newUUID",
				AssetCode:    "asset1",
				Header:       "header",
				Url:          "url",
				Params:       "params",
				Response:     "response",
				Method:       "method",
				ResponseType: "responseType",
			},
		},
	}
	result1 := DealInsertIDAAssetsData(idaContractMap1, idaEventResult1)
	if result1 == nil {
		t.Errorf("Test case 1 failed: Expected non-nil result, got nil")
	} else if len(result1.IDAAssetDetail) != len(expectedResult1.IDAAssetDetail) {
		t.Errorf("Test case 1 failed: Expected %d asset details, got %d", len(expectedResult1.IDAAssetDetail), len(result1.IDAAssetDetail))
	} else if len(result1.IDAAssetAttachment) != len(expectedResult1.IDAAssetAttachment) {
		t.Errorf("Test case 1 failed: Expected %d asset attachments, got %d", len(expectedResult1.IDAAssetAttachment), len(result1.IDAAssetAttachment))
	} else if len(result1.IDAAssetApi) != len(expectedResult1.IDAAssetApi) {
		t.Errorf("Test case 1 failed: Expected %d asset apis, got %d", len(expectedResult1.IDAAssetApi), len(result1.IDAAssetApi))
	}
}

func TestGetIDAAssetUpdateTimeSpan(t *testing.T) {
	// Test case 1: Static update cycle
	updateCycle1 := standard.UpdateCycle{
		UpdateCycleType: common.UpdateCycleTypeStatic,
	}
	details1 := standard.Details{
		TimeSpan: "1 day",
	}
	expectedResult1 := "1 day"
	result1 := GetIDAAssetUpdateTimeSpan(updateCycle1, details1)
	if result1 != expectedResult1 {
		t.Errorf("Test case 1 failed")
	}

	// Test case 2: Periodic update cycle with minutes
	updateCycle2 := standard.UpdateCycle{
		UpdateCycleType: common.UpdateCycleTypePeriodic,
		UpdateCycleUnit: common.IDAUpdateCycleMinute,
		Cycle:           30,
	}
	expectedResult2 := "30分钟"
	result2 := GetIDAAssetUpdateTimeSpan(updateCycle2, standard.Details{})
	if result2 != expectedResult2 {
		t.Errorf("Test case 2 failed")
	}

	// Test case 3: Periodic update cycle with hours
	updateCycle3 := standard.UpdateCycle{
		UpdateCycleType: common.UpdateCycleTypePeriodic,
		UpdateCycleUnit: common.IDAUpdateCycleHour,
		Cycle:           24,
	}
	expectedResult3 := "24小时"
	result3 := GetIDAAssetUpdateTimeSpan(updateCycle3, standard.Details{})
	if result3 != expectedResult3 {
		t.Errorf("Test case 3 failed")
	}

	// Test case 4: Periodic update cycle with days
	updateCycle4 := standard.UpdateCycle{
		UpdateCycleType: common.UpdateCycleTypePeriodic,
		UpdateCycleUnit: common.IDAUpdateCycleday,
		Cycle:           7,
	}
	expectedResult4 := "7天"
	result4 := GetIDAAssetUpdateTimeSpan(updateCycle4, standard.Details{})
	if result4 != expectedResult4 {
		t.Errorf("Test case 4 failed")
	}
}

func TestGetIDAAssetDataScale(t *testing.T) {
	// Test case 1: Data scale type number
	dataScale1 := standard.DataScale{
		Type:  common.IDADataScaleTypeNum,
		Scale: 100,
	}
	expectedResult1 := "100条"
	result1 := GetIDAAssetDataScale(dataScale1)
	if result1 != expectedResult1 {
		t.Errorf("Test case 1 failed")
	}

	// Test case 2: Data scale type M
	dataScale2 := standard.DataScale{
		Type:  common.IDADataScaleTypeM,
		Scale: 1024,
	}
	expectedResult2 := "1024M"
	result2 := GetIDAAssetDataScale(dataScale2)
	if result2 != expectedResult2 {
		t.Errorf("Test case 2 failed")
	}

	// Test case 3: Data scale type G
	dataScale3 := standard.DataScale{
		Type:  common.IDADataScaleTypeG,
		Scale: 1,
	}
	expectedResult3 := "1G"
	result3 := GetIDAAssetDataScale(dataScale3)
	if result3 != expectedResult3 {
		t.Errorf("Test case 3 failed")
	}
}

func TestGetIDAUserCategories(t *testing.T) {
	// Test case 1: Single user category
	userCategories1 := []int{1}
	expectedResult1 := "政府用户"
	result1 := GetIDAUserCategories(userCategories1)
	if result1 != expectedResult1 {
		t.Errorf("Test case 1 failed")
	}

	// Test case 2: Multiple user categories
	userCategories2 := []int{1, 2, 3}
	expectedResult2 := "政府用户,企业用户,个人用户"
	result2 := GetIDAUserCategories(userCategories2)
	if result2 != expectedResult2 {
		t.Errorf("Test case 2 failed, result2:%v, expectedResult2:%v", result2, expectedResult2)
	}

	// Test case 3: No user categories
	userCategories3 := []int{}
	expectedResult3 := ""
	result3 := GetIDAUserCategories(userCategories3)
	if result3 != expectedResult3 {
		t.Errorf("Test case 3 failed")
	}
}

func TestDealUpdateIDAAssetsData(t *testing.T) {

	// Test case 1: Update IDA assets data with new fields
	assetDetailMap1 := map[string]*db.IDAAssetDetail{
		"asset1": {
			ID:                "asset1",
			AssetCode:         "asset1",
			ContractName:      "Contract1",
			ContractAddr:      "contractAddr1",
			AssetName:         "Asset1",
			AssetEnName:       "Asset1",
			Category:          1,
			ImmediatelySupply: true,
			DataScale:         "100条",
			IndustryTitle:     "",
			Summary:           "Summary",
			Creator:           "Creator",
			Holder:            "Holder",
			TxID:              "txID",
			UserCategories:    "政府用户, 企业用户, 个人用户",
			UpdateCycleType:   1,
		},
	}
	idaEventResult1 := &model.IDAEventData{
		IDAUpdatedMap: map[string][]*db.EventIDAUpdatedData{
			"asset1": {
				{
					Field:  db.KeyIDABasic,
					Update: "New Basic",
				},
				{
					Field:  db.KeyIDASUppLy,
					Update: "New Supply",
				},
				{
					Field:  db.KeyIDADetails,
					Update: "New Details",
				},
				{
					Field:  db.KeyIDAOwnership,
					Update: "New Ownership",
				},
				{
					Field:  db.KeyIDAColumns,
					Update: "New Columns",
				},
				{
					Field:  db.KeyIDAAPI,
					Update: "New API",
				},
			},
		},
	}
	expectedResult1 := &db.IDAAssetsUpdateDB{
		UpdateAssetDetails: []*db.IDAAssetDetail{
			{
				ID:                "asset1",
				AssetCode:         "asset1",
				ContractName:      "Contract1",
				ContractAddr:      "contractAddr1",
				AssetName:         "Asset1",
				AssetEnName:       "Asset1",
				Category:          1,
				ImmediatelySupply: true,
				DataScale:         "100条",
				IndustryTitle:     "",
				Summary:           "Summary",
				Creator:           "Creator",
				Holder:            "Holder",
				TxID:              "txID",
				UserCategories:    "政府用户, 企业用户, 个人用户",
				UpdateCycleType:   1,
				UpdatedTime:       idaEventResult1.EventTime,
			},
		},
		InsertAttachment:      nil,
		InsertIDAAssetData:    nil,
		InsertIDAAssetApi:     nil,
		DeleteAttachmentCodes: nil,
		DeleteAssetDataCodes:  nil,
		DeleteAssetApiCodes:   nil,
	}
	result1 := DealUpdateIDAAssetsData(idaEventResult1, assetDetailMap1)
	if len(result1.UpdateAssetDetails) != len(expectedResult1.UpdateAssetDetails) {
		t.Errorf("Test case 1 failed")
	}

	// Test case 2: Update IDA assets data with deleted assets
	assetDetailMap2 := map[string]*db.IDAAssetDetail{
		"asset1": {
			ID:                "asset1",
			AssetCode:         "asset1",
			ContractName:      "Contract1",
			ContractAddr:      "contractAddr1",
			AssetName:         "Asset1",
			AssetEnName:       "Asset1",
			Category:          1,
			ImmediatelySupply: true,
			DataScale:         "100条",
			IndustryTitle:     "",
			Summary:           "Summary",
			Creator:           "Creator",
			Holder:            "Holder",
			TxID:              "txID",
			UserCategories:    "政府用户, 企业用户, 个人用户",
			UpdateCycleType:   1,
		},
	}
	idaEventResult2 := &model.IDAEventData{
		IDADeletedCodeMap: map[string][]string{
			"contractAddr1": {"asset1"},
		},
	}
	expectedResult2 := &db.IDAAssetsUpdateDB{
		UpdateAssetDetails: []*db.IDAAssetDetail{
			{
				ID:                "asset1",
				AssetCode:         "asset1",
				ContractName:      "Contract1",
				ContractAddr:      "contractAddr1",
				AssetName:         "Asset1",
				AssetEnName:       "Asset1",
				Category:          1,
				ImmediatelySupply: true,
				DataScale:         "100条",
				IndustryTitle:     "",
				Summary:           "Summary",
				Creator:           "Creator",
				Holder:            "Holder",
				TxID:              "txID",
				UserCategories:    "政府用户, 企业用户, 个人用户",
				UpdateCycleType:   1,
				IsDeleted:         true,
				UpdatedTime:       idaEventResult2.EventTime,
			},
		},
	}
	result2 := DealUpdateIDAAssetsData(idaEventResult2, assetDetailMap2)
	if len(result2.UpdateAssetDetails) != len(expectedResult2.UpdateAssetDetails) {
		t.Errorf("Test case 2 failed")
	}

	// Test case 3: Update IDA assets data with no changes
	assetDetailMap3 := map[string]*db.IDAAssetDetail{
		"asset1": {
			ID:                "asset1",
			AssetCode:         "asset1",
			ContractName:      "Contract1",
			ContractAddr:      "contractAddr1",
			AssetName:         "Asset1",
			AssetEnName:       "Asset1",
			Category:          1,
			ImmediatelySupply: true,
			DataScale:         "100条",
			IndustryTitle:     "",
			Summary:           "Summary",
			Creator:           "Creator",
			Holder:            "Holder",
			TxID:              "txID",
			UserCategories:    "政府用户, 企业用户, 个人用户",
			UpdateCycleType:   1,
		},
	}
	idaEventResult3 := &model.IDAEventData{
		IDAUpdatedMap:     map[string][]*db.EventIDAUpdatedData{},
		IDADeletedCodeMap: map[string][]string{},
	}
	expectedResult3 := &db.IDAAssetsUpdateDB{
		UpdateAssetDetails:    nil,
		InsertAttachment:      nil,
		InsertIDAAssetData:    nil,
		InsertIDAAssetApi:     nil,
		DeleteAttachmentCodes: nil,
		DeleteAssetDataCodes:  nil,
		DeleteAssetApiCodes:   nil,
	}
	result3 := DealUpdateIDAAssetsData(idaEventResult3, assetDetailMap3)
	if len(result3.UpdateAssetDetails) != len(expectedResult3.UpdateAssetDetails) {
		t.Errorf("Test case 3 failed")
	}

	// Test case 4: Update IDA assets data with multiple assets
	assetDetailMap4 := map[string]*db.IDAAssetDetail{
		"asset1": {
			ID:                "asset1",
			AssetCode:         "asset1",
			ContractName:      "Contract1",
			ContractAddr:      "contractAddr1",
			AssetName:         "Asset1",
			AssetEnName:       "Asset1",
			Category:          1,
			ImmediatelySupply: true,
			DataScale:         "100条",
			IndustryTitle:     "",
			Summary:           "Summary",
			Creator:           "Creator",
			Holder:            "Holder",
			TxID:              "txID",
			UserCategories:    "政府用户, 企业用户, 个人用户",
			UpdateCycleType:   1,
		},
		"asset2": {
			ID:                "asset2",
			AssetCode:         "asset2",
			ContractName:      "Contract2",
			ContractAddr:      "contractAddr2",
			AssetName:         "Asset2",
			AssetEnName:       "Asset2",
			Category:          1,
			ImmediatelySupply: true,
			DataScale:         "200条",
			IndustryTitle:     "",
			Summary:           "Summary",
			Creator:           "Creator",
			Holder:            "Holder",
			TxID:              "txID",
			UserCategories:    "政府用户, 企业用户, 个人用户",
			UpdateCycleType:   1,
		},
	}
	idaEventResult4 := &model.IDAEventData{
		IDAUpdatedMap: map[string][]*db.EventIDAUpdatedData{
			"asset1": {
				{
					Field:  db.KeyIDABasic,
					Update: "New Basic",
				},
			},
			"asset2": {
				{
					Field:  db.KeyIDABasic,
					Update: "New Basic",
				},
			},
		},
	}
	expectedResult4 := &db.IDAAssetsUpdateDB{
		UpdateAssetDetails: []*db.IDAAssetDetail{
			{
				ID:                "asset1",
				AssetCode:         "asset1",
				ContractName:      "Contract1",
				ContractAddr:      "contractAddr1",
				AssetName:         "Asset1",
				AssetEnName:       "Asset1",
				Category:          1,
				ImmediatelySupply: true,
				DataScale:         "100条",
				IndustryTitle:     "",
				Summary:           "Summary",
				Creator:           "Creator",
				Holder:            "Holder",
				TxID:              "txID",
				UserCategories:    "政府用户, 企业用户, 个人用户",
				UpdateCycleType:   1,
				UpdatedTime:       idaEventResult4.EventTime,
			},
			{
				ID:                "asset2",
				AssetCode:         "asset2",
				ContractName:      "Contract2",
				ContractAddr:      "contractAddr2",
				AssetName:         "Asset2",
				AssetEnName:       "Asset2",
				Category:          1,
				ImmediatelySupply: true,
				DataScale:         "200条",
				IndustryTitle:     "",
				Summary:           "Summary",
				Creator:           "Creator",
				Holder:            "Holder",
				TxID:              "txID",
				UserCategories:    "政府用户, 企业用户, 个人用户",
				UpdateCycleType:   1,
				UpdatedTime:       idaEventResult4.EventTime,
			},
		},
	}
	result4 := DealUpdateIDAAssetsData(idaEventResult4, assetDetailMap4)
	if len(result4.UpdateAssetDetails) != len(expectedResult4.UpdateAssetDetails) {
		t.Errorf("Test case 4 failed")
	}
}

func TestHandleIDAAPI(t *testing.T) {
	// 准备输入
	updateData := `[{
		"header": "Authorization",
		"url": "https://api.example.com/data",
		"params": "id=1",
		"response": "{\"ok\":true}",
		"method": "GET",
		"resp_type": "json"
	}]`
	assetCode := "ASSET123"

	var insertAssetApis []*db.IDAApiAsset
	var deleteAssetApiCodes []string

	// 调用被测函数
	handleIDAAPI(updateData, assetCode, &insertAssetApis, &deleteAssetApiCodes)

	// 验证 insertAssetApis 是否正确
	assert.Len(t, insertAssetApis, 1)
	api := insertAssetApis[0]
	assert.Equal(t, assetCode, api.AssetCode)
	assert.Equal(t, "Authorization", api.Header)
	assert.Equal(t, "https://api.example.com/data", api.Url)
	assert.Equal(t, "GET", api.Method)
	assert.Equal(t, "json", api.ResponseType)

	// 验证 ID 是合法 UUID
	_, err := uuid.Parse(api.ID)
	assert.NoError(t, err)

	// 验证 deleteAssetApiCodes 是否追加了 assetCode
	assert.Contains(t, deleteAssetApiCodes, assetCode)
}

func TestHandleIDAColumns(t *testing.T) {
	// 准备输入 JSON
	updateData := `[{
		"name": "user_id",
		"data_type": "VARCHAR",
		"data_length": 36,
		"description": "用户ID",
		"is_primary_key": 1,
		"is_not_null": 1,
		"privacy_query": 0
	}]`
	assetCode := "ASSET456"

	var insertAssetDatas []*db.IDADataAsset
	var deleteAssetDataCodes []string

	// 调用被测函数
	handleIDAColumns(updateData, assetCode, &insertAssetDatas, &deleteAssetDataCodes)

	// 验证 insertAssetDatas
	assert.Len(t, insertAssetDatas, 1)
	col := insertAssetDatas[0]

	assert.Equal(t, assetCode, col.AssetCode)
	assert.Equal(t, "user_id", col.FieldName)
	assert.Equal(t, "VARCHAR", col.FieldType)
	assert.Equal(t, 36, col.FieldLength)
	assert.Equal(t, 1, col.IsPrimaryKey)
	assert.Equal(t, 1, col.IsNotNull)
	assert.Equal(t, 0, col.PrivacyQuery)

	// 验证 ID 是合法 UUID
	_, err := uuid.Parse(col.ID)
	assert.NoError(t, err)

	// 验证 deleteAssetDataCodes
	assert.Contains(t, deleteAssetDataCodes, assetCode)
}

func TestHandleIDABasic(t *testing.T) {
	updateData := `{
		"id": "BASIC123",
		"name": "用户资产",
		"enName": "UserAsset",
		"tags": ["金融", "数据"],
		"attachments": [
			{
				"hash": "hash123",
				"url": "http://example.com/file1.pdf",
				"type": 2,
				"size": 12345,
				"auditor": "Auditor1"
			},
			{
				"hash": "hash456",
				"url": "http://example.com/file2.png",
				"type": 1,
				"size": 67890,
				"auditor": "Auditor2"
			}
		],
		"category": 1,
		"industry": {
			"id": 100,
			"code": "BANK",
			"title": "银行业"
		},
		"summary": "这是一个用户资产的摘要",
		"creator": "Alice",
		"txID": "TX123456"
	}`

	assetCode := "ASSET789"
	assetDetail := &db.IDAAssetDetail{}
	var insertIDAAttachments []*db.IDAAssetAttachment
	var deleteAttachmentCodes []string

	// 调用被测函数
	handleIDABasic(updateData, assetDetail, assetCode, &insertIDAAttachments, &deleteAttachmentCodes)

	// 验证 assetDetail
	assert.Equal(t, "用户资产", assetDetail.AssetName)
	assert.Equal(t, "UserAsset", assetDetail.AssetEnName)
	assert.Equal(t, 1, assetDetail.Category)          // category 是 int
	assert.Equal(t, "银行业", assetDetail.IndustryTitle) // 来自 Industry.Title
	assert.Equal(t, "这是一个用户资产的摘要", assetDetail.Summary)
	assert.Equal(t, "Alice", assetDetail.Creator)
	assert.Equal(t, "TX123456", assetDetail.TxID)

	// 验证 insertIDAAttachments
	assert.Len(t, insertIDAAttachments, 2)

	att1 := insertIDAAttachments[0]
	assert.Equal(t, assetCode, att1.AssetCode)
	assert.Equal(t, "http://example.com/file1.pdf", att1.Url)
	assert.Equal(t, 2, att1.ContextType) // 注意这里 Type 是 int
	_, err := uuid.Parse(att1.ID)
	assert.NoError(t, err)

	att2 := insertIDAAttachments[1]
	assert.Equal(t, "http://example.com/file2.png", att2.Url)
	assert.Equal(t, 1, att2.ContextType)
	_, err = uuid.Parse(att2.ID)
	assert.NoError(t, err)

	// 验证 deleteAttachmentCodes
	assert.Contains(t, deleteAttachmentCodes, assetCode)
}

func TestProcessDataCategory(t *testing.T) {
	ida := &standard.IDAInfo{
		Basic: standard.Basic{
			ID: "ASSET123",
		},
		Columns: []standard.ColumnInfo{
			{
				Name:       "user_id",
				DataType:   "int",
				DataLength: 11,
			},
			{
				Name:       "username",
				DataType:   "varchar",
				DataLength: 50,
			},
		},
	}

	var insetAssetDatas []*db.IDADataAsset

	// 调用被测函数
	processDataCategory(ida, &insetAssetDatas)
}
