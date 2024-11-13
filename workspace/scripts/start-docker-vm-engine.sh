IMAGE_VERSION=v2.3.5

LOG_PATH=/root/chainmaker/workspace/build/log/wx-org.chainmaker.org/go
MOUNT_PATH=/root/chainmaker/workspace/build/data/wx-org.chainmaker.org/go

docker run -itd --rm \
  --net=host \
  -v ${LOG_PATH}:/log \
  -v ${MOUNT_PATH}:/mount \
  -e CHAIN_RPC_PROTOCOL="1" \
  -e CHAIN_RPC_PORT="22351" \
  -e SANDBOX_RPC_PORT="32351" \
  -e MAX_SEND_MSG_SIZE="100" \
  -e MAX_RECV_MSG_SIZE="100" \
  -e MAX_CONN_TIMEOUT="10" \
  -e MAX_ORIGINAL_PROCESS_NUM="20" \
  -e DOCKERVM_CONTRACT_ENGINE_LOG_LEVEL="DEBUG" \
  -e DOCKERVM_SANDBOX_LOG_LEVEL="DEBUG" \
  -e DOCKERVM_LOG_IN_CONSOLE="false" \
	--privileged \
	--name chainmaker_vm_solo \
	chainmakerofficial/chainmaker-vm-engine:${IMAGE_VERSION}
