//go:build crypto
// +build crypto

/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

const (
	BulletproofsOpTypePedersenAddNum        = "pedersenAddNum"
	BulletproofsOpTypePedersenAddCommitment = "pedersenAddCommitment"
	BulletproofsOpTypePedersenSubNum        = "pedersenSubNum"
	BulletproofsOpTypePedersenSubCommitment = "pedersenSubCommitment"
	BulletproofsOpTypePedersenMulNum        = "pedersenMulNum"
	BulletproofsVerify                      = "bulletproofsVerify"
)

type ContractBulletproofs struct {
}

func (c *ContractBulletproofs) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (c *ContractBulletproofs) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (c *ContractBulletproofs) InvokeContract(method string) protogo.Response {

	switch method {
	case BulletproofsOpTypePedersenAddNum:
		return c.PedersenAddNum()
	case BulletproofsOpTypePedersenAddCommitment:
		return c.PedersenAddCommitment()
	case BulletproofsOpTypePedersenSubNum:
		return c.PedersenSubNum()
	case BulletproofsOpTypePedersenSubCommitment:
		return c.PedersenSubCommitment()
	case BulletproofsOpTypePedersenMulNum:
		return c.PedersenMulNum()
	case BulletproofsVerify:
		return c.Verify()
	default:
		return sdk.Error("invalid method")
	}
}

// PedersenAddNum Compute a commitment to x + y from a commitment to x without revealing the value x, where y is a scalar
// commitment: C = xB + rB'
// value: the value y
// return1: the new commitment to x + y: C' = (x + y)B + rB'
func (c *ContractBulletproofs) PedersenAddNum() protogo.Response {
	commitment, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["commitment"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	numStr := sdk.Instance.GetArgs()["num"]
	num, err := strconv.Atoi(string(numStr))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to atoi string to int: %v", err))
	}
	res, err := sdk.Bulletproofs.PedersenAddNum(commitment, uint64(num))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to add num, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Bulletproofs", BulletproofsOpTypePedersenAddNum, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save bulletproofs result, %v", err))
	}
	return sdk.Success(res)
}

// PedersenAddCommitment Compute a commitment to x + y from commitments to x and y, without revealing the value x and y
// commitment1: commitment to x: Cx = xB + rB'
// commitment2: commitment to y: Cy = yB + sB'
// return: commitment to x + y: C = (x + y)B + (r + s)B'
func (c *ContractBulletproofs) PedersenAddCommitment() protogo.Response {
	commitment1, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["commitment1"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	commitment2, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["commitment2"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	res, err := sdk.Bulletproofs.PedersenAddCommitment(commitment1, commitment2)
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to add commitment, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Bulletproofs", BulletproofsOpTypePedersenAddCommitment, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save bulletproofs result, %v", err))
	}
	return sdk.Success(res)
}

// PedersenSubNum Compute a commitment to x - y from a commitment to x without revealing the value x, where y is a scalar
// commitment: C = xB + rB'
// value: the value y
// return1: the new commitment to x - y: C' = (x - y)B + rB'
func (c *ContractBulletproofs) PedersenSubNum() protogo.Response {
	commitment, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["commitment"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	numStr := sdk.Instance.GetArgs()["num"]
	num, err := strconv.Atoi(string(numStr))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to atoi string to int: %v", err))
	}
	res, err := sdk.Bulletproofs.PedersenSubNum(commitment, uint64(num))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to sub num, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Bulletproofs", BulletproofsOpTypePedersenSubNum, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save bulletproofs result, %v", err))
	}
	return sdk.Success(res)
}

// PedersenSubCommitment Compute a commitment to x - y from commitments to x and y, without revealing the value x and y
// commitment1: commitment to x: Cx = xB + rB'
// commitment2: commitment to y: Cy = yB + sB'
// return: commitment to x - y: C = (x - y)B + (r - s)B'
func (c *ContractBulletproofs) PedersenSubCommitment() protogo.Response {
	commitment1, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["commitment1"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	commitment2, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["commitment2"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	res, err := sdk.Bulletproofs.PedersenSubCommitment(commitment1, commitment2)
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to sub commitment, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Bulletproofs", BulletproofsOpTypePedersenSubNum, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save bulletproofs result, %v", err))
	}
	return sdk.Success(res)
}

// PedersenMulNum Compute a commitment to x * y from a commitment to x and an integer y, without revealing the value x and y
// commitment1: commitment to x: Cx = xB + rB'
// value: integer value y
// return: commitment to x * y: C = (x * y)B + (r * y)B'
func (c *ContractBulletproofs) PedersenMulNum() protogo.Response {
	commitment, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["commitment"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	numStr := sdk.Instance.GetArgs()["num"]
	num, err := strconv.Atoi(string(numStr))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to atoi string to int: %v", err))
	}
	res, err := sdk.Bulletproofs.PedersenMulNum(commitment, uint64(num))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to mul num, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Bulletproofs", BulletproofsOpTypePedersenMulNum, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save bulletproofs result, %v", err))
	}
	return sdk.Success(res)
}

// Verify Verify the validity of a proof
// proof: the zero-knowledge proof proving the number committed in commitment is in the range [0, 2^64)
// commitment: commitment bindingly hiding the number x
// return: true on valid proof, false otherwise
func (c *ContractBulletproofs) Verify() protogo.Response {
	proof, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["proof"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	commitment, err := base64.StdEncoding.DecodeString(string(sdk.Instance.GetArgs()["commitment"]))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string: %v", err))
	}
	ok, err := sdk.Bulletproofs.Verify(proof, commitment)
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to verify, %v", err))
	}
	if err = sdk.Instance.PutState("Bulletproofs", BulletproofsVerify, strconv.FormatBool(ok)); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save bulletproofs result, %v", err))
	}
	return sdk.Success([]byte(strconv.FormatBool(ok)))
}

func main() {
	err := sandbox.Start(new(ContractBulletproofs))
	if err != nil {
		log.Fatal(err)
	}
}
