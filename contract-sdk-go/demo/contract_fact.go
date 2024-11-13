/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"encoding/json"
	"log"
	"strconv"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type FactContract struct {
}

func (f *FactContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (f *FactContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (f *FactContract) InvokeContract(method string) protogo.Response {

	switch method {
	case "save":
		return f.Save()
	case "findByFileHash":
		return f.FindByFileHash()
	default:
		return sdk.Error("invalid method")
	}
}

// 存证对象
type Fact struct {
	FileHash string `json:"FileHash"`
	FileName string `json:"FileName"`
	Time     int32  `json:"time"`
}

// 新建存证对象
func NewFact(FileHash string, FileName string, time int32) *Fact {
	fact := &Fact{
		FileHash: FileHash,
		FileName: FileName,
		Time:     time,
	}
	return fact
}

func (f *FactContract) Save() protogo.Response {
	params := sdk.Instance.GetArgs()

	// 获取参数
	fileHash := string(params["file_hash"])
	fileName := string(params["file_name"])
	timeStr := string(params["time"])
	time, err := strconv.Atoi(timeStr)
	if err != nil {
		msg := "time is [" + timeStr + "] not int"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}

	// 构建结构体
	fact := NewFact(fileHash, fileName, int32(time))

	// 序列化
	factBytes, _ := json.Marshal(fact)

	// 发送事件
	sdk.Instance.EmitEvent("topic_vx", []string{fact.FileHash, fact.FileName})

	// 存储数据
	err = sdk.Instance.PutStateByte("fact_bytes", fact.FileHash, factBytes)
	if err != nil {
		return sdk.Error("fail to save fact bytes")
	}

	// 记录日志
	sdk.Instance.Debugf("[save] FileHash=" + fact.FileHash)
	sdk.Instance.Debugf("[save] FileName=" + fact.FileName)

	// 返回结果
	return sdk.Success([]byte(fact.FileName + fact.FileHash))

}

func (f *FactContract) FindByFileHash() protogo.Response {
	// 获取参数
	FileHash := string(sdk.Instance.GetArgs()["file_hash"])

	// 查询结果
	result, err := sdk.Instance.GetStateByte("fact_bytes", FileHash)
	if err != nil {
		return sdk.Error("failed to call get_state")
	}

	// 反序列化
	var fact Fact
	_ = json.Unmarshal(result, &fact)

	// 记录日志
	sdk.Instance.Debugf("[find_by_file_hash] FileHash=" + fact.FileHash)
	sdk.Instance.Debugf("[find_by_file_hash] FileName=" + fact.FileName)

	// 返回结果
	return sdk.Success(result)
}

func main() {

	err := sandbox.Start(new(FactContract))
	if err != nil {
		log.Fatal(err)
	}
}
