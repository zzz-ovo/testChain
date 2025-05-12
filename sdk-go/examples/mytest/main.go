package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	//"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	createContractTimeout = 10000
	claimVersion          = "1.0.0"
	claimByteCodePath     = "bank.7z"
	sdkConfPath           = "sdk_config_solo.yml"
)

var (
	// bankContractName = "bank" + strconv.FormatInt(time.Now().UnixNano(), 10)
	bankContractName = "bank" // 合约名称
	DataSkew         = 1.00001
	AccountCount     = uint64(TxCount)
	TxCount, _       = strconv.Atoi(os.Args[2])
)

type Transaction struct {
	from string
	to   string
}

func main() {
	// 初始化链客户端（需要配置）
	client, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
	)

	if err != nil {
		panic(err)
	}
	switch os.Args[1] {
	case "test":
		easyTest(client)
	case "prepare":
		prepare(client)
	case "benchmark":
		benchmark(client)
	case "line":
		lineTx(client)
	}

}

func prepare(client *sdk.ChainClient) {
	startime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < TxCount; i++ {
		wg.Add(1)
		go func(index uint32) {
			defer wg.Done()
			account := fmt.Sprintf("%09d", index)
			err := testDeposit(client, account, 10000000, false)
			if err != nil {
				fmt.Printf("transfer failed: %v", err)
			}
		}(uint32(i))
	}
	wg.Wait()
	timecost := time.Since(startime).Seconds()
	fmt.Printf("\n %v DepositTx cost %v seconds \n tps is %v", TxCount, timecost, float64(TxCount)/timecost)
}
func benchmark(client *sdk.ChainClient) {
	r := rand.New(rand.NewSource(1))
	zipf := rand.NewZipf(r, DataSkew, 1.0, AccountCount)
	txSet := make([]Transaction, TxCount)
	for i := 0; i < TxCount; i++ {
		from := zipf.Uint64()
		to := from
		for to == from {
			to = zipf.Uint64()
		}
		toAccount := fmt.Sprintf("%09d", to)
		fromAccount := fmt.Sprintf("%09d", from)
		txSet[i] = Transaction{from: fromAccount, to: toAccount}
	}
	startime := time.Now()
	var wg sync.WaitGroup
	//failTxCount := uint32(0)
	for _, tx := range txSet {
		wg.Add(1)
		go func(tx Transaction) {
			defer wg.Done()
			err := testTransfer(client, tx.from, tx.to, 1, false)
			if err != nil {
				fmt.Printf("transfer failed: %v", err)
				return
			}

		}(tx)
	}
	wg.Wait()
	timecost := time.Since(startime).Seconds()
	fmt.Printf("\n %v TransferTx cost %v seconds \n tps is %v", TxCount, timecost, float64(TxCount)/timecost)

}

// func easyTest(client *sdk.ChainClient) {
// 	//fmt.Println("\n====================== 创建合约 ======================")
// 	//usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
// 	//testUserContractClaimCreate(client, true, usernames...)

// 	fmt.Println("\n====================== 测试存款 ======================")
// 	if err := testDeposit(client, "alice", 1000, true); err != nil {
// 		fmt.Printf("Deposit test failed: %v\n", err)
// 	}

// 	fmt.Println("\n====================== 测试转账 ======================")
// 	if err := testTransfer(client, "alice", "bob", 500, true); err != nil {
// 		fmt.Printf("Transfer test failed: %v\n", err)
// 	}

// 	fmt.Println("\n====================== 测试查询 ======================")
// 	if err := testGetBalance(client, "alice", true); err != nil {
// 		fmt.Printf("GetBalance test failed: %v\n", err)
// 	}

// 	fmt.Println("\n====================== 创建取款 ======================")
// 	if err := testWithdraw(client, "alice", 300, true); err != nil {
// 		fmt.Printf("Withdraw test failed: %v\n", err)
// 	}
// 	fmt.Println("\n====================== 测试结束 ======================\n")
// }

