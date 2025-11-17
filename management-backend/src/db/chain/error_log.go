package chain

import (
	"management_backend/src/db/common"
	"management_backend/src/db/connection"
	loggers "management_backend/src/logger"
)

// CreateErrorLogRecord create error log record
func CreateErrorLogRecord(errorLog *common.ChainErrorLog) error {
	if err := connection.DB.Create(&errorLog).Error; err != nil {
		loggers.DBLogger.Error("save chain error log record failed: " + err.Error())
		return err
	}
	return nil
}

// GetLogInfoById getLogInfoById
func GetLogInfoById(id int64) (*common.ChainErrorLog, error) {
	var errLog common.ChainErrorLog
	if err := connection.DB.Model(errLog).Where("id = ?", id).Find(&errLog).Error; err != nil {
		loggers.DBLogger.Error("QueryLogInfoById Failed: " + err.Error())
		return nil, err
	}
	return &errLog, nil
}

// GetLogList getLogList
func GetLogList(chainId string, offset int, limit int) (int64, []*common.ChainErrorLog, error) {
	var (
		count   int64
		logList []*common.ChainErrorLog
		err     error
	)

	if err = connection.DB.Model(&common.ChainErrorLog{}).Where("chain_id = ?", chainId).Count(&count).Error; err != nil {
		loggers.DBLogger.Error("GetLogList Failed: " + err.Error())
		return count, logList, err
	}

	if err = connection.DB.Model(&common.ChainErrorLog{}).Where("chain_id = ?", chainId).
		Order("log_time desc").
		Offset(offset).Limit(limit).Find(&logList).Error; err != nil {
		loggers.DBLogger.Error("GetLogList Failed: " + err.Error())
		return count, logList, err
	}
	return count, logList, err
}
