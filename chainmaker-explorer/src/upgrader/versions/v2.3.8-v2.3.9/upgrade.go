package v238_239

import (
	"chainmaker_web/src/chain"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/logic"
	"chainmaker_web/src/sync/model"
	"chainmaker_web/src/sync/saveTasks"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"

	"log"
)

// 实际升级逻辑
func Upgrade(chainId string) {
	startTime := time.Now()
	//删除表
	err := DeleteTablesByChainID(chainId)
	if err != nil {
		log.Printf("======fail %s====", chainId)
		return
	}

	//创建表
	chainList := make([]*config.ChainInfo, 0)
	chainList = append(chainList, &config.ChainInfo{ChainId: chainId})
	// 初始化链列表表
	chain.InitChainListTable(chainList)

	//最新高度
	blockHeight := GetMaxBlockHeight(chainId)
	log.Printf("Max Block Height for chain %s: %d", chainId, blockHeight)

	//处理数据
	err = BatchFetchTxInfo(chainId, blockHeight)
	endTime := time.Now()
	if err != nil {
		log.Printf("======fail %s====blockHeight %d======duration %v====",
			chainId, blockHeight, endTime.Sub(startTime))
		return
	}
	log.Printf("======success %s====blockHeight %d======duration %v====",
		chainId, blockHeight, endTime.Sub(startTime))

	// 更新统计信息// 更新统计信息
	totalCount, err := dbhandle.GetCrossTxCount(chainId)
	if err != nil {
		log.Printf("======fail %s====TotalCrossTx %d=======",
			chainId, totalCount)
		return
	}
	statistics, err := dbhandle.GetChainStatistics(chainId)
	if err != nil {
		log.Printf("======fail %s====TotalCrossTx %d=======",
			chainId, totalCount)
		return
	}

	statistics.TotalCrossTx = totalCount
	err = dbhandle.UpdateStatisticsRealtime(chainId, statistics)
	if err != nil {
		log.Printf("======fail %s====TotalCrossTx %d=======",
			chainId, totalCount)
		return
	}
	log.Printf("======success %s====TotalCrossTx %d=======",
		chainId, totalCount)
}

func DeleteTablesByChainID(chainId string) error {
	// 根据链ID和表名获取表名
	tableName1 := db.GetTableName(chainId, db.TableCrossSubChainData)
	tableName2 := db.GetTableName(chainId, db.TableCrossMainTransaction)
	tableName3 := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	tableName4 := db.GetTableName(chainId, db.TableCrossBusinessTransaction)
	tableName5 := db.GetTableName(chainId, db.TableCrossSubChainCrossChain)
	tableName6 := db.GetTableName(chainId, db.TableCrossChainContract)
	// 执行删除表的SQL语句
	err := db.GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName1)).Error
	// 如果有错误，返回错误
	if err != nil {
		return err
	}
	err = db.GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName2)).Error
	// 如果有错误，返回错误
	if err != nil {
		return err
	}
	err = db.GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName3)).Error
	// 如果有错误，返回错误
	if err != nil {
		return err
	}
	err = db.GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName4)).Error
	// 如果有错误，返回错误
	if err != nil {
		return err
	}
	err = db.GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName5)).Error
	// 如果有错误，返回错误
	if err != nil {
		return err
	}
	err = db.GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName6)).Error
	// 如果有错误，返回错误
	if err != nil {
		return err
	}
	// 没有错误，返回nil
	return nil
}

func BatchFetchTxInfo(chainId string, maxBlockHeight int64) error {
	blockHeight := int64(0) // 明确初始化从高度0开始
	// 串行遍历每个区块（0 -> maxBlockHeight）
	for blockHeight <= int64(maxBlockHeight) {
		crossResult := model.NewCrossChainResult()
		// 1. 批量查询当前高度交易（支持单高度批量查询）
		txList, err := dbhandle.GetTxInfoByBlockHeight(chainId, []int64{blockHeight})
		if err != nil {
			log.Printf("区块 %d 交易查询失败: %v", blockHeight, err)
			blockHeight++ // 跳过当前高度继续执行
			continue
		}

		// 2. 逐笔交易处理
		for _, tx := range txList {
			if tx.WriteSet == "" {
				continue // 跳过无写集的交易
			}

			// 解析交易写集
			txWrites := make([]config.RwSet, 0)
			if err1 := json.Unmarshal([]byte(tx.WriteSet), &txWrites); err1 != nil {
				log.Printf("交易 %s 写集解析失败: %v", tx.TxId, err1)
				continue
			}

			// 提取跨链信息
			crossChainInfo, err2 := GetCrossChainInfoByWriteSet(txWrites)
			if err2 != nil || crossChainInfo == nil {
				continue // 非跨链交易或解析失败
			}

			// 3. 处理跨链数据（严格按顺序执行）
			crossTxTransferList := logic.GetCrossTxTransfer(chainId, tx.TxId, blockHeight, tx.Timestamp, crossChainInfo)
			mainTransaction := logic.GetMainCrossTransaction(blockHeight, crossChainInfo, tx.TxId, tx.Timestamp)
			businessTxMap := logic.GetBusinessTransaction(chainId, crossChainInfo)
			logic.UpdateCrossChainResultTx(crossResult, mainTransaction, crossTxTransferList, businessTxMap)
		}

		// 4. 按区块顺序执行校验（关键：不可并发！）
		if err1 := sync.CheckCrossChainTransactionData(chainId, crossResult); err1 != nil {
			log.Printf("区块 %d 跨链校验失败: %v", blockHeight, err1)
			return err1 // 严重错误时终止
		}

		err = SaveCrossChainResult(chainId, crossResult)
		if err != nil {
			log.Printf("区块 %d 跨链结果保存失败: %v", blockHeight, err)
			return err // 严重错误时终止
		}

		err = DelayCrossTransfer(chainId, crossResult, blockHeight)
		if err != nil {
			log.Printf("区块 %d 跨链结果保存失败: %v", blockHeight, err)
		}

		log.Printf("===========区块 %d 跨链处理成功=============", blockHeight)
		blockHeight++ // 移动到下一高度
	}
	return nil
}

