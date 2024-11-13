VERSION=v2.3.5
IMAGE_VERSION=v2.3.5

BUILD_TIME = $(shell date "+%Y%m%d%H%M%S")
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT = $(shell git log --pretty=format:'%h' -n 1)

build-test:
	cd test/scripts && ./prepare.sh

build-image:
	cd vm_mgr && go mod vendor
	cd vm_mgr && docker build -t chainmaker-vm-engine \
	--build-arg BUILD_TIME=${BUILD_TIME} \
	--build-arg GIT_BRANCH=${GIT_BRANCH} \
	--build-arg GIT_COMMIT=${GIT_COMMIT} \
	-f Dockerfile ./
	docker tag chainmaker-vm-engine chainmakerofficial/chainmaker-vm-engine:${IMAGE_VERSION}
	docker images | grep chainmaker-vm-engine

image-push:
	docker push chainmakerofficial/chainmaker-vm-engine:${IMAGE_VERSION}

update-gomod:
	cd vm_mgr && rm -rf vendor
	bash .vm_mgr/scripts/gomod_update.sh

gen-dockervm-pb:
	cd pb/proto && protoc -I=. --gogofaster_out=plugins=grpc:../protogo --gogofaster_opt=paths=source_relative dockervm_message.proto
	cd vm_mgr/pb/proto && protoc -I=. --gogofaster_out=plugins=grpc:../protogo --gogofaster_opt=paths=source_relative dockervm_message.proto

clean-test:
	cd test/scripts && ./dockerclean.sh

clean:
	cd vm_mgr && rm -rf vendor
	cd test/scripts && ./dockerclean.sh
	docker image rm chainmakerofficial/chainmaker-vm-engine:${IMAGE_VERSION}
	docker image prune -f

ci:
	golangci-lint run ./...

gomod:
	go get chainmaker.org/chainmaker/common/v2@$(VERSION)
	go get chainmaker.org/chainmaker/localconf/v2@$(VERSION)
	go get chainmaker.org/chainmaker/pb-go/v2@v2.3.6
	go get chainmaker.org/chainmaker/protocol/v2@v2.3.6
	go get chainmaker.org/chainmaker/utils/v2@$(VERSION)
	go get chainmaker.org/chainmaker/vm/v2@v2.3.6
	go mod tidy

ut:
	./test/scripts/prepare.sh
	make build-image
	# UDS: docker run -itd --rm -v $(shell pwd)/data/org1/docker-go:/mount -v $(shell pwd)/log/org1/dockervm:/log --privileged --name chaimaker_vm_test chainmakerofficial/chainmaker-vm-engine:${IMAGE_VERSION}
	docker run -itd --net=host --privileged --name chaimaker_vm_test chainmakerofficial/chainmaker-vm-engine:${IMAGE_VERSION}
	sh ./ut_cover.sh
	docker stop chaimaker_vm_test

version:
	docker inspect chainmakerofficial/chainmaker-vm-engine:${IMAGE_VERSION} | jq '.[].ContainerConfig.Labels'