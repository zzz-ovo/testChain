/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	sdkConfigPKUser1Path = "../sdk_configs/sdk_config_pk_user1.yml"
	amount               = 900000000

	admin2KeyPath = "../../testdata/crypto-config-pk/public/admin/admin2/admin2.key"
	admin2PemPath = "../../testdata/crypto-config-pk/public/admin/admin2/admin2.pem"

	createContractTimeout = 5
	claimVersion          = "2.0.0"
	claimByteCodePath     = "../../testdata/claim-wasm-demo/rust-fact-2.0.0.wasm"
)

var (
	claimContractName = "claim" + strconv.FormatInt(time.Now().UnixNano(), 10)
)

/**
测试方法：
1.chainmaker-go中启动pk模式4节点区块链网络
	./prepare_pk.sh 4 1
	./build_release.sh
	./cluster_quick_start.sh normal
2.停止节点，开启gas
	修改4个节点by1.yml, enable_gas: true
	./cluster_quick_stop.sh clean
	./cluster_quick_start.sh normal
3.拷贝公私钥
	cp -rf chainmaker-go/build/crypto-config/node1/admin sdk-go/testdata/crypto-config-pk/public/admin
4.执行
	cd examples/payer
	go run main.go
*/

func main() {
	//测试对合约名+方法名设置Payer
	testContractAndMethod()
	//测试对合约名名设置Payer
	testContractNoMethod()
}
func testContractAndMethod() {
	methodSave := "save"
	claimContractName = "claim" + strconv.FormatInt(time.Now().UnixNano(), 10)

	client, err := examples.CreateChainClientWithSDKConf(sdkConfigPKUser1Path)
	if err != nil {
		log.Fatal(err)
	}
	pk, err := client.GetPublicKey().String()
	if err != nil {
		log.Fatal(err)
	}

	addr, err := sdk.GetEVMAddressFromPKPEM(pk, crypto.HASH_TYPE_SHA256)
	fmt.Printf("addr: %s\n", addr)

	fmt.Println("====================== 设置client自己为 gas admin ======================")
	if err := setGasAdmin(client, addr); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 获取 gas admin ======================")
	if err := getGasAdmin(client); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("====================== 充值gas账户%d个gas ======================\n", amount*2)
	rechargeGasList := []*syscontract.RechargeGas{
		{
			Address:   addr,
			GasAmount: amount * 2,
		},
	}
	if err := rechargeGas(client, rechargeGasList); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}

	//fmt.Println("====================== 查询admin2的地址 ======================")
	admin2Key, err := ioutil.ReadFile(admin2KeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	admin2Pem, err := ioutil.ReadFile(admin2PemPath)
	if err != nil {
		log.Fatalln(err)
	}
	addr2, err := sdk.GetEVMAddressFromPKPEM(string(admin2Pem), crypto.HASH_TYPE_SHA256)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("addr2: %s\n", addr2)

	fmt.Printf("====================== 充值gas账户%d个gas ======================\n", amount*2)
	rechargeGasList2 := []*syscontract.RechargeGas{
		{
			Address:   addr2,
			GasAmount: amount * 2,
		},
	}
	if err := rechargeGas(client, rechargeGasList2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 安装合约 ======================")
	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	testUserContractClaimCreate(client, true, usernames...)

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 调用合约 ======================")
	_, err = testUserContractClaimInvoke(client, "save", true, uuid.GetUUID())
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 设置payer ======================")
	resp, err := client.SetContractMethodPayer(addr2, claimContractName, methodSave, "00", "public", admin2Key, admin2Pem, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 查询payer ======================")
	resp, err = client.QueryContractMethodPayer(claimContractName, methodSave, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 调用合约 ======================")
	txId := uuid.GetUUID()
	_, err = testUserContractClaimInvoke(client, "save", true, txId)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("====================== 查询交易 txId=%s ======================\n", txId)
	resp, err = client.QueryTxPayer(txId, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 取消payer ======================")
	resp, err = client.UnsetContractMethodPayer(claimContractName, methodSave, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 查询payer ======================")
	resp, err = client.QueryContractMethodPayer(claimContractName, methodSave, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 调用合约 ======================")
	txId = uuid.GetUUID()
	_, err = testUserContractClaimInvoke(client, "save", true, txId)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("====================== 查询交易 txId=%s ======================\n", txId)
	resp, err = client.QueryTxPayer(txId, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}
}

func testContractNoMethod() {
	methodSave := ""
	claimContractName = "claim" + strconv.FormatInt(time.Now().UnixNano(), 10)

	client, err := examples.CreateChainClientWithSDKConf(sdkConfigPKUser1Path)
	if err != nil {
		log.Fatal(err)
	}
	pk, err := client.GetPublicKey().String()
	if err != nil {
		log.Fatal(err)
	}

	addr, err := sdk.GetEVMAddressFromPKPEM(pk, crypto.HASH_TYPE_SHA256)
	fmt.Printf("addr: %s\n", addr)

	fmt.Println("====================== 设置client自己为 gas admin ======================")
	if err := setGasAdmin(client, addr); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 获取 gas admin ======================")
	if err := getGasAdmin(client); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("====================== 充值gas账户%d个gas ======================\n", amount*2)
	rechargeGasList := []*syscontract.RechargeGas{
		{
			Address:   addr,
			GasAmount: amount * 2,
		},
	}
	if err := rechargeGas(client, rechargeGasList); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}

	//fmt.Println("====================== 查询admin2的地址 ======================")
	admin2Key, err := ioutil.ReadFile(admin2KeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	admin2Pem, err := ioutil.ReadFile(admin2PemPath)
	if err != nil {
		log.Fatalln(err)
	}
	addr2, err := sdk.GetEVMAddressFromPKPEM(string(admin2Pem), crypto.HASH_TYPE_SHA256)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("addr2: %s\n", addr2)

	fmt.Printf("====================== 充值gas账户%d个gas ======================\n", amount*2)
	rechargeGasList2 := []*syscontract.RechargeGas{
		{
			Address:   addr2,
			GasAmount: amount * 2,
		},
	}
	if err := rechargeGas(client, rechargeGasList2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 安装合约 ======================")
	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	testUserContractClaimCreate(client, true, usernames...)

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 调用合约 ======================")
	_, err = testUserContractClaimInvoke(client, "save", true, uuid.GetUUID())
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 设置payer ======================")
	resp, err := client.SetContractMethodPayer(addr2, claimContractName, methodSave, "00", "public", admin2Key, admin2Pem, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 查询payer ======================")
	resp, err = client.QueryContractMethodPayer(claimContractName, methodSave, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 调用合约 ======================")
	txId := uuid.GetUUID()
	_, err = testUserContractClaimInvoke(client, "save", true, txId)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("====================== 查询交易 txId=%s ======================\n", txId)
	resp, err = client.QueryTxPayer(txId, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 取消payer ======================")
	resp, err = client.UnsetContractMethodPayer(claimContractName, methodSave, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 查询payer ======================")
	resp, err = client.QueryContractMethodPayer(claimContractName, methodSave, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 调用合约 ======================")
	txId = uuid.GetUUID()
	_, err = testUserContractClaimInvoke(client, "save", true, txId)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("====================== 查询交易 txId=%s ======================\n", txId)
	resp, err = client.QueryTxPayer(txId, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, addr); err != nil {
		log.Fatal(err)
	}
	if err := getGasBalance(client, addr2); err != nil {
		log.Fatal(err)
	}
}

func setGasAdmin(cc *sdk.ChainClient, address string) error {
	payload, err := cc.CreateSetGasAdminPayload(address)
	if err != nil {
		return err
	}
	endorsers, err := examples.GetEndorsersWithAuthType(cc.GetHashType(),
		cc.GetAuthType(), payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, endorsers, -1, true)
	if err != nil {
		return err
	}

	fmt.Printf("set gas admin resp: %+v\n", resp)
	return nil
}

func getGasAdmin(cc *sdk.ChainClient) error {
	adminAddr, err := cc.GetGasAdmin()
	if err != nil {
		return err
	}
	fmt.Printf("get gas admin address: %+v\n", adminAddr)
	return nil
}

func rechargeGas(cc *sdk.ChainClient, rechargeGasList []*syscontract.RechargeGas) error {
	payload, err := cc.CreateRechargeGasPayload(rechargeGasList)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, nil, -1, true)
	if err != nil {
		return err
	}

	fmt.Printf("recharge gas resp: %+v\n", resp)
	return nil
}

func getGasBalance(cc *sdk.ChainClient, address string) error {
	balance, err := cc.GetGasBalance(address)
	if err != nil {
		return err
	}

	fmt.Printf("get gas balance: %+v\n", balance)
	return nil
}

func testUserContractClaimCreate(client *sdk.ChainClient, withSyncResult bool, usernames ...string) {

	resp, err := createUserContract(client, claimContractName, claimVersion, claimByteCodePath,
		common.RuntimeType_WASMER, []*common.KeyValuePair{}, withSyncResult, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CREATE claim contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string, runtime common.RuntimeType,
	kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

	var (
		err       error
		codeBytes []byte
	)

	//payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	//if err != nil {
	//	return nil, err
	//}
	isFile := utils.Exists(byteCodePath)
	if isFile {
		bz, err := ioutil.ReadFile(byteCodePath)
		if err != nil {
			return nil, fmt.Errorf("read from byteCode file %s failed, %s", byteCodePath, err)
		}

		if runtime == common.RuntimeType_EVM { // evm contract hex need decode to bytes
			codeBytesStr := strings.TrimSpace(string(bz))
			if codeBytes, err = hex.DecodeString(codeBytesStr); err != nil {
				return nil, fmt.Errorf("decode evm contract hex to bytes failed, %s", err)
			}
		} else { // wasm bin file no need decode
			codeBytes = bz
		}
	} else {
		if runtime == common.RuntimeType_EVM {
			byteCodePath = strings.TrimSpace(byteCodePath)
		}

		if codeBytes, err = hex.DecodeString(byteCodePath); err != nil {
			if codeBytes, err = base64.StdEncoding.DecodeString(byteCodePath); err != nil {
				return nil, fmt.Errorf("decode byteCode string failed, %s", err)
			}
		}
	}
	payload := client.CreatePayload("", common.TxType_INVOKE_CONTRACT,
		syscontract.SystemContract_CONTRACT_MANAGE.String(), syscontract.ContractManageFunction_INIT_CONTRACT.String(), kvs, 0, &common.Limit{GasLimit: amount})
	payload.Parameters = append(payload.Parameters, &common.KeyValuePair{
		Key:   syscontract.InitContract_CONTRACT_NAME.String(),
		Value: []byte(contractName),
	})

	payload.Parameters = append(payload.Parameters, &common.KeyValuePair{
		Key:   syscontract.InitContract_CONTRACT_VERSION.String(),
		Value: []byte(version),
	})

	payload.Parameters = append(payload.Parameters, &common.KeyValuePair{
		Key:   syscontract.InitContract_CONTRACT_RUNTIME_TYPE.String(),
		Value: []byte(runtime.String()),
	})

	payload.Parameters = append(payload.Parameters, &common.KeyValuePair{
		Key:   syscontract.InitContract_CONTRACT_BYTECODE.String(),
		Value: codeBytes,
	})

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
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

func testUserContractClaimInvoke(client *sdk.ChainClient,
	method string, withSyncResult bool, txId string) (string, error) {

	curTime := strconv.FormatInt(time.Now().Unix(), 10)

	fileHash := uuid.GetUUID()
	kvs := []*common.KeyValuePair{
		{
			Key:   "time",
			Value: []byte(curTime),
		},
		{
			Key:   "file_hash",
			Value: []byte(fileHash),
		},
		{
			Key:   "file_name",
			Value: []byte(fmt.Sprintf("file_%s", curTime)),
		},
	}

	err := invokeUserContract(client, claimContractName, method, txId, kvs, withSyncResult)
	if err != nil {
		return "", err
	}

	return fileHash, nil
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string,
	kvs []*common.KeyValuePair, withSyncResult bool) error {

	payload := client.CreatePayload(txId, common.TxType_INVOKE_CONTRACT, contractName, method, kvs, 0, &common.Limit{GasLimit: amount})

	resp, err := client.SendContractManageRequest(payload, nil, createContractTimeout, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}