func DelayCrossTransfer(chainId string, crossResult *model.CrossChainResult, blockHeight int64) error {
	//计算主子链跨链次数
	crossSubChainIdMap := logic.ParseCrossCycleTxTransfer(crossResult.InsertCrossTransfer)
	//主子链-获取子链信息
	var subChainIds []string
	for crossSubChainId := range crossSubChainIdMap {
		subChainIds = append(subChainIds, crossSubChainId)
	}
	crossSubChainDBMap, err := dbhandle.GetCrossSubChainById(chainId, subChainIds)
	if err != nil {
		log.Printf("GetCrossSubChainById error: %v", err)
		return err
	}
	//主子链-获取子链跨链列表交易数据
	crossSubChainCrossDB, err := dbhandle.GetCrossSubChainCrossNum(chainId, subChainIds)
	if err != nil {
		log.Printf("GetCrossSubChainCrossNum error: %v", err)
		return err
	}

	insertSubChainData, updateSubChainData, crossChainContracts := logic.DealCrossSubChainData(
		crossResult.InsertCrossTransfer, crossSubChainDBMap)

	insertSubChainCross, updateSubChainCross, err := logic.DealSubChainCrossChainNum(chainId, crossSubChainIdMap,
		crossSubChainCrossDB, blockHeight)
	if err != nil {
		log.Printf("DealSubChainCrossChainNum error: %v", err)
		return err
	}

	delayCrossChain := &model.DelayCrossChain{
		InsertSubChainData:  insertSubChainData,
		UpdateSubChainData:  updateSubChainData,
		InsertSubChainCross: insertSubChainCross,
		UpdateSubChainCross: updateSubChainCross,
		CrossChainContracts: crossChainContracts,
	}
	err = saveTasks.SaveDelayCrossChainToDB(chainId, delayCrossChain)
	if err != nil {
		log.Printf("SaveDelayCrossChainToDB error: %v", err)
		return err
	}
	return nil
}

func SaveCrossChainResult(chainId string, crossResult *model.CrossChainResult) error {
	err := saveTasks.UpdateCrossTransfer(chainId, crossResult)
	if err != nil {
		return err
	}

	//存储跨链交易
	err = dbhandle.InsertCrossSubTransaction(chainId, crossResult.CrossMainTransaction)
	if err != nil {
		return err
	}

	//业务交易数据
	err = saveTasks.SaveBusinessTransaction(chainId, crossResult.BusinessTxMap)
	if err != nil {
		return err
	}
	return nil
}

func GetCrossChainInfoByWriteSet(txWrites []config.RwSet) (
	*tcipCommon.CrossChainInfo, error) {
	crossChainInfo := &tcipCommon.CrossChainInfo{}
	if len(txWrites) == 0 {
		return nil, nil
	}

	for _, write := range txWrites {
		//以c开头的key可以解析成CrossChainInfo
		isCrossChainTx := common.IsRelayCrossChainInfo(write.Key, write.ContractName)
		if !isCrossChainTx {
			continue
		}

		// Base64 解码
		decoded, err := base64.StdEncoding.DecodeString(string(write.Value))
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(decoded, crossChainInfo)
		if err != nil {
			return nil, err
		}
		break
	}

	if crossChainInfo.CrossChainId != "" {
		return crossChainInfo, nil
	}
	return nil, nil
}

// 获取最大区块高度（复用主项目代码）
func GetMaxBlockHeight(chainId string) int64 {
	statistics, err := dbhandle.GetChainStatistics(chainId)
	if err != nil {
		return 0 // 错误处理
	}
	if statistics == nil {
		return 0 // 没有统计数据
	}
	return statistics.BlockHeight
}