func easyTest(client *sdk.ChainClient) {
	//fmt.Println("\n====================== 创建合约 ======================")
	//usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	//testUserContractClaimCreate(client, true, usernames...)
	fmt.Println("\n====================== 测试查询 ======================")
	if err := testGetBalance(client, "alice", true); err != nil {
		fmt.Printf("GetBalance test failed: %v\n", err)
	}
	fmt.Println("\n====================== 测试查询 ======================")
	if err := testGetBalance(client, "bob", true); err != nil {
		fmt.Printf("GetBalance test failed: %v\n", err)
	}
	fmt.Println("\n====================== 测试存款 ======================")
	if err := testDeposit(client, "alice", 1000, true); err != nil {
		fmt.Printf("Deposit test failed: %v\n", err)
	}

	fmt.Println("\n====================== 测试转账 ======================")
	if err := testTransfer(client, "alice", "bob", 500, true); err != nil {
		fmt.Printf("Transfer test failed: %v\n", err)
	}

	fmt.Println("\n====================== 测试查询 ======================")
	if err := testGetBalance(client, "alice", true); err != nil {
		fmt.Printf("GetBalance test failed: %v\n", err)
	}

	fmt.Println("\n====================== 测试查询 ======================")
	if err := testGetBalance(client, "bob", true); err != nil {
		fmt.Printf("GetBalance test failed: %v\n", err)
	}
	fmt.Println("\n====================== 测试结束 ======================\n")
}

// 测试存款操作
func testDeposit(client *sdk.ChainClient, account string, amount int, withSyncResult bool) error {
	kvs := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("deposit"),
		},
		{
			Key:   "account",
			Value: []byte(account),
		},
		{
			Key:   "amount",
			Value: []byte(strconv.Itoa(amount)),
		},
	}

	_, err := invokeBankContract(client, bankContractName, "invoke_contract", "", kvs, withSyncResult)
	if err != nil {
		return fmt.Errorf("deposit failed: %v", err)
	}
	//fmt.Printf("Deposit %d to %s success", amount, account)
	return nil
}

// 测试取款操作
func testWithdraw(client *sdk.ChainClient, account string, amount int, withSyncResult bool) error {
	kvs := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("withdraw"),
		},
		{
			Key:   "account",
			Value: []byte(account),
		},
		{
			Key:   "amount",
			Value: []byte(strconv.Itoa(amount)),
		},
	}

	_, err := invokeBankContract(client, bankContractName, "invoke_contract", "", kvs, withSyncResult)
	if err != nil {
		return fmt.Errorf("withdraw failed: %v", err)
	}
	fmt.Printf("Withdraw %d from %s success", amount, account)
	return nil
}

// 测试转账操作
func testTransfer(client *sdk.ChainClient, from, to string, amount int, withSyncResult bool) error {
	kvs := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("transfer"),
		},
		{
			Key:   "from",
			Value: []byte(from),
		},
		{
			Key:   "to",
			Value: []byte(to),
		},
		{
			Key:   "amount",
			Value: []byte(strconv.Itoa(amount)),
		},
	}

	_, err := invokeBankContract(client, bankContractName, "invoke_contract", "", kvs, withSyncResult)
	if err != nil {
		return fmt.Errorf("transfer failed: %v", err)
	}
	//fmt.Printf("Transfer %d from %s to %s success", amount, from, to)
	return nil
}

// 测试查询余额
func testGetBalance(client *sdk.ChainClient, account string, withSyncResult bool) error {
	kvs := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("getBalance"),
		},
		{
			Key:   "account",
			Value: []byte(account),
		},
	}

	resp, err := invokeBankContract(client, bankContractName, "invoke_contract", "", kvs, withSyncResult)
	if err != nil {
		return fmt.Errorf("get balance failed: %v", err)
	}

	// 实际使用中需要解析合约返回结果
	// 这里假设直接返回余额数值
	balanceStr := string(resp.ContractResult.Result)
	balance, _ := strconv.Atoi(balanceStr)
	fmt.Printf("Balance of %s is %d", account, balance)
	return nil
}

// 通用合约调用方法
func invokeBankContract(client *sdk.ChainClient, contractName, method, txId string, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {
	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
	if err != nil {
		return resp, fmt.Errorf("invoke contract failed: %v", err)
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return resp, fmt.Errorf("contract response error: [code:%d]/[msg:%s]", resp.Code, resp.Message)
	}

	//fmt.Printf("Operation success! Contract result: %s\n", resp.ContractResult.Result)
	return resp, nil
}

func testUserContractClaimCreate(client *sdk.ChainClient, withSyncResult bool, usernames ...string) {

	resp, err := createUserContract(client, bankContractName, claimVersion, claimByteCodePath,
		common.RuntimeType_DOCKER_GO, []*common.KeyValuePair{}, withSyncResult, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CREATE claim contract resp: %+v\n", resp)
}
func createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string, runtime common.RuntimeType,
	kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	endorsers, err := examples.GetEndorsersWithAuthType(client.GetHashType(),
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func lineTx(client *sdk.ChainClient) {

	for i := 0; i < TxCount; i++ {
		if err := testDeposit(client, "alice", 1, false); err != nil {
			fmt.Printf("Deposit test failed: %v\n", err)
		}
	}
}
