/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"fmt"
	"log"
	"time"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	vmPb "chainmaker.org/chainmaker/pb-go/v2/vm"
)

type TestContract struct {
}

func (t *TestContract) InitContract() protogo.Response {

	return sdk.Success([]byte("Init Success"))
}

func (t *TestContract) UpgradeContract() protogo.Response {

	return sdk.Success([]byte("Upgrade Success"))
}

func (t *TestContract) InvokeContract(method string) protogo.Response {

	switch method {
	case "display":
		return t.display()
	case "put_state":
		return t.putState()
	case "put_state_byte":
		return t.putStateByte()
	case "put_state_from_key":
		return t.putStateFromKey()
	case "put_state_from_key_byte":
		return t.putStateFromKeyByte()
	case "get_state":
		return t.getState()
	case "get_state_byte":
		return t.getStateByte()
	case "get_state_from_key":
		return t.getStateFromKey()
	case "get_state_from_key_byte":
		return t.getStateFromKeyByte()
	case "del_state":
		return t.delState()
	case "time_out":
		return t.timeOut()
	case "out_of_range":
		return t.outOfRange()
	case "cross_contract":
		return t.crossContract()
	case "cross_contract_self":
		return t.callSelf()

	// kvIterator
	case "construct_data":
		return t.constructData()
	case "kv_iterator_test":
		return t.kvIterator()

	// keyHistoryIterator
	case "key_history_kv_iter":
		return t.keyHistoryIter()

	// getSenderAddress
	case "get_sender_address":
		return t.getSenderAddr()
	case "get_sender":
		return t.getSender()
	case "get_origin":
		return t.getOrigin()

		// performance
	case "save":
		return t.Save()
	case "findByFileHash":
		return t.FindByFileHash()

	default:
		msg := fmt.Sprintf("unknown method")
		return sdk.Error(msg)
	}
}

func (t *TestContract) Save() protogo.Response {
	key := string(sdk.Instance.GetArgs()["file_key"])
	name := sdk.Instance.GetArgs()["file_name"]

	err := sdk.Instance.PutStateByte(key, "", name)
	if err != nil {
		return sdk.Error("fail to save")
	}
	return sdk.Success([]byte("success"))
}

func (t *TestContract) FindByFileHash() protogo.Response {
	key := string(sdk.Instance.GetArgs()["file_key"])

	_, err := sdk.Instance.GetStateByte(key, "")
	if err != nil {
		return sdk.Error("fail to find")
	}
	return sdk.Success([]byte(""))
}

func (t *TestContract) getSenderAddr() protogo.Response {
	senderAddr, err := sdk.Instance.GetSenderAddr()
	if err != nil {
		msg := "GetSenderAddr failed"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}

	l := len([]byte(senderAddr))
	msg := fmt.Sprintf("=== sender address: [%s] len: %d===", senderAddr, l)
	sdk.Instance.Debugf(msg)
	return sdk.Success([]byte(senderAddr))
}

func (t *TestContract) getSender() protogo.Response {
	sender, err := sdk.Instance.Sender()
	if err != nil {
		msg := "get sender failed"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}

	l := len([]byte(sender))
	msg := fmt.Sprintf("=== sender: [%s] len: %d===", sender, l)
	sdk.Instance.Debugf(msg)
	return sdk.Success([]byte(sender))
}

func (t *TestContract) getOrigin() protogo.Response {
	origin, err := sdk.Instance.Origin()
	if err != nil {
		msg := "get origin failed"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}

	l := len([]byte(origin))
	msg := fmt.Sprintf("=== origin: [%s] len: %d===", origin, l)
	sdk.Instance.Debugf(msg)
	return sdk.Success([]byte(origin))
}

