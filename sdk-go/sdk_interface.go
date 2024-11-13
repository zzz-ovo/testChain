/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"context"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/config"
	"chainmaker.org/chainmaker/pb-go/v2/discovery"
	"chainmaker.org/chainmaker/pb-go/v2/store"
	"chainmaker.org/chainmaker/pb-go/v2/sync"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/pb-go/v2/txpool"
)

// SDKInterface # ChainMaker Go SDK 接口说明
type SDKInterface interface {
	// ## 1 用户合约接口
	// ### 1.1 创建合约待签名payload生成
	// **参数说明**
	//   - contractName: 合约名
	//   - version: 版本号
	//   - byteCodeStringOrFilePath: 支持传入合约二进制文件路径或Hex或Base64编码的string
	//   - runtime: 合约运行环境
	//   - kvs: 合约初始化参数
	// ```go
	CreateContractCreatePayload(contractName, version, byteCodeStringOrFilePath string, runtime common.RuntimeType,
		kvs []*common.KeyValuePair) (*common.Payload, error)
	// ```

	// ### 1.2 升级合约待签名payload生成
	// **参数说明**
	//   - contractName: 合约名
	//   - version: 版本号
	//   - byteCodeStringOrFilePath: 支持传入合约二进制文件路径或Hex或Base64编码的string
	//   - runtime: 合约运行环境
	//   - kvs: 合约升级参数
	// ```go
	CreateContractUpgradePayload(contractName, version, byteCodeStringOrFilePath string, runtime common.RuntimeType,
		kvs []*common.KeyValuePair) (*common.Payload, error)
	// ```

	// ### 1.3 冻结合约payload生成
	// **参数说明**
	//   - contractName: 合约名
	// ```go
	CreateContractFreezePayload(contractName string) (*common.Payload, error)
	// ```

	// ### 1.4 解冻合约payload生成
	// **参数说明**
	//   - contractName: 合约名
	// ```go
	CreateContractUnfreezePayload(contractName string) (*common.Payload, error)
	// ```

	// ### 1.5 吊销合约payload生成
	// **参数说明**
	//   - contractName: 合约名
	// ```go
	CreateContractRevokePayload(contractName string) (*common.Payload, error)
	// ```

	// ### 1.6 合约管理获取Payload签名
	// **参数说明**
	//   - payload: 待签名payload
	// ```go
	SignContractManagePayload(payload *common.Payload) (*common.EndorsementEntry, error)
	// ```

	// ### 1.7 发送合约管理请求（创建、更新、冻结、解冻、吊销）
	// **参数说明**
	//   - payload: 交易payload
	//   - endorsers: 背书签名信息列表
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	SendContractManageRequest(payload *common.Payload, endorsers []*common.EndorsementEntry, timeout int64,
		withSyncResult bool) (*common.TxResponse, error)
	// ### 1.7 发送合约管理请求（创建、更新、冻结、解冻、吊销）
	// **参数说明**
	//   - payer: gas扣费代扣人
	//   - 其他同SendContractManageRequest
	SendContractManageRequestWithPayer(payload *common.Payload, endorsers []*common.EndorsementEntry,
		payer *common.EndorsementEntry, timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 1.8 合约调用
	// **参数说明**
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - txId: 交易ID
	//           格式要求：长度为64字节，字符在a-z0-9
	//           可为空，若为空字符串，将自动生成txId
	//   - kvs: 合约参数
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	//   - limit: transaction limitation，执行交易时的资源消耗上限，设为nil则不设置上限
	// ```go
	InvokeContract(contractName, method, txId string, kvs []*common.KeyValuePair, timeout int64,
		withSyncResult bool) (*common.TxResponse, error)
	InvokeContractWithLimit(contractName, method, txId string, kvs []*common.KeyValuePair, timeout int64,
		withSyncResult bool, limit *common.Limit) (*common.TxResponse, error)
	// ```

	// ### 1.9 合约查询接口调用
	// **参数说明**
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - kvs: 合约参数
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	// ```go
	QueryContract(contractName, method string, kvs []*common.KeyValuePair, timeout int64) (*common.TxResponse, error)
	// ```

	// ### 1.10 构造待发送交易体
	// **参数说明**
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - txId: 交易ID
	//           格式要求：长度为64字节，字符在a-z0-9
	//           可为空，若为空字符串，将自动生成txId
	//   - kvs: 合约参数
	// ```go
	GetTxRequest(contractName, method, txId string, kvs []*common.KeyValuePair) (*common.TxRequest, error)
	// ```

	// ### 1.11 发送已构造好的交易体
	// **参数说明**
	//   - txRequest: 已构造好的交易体
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	SendTxRequest(txRequest *common.TxRequest, timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ## 2 系统合约接口
	// ### 2.1 根据交易Id查询交易
	// **参数说明**
	//   - txId: 交易ID
	// ```go
	GetTxByTxId(txId string) (*common.TransactionInfo, error)
	// ```

	// ### 2.2 根据交易Id查询包含rwset的交易
	// **参数说明**
	//   - txId: 交易ID
	// ```go
	GetTxWithRWSetByTxId(txId string) (*common.TransactionInfoWithRWSet, error)
	// ```

	// ### 2.3 根据区块高度查询区块
	// **参数说明**
	//   - blockHeight: 指定区块高度，若为-1，将返回最新区块
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByHeight(blockHeight uint64, withRWSet bool) (*common.BlockInfo, error)
	// ```

	//Deprecated, 请使用GetBlockByHeight
	//### 2.4 根据区块高度查询完整区块
	// **参数说明**
	//   - blockHeight: 指定区块高度，若为-1，将返回最新区块
	// ```go
	GetFullBlockByHeight(blockHeight uint64) (*store.BlockWithRWSet, error)
	// ```

	// ### 2.5 根据区块哈希查询区块
	// **参数说明**
	//   - blockHash: 指定区块Hash
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.6 根据交易Id查询区块
	// **参数说明**
	//   - txId: 交易ID
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.7 查询最新的配置块
	// **参数说明**
	//   - withRWSet: 是否返回读写集
	// ```go
	GetLastConfigBlock(withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.8 查询最新区块
	// **参数说明**
	//   - withRWSet: 是否返回读写集
	// ```go
	GetLastBlock(withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.9 查询节点加入的链信息
	//    - 返回ChainId清单
	// ```go
	GetNodeChainList() (*discovery.ChainList, error)
	// ```

	// ### 2.10 查询链信息
	//   - 包括：当前链最新高度，链节点信息
	// ```go
	GetChainInfo() (*discovery.ChainInfo, error)
	// ```

	// ### 2.11 根据交易Id获取区块高度
	// **参数说明**
	//   - txId: 交易ID
	// ```go
	GetBlockHeightByTxId(txId string) (uint64, error)
	// ```

	// ### 2.12 根据区块Hash获取区块高度
	// **参数说明**
	//   - blockHash: 指定区块Hash
	// ```go
	GetBlockHeightByHash(blockHash string) (uint64, error)
	// ```

	// ### 2.13 查询当前最新区块高度
	// ```go
	GetCurrentBlockHeight() (uint64, error)
	// ```

	// ### 2.14 根据区块高度查询区块头
	// **参数说明**
	//   - blockHeight: 指定区块高度，若为-1，将返回最新区块头
	// ```go
	GetBlockHeaderByHeight(blockHeight uint64) (*common.BlockHeader, error)
	// ```

	// ### 2.15 系统合约调用
	// **参数说明**
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - txId: 交易ID
	//           格式要求：长度为64字节，字符在a-z0-9
	//           可为空，若为空字符串，将自动生成txId
	//   - kvs: 合约参数
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	InvokeSystemContract(contractName, method, txId string, kvs []*common.KeyValuePair, timeout int64,
		withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 2.16 系统合约查询接口调用
	// **参数说明**
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - kvs: 合约参数
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	// ```go
	QuerySystemContract(contractName, method string, kvs []*common.KeyValuePair, timeout int64) (*common.TxResponse, error)
	// ```

	// ### 2.17 根据交易Id获取Merkle路径
	// **参数说明**
	//   - txId: 交易ID
	// ```go
	GetMerklePathByTxId(txId string) ([]byte, error)
	// ```

	// ### 2.18 开放系统合约
	// **参数说明**
	//   - grantContractList: 需要开放的系统合约字符串数组
	// ```go
	CreateNativeContractAccessGrantPayload(grantContractList []string) (*common.Payload, error)
	// ```

	// ### 2.19 弃用系统合约
	// **参数说明**
	//   - revokeContractList: 需要弃用的系统合约字符串数组
	// ```go
	CreateNativeContractAccessRevokePayload(revokeContractList []string) (*common.Payload, error)
	// ```

	// ### 2.20 查询指定合约的信息，包括系统合约和用户合约
	// **参数说明**
	//   - contractName: 指定查询的合约名字，包括系统合约和用户合约
	// ```go
	GetContractInfo(contractName string) (*common.Contract, error)
	// ```

	// ### 2.21 查询所有的合约名单，包括系统合约和用户合约
	// **返回值说明**
	//   - []*common.Contract: 链上所有的合约列表，包括系统合约和用户合约
	// ```go
	GetContractList() ([]*common.Contract, error)
	// ```

	// ### 2.22 查询已禁用的系统合约名单
	// **返回值说明**
	//   - []string: 链上已禁用的系统合约名字列表
	// ```go
	GetDisabledNativeContractList() ([]string, error)
	// ```

	// ### 2.23 根据区块高度获得区块，应对超大block，可设置返回block的截断规则。
	// **参数说明**
	//   - blockHeight: 指定区块高度
	//   - withRWSet: 是否返回读写集
	//   - truncateLength 截断的长度设置，如果此参数为>0，超过此长度的剩余数据将被丢弃。如果<=0则不截断
	//   - truncateModel 截断的模式设置:hash,truncate,empty
	// ```go
	GetBlockByHeightTruncate(blockHeight uint64, withRWSet bool, truncateLength int,
		truncateModel string) (*common.BlockInfo, error)
	// ```

	// ## 3 链配置接口
	// ### 3.1 查询最新链配置
	// ```go
	GetChainConfig() (*config.ChainConfig, error)
	// ```

	// ### 3.2 根据指定区块高度查询最近链配置
	// **参数说明**
	//   - blockHeight: 指定区块高度
	//     如果当前区块就是配置块，直接返回当前区块的链配置
	// ```go
	GetChainConfigByBlockHeight(blockHeight uint64) (*config.ChainConfig, error)
	// ```

	// ### 3.3 查询最新链配置序号Sequence
	//   - 用于链配置更新
	// ```go
	GetChainConfigSequence() (uint64, error)
	// ```

	// ### 3.4 链配置更新获取Payload签名
	// **参数说明**
	//   - payload: 待签名payload
	// ```go
	SignChainConfigPayload(payload *common.Payload) (*common.EndorsementEntry, error)
	// ```

	// ### 3.5 发送链配置更新请求
	// **参数说明**
	//   - payload: 待签名payload
	//   - endorsers: 背书签名信息列表
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	SendChainConfigUpdateRequest(payload *common.Payload, endorsers []*common.EndorsementEntry, timeout int64,
		withSyncResult bool) (*common.TxResponse, error)
	// ```

	// > 以下CreateChainConfigXXXXXXPayload方法，用于生成链配置待签名payload，在进行多签收集后(需机构Admin权限账号签名)，用于链配置的更新

	// ### 3.6 更新Core模块待签名payload生成
	// **参数说明**
	//   - txSchedulerTimeout: 交易调度器从交易池拿到交易后, 进行调度的时间，其值范围为(0, 60], 若设置为0，则抛出错误
	//   - txSchedulerValidateTimeout: 交易调度器从区块中拿到交易后, 进行验证的超时时间，其值范围为(0, 60], 若设置为0，则抛出错误
	// ```go
	CreateChainConfigCoreUpdatePayload(txSchedulerTimeout, txSchedulerValidateTimeout uint64) (*common.Payload, error)
	// ```

	// ### 3.7 更新链配置的区块相关参数待签名payload生成
	// **参数说明**
	//   - txTimestampVerify: 是否需要开启交易时间戳校验
	//   - blockTimestampVerify: 是否需要开启区块头时间戳校验
	//   - txTimeout: 交易时间戳的过期时间(秒)，其值范围为[600, +∞)
	//   - blockTimeout: 区块头时间戳的过期时间(秒)，其值范围为[10, +∞)
	//   - blockTxCapacity: 区块中最大交易数，其值范围为(0, +∞]
	//   - blockSize: 区块最大限制，单位MB，其值范围为(0, +∞]
	//   - blockInterval: 出块间隔，单位:ms，其值范围为[10, +∞]
	//   - txParamterSize: 交易的参数的最大值限制，单位：MB，其值范围为[0,100]
	// ```go
	CreateChainConfigBlockUpdatePayload(txTimestampVerify, blockTimestampVerify bool, txTimeout,
		blockTimeout, blockTxCapacity, blockSize, blockInterval, txParamterSize uint32) (*common.Payload, error)
	// ```

	// ### 3.8 添加信任组织根证书待签名payload生成
	// **参数说明**
	//   - trustRootOrgId: 组织Id
	//   - trustRootCrt: 根证书
	// ```go
	CreateChainConfigTrustRootAddPayload(trustRootOrgId string, trustRootCrt []string) (*common.Payload, error)
	// ```

	// ### 3.9 更新信任组织根证书待签名payload生成
	// **参数说明**
	//   - trustRootOrgId: 组织Id
	//   - trustRootCrt: 根证书
	// ```go
	CreateChainConfigTrustRootUpdatePayload(trustRootOrgId string, trustRootCrt []string) (*common.Payload, error)
	// ```

	// ### 3.10 删除信任组织根证书待签名payload生成
	// **参数说明**
	//   - orgIdOrPKPubkeyPEM: cert模式下传组织Id，public模式下传公钥PEM
	// ```go
	CreateChainConfigTrustRootDeletePayload(orgIdOrPKPubkeyPEM string) (*common.Payload, error)
	// ```

	// ### 3.11 添加信任成员证书待签名payload生成
	// **参数说明**
	//   - trustMemberOrgId: 组织Id
	//   - trustMemberNodeId: 节点Id
	//   - trustMemberRole: 成员角色
	//   - trustMemberInfo: 成员信息内容
	// ```go
	CreateChainConfigTrustMemberAddPayload(trustMemberOrgId, trustMemberNodeId,
		trustMemberRole, trustMemberInfo string) (*common.Payload, error)
	// ```

	// ### 3.12 删除信任成员证书待签名payload生成
	// **参数说明**
	//   - trustMemberInfo: 成员信息内容
	// ```go
	CreateChainConfigTrustMemberDeletePayload(trustMemberInfo string) (*common.Payload, error)
	// ```

	// ### 3.13 添加权限配置待签名payload生成
	// **参数说明**
	//   - permissionResourceName: 权限名
	//   - policy: 权限规则
	// ```go
	CreateChainConfigPermissionAddPayload(permissionResourceName string,
		policy *accesscontrol.Policy) (*common.Payload, error)
	// ```

	// ### 3.14 更新权限配置待签名payload生成
	// **参数说明**
	//   - permissionResourceName: 权限名
	//   - policy: 权限规则
	// ```go
	CreateChainConfigPermissionUpdatePayload(permissionResourceName string,
		policy *accesscontrol.Policy) (*common.Payload, error)
	// ```

	// ### 3.15 删除权限配置待签名payload生成
	// **参数说明**
	//   - permissionResourceName: 权限名
	// ```go
	CreateChainConfigPermissionDeletePayload(permissionResourceName string) (*common.Payload, error)
	// ```

	// ### 3.16 添加共识节点地址待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeIds: 节点Id
	// ```go
	CreateChainConfigConsensusNodeIdAddPayload(nodeOrgId string, nodeIds []string) (*common.Payload, error)
	// ```

	// ### 3.17 更新共识节点地址待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeOldNodeId: 节点原Id
	//   - nodeNewNodeId: 节点新Id
	// ```go
	CreateChainConfigConsensusNodeIdUpdatePayload(nodeOrgId, nodeOldNodeId, nodeNewNodeId string) (*common.Payload, error)
	// ```

	// ### 3.18 删除共识节点地址待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeId: 节点Id
	// ```go
	CreateChainConfigConsensusNodeIdDeletePayload(nodeOrgId, nodeId string) (*common.Payload, error)
	// ```

	// ### 3.19 添加共识节点待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeIds: 节点Id
	// ```go
	CreateChainConfigConsensusNodeOrgAddPayload(nodeOrgId string, nodeIds []string) (*common.Payload, error)
	// ```

	// ### 3.20 更新共识节点待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeIds: 节点Id
	// ```go
	CreateChainConfigConsensusNodeOrgUpdatePayload(nodeOrgId string, nodeIds []string) (*common.Payload, error)
	// ```

	// ### 3.21 删除共识节点待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	// ```go
	CreateChainConfigConsensusNodeOrgDeletePayload(nodeOrgId string) (*common.Payload, error)
	// ```

	// ### 3.22 添加共识扩展字段待签名payload生成
	// **参数说明**
	//   - kvs: 字段key、value对
	// ```go
	CreateChainConfigConsensusExtAddPayload(kvs []*common.KeyValuePair) (*common.Payload, error)
	// ```

	// ### 3.23 更新共识扩展字段待签名payload生成
	// **参数说明**
	//   - kvs: 字段key、value对
	// ```go
	CreateChainConfigConsensusExtUpdatePayload(kvs []*common.KeyValuePair) (*common.Payload, error)
	// ```

	// ### 3.24 删除共识扩展字段待签名payload生成
	// **参数说明**
	//   - keys: 待删除字段
	// ```go
	CreateChainConfigConsensusExtDeletePayload(keys []string) (*common.Payload, error)
	// ```

	// ### 3.25 修改地址类型payload生成
	// **参数说明**
	//   - addrType: 地址类型，0-ChainMaker; 1-ZXL
	// ```go
	CreateChainConfigAlterAddrTypePayload(addrType string) (*common.Payload, error)
	// ```

	// ### 3.26 启用或停用Gas计费开关payload生成
	// ```go
	CreateChainConfigEnableOrDisableGasPayload() (*common.Payload, error)
	// ```

	// ### 3.27 开启或关闭链配置的Gas优化payload生成
	// ```go
	CreateChainConfigOptimizeChargeGasPayload(enable bool) (*common.Payload, error)
	// ```

	// ### 3.28  查询最新权限配置列表
	// ```go
	GetChainConfigPermissionList() ([]*config.ResourcePolicy, error)
	// ```

	// ## 4 证书管理接口
	// ### 4.1 用户证书添加
	// **参数说明**
	//   - 在common.TxResponse.ContractResult.Result字段中返回成功添加的certHash
	// ```go
	AddCert() (*common.TxResponse, error)
	// ```

	// ### 4.2 用户证书删除
	// **参数说明**
	//   - certHashes: 证书Hash列表
	// ```go
	DeleteCert(certHashes []string) (*common.TxResponse, error)
	// ```

	// ### 4.3 用户证书查询
	// **参数说明**
	//   - certHashes: 证书Hash列表
	// **返回值说明**
	//   - *common.CertInfos: 包含证书Hash和证书内容的列表
	// ```go
	QueryCert(certHashes []string) (*common.CertInfos, error)
	// ```

	// ### 4.4 获取用户证书哈希
	// ```go
	GetCertHash() ([]byte, error)
	// ```

	// ### 4.5 生成证书管理操作Payload（三合一接口）
	// **参数说明**
	//   - method: CERTS_FROZEN(证书冻结)/CERTS_UNFROZEN(证书解冻)/CERTS_REVOCATION(证书吊销)
	//   - kvs: 证书管理操作参数
	// ```go
	CreateCertManagePayload(method string, kvs []*common.KeyValuePair) *common.Payload
	// ```

	// ### 4.6 生成证书冻结操作Payload
	// **参数说明**
	//   - certs: X509证书列表
	// ```go
	CreateCertManageFrozenPayload(certs []string) *common.Payload
	// ```

	// ### 4.7 生成证书解冻操作Payload
	// **参数说明**
	//   - certs: X509证书列表
	// ```go
	CreateCertManageUnfrozenPayload(certs []string) *common.Payload
	// ```

	// ### 4.8 生成证书吊销操作Payload
	// **参数说明**
	//   - certs: X509证书列表
	// ```go
	CreateCertManageRevocationPayload(certCrl string) *common.Payload
	// ```

	// ### 4.9 待签payload签名
	//  *一般需要使用具有管理员权限账号进行签名*
	// **参数说明**
	//   - payload: 待签名payload
	// ```go
	SignCertManagePayload(payload *common.Payload) (*common.EndorsementEntry, error)
	// ```

	// ### 4.10 发送证书管理请求（证书冻结、解冻、吊销）
	// **参数说明**
	//   - payload: 交易payload
	//   - endorsers: 背书签名信息列表
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	SendCertManageRequest(payload *common.Payload, endorsers []*common.EndorsementEntry, timeout int64,
		withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ## 5 消息订阅接口
	// ### 5.1 区块订阅
	// **参数说明**
	//   - startBlock: 订阅起始区块高度，包含起始区块。若为-1，表示订阅实时最新区块
	//   - endBlock: 订阅结束区块高度，包含结束区块。若为-1，表示订阅实时最新区块
	//   - withRwSet: 是否返回读写集
	//   - onlyHeader: 若设置为true，将忽略withRwSet选项，仅返回区块头（common.BlockHeader）,若设置为false，将返回common.BlockInfo
	// ```go
	SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRWSet, onlyHeader bool) (<-chan interface{}, error)
	// ```

	// ### 5.2 交易订阅
	// **参数说明**
	//   - startBlock: 订阅起始区块高度，包含起始区块。若为-1，表示订阅实时最新区块
	//   - endBlock: 订阅结束区块高度，包含结束区块。若为-1，表示订阅实时最新区块
	//   - contractName ：指定订阅指定合约的交易，可以传用户合约名称或系统合约名称，若为空，表示订阅所有合约的交易
	//   - txIds: 订阅txId列表，若为空，表示订阅所有txId
	// ```go
	SubscribeTx(ctx context.Context, startBlock, endBlock int64, contractName string,
		txIds []string) (<-chan interface{}, error)
	// ```

	// ### 5.3 合约事件订阅
	// **参数说明**
	//   - startBlock: 订阅起始区块高度，包含起始区块。若为-1，表示订阅实时最新区块
	//   - endBlock: 订阅结束区块高度，包含结束区块。若为-1，表示订阅实时最新区块
	//   - contractName ：指定订阅的合约名称
	//   - topic ：指定订阅主题
	// ```go
	SubscribeContractEvent(ctx context.Context, startBlock, endBlock int64, contractName,
		topic string) (<-chan interface{}, error)
	// ```

	// ### 5.4 多合一订阅
	// **参数说明**
	//   - txType: 订阅交易类型，目前已支持：区块消息订阅(common.TxType_SUBSCRIBE_BLOCK_INFO)、交易消息订阅(common.TxType_SUBSCRIBE_TX_INFO)
	//   - payloadBytes: 消息订阅参数payload
	// ```go
	Subscribe(ctx context.Context, payloadBytes *common.Payload) (<-chan interface{}, error)
	// ```

	//### 5.5 根据别名前缀订阅交易(仅支持证书模式链)
	//**参数说明**
	//- startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
	//- endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
	//- aliasPrefix: 证书别名前缀，不允许为空
	// ```go
	SubscribeTxByPreAlias(ctx context.Context, startBlock, endBlock int64,
		aliasPrefix string) (<-chan interface{}, error)
	//```

	//### 5.6 根据交易id前缀订阅交易
	//**参数说明**
	//- startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
	//- endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
	//- txIdPrefix: 交易id前缀，不允许为空
	//```go
	SubscribeTxByPreTxId(ctx context.Context, startBlock, endBlock int64,
		txIdPrefix string) (<-chan interface{}, error)
	//```

	//### 5.7 根据组织id前缀订阅交易
	//**参数说明**
	//- startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
	//- endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
	//- orgIdPrefix: 组织id前缀，不允许为空
	//```go
	SubscribeTxByPreOrgId(ctx context.Context, startBlock, endBlock int64,
		orgIdPrefix string) (<-chan interface{}, error)
	//```

	// ## 6 证书压缩
	// *开启证书压缩可以减小交易包大小，提升处理性能*
	// ### 6.1 启用压缩证书功能
	// ```go
	EnableCertHash() error
	// ```

	// ### 6.2 停用压缩证书功能
	// ```go
	DisableCertHash() error
	// ```

	// ## 7 层级属性加密类接口
	// > 注意：层级属性加密模块 `Id` 使用 `/` 作为分隔符，例如： Org1/Ou1/Member1
	// ### 7.1 生成层级属性参数初始化交易 payload
	// **参数说明**
	//   - orgId: 参与方组织 id
	//   - hibeParams: 传入序列化后的hibeParams byte数组
	// ```go
	CreateHibeInitParamsTxPayloadParams(orgId string, hibeParams []byte) ([]*common.KeyValuePair, error)
	// ```

	// ### 7.2 生成层级属性加密交易 payload，加密参数已知
	// **参数说明**
	//   - plaintext: 待加密交易消息明文
	//   - receiverIds: 消息接收者 id 列表，需和 paramsList 一一对应
	//   - paramsBytesList: 消息接收者对应的加密参数，需和 receiverIds 一一对应
	//   - txId: 以交易 Id 作为链上存储 hibeMsg 的 Key, 如果不提供存储的信息可能被覆盖
	//   - keyType: 对明文进行对称加密的方法，请传入 common 中 crypto 包提供的方法，目前提供AES和SM4两种方法
	// ```go
	CreateHibeTxPayloadParamsWithHibeParams(plaintext []byte, receiverIds []string, paramsBytesList [][]byte, txId string,
		keyType crypto.KeyType) ([]*common.KeyValuePair, error)
	// ```

	// ### 7.3 生成层级属性加密交易 payload，参数由链上查询得出
	// **参数说明**
	//   - contractName: 合约名
	//   - queryParamsMethod: 链上查询 hibe.Params 的合约方法
	//   - plaintext: 待加密交易消息明文
	//   - receiverIds: 消息接收者 id 列表，需和 paramsList 一一对应
	//   - paramsList: 消息接收者对应的加密参数，需和 receiverIds 一一对应
	//   - receiverOrgIds: 链上查询 hibe Params 的 Key 列表，需要和 receiverIds 一一对应
	//   - txId: 以交易 Id 作为链上存储 hibeMsg 的 Key, 如果不提供存储的信息可能被覆盖
	//   - keyType: 对明文进行对称加密的方法，请传入 common 中 crypto 包提供的方法，目前提供AES和SM4两种方法
	//   - timeout: （内部查询 HibeParams 的）超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	// ```go
	CreateHibeTxPayloadParamsWithoutHibeParams(contractName, queryParamsMethod string, plaintext []byte,
		receiverIds []string, receiverOrgIds []string, txId string, keyType crypto.KeyType,
		timeout int64) ([]*common.KeyValuePair, error)
	// ```

	// ### 7.4 查询某一组织的加密公共参数，返回其序列化后的byte数组
	// **参数说明**
	//   - contractName: 合约名
	//   - method: 查询的合约方法名
	//   - orgId: 参与方 id
	//   - timeout: 查询超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	// ```go
	QueryHibeParamsWithOrgId(contractName, method, orgId string, timeout int64) ([]byte, error)
	// ```

	// ### 7.5 已知交易id，根据私钥解密密文交易
	// **参数说明**
	//   - localId: 本地层级属性加密 id
	//   - hibeParams: hibeParams 序列化后的byte数组
	//   - hibePrivKey: hibe私钥序列化后的byte数组
	//   - txId: 层级属性加密交易 id
	//   - keyType: 对加密信息进行对称解密的方法，请和加密时使用的方法保持一致，请传入 common 中 crypto 包提供的方法，目前提供AES和SM4两种方法
	// ```go
	DecryptHibeTxByTxId(localId string, hibeParams []byte, hibePrvKey []byte, txId string,
		keyType crypto.KeyType) ([]byte, error)
	// ```

	// ## 8 数据归档接口
	// **（注意：请使用归档工具cmc进行归档操作，以下接口是归档原子接口，并不包括归档完整流程）**
	// ### 8.1 获取已归档区块高度，该方法已弃用，请优先使用GetArchiveStatus
	// **参数说明**
	//   - 输出已归档的区块高度
	// ```go
	//Deprecated
	GetArchivedBlockHeight() (uint64, error)
	// ```
	// ### 8.2 获取节点归档状态详细信息
	// **参数说明**
	//   - 输出：输出节点支持的归档类型：ExtraDB(外部数据库)，ArchiveCenter(归档中心)
	//   - 输出：已归档的区块高度
	//   - 输出：节点允许的归档高度: 高度必须小于(CurrBlockHeight-UnarchivedHeight)，UnarchivedHeight在chainmaker.yml storage中设置
	//   - 输出：当前节点的归档状态：Normal(正常), Archiving(正在归档), Restoring(正在恢复)
	//   - 输出：当前节点上可恢复区块的暂存文件信息
	// ```go
	GetArchiveStatus() (*store.ArchiveStatus, error)

	// ```
	// ### 8.3 构造数据归档区块Payload
	// **参数说明**
	//   - targetBlockHeight: 归档目标区块高度，节点ArchiveStatus.Process为Normal时才可发送成功
	// ```go
	CreateArchiveBlockPayload(targetBlockHeight uint64) (*common.Payload, error)
	// ```

	// ### 8.4 构造归档数据恢复Payload
	// **参数说明**
	//   - fullBlock: 完整区块数据（对应结构：store.BlockWithRWSet）
	// ```go
	CreateRestoreBlockPayload(fullBlock []byte) (*common.Payload, error)
	// ```

	// ### 8.6 发送归档请求
	// **参数说明**
	//   - payload: 指向payload对象的指针
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//	 - 备注：必须在ArchiveStatus.Process归档状态为Normal时才可调用成功
	// ```go
	SendArchiveBlockRequest(payload *common.Payload, timeout int64) (*common.TxResponse, error)
	// ```

	// ### 8.7 归档数据恢复
	// **参数说明**
	//   - payload: 指向payload对象的指针
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//	 - 备注：必须在ArchiveStatus.Process归档状态为Normal时才可调用成功
	// ```go
	SendRestoreBlockRequest(payload *common.Payload, timeout int64) (*common.TxResponse, error)
	// ```

	// ### 8.8 根据交易Id查询已归档交易
	// **参数说明**
	//   - txId: 交易ID
	// ```go
	GetArchivedTxByTxId(txId string) (*common.TransactionInfo, error)
	// ```

	// ### 8.9 根据区块高度查询已归档区块
	// **参数说明**
	//   - blockHeight: 指定区块高度
	//   - withRWSet: 是否返回读写集
	// ```go
	GetArchivedBlockByHeight(blockHeight uint64, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 8.11 根据区块哈希查询已归档区块
	// **参数说明**
	//   - blockHash: 指定区块Hash
	//   - withRWSet: 是否返回读写集
	// ```go
	GetArchivedBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 8.12 根据交易Id查询已归档区块
	// **参数说明**
	//   - txId: 交易ID
	//   - withRWSet: 是否返回读写集
	// ```go
	GetArchivedBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ## 9 隐私计算系统合约接口
	// ### 9.1 保存隐私合约计算结果，包括合约部署
	// **参数说明**
	//   - contractName: 合约名称
	//   - contractVersion: 合约版本号
	//   - isDeployment: 是否是部署合约
	//   - codeHash: 合约字节码hash值
	//   - reportHash: Enclave report hash值
	//   - result: 隐私合约执行结果
	//   - codeHeader: solodity合部署合约时合约字节码的header数据
	//   - txId: 交易Id
	//   - rwSet: 隐私合约执行产生的读写集
	//   - sign: Enclave对执行结果数据的结果签名
	//   - events: 合约执行产生的事件
	//   - privateReq: 用户调用隐私计算请求时的request序列化字节数组
	//   - withSyncResult: 是否同步返回调用结果
	//   - timeout: 发送交易的超时时间
	// ```go
	SaveData(contractName string, contractVersion string, isDeployment bool, codeHash []byte, reportHash []byte,
		result *common.ContractResult, codeHeader []byte, txId string, rwSet *common.TxRWSet, sign []byte,
		events *common.StrSlice, privateReq []byte, withSyncResult bool, timeout int64) (*common.TxResponse, error)
	// ```

	// ### 9.2 保存远程证明
	// **参数说明**
	//   - proof: 远程证明
	//   - txId: 交易Id
	//   - withSyncResult: 是否同步返回调用结果
	//   - timeout: 交易发送超时时间
	// ```go
	SaveRemoteAttestationProof(proof, txId string, withSyncResult bool, timeout int64) (*common.TxResponse, error)
	// ```

	// ### 9.3 构建上传Enclave CA证书的报文
	// **参数说明**
	//   - caCert: Enclave CA证书
	//   - txId: 交易Id
	// ```go
	CreateSaveEnclaveCACertPayload(caCert, txId string) (*common.Payload, error)
	// ```

	// ### 9.4 获取Enclave CA证书
	// ```go
	GetEnclaveCACert() ([]byte, error)
	// ```

	// ### 9.5 隐私计算调用者权限验证
	// **参数说明**
	//   - payload: 用户签名验证的payload内容
	//   - orgIds: 组织Id的slice，注意和signPairs里面SignInfo的证书顺序一致
	//   - signPairs: 用户多签的签名和证书slice
	// ```go
	CheckCallerCertAuth(payload string, orgIds []string, signPairs []*syscontract.SignInfo) (*common.TxResponse, error)
	// ```

	// ### 9.6 获取Enclave的report
	// **参数说明**
	//   - enclaveId: Enclave的Id，当前固定为"global_enclave_id"
	// ```go
	GetEnclaveReport(enclaveId string) ([]byte, error)
	// ```

	// ### 9.7 获取隐私证明材料
	// **参数说明**
	//   - enclaveId: Enclave的Id，当前固定为"global_enclave_id"
	// ```go
	GetEnclaveProof(enclaveId string) ([]byte, error)
	// ```

	// ### 9.8 获取隐私合约计算结果
	// **参数说明**
	//   - contractName: 合约名称
	//   - key: 计算结果对应的键值
	// ```go
	GetData(contractName, key string) ([]byte, error)
	// ```

	// ### 9.9 保存隐私目录
	// **参数说明**
	//   - orderId: 隐私目录的主键，供以后查询使用
	//   - txId: 交易ID
	//   - privateDir:
	//   - withSyncResult: 是否同步等待交易结果
	//   - timeout: 等待交易结果的超时时间
	// ```go
	SaveDir(orderId, txId string, privateDir *common.StrSlice, withSyncResult bool,
		timeout int64) (*common.TxResponse, error)
	// ```

	// ### 9.10 获取用户部署的隐私合约
	// **参数说明**
	//   - contractName: 合约名称
	//   - codeHash: 代码哈希
	// ```go
	GetContract(contractName, codeHash string) (*common.PrivateGetContract, error)
	// ```

	// ### 9.11 获取用户的隐私目录
	// **参数说明**
	//   - orderId: 隐私目录的主键
	// ```go
	GetDir(orderId string) ([]byte, error)
	// ```

	// ### 9.12 构建上传隐私计算环境的report的报文
	// **参数说明**
	//   - enclaveId: 隐私计算环境的标识
	//   - report: 隐私计算环境的report
	//   - txId: 交易ID
	// ```go
	CreateSaveEnclaveReportPayload(enclaveId, report, txId string) (*common.Payload, error)
	// ```

	// ### 9.13 获取隐私计算环境的加密公钥
	// **参数说明**
	//   - enclaveId: 隐私计算环境的标识
	// ```go
	GetEnclaveEncryptPubKey(enclaveId string) ([]byte, error)
	// ```

	// ### 9.14 获取隐私计算环境的验签公钥
	// **参数说明**
	//   - enclaveId: 隐私计算环境的标识
	// ```go
	GetEnclaveVerificationPubKey(enclaveId string) ([]byte, error)
	// ```

	// ### 9.15 获取隐私证明材料中的Challenge
	// **参数说明**
	//   - enclaveId: 隐私计算环境的标识
	// ```go
	GetEnclaveChallenge(enclaveId string) ([]byte, error)
	// ```

	// ### 9.16 获取隐私证明材料中的Signature
	// **参数说明**
	//   - enclaveId: 隐私计算环境的标识
	// ```go
	GetEnclaveSignature(enclaveId string) ([]byte, error)
	// ```

	// ## 10 系统类接口
	// ### 10.1 SDK停止接口
	// *关闭连接池连接，释放资源*
	// ```go
	Stop() error
	// ```

	// ### 10.2 获取链版本
	// ```go
	GetChainMakerServerVersion() (string, error)

	// ### 10.2 获取链版本
	// ```go
	GetChainMakerServerVersionCustom(ctx context.Context) (string, error)
	// ```

	// ## 11 公钥身份类接口
	// ### 11.1 构造添加公钥身份请求
	// **参数说明**
	//   - pubkey: 公钥信息
	//   - orgId: 组织id
	//   - role:   角色，支持client,light,common
	// ```go
	CreatePubkeyAddPayload(pubkey string, orgId string, role string) (*common.Payload, error)
	// ```

	// ### 11.2 构造删除公钥身份请求
	// **参数说明**
	//   - pubkey: 公钥信息
	//   - orgId: 组织id
	// ```go
	CreatePubkeyDelPayload(pubkey string, orgId string) (*common.Payload, error)
	// ```

	// ### 11.3 构造查询公钥身份请求
	// **参数说明**
	//   - pubkey: 公钥信息
	// ```go
	CreatePubkeyQueryPayload(pubkey string) (*common.Payload, error)
	// ```

	// ### 11.4 发送公钥身份管理请求（添加、删除）
	// **参数说明**
	//   - payload: 交易payload
	//   - endorsers: 背书签名信息列表
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	SendPubkeyManageRequest(payload *common.Payload, endorsers []*common.EndorsementEntry, timeout int64,
		withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ## 12 多签类接口
	// ### 12.1 发起多签请求
	// **参数说明**
	//   - payload: 待签名payload
	//   - endorsers: 背书签名信息列表
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	MultiSignContractReq(payload *common.Payload, endorsers []*common.EndorsementEntry, timeout int64,
		withSyncResult bool) (*common.TxResponse, error)
	// **参数说明**
	// payer: 指定gas代扣账户
	MultiSignContractReqWithPayer(payload *common.Payload, endorsers []*common.EndorsementEntry,
		payer *common.EndorsementEntry, timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 12.2 发起多签投票
	// **参数说明**
	//   - payload: 待签名payload
	//   - endorser: 投票人对多签请求 payload 的签名信息
	//   - isAgree: 投票人对多签请求是否同意，true为同意，false则反对
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	MultiSignContractVote(payload *common.Payload, endorser *common.EndorsementEntry, isAgree bool,
		timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 12.3 根据txId查询多签状态
	// **参数说明**
	//   - txId: 需要查询的多签请求交易Id
	// ```go
	MultiSignContractQuery(txId string) (*common.TxResponse, error)
	// ```

	// ### 12.4 触发执行多签请求
	// **参数说明**
	//   - payload: 待签名payload
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//	 - limit: 本次执行多签请求支付的 gas 上限
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	MultiSignContractTrig(multiSignReqPayload *common.Payload,
		timeout int64, limit *common.Limit, withSyncResult bool) (*common.TxResponse, error)
	// **参数说明**
	// payer: 指定gas代扣账户
	MultiSignContractTrigWithPayer(multiSignReqPayload *common.Payload, payer *common.EndorsementEntry,
		timeout int64, limit *common.Limit, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 12.5 根据发起多签请求所需的参数构建payload
	// **参数说明**
	//   - pairs: 发起多签请求所需的参数
	// ```go
	CreateMultiSignReqPayload(pairs []*common.KeyValuePair) *common.Payload
	// ```

	// ### 12.6 根据发起多签请求所需的参数构建payload
	// **参数说明**
	//   - pairs: 发起多签请求所需的参数
	//   - gasLimit: 本次交易允许消费的最大gas, 0代表无限制
	// ```go
	CreateMultiSignReqPayloadWithGasLimit(pairs []*common.KeyValuePair, gasLimit uint64) *common.Payload
	// ```

	// ### 12.7 发起多签投票
	// **参数说明**
	//   - payload: 待签名payload
	//   - endorser: 投票人对多签请求 payload 的签名信息
	//   - isAgree: 投票人对多签请求是否同意，true为同意，false则反对
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - gasLimit: 交易允许消费的最大 gas, 0 代表无限制
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	MultiSignContractVoteWithGasLimit(payload *common.Payload, endorser *common.EndorsementEntry, isAgree bool,
		timeout int64, gasLimit uint64, withSyncResult bool) (*common.TxResponse, error)
	// **参数说明**
	// payer: 指定gas代扣账户
	MultiSignContractVoteWithGasLimitAndPayer(payload *common.Payload, endorser *common.EndorsementEntry,
		payer *common.EndorsementEntry, isAgree bool, timeout int64, gasLimit uint64,
		withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 12.8 根据txId查询多签状态
	// **参数说明**
	//   - txId: 需要查询的多签请求交易Id
	//   - params: 格式化参数
	// ```go
	MultiSignContractQueryWithParams(txId string, params []*common.KeyValuePair) (*common.TxResponse, error)
	// ```

	// ## 13 gas管理相关接口
	// ### 13.1 构造设置gas管理员payload
	// **参数说明**
	//   - address: gas管理员的地址
	// ```go
	CreateSetGasAdminPayload(address string) (*common.Payload, error)
	// ```

	// ### 13.2 查询gas管理员
	// **返回值说明**
	//   - string: gas管理员的账号地址
	// ```go
	GetGasAdmin() (string, error)
	// ```

	// ### 13.3 构造 充值gas账户 payload
	// **参数说明**
	//   - rechargeGasList: 一个或多个gas账户充值指定gas数量
	// ```go
	CreateRechargeGasPayload(rechargeGasList []*syscontract.RechargeGas) (*common.Payload, error)
	// ```

	// ### 13.4 查询gas账户余额
	// **参数说明**
	//   - address: 查询gas余额的账户地址
	// ```go
	GetGasBalance(address string) (int64, error)
	// ```

	// ### 13.5 构造 退还gas账户的gas payload
	// **参数说明**
	//   - address: 退还gas的账户地址
	//   - amount: 退还gas的数量
	// ```go
	CreateRefundGasPayload(address string, amount int64) (*common.Payload, error)
	// ```

	// ### 13.6 构造 冻结指定gas账户 payload
	// **参数说明**
	//   - address: 冻结指定gas账户的账户地址
	// ```go
	CreateFrozenGasAccountPayload(address string) (*common.Payload, error)
	// ```

	// ### 13.7 构造 解冻指定gas账户 payload
	// **参数说明**
	//   - address: 解冻指定gas账户的账户地址
	// ```go
	CreateUnfrozenGasAccountPayload(address string) (*common.Payload, error)
	// ```

	// ### 13.8 查询gas账户的状态
	// **参数说明**
	//   - address: 解冻指定gas账户的账户地址
	// **返回值说明**
	//   - bool: true表示账号未被冻结，false表示账号已被冻结
	// ```go
	GetGasAccountStatus(address string) (bool, error)
	// ```

	// ### 13.9 发送gas管理类请求
	// **参数说明**
	//   - payload: 交易payload
	//   - endorsers: 背书签名信息列表
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	SendGasManageRequest(payload *common.Payload, endorsers []*common.EndorsementEntry, timeout int64,
		withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 13.10 为payload添加gas limit
	// **参数说明**
	//   - payload: 交易payload
	//   - limit: transaction limitation，执行交易时的资源消耗上限
	// ```go
	AttachGasLimit(payload *common.Payload, limit *common.Limit) *common.Payload
	// ```

	// ### 13.11 估算交易的gas消耗量
	// **参数说明**
	//   - payload: 待估算gas消耗量的交易payload
	// **返回值说明**
	//   - uint64: 估算出的gas消耗量
	// ```go
	EstimateGas(payload *common.Payload) (uint64, error)
	// ```

	// ### 13.12 构造 配置账户基础gas消耗数量 payload
	// **参数说明**
	//   - amount: 调用合约消耗的基础的 gas 数量
	// ```go
	CreateSetInvokeBaseGasPayload(invokeBaseGas int64) (*common.Payload, error)
	// ```

	// ### 13.13 设置链配置的 default gas_price 参数
	// **参数说明**
	//   - invokeGasPrice: 调用合约中所带参数消耗的每个字节的 gas 价格，支持格式为10进制整数｜小数
	// ```go
	CreateSetInvokeGasPricePayload(invokeGasPrice string) (*common.Payload, error)
	// ```

	// ### 13.14 设置链配置的 install_base gas 参数
	// **参数说明**
	//   - amount: 安装｜升级合约消耗的基础的 gas 数量
	// ```go
	CreateSetInstallBaseGasPayload(installBaseGas int64) (*common.Payload, error)
	// ```

	// ### 13.15 设置链配置的 install gas_price 参数
	// **参数说明**
	//   - invokeGasPrice: 安装｜升级合约中所带参数消耗的每个字节的 gas 价格，支持格式为10进制整数｜小数
	// ```go
	CreateSetInstallGasPricePayload(installGasPrice string) (*common.Payload, error)
	// ```

	// ## 14 别名相关接口
	// ### 14.1 添加别名
	// ```go
	AddAlias() (*common.TxResponse, error)
	// ```

	// ### 14.2 构造`更新别名的证书`payload
	// **参数说明**
	//   - alias: 带更新证书的别名
	//   - newCertPEM: 新的证书，此新证书将替换掉alias关联的证书
	// ```go
	CreateUpdateCertByAliasPayload(alias, newCertPEM string) *common.Payload
	// ```

	// ### 14.3 签名`更新别名的证书`payload
	// **参数说明**
	//   - payload: 交易payload
	// ```go
	SignUpdateCertByAliasPayload(payload *common.Payload) (*common.EndorsementEntry, error)
	// ```

	// ### 14.4 发起`更新别名的证书`交易
	// **参数说明**
	//   - payload: 交易payload
	//   - endorsers: 背书签名信息列表
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	UpdateCertByAlias(payload *common.Payload, endorsers []*common.EndorsementEntry,
		timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 14.5 查询别名详细信息
	// **参数说明**
	//   - aliases: 带查询的证书别名切片，根据这些别名查询返回AliasInfos
	// ```go
	QueryCertsAlias(aliases []string) (*common.AliasInfos, error)
	// ```

	// ### 14.6 构造`删除别名`payload
	// **参数说明**
	//   - aliases: 带删除的证书别名切片
	// ```go
	CreateDeleteCertsAliasPayload(aliases []string) *common.Payload
	// ```

	// ### 14.7 签名`删除别名`payload
	// **参数说明**
	//   - payload: 交易payload
	// ```go
	SignDeleteAliasPayload(payload *common.Payload) (*common.EndorsementEntry, error)
	// ```

	// ### 14.8 发起`删除别名`交易
	// **参数说明**
	//   - payload: 交易payload
	//   - endorsers: 背书签名信息列表
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult为空，可以通过common.TxResponse.TxId查询交易结果
	// ```go
	DeleteCertsAlias(payload *common.Payload, endorsers []*common.EndorsementEntry,
		timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ## 15 交易池相关接口
	// ### 15.1 获取交易池状态
	// ```go
	GetPoolStatus() (*txpool.TxPoolStatus, error)
	// ```

	// ### 15.2 获取不同交易类型和阶段中的交易Id列表。
	// **参数说明**
	//   - txType: 交易类型 在pb的txpool包中进行了定义
	//   - txStage: 交易阶段 在pb的txpool包中进行了定义
	// **返回值说明**
	//   - []string: 交易Id列表
	// ```go
	GetTxIdsByTypeAndStage(txType txpool.TxType, txStage txpool.TxStage) ([]string, error)
	// ```

	// ### 15.3 根据txIds获取交易池中存在的txs，并返回交易池缺失的tx的txIds
	// **参数说明**
	//   - txIds: 交易Id列表
	// **返回值说明**
	//   - []*common.Transaction: 交易池中存在的txs
	//   - []string: 交易池缺失的tx的txIds
	// ```go
	GetTxsInPoolByTxIds(txIds []string) ([]*common.Transaction, []string, error)

	// ```

	// ## 16 Payer相关接口
	// ### 16.1 设置Payer
	// ```go
	SetContractMethodPayer(
		payerAddress string, contractName string, methodName string, requestId string,
		payerOrgId string, payerKeyPem []byte, payerCertPem []byte,
		gasLimit uint64) (
		*common.TxResponse, error)
	// ```

	// ### 16.2 取消Payer
	// ```go
	UnsetContractMethodPayer(
		contractName string, methodName string, gasLimit uint64) (*common.TxResponse, error)
	// ```

	// ### 16.3 查询Payer
	// ```go
	QueryContractMethodPayer(
		contractName string, methodName string, gasLimit uint64) (*common.TxResponse, error)
	// ```

	// ### 16.3 查询交易Payer
	// ```go
	QueryTxPayer(tsId string, gasLimit uint64) (*common.TxResponse, error)
	// ```

	// ArchiveBlocks 归档指定区块高度范围的区块到归档服务
	// @param beginHeight
	// @param endHeight
	// @param mode
	// @param heightNoticeCallback 归档进度通知
	// @return error
	ArchiveBlocks(archiveHeight uint64, mode string, heightNoticeCallback func(ProcessMessage) error) error

	// RestoreBlocks 从归档服务中恢复到指定高度的区块到节点中
	// @param restoreHeight 指定恢复高度
	// @param heightNoticeCallback 恢复进度通知
	// @return error
	RestoreBlocks(restoreHeight uint64, mode string, heightNoticeCallback func(ProcessMessage) error) error

	// ## 17 共识相关接口
	// ### 17.1 获取所有共识节点的身份标识
	// **返回值说明**
	//   - []string: 共识节点身份标识
	//   - err: 当查询的节点非共识节点时或共识节点内部查询中出现错误，返回error
	// ```go
	GetConsensusValidators() ([]string, error)
	// ```

	// ### 17.2 获取节点正在共识的区块高度
	// **返回值说明**
	//   - uint64: 节点正在共识的区块高度
	//	 - err: 当查询的节点非共识节点时或共识节点内部查询中出现错误，返回error
	// ```go
	GetConsensusHeight() (uint64, error)
	// ```

	// ### 17.3 获取共识节点的状态
	// **返回值说明**
	//   - []byte: 查询的共识节点状态
	//	 - err: 当查询的节点非共识节点时或共识节点内部查询中出现错误，返回error
	// ```go
	GetConsensusStateJSON() ([]byte, error)
	// ```

	// ### 17.4 获取节点同步的状态
	// **参数说明**
	//   - withOthersState: 是否连带获取其他节点的高度状态
	// **返回值说明**
	//   - *sync.SyncState,: 查询的节点同步的状态
	//	 - err: 异常，返回error
	// ```go
	GetSyncState(withOthersState bool) (*sync.SyncState, error)
	// ```
}

// ProcessMessage 处理过程通知消息
type ProcessMessage struct {
	CurrentHeight uint64
	Total         uint64
	Error         error
}

// ArchiveService 归档服务
type ArchiveService interface {

	// ## 2 系统合约接口
	// ### 2.1 根据交易Id查询交易
	// **参数说明**
	//   - txId: 交易ID
	// ```go
	GetTxByTxId(txId string) (*common.TransactionInfo, error)
	// ```

	// ### 2.2 根据交易Id查询包含rwset的交易
	// **参数说明**
	//   - txId: 交易ID
	// ```go
	GetTxWithRWSetByTxId(txId string) (*common.TransactionInfoWithRWSet, error)
	// ```

	// ### 2.3 根据区块高度查询区块
	// **参数说明**
	//   - blockHeight: 指定区块高度，若为-1，将返回最新区块
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByHeight(blockHeight uint64, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.5 根据区块哈希查询区块
	// **参数说明**
	//   - blockHash: 指定区块Hash
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.6 根据交易Id查询区块
	// **参数说明**
	//   - txId: 交易ID
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error)

	// GetChainConfigByBlockHeight 获取指定高度的区块的配置
	// @param blockHeight
	// @return *config.ChainConfig
	// @return error
	GetChainConfigByBlockHeight(blockHeight uint64) (*config.ChainConfig, error)
	// Register 注册一个新的链到归档中心
	// @param genesis 新链的创世区块
	// @return error
	Register(genesis *common.BlockInfo) error

	// ArchiveBlock 归档一个区块
	// @param block
	// @return error
	ArchiveBlock(block *common.BlockInfo) error

	// ArchiveBlocks 传入一个区块查询的迭代器，按顺序将区块归档
	// @param bi 区块迭代器
	// @return error
	ArchiveBlocks(bi BlockIterator, heightNoticeCallback func(ProcessMessage) error) error

	// GetArchivedStatus 获得归档的状态
	// @return archivedHeight
	// @return inArchive
	// @return code
	// @return err
	GetArchivedStatus() (archivedHeight uint64, inArchive bool, code uint32, err error)

	// Close 关闭与归档数据源的链接
	// @return err
	Close() error
}

// BlockIterator 区块查询迭代器
type BlockIterator interface {
	// Next 是否有下一个值
	// @return bool
	Next() bool
	// Value 取值
	// @return *common.BlockInfo
	// @return error
	Value() (*common.BlockInfo, error)
	// Release 释放资源
	Release()
	// Total 总数
	// @return uint64
	Total() uint64
	// Current 当前索引
	// @return uint64
	Current() uint64
}
