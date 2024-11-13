/*
Copyright (C) BABEC. All rights reserved.
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

type Contract2 struct {
}

func (c *Contract2) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (c *Contract2) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (c *Contract2) InvokeContract(method string) protogo.Response {

	switch method {
	case "display":
		return c.Display()
	case "crossCall":
		return c.CrossCall()
	default:
		return sdk.Error("invalid method")
	}
}

func (c *Contract2) Display() protogo.Response {
	return sdk.Success([]byte("successfully display"))
}

func (c *Contract2) CrossCall() protogo.Response {

	contractName := "contract1"
	contractMethod := "find"

	crossContractArgs := make(map[string][]byte)
	crossContractArgs["key"] = []byte("key")

	result := sdk.Instance.CallContract(contractName, contractMethod, crossContractArgs)
	return result
}

func main() {
	err := sandbox.Start(new(Contract2))
	if err != nil {
		log.Fatal(err)
	}
}
