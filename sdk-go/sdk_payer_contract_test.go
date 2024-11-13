package chainmaker_sdk_go

import (
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/stretchr/testify/require"
)

func TestChainClient_SetContractMethodPayer(t *testing.T) {
	//nolint:gosec
	privKey := `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIEAmnFvzHhYsqScEJJj2jlPFRsYRbYpNlr9LVan1xN4GoAoGCCqGSM49
AwEHoUQDQgAE5Gqhg0DDxrH1A4P3hVperV0IeC5qtIqoCu4Qjf11GSbxW7aBIeR5
aXHO9KxVWqg3oBM5uLZmiFgphpvaz79ADQ==
-----END EC PRIVATE KEY-----
`
	//nolint:gosec
	pubKey := `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEnV76M1KpWj47q58vwx6+boBWRcGL
C/Hsok7mM3Sh7hVm0P1kfPlSH9D9cRk3b0xYOZaaPuucgx1GncrHkzzzxA==
-----END PUBLIC KEY-----
`
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.SetContractMethodPayer(
				"address1",
				"contractName1",
				"method1",
				"requestId1",
				"org1",
				[]byte(privKey),
				[]byte(pubKey),
				1000000)
			require.NotNil(t, err)
		})
	}
}

func TestChainClient_UnsetContractMethodPayer(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.UnsetContractMethodPayer(
				"contractName1",
				"method1",
				1000000)
			require.NotNil(t, err)
		})
	}
}

func TestChainClient_QueryContractMethodPayer(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.QueryContractMethodPayer(
				"contractName1",
				"method1",
				1000000)
			require.NotNil(t, err)
		})
	}
}

func TestChainClient_QueryTxPayer(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.QueryTxPayer("txId", 1000000)
			require.NotNil(t, err)
		})
	}
}
