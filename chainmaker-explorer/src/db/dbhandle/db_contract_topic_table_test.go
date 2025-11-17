package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 InsertDecodeEventByABI
func TestInsertDecodeEventByABI(t *testing.T) {
	events := []map[string]interface{}{
		{"field1": "value1"},
		{"field2": "value2"},
	}

	// 执行测试
	_ = InsertDecodeEventByABI(db.UTchainID, "table1", events)
}

func TestCreateTopicTable(t *testing.T) {
	type tableStructure struct {
		Field1 string `json:"field1"`
	}
	err := db.CreateTopicTable(db.UTchainID, "tablename111", tableStructure{})
	assert.NoError(t, err)
}

// 测试 DeleteTopicDataRecord
func TestDeleteTopicDataRecord(t *testing.T) {
	DeleteTopicDataRecord(db.UTchainID, "", "table1")

	// 执行测试
	DeleteTopicDataRecord(db.UTchainID, "version", "table1")
	DeleteTopicDataRecord(db.UTchainID, "version", "tablename111")
}
