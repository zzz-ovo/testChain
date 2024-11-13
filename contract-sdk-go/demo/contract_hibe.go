/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package demo

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"encoding/json"
	"fmt"
	"log"
)

type ContractHIBE struct {
}

func (c *ContractHIBE) InitContract() protogo.Response {
	return sdk.Success([]byte("Init success"))
}

func (c *ContractHIBE) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade success"))
}

func (c *ContractHIBE) InvokeContract(method string) protogo.Response {

	switch method {
	case "save_hibe_params":
		return c.SaveParams()
	case "save_hibe_msg":
		return c.SaveMsg()
	case "find_params_by_org_id":
		return c.FindParamsByOrgId()
	default:
		return sdk.Error("invalid method")
	}
}

// SaveParams saves HIBE params
func (c *ContractHIBE) SaveParams() protogo.Response {
	params := sdk.Instance.GetArgs()
	if orgId, ok := params["org_id"]; !ok {
		return sdk.Error("key org_id not exist")
	} else if hibeParams, ok := params["params"]; !ok {
		return sdk.Error("key params not exist")
	} else {
		if err := sdk.Instance.PutStateByte("HIBE_params", string(orgId), hibeParams); err != nil {
			return sdk.Error(fmt.Sprintf("failed to save hibe params, %v", err))
		}
	}
	return sdk.Success([]byte("success"))
}

// SaveMsg saves HIBE msg
func (c *ContractHIBE) SaveMsg() protogo.Response {
	params := sdk.Instance.GetArgs()
	paramsBytes, _ := json.Marshal(params)
	if txId, ok := params["tx_id"]; !ok {
		return sdk.Error("key tx_id not exist")
	} else {
		if err := sdk.Instance.PutStateByte("HIBE_msg", string(txId), paramsBytes); err != nil {
			return sdk.Error(fmt.Sprintf("failed to save hibe msg, %v", err))
		}
	}
	return sdk.Success([]byte("success"))
}

// FindParamsByOrgId return params by org id
func (c *ContractHIBE) FindParamsByOrgId() protogo.Response {
	params := sdk.Instance.GetArgs()
	if orgId, ok := params["org_id"]; !ok {
		return sdk.Error("key org_id not exist")
	} else {
		if val, err := sdk.Instance.GetStateByte("HIBE_params", string(orgId)); err != nil {
			return sdk.Error(fmt.Sprintf("failed to get hibe params, %v", err))
		} else {
			return sdk.Success(val)
		}
	}
}

func main() {
	err := sandbox.Start(new(ContractHIBE))
	if err != nil {
		log.Fatal(err)
	}
}
