/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"log"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type StoreMapContract struct {
}

func (s *StoreMapContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (s *StoreMapContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (s *StoreMapContract) InvokeContract(method string) protogo.Response {

	switch method {
	case "set":
		return s.Set()
	case "get":
		return s.Get()
	case "del":
		return s.Del()
	default:
		return sdk.Error("invalid method")
	}
}

func (s *StoreMapContract) Set() protogo.Response {
	var err error
	var storeMap *sdk.StoreMap

	var deep int64 = 3
	m := "m1"

	storeMap, err = sdk.NewStoreMap(m, deep)
	if err != nil {
		return sdk.Error(err.Error())
	}

	key := []string{"key1", "key2", "key3"}
	value := []byte("value")

	err = storeMap.Set(key, value)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("success"))
}

func (s *StoreMapContract) Get() protogo.Response {
	var err error
	var storeMap *sdk.StoreMap

	var deep int64 = 3
	m := "m1"

	storeMap, err = sdk.NewStoreMap(m, deep)
	if err != nil {
		return sdk.Error(err.Error())
	}

	key := []string{"key1", "key2", "key3"}

	_, err = storeMap.Get(key)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("success"))
}

func (s *StoreMapContract) Del() protogo.Response {
	var err error
	var storeMap *sdk.StoreMap

	var deep int64 = 3
	m := "m1"

	storeMap, err = sdk.NewStoreMap(m, deep)
	if err != nil {
		return sdk.Error(err.Error())
	}

	key := []string{"key1", "key2", "key3"}

	err = storeMap.Del(key)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("success"))
}

func (s *StoreMapContract) Exist() protogo.Response {
	var err error
	var storeMap *sdk.StoreMap

	var deep int64 = 3
	m := "m1"

	storeMap, err = sdk.NewStoreMap(m, deep)
	if err != nil {
		return sdk.Error(err.Error())
	}

	key := []string{"key1", "key2", "key3"}

	_, err = storeMap.Exist(key)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("success"))
}

func main() {

	err := sandbox.Start(new(StoreMapContract))
	if err != nil {
		log.Fatal(err)
	}
}
