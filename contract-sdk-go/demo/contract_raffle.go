/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"encoding/json"
	"fmt"
	"log"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type RaffleContract struct {
}

func (r *RaffleContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (r *RaffleContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (r *RaffleContract) InvokeContract(method string) protogo.Response {

	switch method {
	case "register_all":
		return r.RegisterAll()
	case "raffle":
		return r.Raffle()
	case "query":
		return r.query()
	default:
		return sdk.Error("invalid method")
	}
}

//type Peoples struct {
//	Peoples map[string]int `json:"peoples"`
//}

type Peoples struct {
	Peoples map[int]string `json:"peoples"`
}

func (r *RaffleContract) RegisterAll() protogo.Response {
	params := sdk.Instance.GetArgs()

	// 获取参数
	value := params["peoples"]
	var errMsg string
	if len(value) == 0 {
		errMsg = "value should not be empty!"
		sdk.Instance.Debugf(errMsg)
		return sdk.Error(errMsg)
	}

	var peoples Peoples
	err := json.Unmarshal(value, &peoples)
	for i := 1; i < len(peoples.Peoples); i++ {
		if name, ok := peoples.Peoples[i]; !ok || len(name) == 0 {
			errMsg = fmt.Sprintf("[registerAll] name should not be empty for number %d", i)
			sdk.Instance.Debugf(errMsg)
			return sdk.Error(errMsg)
		}
	}
	err = sdk.Instance.PutStateByte("peoples", "", value)
	if err != nil {
		errMsg = fmt.Sprintf("[register] put state bytes failed, %s", err)
		sdk.Instance.Debugf(errMsg)
		return sdk.Error(errMsg)
	}
	// 返回结果
	return sdk.Success([]byte("ok"))
}

func (r *RaffleContract) Raffle() protogo.Response {
	params := sdk.Instance.GetArgs()

	var errMsg string
	argTimestamp := string(params["timestamp"])
	if len(argTimestamp) == 0 {
		errMsg = "argTimestamp should not be empty!"
		sdk.Instance.Debugf(errMsg)
		return sdk.Error(errMsg)
	}

	peoplesData, err := sdk.Instance.GetStateByte("peoples", "")
	if err != nil {
		errMsg = "get peoples data from store failed!"
		sdk.Instance.Debugf(errMsg)
		return sdk.Error(errMsg)
	}

	var peoples Peoples
	err = json.Unmarshal(peoplesData, &peoples)
	if err != nil {
		errMsg = fmt.Sprintf("unmarshal peoples data failed, %s", err)
		sdk.Instance.Debugf(errMsg)
		return sdk.Error(errMsg)
	}
	num := r.BkdrHash(argTimestamp)
	num = num % len(peoples.Peoples)

	result := fmt.Sprintf("num: %d, name: %s", num, peoples.Peoples[num])
	delete(peoples.Peoples, num)
	newPeoplesData, err := json.Marshal(peoples.Peoples)
	if err != nil {
		errMsg = fmt.Sprintf("marshal new peoples data failed, %s", err)
		sdk.Instance.Debugf(errMsg)
		return sdk.Error(errMsg)
	}
	err = sdk.Instance.PutStateByte("peoples", "", newPeoplesData)
	if err != nil {
		errMsg = fmt.Sprintf("put new peoples data failed, %s", err)
		sdk.Instance.Debugf(errMsg)
		return sdk.Error(errMsg)
	}

	return sdk.Success([]byte(result))
}

func (r *RaffleContract) BkdrHash(timestamp string) int {
	hash := 0
	seed := 131
	for x := range timestamp {
		hash = hash*seed + x
	}
	return hash & 0x7FFFFFFF
}

func (r *RaffleContract) query() protogo.Response {
	peoplesData, err := sdk.Instance.GetStateByte("peoples", "")
	if err != nil {
		errMsg := "get peoples data from store failed!"
		sdk.Instance.Debugf(errMsg)
		return sdk.Error(errMsg)
	}

	return sdk.Success(peoplesData)
}

//func (f *RaffleContract) register() protogo.Response {
//	params := f.Sdk.GetArgs()
//
//	// 获取参数
//	name := string(params["name"])
//	if len(name) == 0 {
//		msg := "name should not be empty!"
//		f.Sdk.Debugf(msg)
//		return sdk.Error(msg)
//	}
//
//	index := 0
//	var errMsg string
//	result, err := f.Sdk.GetState("index", "")
//	if err != nil {
//		errMsg = fmt.Sprintf("[register] get index from store failed, %s", err)
//		f.Sdk.Debugf(errMsg)
//		return sdk.Error(errMsg)
//	}
//	if len(result) != 0 {
//		index, err = strconv.Atoi(result)
//	}
//	if err != nil {
//		errMsg = fmt.Sprintf("[register] convert index failed, %s", err)
//		f.Sdk.Debugf(errMsg)
//		return sdk.Error(errMsg)
//	}
//	if index == 0 {
//		index = 100
//	} else {
//		index++
//	}
//	resultBytes, err := f.Sdk.GetStateByte("peoples", "")
//	if err != nil {
//		errMsg = fmt.Sprintf("[register] get peoples data failed, %s", err)
//		f.Sdk.Debugf(errMsg)
//		return sdk.Error(errMsg)
//	}
//	var peoples Peoples
//	if err = json.Unmarshal(resultBytes, &peoples); err != nil {
//		errMsg = fmt.Sprintf("[register] unmarshal peoples failed, %s", err)
//		f.Sdk.Debugf(errMsg)
//		return sdk.Error(errMsg)
//	}
//	_, ok := peoples.Peoples[name]
//	if ok {
//		errMsg = fmt.Sprintf("[register] %s has already been register", name)
//		f.Sdk.Debugf(errMsg)
//		return sdk.Error(errMsg)
//	}
//	peoples.Peoples[name] = index
//	peoplesBytes, err := json.Marshal(peoples)
//	if err != nil {
//		errMsg = fmt.Sprintf("[register] mashal peoples failed, %s", err)
//		f.Sdk.Debugf(errMsg)
//		return sdk.Error(errMsg)
//	}
//	err = f.Sdk.PutStateByte("peoples", "", peoplesBytes)
//	if err != nil {
//		errMsg = fmt.Sprintf("[register] put state bytes failed, %s", err)
//		f.Sdk.Debugf(errMsg)
//		return sdk.Error(errMsg)
//	}
//	indexStr := strconv.Itoa(index)
//	err = f.Sdk.PutState("index", "", indexStr)
//	if err != nil {
//		errMsg = fmt.Sprintf("[register] put state bytes failed, %s", err)
//		f.Sdk.Debugf(errMsg)
//		return sdk.Error(errMsg)
//	}
//
//	// 返回结果
//	return sdk.Success([]byte(indexStr))
//}

func main() {

	err := sandbox.Start(new(RaffleContract))
	if err != nil {
		log.Fatal(err)
	}
}
