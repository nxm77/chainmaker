// Package logic 逻辑处理
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	loggers "chainmaker_web/src/logger"
	"chainmaker_web/src/utils"
	"encoding/json"
	"reflect"

	syncLogic "chainmaker_web/src/sync/logic"

	"github.com/google/uuid"
)

var log = loggers.GetLogger(loggers.MODULE_WEB)

// validateABIFile 验证ABI文件
func validateABIFile(parsedABI []*utils.ContractABI) error {
	// 遍历解析后的ABI
	for _, abi := range parsedABI {
		// 如果ABI类型不是函数类型、构造函数类型或事件类型，则返回错误
		if abi.Type != utils.ABIFunctionType &&
			abi.Type != utils.ABIConstructorType &&
			abi.Type != utils.ABIEventType {
			return entity.GetErrorMsg(entity.ErrorContractAbiFunctionType)
		}

		// 如果ABI类型是事件类型
		if abi.Type == utils.ABIEventType {
			//检查是否包含系统字段
			for _, input := range abi.Inputs {
				// 如果包含系统字段，则返回错误
				if _, exists := db.ABISystemFields[input.Name]; exists {
					return entity.GetErrorMsg(entity.ErrorContractAbiEventInputs)
				}
			}
		}
	}
	// 如果没有错误，则返回nil
	return nil
}

func SaveContractABI(params *entity.UploadContractAbiParams, upgradeContract *db.UpgradeContractTransaction) error {
	// 1. 解析并验证基础数据
	parsedABI, err := utils.ParseContractABI(params.AbiJson)
	if err != nil {
		// 记录解析失败日志
		log.Errorf("ParseContractABI ABI解析失败: %v", err)
		return err
	}

	// 2. 检验ABI文件是否合法
	if errABI := validateABIFile(parsedABI); err != nil {
		return errABI
	}

	// 2. 保存/更新ABI文件
	_, err = saveOrUpdateABIFile(params, upgradeContract, parsedABI)
	if err != nil {
		return err
	}

	// 3. 处理事件表逻辑
	return processEventTables(params, upgradeContract, parsedABI)
}

func processEventTables(params *entity.UploadContractAbiParams, upgradeContract *db.UpgradeContractTransaction,
	parsedABI []*utils.ContractABI) error {
	// 处理每个事件
	for _, event := range parsedABI {
		if event.Type != utils.ABIEventType {
			continue
		}

		// 生成表名（新格式）
		tableName := utils.GenerateTableName(params.ContractAddr, event.Inputs)
		if tableName == "" {
			log.Errorf("processEventTables 生成表名失败: %s", event.Name)
			return entity.GetErrorMsg(entity.ErrorGenerateTableNameFailed)
		}

		// 处理当前版本记录
		if err := handleEventVersion(params, upgradeContract, event, tableName); err != nil {
			return err
		}
	}

	return nil
}

// 3.2 处理单个事件版本
// 处理合约ABI版本
func handleEventVersion(params *entity.UploadContractAbiParams, upgradeContract *db.UpgradeContractTransaction,
	event *utils.ContractABI, tableName string) error {
	isNew := true
	// 查询当前版本是否已有记录
	existing, err := dbhandle.GetContractABITopic(params.ChainId, params.ContractAddr, params.ContractVersion, event.Name)
	if err != nil {
		log.Errorf("查询当前版本记录失败: %v", err)
		return err
	}

	// 如果存在记录，判断表名是否一致，不一致则更新表名
	if existing != nil {
		oldTableName := existing.TopicTableName // 先保存旧表名

		// 如果表名不一致，需要更新表名称，并删除历史解析记录
		if oldTableName != tableName {
			existing.TopicTableName = tableName // 更新为新表名

			// 1. 先更新数据库记录
			if err := dbhandle.UpdateContractABITopic(params.ChainId, existing); err != nil {
				return err
			}

			// 2. 删除旧表的数据（使用旧表名）
			dbhandle.DeleteTopicDataRecord(params.ChainId, params.ContractVersion, oldTableName)
			log.Infof("表名更新完成 | old_table=%s -> new_table=%s", oldTableName, tableName)
		} else {
			isNew = false
		}
	}

	if isNew {
		// 如果是新版本，直接创建新记录
		err := createTopicRecord(params, upgradeContract, event, tableName)
		return err
	}
	return nil
}

