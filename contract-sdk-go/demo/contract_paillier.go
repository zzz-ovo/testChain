//go:build crypto
// +build crypto

/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"encoding/base64"
	"fmt"
	"log"
)

const (
	PaillierOpTypeAddCiphertext = "AddCiphertext"
	PaillierOpTypeAddPlaintext  = "AddPlaintext"
	PaillierOpTypeSubCiphertext = "SubCiphertext"
	PaillierOpTypeSubPlaintext  = "SubPlaintext"
	PaillierOpTypeNumMul        = "NumMul"
	PaillierGet                 = "PaillierGet"
)

type ContractPaillier struct {
}

func (c *ContractPaillier) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (c *ContractPaillier) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (c *ContractPaillier) InvokeContract(method string) protogo.Response {

	switch method {
	case PaillierOpTypeAddCiphertext:
		return c.AddCiphertext()
	case PaillierOpTypeAddPlaintext:
		return c.AddPlaintext()
	case PaillierOpTypeSubCiphertext:
		return c.SubCiphertext()
	case PaillierOpTypeSubPlaintext:
		return c.SubPlaintext()
	case PaillierOpTypeNumMul:
		return c.NumMul()
	case PaillierGet:
		return c.PaillierGet()
	default:
		return sdk.Error("invalid method")
	}
}

// AddCiphertext homomorphically adds together two cipher texts.
func (c *ContractPaillier) AddCiphertext() protogo.Response {
	params := sdk.Instance.GetArgs()
	pubkey, ok := params["pubkey"]
	if !ok {
		return sdk.Error("pubkey is nil")
	}
	para1, ok := params["para1"]
	if !ok {
		return sdk.Error("para1 is nil")
	}
	ctBytes1, err := base64.StdEncoding.DecodeString(string(para1))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string para1: %v", err))
	}
	para2, ok := params["para2"]
	if !ok {
		return sdk.Error("para2 is nil")
	}
	ctBytes2, err := base64.StdEncoding.DecodeString(string(para2))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string para2: %v", err))
	}
	res, err := sdk.Paillier.AddCiphertext(pubkey, ctBytes1, ctBytes2)
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to add cipher text, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Paillier", PaillierOpTypeAddCiphertext, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save paillier result, %v", err))
	}
	return sdk.Success(res)
}

// AddPlaintext homomorphically adds a passed constant to the encrypted integer
func (c *ContractPaillier) AddPlaintext() protogo.Response {
	params := sdk.Instance.GetArgs()
	pubkey, ok := params["pubkey"]
	if !ok {
		return sdk.Error("pubkey is nil")
	}
	para1, ok := params["para1"]
	if !ok {
		return sdk.Error("para1 is nil")
	}
	ctBytes, err := base64.StdEncoding.DecodeString(string(para1))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string para1: %v", err))
	}
	para2, ok := params["para2"]
	if !ok {
		return sdk.Error("para2 is nil")
	}
	ptStr := string(para2)
	res, err := sdk.Paillier.AddPlaintext(pubkey, ctBytes, ptStr)
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to add plaintext, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Paillier", PaillierOpTypeAddPlaintext, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save paillier result, %v", err))
	}
	return sdk.Success(res)
}

// SubCiphertext homomorphically subs two cipher texts.
func (c *ContractPaillier) SubCiphertext() protogo.Response {
	params := sdk.Instance.GetArgs()
	pubkey, ok := params["pubkey"]
	if !ok {
		return sdk.Error("pubkey is nil")
	}
	para1, ok := params["para1"]
	if !ok {
		return sdk.Error("para1 is nil")
	}
	ctBytes1, err := base64.StdEncoding.DecodeString(string(para1))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string para1: %v", err))
	}
	para2, ok := params["para2"]
	if !ok {
		return sdk.Error("para2 is nil")
	}
	ctBytes2, err := base64.StdEncoding.DecodeString(string(para2))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string para2: %v", err))
	}
	res, err := sdk.Paillier.SubCiphertext(pubkey, ctBytes1, ctBytes2)
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to sub cipher text, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Paillier", PaillierOpTypeSubCiphertext, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save paillier result, %v", err))
	}
	return sdk.Success(res)
}

// SubPlaintext homomorphically subs a passed constant to the encrypted integer
func (c *ContractPaillier) SubPlaintext() protogo.Response {
	params := sdk.Instance.GetArgs()
	pubkey, ok := params["pubkey"]
	if !ok {
		return sdk.Error("pubkey is nil")
	}
	para1, ok := params["para1"]
	if !ok {
		return sdk.Error("para1 is nil")
	}
	ctBytes, err := base64.StdEncoding.DecodeString(string(para1))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string para1: %v", err))
	}
	para2, ok := params["para2"]
	if !ok {
		return sdk.Error("para2 is nil")
	}
	ptStr := string(para2)
	res, err := sdk.Paillier.SubPlaintext(pubkey, ctBytes, ptStr)
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to sub plaintext, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Paillier", PaillierOpTypeSubPlaintext, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save paillier result, %v", err))
	}
	return sdk.Success(res)
}

// NumMul homomorphically multiplies an encrypted integer (cipher text) by a constant
func (c *ContractPaillier) NumMul() protogo.Response {
	params := sdk.Instance.GetArgs()
	pubkey, ok := params["pubkey"]
	if !ok {
		return sdk.Error("pubkey is nil")
	}
	para1, ok := params["para1"]
	if !ok {
		return sdk.Error("para1 is nil")
	}
	ctBytes, err := base64.StdEncoding.DecodeString(string(para1))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to base64 decode string para1: %v", err))
	}
	para2, ok := params["para2"]
	if !ok {
		return sdk.Error("para2 is nil")
	}
	ptStr := string(para2)
	res, err := sdk.Paillier.NumMul(pubkey, ctBytes, ptStr)
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to mul plaintext, %v", err))
	}
	if err = sdk.Instance.PutStateByte("Paillier", PaillierOpTypeNumMul, res); err != nil {
		return sdk.Error(fmt.Sprintf("failed to save paillier result, %v", err))
	}
	return sdk.Success(res)
}

// PaillierGet return result
func (c *ContractPaillier) PaillierGet() protogo.Response {
	params := sdk.Instance.GetArgs()
	handletype, ok := params["handletype"]
	if !ok {
		return sdk.Error("handletype is nil")
	}
	result, err := sdk.Instance.GetStateByte("Paillier", string(handletype))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to get paillier result, %v", err))
	}
	return sdk.Success(result)
}

func main() {
	err := sandbox.Start(new(ContractPaillier))
	if err != nil {
		log.Fatal(err)
	}
}