func (t *TestContract) keyHistoryIter() protogo.Response {
	sdk.Instance.Debugf("===Key History Iter START===")
	args := sdk.Instance.GetArgs()
	key := string(args["key"])
	field := string(args["field"])

	iter, err := sdk.Instance.NewHistoryKvIterForKey(key, field)
	if err != nil {
		msg := "NewHistoryIterForKey failed"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}

	sdk.Instance.Debugf("===create iter success===")

	count := 0
	for iter.HasNext() {
		sdk.Instance.Debugf("HasNext")
		count++
		km, err := iter.Next()
		if err != nil {
			msg := "iterator failed to get the next element"
			sdk.Instance.Debugf(msg)
			return sdk.Error(msg)
		}

		sdk.Instance.Debugf(fmt.Sprintf("=== Data History [%d] Info:", count))
		sdk.Instance.Debugf(fmt.Sprintf("=== Key: [%s]", km.Key))
		sdk.Instance.Debugf(fmt.Sprintf("=== Field: [%s]", km.Field))
		sdk.Instance.Debugf(fmt.Sprintf("=== Value: [%s]", km.Value))
		sdk.Instance.Debugf(fmt.Sprintf("=== TxId: [%s]", km.TxId))
		sdk.Instance.Debugf(fmt.Sprintf("=== BlockHeight: [%d]", km.BlockHeight))
		sdk.Instance.Debugf(fmt.Sprintf("=== IsDelete: [%t]", km.IsDelete))
		sdk.Instance.Debugf(fmt.Sprintf("=== Timestamp: [%s]", km.Timestamp))
	}

	closed, err := iter.Close()
	if !closed || err != nil {
		msg := fmt.Sprintf("iterator close failed, %s", err.Error())
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}
	sdk.Instance.Debugf("===iter close success===")

	sdk.Instance.Debugf("===Key History Iter END===")

	return sdk.Success([]byte("get key history successfully"))
}

func (t *TestContract) display() protogo.Response {
	return sdk.Success([]byte("display successful"))
}

func (t *TestContract) putState() protogo.Response {
	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])
	getField := string(args["field"])
	getValue := string(args["value"])

	err := sdk.Instance.PutState(getKey, getField, getValue)
	if err != nil {
		return sdk.Error(err.Error())
	}

	return sdk.Success([]byte("put state successfully"))
}

func (t *TestContract) putStateByte() protogo.Response {
	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])
	getField := string(args["field"])
	getValue := args["value"]

	err := sdk.Instance.PutStateByte(getKey, getField, getValue)
	if err != nil {
		return sdk.Error(err.Error())
	}

	return sdk.Success([]byte("put state successfully"))
}

func (t *TestContract) putStateFromKey() protogo.Response {
	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])
	getValue := string(args["value"])

	err := sdk.Instance.PutStateFromKey(getKey, getValue)
	if err != nil {
		return sdk.Error(err.Error())
	}

	return sdk.Success([]byte("put state successfully"))
}

func (t *TestContract) putStateFromKeyByte() protogo.Response {
	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])
	getValue := args["value"]

	err := sdk.Instance.PutStateFromKeyByte(getKey, getValue)
	if err != nil {
		return sdk.Error(err.Error())
	}

	return sdk.Success([]byte("put state successfully"))
}

func (t *TestContract) getState() protogo.Response {

	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])
	field := string(args["field"])

	result, err := sdk.Instance.GetState(getKey, field)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(result))
}

func (t *TestContract) getStateByte() protogo.Response {

	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])
	field := string(args["field"])

	result, err := sdk.Instance.GetStateByte(getKey, field)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success(result)
}

func (t *TestContract) getStateFromKey() protogo.Response {

	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])

	result, err := sdk.Instance.GetStateFromKey(getKey)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(result))
}

func (t *TestContract) getStateFromKeyByte() protogo.Response {
	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])

	result, err := sdk.Instance.GetStateFromKeyByte(getKey)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success(result)
}

func (t *TestContract) delState() protogo.Response {
	args := sdk.Instance.GetArgs()

	getKey := string(args["key"])
	getField := string(args["field"])

	err := sdk.Instance.DelState(getKey, getField)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("delete successfully"))
}

func (t *TestContract) timeOut() protogo.Response {
	time.Sleep(50 * time.Second)
	return sdk.Success([]byte("success finish timeout"))
}

func (t *TestContract) outOfRange() protogo.Response {
	var group []string
	group[0] = "abc"
	fmt.Println(group[0])
	return sdk.Success([]byte("exit out of range"))
}