// 3.3 创建Topic记录
func createTopicRecord(params *entity.UploadContractAbiParams, upgradeContract *db.UpgradeContractTransaction,
	event *utils.ContractABI, tableName string) error {
	fields := make([]*db.DynamicStructField, 0)
	for _, input := range event.Inputs {
		fields = append(fields, &db.DynamicStructField{
			Name:    input.Name,
			Type:    input.Type,
			Indexed: input.Indexed,
		})
	}

	// 需要创建新表
	structType := db.CreateDynamicStructWithSystemFields(fields)
	dynamicStruct := reflect.New(structType).Interface()
	if err := db.CreateTopicTable(params.ChainId, tableName, dynamicStruct); err != nil {
		return err
	}

	// 创建topic记录
	record := &db.ContractABITopic{
		Id:              uuid.New().String(),
		ContractName:    upgradeContract.ContractName,
		ContractAddr:    params.ContractAddr,
		ContractVersion: params.ContractVersion,
		TopicTableName:  tableName,
		Topic:           event.Name,
	}

	err := dbhandle.InsertContractABITopic(params.ChainId, record)
	return err
}

// 2. 保存ABI文件
// 保存ABI文件
func saveOrUpdateABIFile(params *entity.UploadContractAbiParams, upgradeContract *db.UpgradeContractTransaction,
	parsedABI []*utils.ContractABI) (bool, error) {
	// 将解析后的ABI序列化为JSON格式
	abiJSON, err := json.Marshal(parsedABI)
	if err != nil {
		log.Errorf("saveABIFile ABI序列化失败: %v", err)
		return false, nil
	}

	// 查询已存在的ABI文件
	existingABIFile, err := dbhandle.GetContractABIFile(params.ChainId, params.ContractAddr, params.ContractVersion)
	if err != nil {
		log.Errorf("GetContractABIFile ABI文件查询失败: %v", err)
		return false, err
	}

	// 如果不存在已存在的ABI文件，则插入新的ABI文件
	if existingABIFile == nil {
		newABIFile := &db.ContractABIFile{
			Id:              uuid.New().String(),
			ContractName:    upgradeContract.ContractName,
			ContractAddr:    params.ContractAddr,
			ContractVersion: params.ContractVersion,
			ABIJson:         string(abiJSON),
		}
		errDB := dbhandle.InsertContractABIFile(params.ChainId, newABIFile)
		return true, errDB
	}

	// 如果已存在已存在的ABI文件，则更新ABI文件
	existingABIFile.ABIJson = string(abiJSON)
	errDB := dbhandle.UpdateContractABIFile(params.ChainId, existingABIFile)
	return false, errDB
}

// 获取合约主题
func GetContractTopics(chainId, contractAddr, version string) ([]string, error) {
	contractTopics := make([]string, 0)
	abiFuncList, err := GetContractABIUnmarshal(chainId, contractAddr, version)
	if err != nil {
		log.Errorf("GetContractABIUnmarshal err: %v", err)
		return nil, err
	}

	// 遍历abi函数列表，获取事件名称
	for _, abi := range abiFuncList {
		if abi.Type == utils.ABIEventType {
			contractTopics = append(contractTopics, abi.Name)
		}
	}
	return contractTopics, nil
}

func GetContractTopicIndexs(chainId, contractAddr, version, topic string) ([]string, error) {
	topicIndexs := make([]string, 0)
	// 获取合约abi详情
	abiFuncList, err := GetContractABIUnmarshal(chainId, contractAddr, version)
	if err != nil {
		log.Errorf("GetContractABIUnmarshal err: %v", err)
		return nil, err
	}
	// 遍历abi函数列表，获取事件名称
	for _, abi := range abiFuncList {
		if abi.Type == utils.ABIEventType && abi.Name == topic {
			for _, input := range abi.Inputs {
				if input.Indexed {
					topicIndexs = append(topicIndexs, input.Name)
				}
			}
		}
	}

	return topicIndexs, nil
}

