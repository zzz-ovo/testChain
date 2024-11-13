/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/common/v2/evmutils/abi"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type CrossEvmContract struct {
}

func (c *CrossEvmContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (c *CrossEvmContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (c *CrossEvmContract) InvokeContract(method string) protogo.Response {

	switch method {
	case "crossEvmStorageSet":
		return c.CrossEvmStorageSet()
	case "crossEvmStorageGet":
		return c.CrossEvmStorageGet()
	default:
		return sdk.Error("invalid method")
	}
}

func (c *CrossEvmContract) CrossEvmStorageSet() protogo.Response {
	// 获取所有参数
	args := sdk.Instance.GetArgs()

	// 取出 storage abi
	storateABIJson := string(args["storage_abi"])
	storageAbi, err := abi.JSON(strings.NewReader(string(storateABIJson)))
	if err != nil {
		return sdk.Error(err.Error())
	}

	// 取出 data的10进制字符串表示
	dataStr := string(args["storage_set_data"])
	// dataStr 转换成 int64
	data, err := strconv.ParseInt(dataStr, 10, 64)
	if err != nil {
		return sdk.Error(err.Error())
	}

	// 取出跨合约调用的方法
	evmStorageFuncName := string(args["storage_set_func_name"])
	dataByte, err := storageAbi.Pack(evmStorageFuncName, data)
	if err != nil {
		return sdk.Error(err.Error())
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	// 取出 evm storage 合约名
	storageContractName := string(args["storage_contract_name"])

	crossContractArgs := make(map[string][]byte)
	crossContractArgs["data"] = []byte(dataString)

	// response could be correct or error
	response := sdk.Instance.CallContract(storageContractName, method, crossContractArgs)
	sdk.Instance.EmitEvent("cross contract set", []string{"success"})
	return response
}

func (c *CrossEvmContract) CrossEvmStorageGet() protogo.Response {
	// 获取所有参数
	args := sdk.Instance.GetArgs()

	// 取出 storage abi
	storateABIJson := string(args["storage_abi"])
	storageAbi, err := abi.JSON(strings.NewReader(string(storateABIJson)))
	if err != nil {
		return sdk.Error(err.Error())
	}

	// 取出跨合约调用的方法
	evmStorageFuncName := string(args["storage_get_func_name"])

	dataByte, err := storageAbi.Pack(evmStorageFuncName)
	// dataByte, err := storageAbi.Pack("get")
	if err != nil {
		if err.Error() != "contract does not have a constructor" {
			return sdk.Error(err.Error())
		}
	}

	dataString := hex.EncodeToString(dataByte)

	method := dataString[0:8]

	// 取出 evm storage 合约名
	storageContractName := string(args["storage_contract_name"])

	crossContractArgs := make(map[string][]byte)
	crossContractArgs["data"] = []byte(dataString)

	// response could be correct or error
	response := sdk.Instance.CallContract(storageContractName, method, crossContractArgs)

	val, err := storageAbi.Unpack(evmStorageFuncName, response.Payload)
	if err != nil {
		return sdk.Error(err.Error())
	}

	sdk.Instance.EmitEvent("cross contract get", []string{fmt.Sprintf("%s", val)})
	// return response
	return sdk.Success([]byte(fmt.Sprintf("get value from evm: %s", val)))
}

func main() {
	err := sandbox.Start(new(CrossEvmContract))
	if err != nil {
		log.Fatal(err)
	}
}
