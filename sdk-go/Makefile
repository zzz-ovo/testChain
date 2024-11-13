VERSION=v2.3.5

build:
	go mod tidy && go build ./...

lint:
	golangci-lint run ./...

mockgen:
	mockgen -destination ./mock/sdk_mock.go -package mock -source sdk_interface.go

gomod:
	go get chainmaker.org/chainmaker/common/v2@$(VERSION)
	go get chainmaker.org/chainmaker/pb-go/v2@v2.3.6
	go get chainmaker.org/chainmaker/protocol/v2@v2.3.6
	go get chainmaker.org/chainmaker/utils/v2@$(VERSION)
	go mod tidy

.PHONY: build lint
