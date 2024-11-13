PLATFORM=$(shell uname -m)
DATETIME=$(shell date "+%Y%m%d%H%M%S")
VERSION=v2.3.5


build: chainmaker-cryptogen

vendor:
	# 注意：执行此方法前需要切换common项目到对应分支或commit
	# 密码学 gmssl 相关
	@cp -a ../common/opencrypto/gmssl/gmssl/include/openssl/ ./src/vendor/chainmaker.org/chainmaker/common/v2/opencrypto/gmssl/gmssl/
	@cp -a ../common/opencrypto/gmssl/gmssl/lib/linux ./src/vendor/chainmaker.org/chainmaker/common/v2/opencrypto/gmssl/gmssl/lib/
	# 密码学 tencentsm 相关
	@cp -a ../common/opencrypto/tencentsm/tencentsm/include ./src/vendor/chainmaker.org/chainmaker/common/v2/opencrypto/tencentsm/tencentsm/
	@cp -a ../common/opencrypto/tencentsm/tencentsm/lib ./src/vendor/chainmaker.org/chainmaker/common/v2/opencrypto/tencentsm/tencentsm/
	@# 密码学 bulletproofs 相关
	@#cp -a ../common/crypto/bulletproofs/bulletproofs_cgo/c_include ./src/vendor/chainmaker.org/chainmaker/common/v2/crypto/bulletproofs/bulletproofs_cgo/
	@#cp -a ../common/crypto/bulletproofs/bulletproofs_cgo/c_lib ./src/ven/dor/chainmaker.org/chainmaker/common/v2/crypto/bulletproofs/bulletproofs_cgo
	@cd src && go build -mod=vendor -o ../bin/chainmaker-cryptogen
	@mkdir -p ./release
	@rm -rf ./tmp/chainmaker-cryptogen/
	@mkdir -p ./tmp/chainmaker-cryptogen/
	@mkdir ./tmp/chainmaker-cryptogen/bin
	@mkdir ./tmp/chainmaker-cryptogen/config
	@cp bin/chainmaker-cryptogen ./tmp/chainmaker-cryptogen/bin
	@cp -r config/crypto_config_template.yml ./tmp/chainmaker-cryptogen/config
	@cd ./tmp;tar -zcvf chainmaker-cryptogen-$(VERSION).$(DATETIME).$(PLATFORM).tar.gz chainmaker-cryptogen; mv chainmaker-cryptogen-$(VERSION).$(DATETIME).$(PLATFORM).tar.gz ../release
	@rm -rf ./tmp/

chainmaker-cryptogen:
	@cd src && go mod tidy && GOPATH=${GOPATH} go build -o ../bin/chainmaker-cryptogen
	@mkdir -p ./release
	@rm -rf ./tmp/chainmaker-cryptogen/
	@mkdir -p ./tmp/chainmaker-cryptogen/
	@mkdir ./tmp/chainmaker-cryptogen/bin
	@mkdir ./tmp/chainmaker-cryptogen/config
	@cp bin/chainmaker-cryptogen ./tmp/chainmaker-cryptogen/bin
	@cp -r config/crypto_config_template.yml ./tmp/chainmaker-cryptogen/config
	@cd ./tmp;tar -zcvf chainmaker-cryptogen-$(VERSION).$(DATETIME).$(PLATFORM).tar.gz chainmaker-cryptogen; mv chainmaker-cryptogen-$(VERSION).$(DATETIME).$(PLATFORM).tar.gz ../release
	@rm -rf ./tmp/

clean:
	@rm -rf ./bin

gomod:
	@cd src && go get chainmaker.org/chainmaker/common/v2@$(VERSION)
	@cd src && go mod tidy