func (t *TestContract) crossContract() protogo.Response {

	args := sdk.Instance.GetArgs()

	contractName := string(args["contract_name"])
	// contractVersion := string(args["contract_version"])

	calledMethod := string(args["contract_method"])

	crossContractArgs := make(map[string][]byte)
	crossContractArgs["method"] = []byte(calledMethod)

	// response could be correct or error
	response := sdk.Instance.CallContract(contractName, calledMethod, crossContractArgs)
	sdk.Instance.EmitEvent("cross contract", []string{"success"})
	return response
}

func (t *TestContract) callSelf() protogo.Response {

	sdk.Instance.Debugf("testing call self")

	contractName := "contract_test10"
	contractMethod := "CallSelf"

	crossContractArgs := make(map[string][]byte)

	response := sdk.Instance.CallContract(contractName, contractMethod, crossContractArgs)
	return response
}

// constructData 提供Kv迭代器的测试数据
/*
	| Key   | Field   | Value |
	| ---   | ---     | ---   |
	| key1  | field1  | val   |
	| key1  | field2  | val   |
	| key1  | field23 | val   |
	| ey1   | field3  | val   |
	| key2  | field1  | val   |
	| key3  | field2  | val   |
	| key33 | field2  | val   |
	| key4  | field3  | val   |
*/
func (t *TestContract) constructData() protogo.Response {
	dataList := []struct {
		key   string
		field string
		value string
	}{
		{key: "key1", field: "field1", value: "val"},
		{key: "key1", field: "field2", value: "val"},
		{key: "key1", field: "field23", value: "val"},
		{key: "key1", field: "field3", value: "val"},
		{key: "key2", field: "field1", value: "val"},
		{key: "key3", field: "field2", value: "val"},
		{key: "key33", field: "field2", value: "val"},
		{key: "key33", field: "field2", value: "val"},
		{key: "key4", field: "field3", value: "val"},
	}

	for _, data := range dataList {
		err := sdk.Instance.PutState(data.key, data.field, data.value)
		if err != nil {
			msg := fmt.Sprintf("constructData failed, %s", err.Error())
			sdk.Instance.Debugf(msg)
			return sdk.Error(msg)
		}
	}

	return sdk.Success([]byte("construct success!"))
}

