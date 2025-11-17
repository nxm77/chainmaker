/*
Package db comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package db

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestNewDIDSaveData(t *testing.T) {
	didSaveData := NewDIDSaveData()

	assert.NotNil(t, didSaveData)
	assert.NotNil(t, didSaveData.SaveDIDDetails)
	assert.NotNil(t, didSaveData.InsertDIDHistory)
	assert.NotNil(t, didSaveData.UpdateDIDStatus)
	assert.NotNil(t, didSaveData.UpdateDIDIssuer)
	assert.NotNil(t, didSaveData.SaveVCTemplates)
	assert.NotNil(t, didSaveData.InsertVCHistory)
	assert.NotNil(t, didSaveData.UpdateVCStatus)

	assert.Equal(t, 0, len(didSaveData.SaveDIDDetails))
	assert.Equal(t, 0, len(didSaveData.InsertDIDHistory))
	assert.Equal(t, 0, len(didSaveData.UpdateDIDStatus))
	assert.Equal(t, 0, len(didSaveData.UpdateDIDIssuer))
	assert.Equal(t, 0, len(didSaveData.SaveVCTemplates))
	assert.Equal(t, 0, len(didSaveData.InsertVCHistory))
	assert.Equal(t, 0, len(didSaveData.UpdateVCStatus))
}
