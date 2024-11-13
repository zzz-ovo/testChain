/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/stretchr/testify/require"
)

func TestChainClient_PubKeyManage(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreatePubkeyAddPayload("pubkey", "orgId", "role")
			require.Nil(t, err)
			_, err = cc.CreatePubkeyDelPayload("pubkey", "orgId")
			require.Nil(t, err)
			_, err = cc.CreatePubkeyQueryPayload("pubkey")
			require.Nil(t, err)
		})
	}
}