func GetContractTopicColumns(chainId, contractAddr, version, topic string) ([]string, error) {
	topicColumns := make([]string, 0)
	// 获取合约abi详情
	abiFuncList, err := GetContractABIUnmarshal(chainId, contractAddr, version)
	if err != nil {
		log.Errorf("GetContractABIUnmarshal err: %v", err)
		return nil, err
	}
	// 遍历abi函数列表，获取事件名称
	for _, abi := range abiFuncList {
		if abi.Type == utils.ABIEventType && abi.Name == topic {
			for _, input := range abi.Inputs {
				topicColumns = append(topicColumns, input.Name)
			}
		}
	}

	//拼接系统字段
	systemFields := []string{
		db.ABISystemFieldTxID,
		db.ABISystemFieldTimestamp,
	}
	topicColumns = append(topicColumns, systemFields...)
	return topicColumns, nil
}

// 获取合约abi解析的结构体
func GetContractABIUnmarshal(chainId, contractAddr, version string) ([]*utils.ContractABI, error) {
	// 获取合约abi详情
	abiDetail, err := dbhandle.GetContractABIFile(chainId, contractAddr, version)
	if err != nil {
		return nil, err
	}

	if abiDetail == nil {
		return nil, err
	}

	// 解析abi
	abiFuncList := make([]*utils.ContractABI, 0)
	err = json.Unmarshal([]byte(abiDetail.ABIJson), &abiFuncList)
	if err != nil {
		return nil, err
	}
	return abiFuncList, nil
}

func GetDecodeContractEvents(params *entity.GetDecodeContractEventsParams, topicColumns, topicIndex []string) (
	[]map[string]interface{}, int64, error) {
	var (
		offset       = params.Offset
		limit        = params.Limit
		chainId      = params.ChainId
		contractAddr = params.ContractAddr
		version      = params.ContractVersion
		topic        = params.Topic
		searchParams = params.SearchParams
	)
	// 查询当前版本是否已有记录
	abiTopicInfo, err := dbhandle.GetContractABITopic(chainId, contractAddr, version, topic)
	if err != nil {
		return nil, 0, err
	}

	var tableName string
	if abiTopicInfo != nil {
		tableName = abiTopicInfo.TopicTableName
	}

	filterSearchParams := filterSearchParams(searchParams, topicIndex)
	results, total, err := db.DBHandler.GetDecodeEventByABIAndTotal(offset, limit, chainId, contractAddr,
		version, topic, tableName, topicColumns, filterSearchParams)
	return results, total, err
}

// 过滤掉不在topicIndexs里面的searchParams数据
func filterSearchParams(searchParams []entity.SearchParam, topicIndexs []string) []entity.SearchParam {
	// 创建索引集合用于快速查找
	indexSet := make(map[string]bool, len(topicIndexs))
	for _, index := range topicIndexs {
		indexSet[index] = true
	}

	// 过滤符合条件的参数
	filtered := make([]entity.SearchParam, 0, len(searchParams))
	for _, param := range searchParams {
		if indexSet[param.Name] && param.Value != "" {
			filtered = append(filtered, param)
		}
	}
	return filtered
}

func AsyncHandleContractABIEvent(chainId, contractAddr, contractVersion string) error {
	var (
		offset       int
		totalSuccess int64
		totalFail    int64
	)
	//循环获取合约ABI事件，直到数据为空,统计一共处理了多少条数据
	for {
		// 获取合约ABI事件
		events, err := dbhandle.GetABIEventsByVersion(offset, entity.LimitMaxSpec, chainId,
			contractAddr, contractVersion)
		if err != nil || len(events) == 0 {
			break
		}

		// 处理合约ABI事件
		topicTableEvents := syncLogic.BuildEventDataByABI(chainId, events)
		// 遍历topicTableEvents中的每一个tableName和events
		for tableName, values := range topicTableEvents {
			// 调用dbhandle.InsertDecodeEventByABI函数将events插入到tableName表中
			err := dbhandle.InsertDecodeEventByABI(chainId, tableName, values)
			// 如果插入失败，则返回错误
			if err != nil {
				log.Errorf("InsertDecodeEventByABI failed, err: %v", err)
				totalFail += int64(len(values))
				continue
			}
			totalSuccess += int64(len(values))
		}

		offset += entity.LimitMaxSpec
		if len(events) < entity.LimitMaxSpec {
			break
		}
	}

	log.Infof("AsyncHandleContractABIEvent 处理合约ABI事件完成 | chainId=%s | contractAddr=%s | "+
		"contractVersion=%s | totalSuccess=%d | totalFail=%d",
		chainId, contractAddr, contractVersion, totalSuccess, totalFail)
	return nil
}
