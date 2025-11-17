package saveTasks

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"sync"

	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	"github.com/panjf2000/ants/v2"
)

// TaskSaveTransactions
type TaskSaveTransactions struct {
	ChainId            string
	Transactions       map[string]*db.Transaction
	UpgradeContractTxs []*db.UpgradeContractTransaction
}

// NewTaskSaveTransactions
// @param chainId
// @param transactions
// @return *TaskSaveTransactions
func (e TaskSaveTransactions) Execute() error {
	err := SaveTransactionsToDB(e.ChainId, e.Transactions)
	if err != nil {
		return err
	}
	err = SaveUpgradeContractTxToDB(e.ChainId, e.UpgradeContractTxs)
	return err
}

// SaveTransactionsToDB
//
//	@Description: 存储交易数据
//	@param chainId
//	@param transactions 交易列表
//	@return error
func SaveTransactionsToDB(chainId string, transactions map[string]*db.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}
	var (
		goRoutinePool *ants.Pool
		err           error
	)

	// 将交易分割为大小为50的批次
	batches := batchTransactions(transactions)
	errChan := make(chan error, config.MaxDBPoolSize)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.Transaction) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertTransactions(chainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertTransactions submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// SaveUpgradeContractTxToDB
//
//	@Description: 保存合约交易信息
//	@param chainId
//	@param transactions 交易map
//	@param contractTxId 合约升级交易id
//	@return error
func SaveUpgradeContractTxToDB(chainId string, upgradeContractTxs []*db.UpgradeContractTransaction) error {
	if len(upgradeContractTxs) == 0 {
		return nil
	}

	//插入数据
	err := dbhandle.InsertUpgradeContractTx(chainId, upgradeContractTxs)
	if err != nil {
		return err
	}

	return nil
}

// TaskInsertUser
type TaskInsertUser struct {
	ChainId  string
	UserList map[string]*db.User
}

// Execute
// @description: 执行任务-插入用户数据
// @return error 错误信息
func (e TaskInsertUser) Execute() error {
	users := e.UserList
	if len(users) == 0 {
		return nil
	}
	// 将分割为大小为100的批次
	batches := batchUsers(users)
	for _, batch := range batches {
		//插入数据
		err := dbhandle.BatchInsertUser(e.ChainId, batch)
		if err != nil {
			return err
		}
	}
	return nil
}

// TaskSaveContract
type TaskSaveContract struct {
	ChainId      string
	InsertList   []*db.Contract
	UpdateList   []*db.Contract
	ByteCodeList []*db.ContractByteCode
}

// Execute
// @description: 执行任务-保存合约数据
// @return error 错误信息
func (e TaskSaveContract) Execute() error {
	insertContracts := e.InsertList
	updateContracts := e.UpdateList
	byteCodeList := e.ByteCodeList
	for _, contract := range insertContracts {
		//insert
		err := dbhandle.InsertContract(e.ChainId, contract)
		if err != nil {
			log.Errorf("TaskSaveContract InsertContract failed, contract:%v", contract)
			return err
		}
	}

	for _, contract := range updateContracts {
		//update
		err := dbhandle.UpdateContract(e.ChainId, contract)
		if err != nil {
			log.Errorf("TaskSaveContract UpdateContract failed, contract:%v", contract)
			return err
		}
	}

	//存储合约字节码
	err := dbhandle.InsertContractByteCodes(e.ChainId, byteCodeList)
	if err != nil {
		log.Errorf("TaskSaveContract InsertContractByteCodes failed")
		return err
	}
	return nil
}

// TaskSaveStandardContract
type TaskSaveStandardContract struct {
	ChainId            string
	InsertFTContracts  []*db.FungibleContract
	InsertNFTContracts []*db.NonFungibleContract
	InsertIDAContracts []*db.IDAContract
}

// Execute
// @description: 执行任务-保存合约数据
// @return error 错误信息
func (e TaskSaveStandardContract) Execute() error {
	ftContracts := e.InsertFTContracts
	nftContracts := e.InsertNFTContracts
	idaContracts := e.InsertIDAContracts
	var err error
	err = dbhandle.InsertFungibleContract(e.ChainId, ftContracts)
	if err != nil {
		log.Errorf("SaveStandardContractToDB err:%v contract:%v", err, ftContracts)
		return err
	}

	err = dbhandle.InsertNonFungibleContract(e.ChainId, nftContracts)
	if err != nil {
		log.Errorf("SaveStandardContractToDB err:%v contract:%v", err, nftContracts)
		return err
	}

	err = dbhandle.InsertIDAContract(e.ChainId, idaContracts)
	if err != nil {
		log.Errorf("SaveStandardContractToDB err:%v contract:%v", err, idaContracts)
		return err
	}
	return nil
}

// TaskEvidenceContract
type TaskEvidenceContract struct {
	ChainId           string
	EvidenceContracts []*db.EvidenceContract
}

// Execute
// @description: 执行任务-保存合约数据
// @return error 错误信息
func (e TaskEvidenceContract) Execute() error {
	contractList := e.EvidenceContracts
	if len(contractList) == 0 {
		return nil
	}

	//插入数据
	err := dbhandle.InsertEvidenceContract(e.ChainId, contractList)
	if err != nil {
		log.Errorf("SaveEvidenceContractToDB failed, contract:%v", contractList)
		return err
	}
	return nil
}

// TaskInsertContractEvents
type TaskInsertContractEvents struct {
	ChainId        string
	ContractEvents []*db.ContractEvent
}

// Execute
// @description: 执行任务-保存合约数据
// @return error 错误信息
func (e TaskInsertContractEvents) Execute() error {
	contractEvents := e.ContractEvents
	chainId := e.ChainId

	if len(contractEvents) == 0 {
		return nil
	}
	var (
		goRoutinePool *ants.Pool
		err           error
		wg            sync.WaitGroup
		errChan       = make(chan error, config.MaxDBPoolSize)
	)

	// 将交易分割为大小为SaveBatchSize的批次
	batches := batchContractEvents(contractEvents)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.ContractEvent) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertContractEvent(chainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertContractEvent submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// TaskInsertGasRecord
type TaskInsertGasRecord struct {
	ChainId          string
	InsertGasRecords []*db.GasRecord
}

// Execute
// @description: 执行任务-保存合约数据
// @return error 错误信息
// SaveGasRecordToDB 保存gasR
func (e TaskInsertGasRecord) Execute() error {
	if len(e.InsertGasRecords) == 0 {
		return nil
	}
	var (
		goRoutinePool *ants.Pool
		err           error
		wg            sync.WaitGroup
		errChan       = make(chan error, config.MaxDBPoolSize)
	)

	// 将交易分割为大小为10的批次
	batches := batchGasRecords(e.InsertGasRecords)
	if goRoutinePool, err = ants.NewPool(config.MaxDBPoolSize, ants.WithPreAlloc(false)); err != nil {
		log.Error(GoRoutinePoolErr + err.Error())
		return err
	}

	defer goRoutinePool.Release()
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(insertList []*db.GasRecord) func() {
			return func() {
				defer wg.Done()
				//插入数据
				err = dbhandle.InsertGasRecord(e.ChainId, insertList)
				if err != nil {
					errChan <- err
				}

			}
		}(batch))
		if errSub != nil {
			log.Error("InsertGasRecord submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	return nil
}

// TaskSaveChainConfig
type TaskSaveChainConfig struct {
	ChainId            string
	UpdateChainConfigs []*pbConfig.ChainConfig
}

// Execute
// @description: 执行任务-保存合约数据
// @return error 错误信息
// SaveGasRecordToDB 保存gasR
func (e TaskSaveChainConfig) Execute() error {
	if len(e.UpdateChainConfigs) == 0 {
		return nil
	}

	for _, chainConfig := range e.UpdateChainConfigs {
		if chainConfig == nil || chainConfig.ChainId == "" {
			continue
		}

		err := dbhandle.UpdateChainInfoByConfig(e.ChainId, chainConfig)
		if err != nil {
			log.Errorf("TaskSaveChainConfig failed chainConfig:%v ", chainConfig)
			return err
		}
	}

	return nil
}

// TaskSaveContract
type TaskSaveContractCrossCallTxs struct {
	ChainId        string
	InsertCrossTxs []*db.ContractCrossCallTransaction
	InsertCalls    map[string]*db.ContractCrossCall
}

// Execute
// @description: 执行任务-保存合约数据
// @return error 错误信息
func (e TaskSaveContractCrossCallTxs) Execute() error {
	if len(e.InsertCrossTxs) == 0 {
		return nil
	}

	batches := batchContractCallTxs(e.InsertCrossTxs)
	for _, batch := range batches {
		//插入数据
		err := dbhandle.BatchInsertContractCrossCallTxs(e.ChainId, batch)
		if err != nil {
			log.Errorf("BatchInsertContractCrossCallTxs failed, InsertCrossTxs:%v", batch)
			return err
		}
	}

	if len(e.InsertCalls) > 0 {
		insertList := make([]*db.ContractCrossCall, 0, len(e.InsertCalls))
		for _, call := range e.InsertCalls {
			insertList = append(insertList, call)
		}

		//插入数据
		err := dbhandle.BatchInsertContractCrossCalls(e.ChainId, insertList)
		if err != nil {
			log.Errorf("BatchInsertContractCrossCalls failed, InsertCalls:%v", insertList)
			return err
		}
	}
	return nil
}
