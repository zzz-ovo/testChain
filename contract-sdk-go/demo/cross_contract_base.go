/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"fmt"
	"log"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type Contract1 struct {
}

func (c *Contract1) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (c *Contract1) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (c *Contract1) InvokeContract(method string) protogo.Response {

	switch method {
	case "save":
		return c.Save()
	case "find":
		return c.Find()
	default:
		return sdk.Error("invalid method")
	}
}

func (c *Contract1) Save() protogo.Response {
	params := sdk.Instance.GetArgs()

	key := string(params["key"])
	value := string(params["value"])

	err := sdk.Instance.PutStateFromKey(key, value)
	if err != nil {
		errMsg := fmt.Sprintf("fail to save key [%s], value [%s]: err: [%s]",
			key, value, err)
		return sdk.Error(errMsg)
	}
	return sdk.Success([]byte("successfully save"))
}

func (c *Contract1) Find() protogo.Response {
	params := sdk.Instance.GetArgs()

	key := string(params["key"])
	value, err := sdk.Instance.GetStateFromKey(key)
	if err != nil {
		errMsg := fmt.Sprintf("fail to get key [%s], value [%s]: err: [%s]",
			key, value, err)
		return sdk.Error(errMsg)
	}
	return sdk.Success([]byte(value))
}

func main() {
	err := sandbox.Start(new(Contract1))
	if err != nil {
		log.Fatal(err)
	}
}
