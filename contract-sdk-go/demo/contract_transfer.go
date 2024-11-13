/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"fmt"
	"log"
	"strconv"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type TransferContract struct {
}

func (t *TransferContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (t *TransferContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (t *TransferContract) InvokeContract(method string) protogo.Response {

	switch method {
	case "init":
		return t.Init()
	case "transfer":
		return t.Transfer()
	default:
		return sdk.Error("invalid method")
	}
}

func (t *TransferContract) Init() protogo.Response {
	args := sdk.Instance.GetArgs()

	accFrom := string(args["accFrom"])
	accTo := string(args["accTo"])
	fromBal := string(args["from_bal"])
	toBal := string(args["to_bal"])
	startIndex := string(args["start_index"]) // 分批创建账户 - start index
	endIndex := string(args["end_index"])     // 分批创建账户 - end index

	if !t.isNumber(fromBal) {
		return sdk.Error("from_bal is not a number")
	}

	if !t.isNumber(toBal) {
		return sdk.Error("to_bal is not a number")
	}

	start, err := strconv.Atoi(startIndex)
	if err != nil {
		return sdk.Error("start index is not a number")
	}

	end, err := strconv.Atoi(endIndex)
	if err != nil {
		return sdk.Error("end index is not a number")
	}

	if start > end {
		return sdk.Error("start index bigger than end index")
	}

	for i := start; i <= end; i++ {
		newAccFrom := accFrom + strconv.Itoa(i)
		newAccTo := accTo + strconv.Itoa(i)

		err = sdk.Instance.PutStateFromKey(newAccFrom, fromBal)
		if err != nil {
			return sdk.Error(fmt.Sprintf("putState(%s, %s) err: %+v", newAccFrom, fromBal, err))
		}

		err = sdk.Instance.PutStateFromKey(newAccTo, toBal)
		if err != nil {
			return sdk.Error(fmt.Sprintf("putState(%s, %s) err: %+v", newAccTo, toBal, err))
		}
	}

	return sdk.Success([]byte("init success"))
}

func (t *TransferContract) Transfer() protogo.Response {
	args := sdk.Instance.GetArgs()
	accFrom := string(args["acc_from"])
	accTo := string(args["acc_to"])
	amtTrans, err := strconv.Atoi(string(args["amt_trans"]))
	if err != nil {
		return sdk.Error("amt_trans is not a number")
	}

	fromBalStr, err := sdk.Instance.GetStateFromKey(accFrom)
	if err != nil {
		return sdk.Error(fmt.Sprintf("getState(%s) error: %+v", accFrom, err))
	}

	fromBal, err := strconv.Atoi(fromBalStr)
	if err != nil {
		return sdk.Error("from_bal is not a number")
	}

	toBalStr, err := sdk.Instance.GetStateFromKey(accTo)
	if err != nil {
		return sdk.Error(fmt.Sprintf("getState(%s) error: %+v", accTo, err))
	}
	toBal, err := strconv.Atoi(toBalStr)
	if err != nil {
		return sdk.Error("to_bal is not a number")
	}
	if fromBal < amtTrans {
		return sdk.Error(fmt.Sprintf("money doesn't enough, from_bal: %d, amt_trans: %d", fromBal, amtTrans))
	}
	fromBal -= amtTrans
	toBal += amtTrans
	err = sdk.Instance.PutStateFromKey(accFrom, strconv.Itoa(fromBal))
	if err != nil {
		return sdk.Error(fmt.Sprintf("putState(%s, %d) err: %+v", accFrom, fromBal, err))
	}
	err = sdk.Instance.PutStateFromKey(accTo, strconv.Itoa(toBal))
	if err != nil {
		return sdk.Error(fmt.Sprintf("putState(%s, %d) err: %+v", accTo, toBal, err))
	}
	return sdk.Success([]byte("transfer success"))
}

func (t *TransferContract) isNumber(bal string) bool {
	_, err := strconv.Atoi(bal)
	if err != nil {
		return false
	}
	return true
}

func main() {
	err := sandbox.Start(new(TransferContract))
	if err != nil {
		log.Fatal(err)
	}
}
