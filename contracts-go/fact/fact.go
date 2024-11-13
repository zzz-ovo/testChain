/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// FactContract contract for fact
type FactContract struct {
}

// Fact object to store
type Fact struct {
	FileHash string `json:"fileHash"`
	FileName string `json:"fileName"`
	Time     int32  `json:"time"`
}

// NewFact create a fact object
func NewFact(fileHash string, fileName string, time int32) *Fact {
	fact := &Fact{
		FileHash: fileHash,
		FileName: fileName,
		Time:     time,
	}
	return fact
}

// InitContract install contract func
func (f *FactContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrade contract func
func (f *FactContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// InvokeContract the entry func of invoke contract func
func (f *FactContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "save":
		return f.save()
	case "findByFileHash":
		return f.findByFileHash()
	default:
		return sdk.Error("invalid method")
	}
}

func (f *FactContract) save() protogo.Response {
	params := sdk.Instance.GetArgs()

	// 获取参数
	fileHash := string(params["file_hash"])
	fileName := string(params["file_name"])
	timeStr := string(params["time"])
	time, err := strconv.ParseInt(timeStr, 10, 0)
	if err != nil {
		msg := "time is [" + timeStr + "] not int"
		sdk.Instance.Errorf(msg)
		return sdk.Error(msg)
	}

	// 构建结构体
	fact := NewFact(fileHash, fileName, int32(time))

	// 序列化
	factBytes, err := json.Marshal(fact)
	if err != nil {
		return sdk.Error(fmt.Sprintf("marshal fact failed, err: %s", err))
	}
	// 发送事件
	sdk.Instance.EmitEvent("topic_vx", []string{fact.FileHash, fact.FileName})

	// 存储数据
	err = sdk.Instance.PutStateByte("fact_bytes", fact.FileHash, factBytes)
	if err != nil {
		return sdk.Error("fail to save fact bytes")
	}

	// 记录日志
	sdk.Instance.Infof("[save] fileHash=" + fact.FileHash)
	sdk.Instance.Infof("[save] fileName=" + fact.FileName)

	// 返回结果
	return sdk.Success([]byte(fact.FileName + fact.FileHash))

}

func (f *FactContract) findByFileHash() protogo.Response {
	// 获取参数
	fileHash := string(sdk.Instance.GetArgs()["file_hash"])

	// 查询结果
	result, err := sdk.Instance.GetStateByte("fact_bytes", fileHash)
	if err != nil {
		return sdk.Error("failed to call get_state")
	}

	// 反序列化
	var fact Fact
	if err = json.Unmarshal(result, &fact); err != nil {
		return sdk.Error(fmt.Sprintf("unmarshal fact failed, err: %s", err))
	}

	// 记录日志
	sdk.Instance.Infof("[findByFileHash] fileHash=" + fact.FileHash)
	sdk.Instance.Infof("[findByFileHash] fileName=" + fact.FileName)

	// 返回结果
	return sdk.Success(result)
}

func main() {
	err := sandbox.Start(new(FactContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