// kvIterator 前置数据
/*
	| Key   | Field   | Value |
	| ---   | ---     | ---   |
	| key1  | field1  | val   |
	| key1  | field2  | val   |
	| key1  | field23 | val   |
	| ey1   | field3  | val   |
	| key2  | field1  | val   |
	| key3  | field2  | val   |
	| key33 | field2  | val   |
	| key4  | field3  | val   |
*/
func (t *TestContract) kvIterator() protogo.Response {
	sdk.Instance.Debugf("===construct START===")
	dataList := []struct {
		key   string
		field string
		value string
	}{
		{key: "key1", field: "field1", value: "val"},
		{key: "key1", field: "field2", value: "val"},
		{key: "key1", field: "field23", value: "val"},
		{key: "key1", field: "field3", value: "val"},
		{key: "key2", field: "field1", value: "val"},
		{key: "key3", field: "field2", value: "val"},
		{key: "key33", field: "field2", value: "val"},
		{key: "key33", field: "field2", value: "val"},
		{key: "key4", field: "field3", value: "val"},
	}

	for _, data := range dataList {
		err := sdk.Instance.PutState(data.key, data.field, data.value)
		if err != nil {
			msg := fmt.Sprintf("constructData failed, %s", err.Error())
			sdk.Instance.Debugf(msg)
			return sdk.Error(msg)
		}
	}
	sdk.Instance.Debugf("===construct END===")

	sdk.Instance.Debugf("===kvIterator START===")
	iteratorList := make([]sdk.ResultSetKV, 4)

	// 能查询出 key2, key3, key33 三条数据
	iterator, err := sdk.Instance.NewIterator("key2", "key4")
	if err != nil {
		msg := "NewIterator failed"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}
	iteratorList[0] = iterator

	// 能查询出 field1, field2, field23 三条数据
	iteratorWithField, err := sdk.Instance.NewIteratorWithField("key1", "field1", "field3")
	if err != nil {
		// msg := "create with " + string(key1) + string(field1) + string(field3) + " failed"
		msg := "create with " + "key1" + "field1" + "field3" + " failed"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}
	iteratorList[1] = iteratorWithField

	// 能查询出 key3, key33 两条数据
	preWithKeyIterator, err := sdk.Instance.NewIteratorPrefixWithKey("key3")
	if err != nil {
		msg := "NewIteratorPrefixWithKey failed"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}
	iteratorList[2] = preWithKeyIterator

	// 能查询出 field2, field23 三条数据
	preWithKeyFieldIterator, err := sdk.Instance.NewIteratorPrefixWithKeyField("key1", "field2")
	if err != nil {
		msg := "NewIteratorPrefixWithKeyField failed"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}
	iteratorList[3] = preWithKeyFieldIterator

	for index, iter := range iteratorList {
		index++
		sdk.Instance.Debugf(fmt.Sprintf("===iterator %d START===", index))
		for iter.HasNext() {
			sdk.Instance.Debugf("HasNext Success")
			key, field, value, err := iter.Next()
			if err != nil {
				msg := "iterator failed to get the next element"
				sdk.Instance.Debugf(msg)
				return sdk.Error(msg)
			}

			sdk.Instance.Debugf(fmt.Sprintf("===[key: %s]===", key))
			sdk.Instance.Debugf(fmt.Sprintf("===[field: %s]===", field))
			sdk.Instance.Debugf(fmt.Sprintf("===[value: %s]===", value))
		}

		closed, err := iter.Close()
		if !closed || err != nil {
			msg := fmt.Sprintf("iterator %d close failed, %s", index, err.Error())
			sdk.Instance.Debugf(msg)
			return sdk.Error(msg)
		}
		sdk.Instance.Debugf(fmt.Sprintf("===iterator %d END===", index))
	}
	sdk.Instance.Debugf("===kvIterator END===")

	return sdk.Success([]byte("SUCCESS"))
}

// GetBatchState
/*
	| Key   | Field   | Value |
	| ---   | ---     | ---   |
	| key1  | field1  | val   |
	| key1  | field2  | val   |
	| key1  | field23 | val   |
	| ey1   | field3  | val   |
	| key2  | field1  | val   |
	| key3  | field2  | val   |
	| key33 | field2  | val   |
	| key4  | field3  | val   |
*/
func (t *TestContract) getBatchState() protogo.Response {
	sdk.Instance.Debugf("===construct START===")
	dataList := []struct {
		key   string
		field string
		value string
	}{
		{key: "key1", field: "field1", value: "key1-field1"},
		{key: "key1", field: "field2", value: "key1-field2"},
		{key: "key1", field: "field23", value: "key1-field23"},
		{key: "key1", field: "field3", value: "key1-field3"},
		{key: "key2", field: "field1", value: "key2-field1"},
		{key: "key3", field: "field2", value: "key3-field2"},
		{key: "key33", field: "field2", value: "key33-field2"},
		{key: "key33", field: "field2", value: "key33-field2"},
		{key: "key4", field: "field3", value: "key4-field3"},
	}

	for _, data := range dataList {
		err := sdk.Instance.PutState(data.key, data.field, data.value)
		if err != nil {
			msg := fmt.Sprintf("constructData failed, %s", err.Error())
			sdk.Instance.Debugf(msg)
			return sdk.Error(msg)
		}
	}
	sdk.Instance.Debugf("===construct END===")

	batchKeys := []*vmPb.BatchKey{
		{"key1", "field1", nil, ""},
		{"key1", "field2", nil, ""},
		{"key3", "field2", nil, ""},
		{"key4", "field2", nil, ""},
	}
	gotValues, err := sdk.Instance.GetBatchState(batchKeys)
	if err != nil {
		msg := fmt.Sprintf("constructData failed, %s", err.Error())
		return sdk.Error(msg)
	}
	result := fmt.Sprintf("%v", gotValues)

	return sdk.Success([]byte(result))
}

func main() {
	err := sandbox.Start(new(TestContract))
	if err != nil {
		log.Fatal(err)
	}
}
