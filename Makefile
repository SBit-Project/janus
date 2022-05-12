ifndef GOBIN
GOBIN := $(GOPATH)/bin
endif

ifdef JANUS_PORT
JANUS_PORT := $(JANUS_PORT)
else
JANUS_PORT := 22402
endif


# Latest commit hash
GIT_SHA=$(shell git rev-parse HEAD)

# If working copy has changes, append `-local` to hash
GIT_DIFF=$(shell git diff -s --exit-code || echo "-local")
GIT_REV=$(GIT_SHA)$(GIT_DIFF)
GIT_TAG=$(shell git describe --tags 2>/dev/null)

ifeq ($(GIT_TAG),)
GIT_TAG := $(GIT_REV)
else
GIT_TAG := $(GIT_TAG)$(GIT_DIFF)
endif

check-env:
ifndef GOPATH
	$(error GOPATH is undefined)
endif

.PHONY: install
install: 
	go install \
		-ldflags "-X 'github.com/qtumproject/janus/pkg/params.GitSha=`./sha.sh``git diff -s --exit-code || echo \"-local\"`'" \
		github.com/SBit-Project/janus

.PHONY: release
release: darwin linux

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./build/janus-darwin-amd64 github.com/SBit-Project/janus/cli/janus

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -o ./build/janus-linux-amd64 github.com/SBit-Project/janus/cli/janus

.PHONY: quick-start
quick-start-regtest:
	cd docker && ./spin_up.regtest.sh && cd ..

.PHONY: quick-start-testnet
quick-start-testnet:
	cd docker && ./spin_up.testnet.sh && cd ..

.PHONY: quick-start-mainnet
quick-start-mainnet:
	cd docker && ./spin_up.mainnet.sh && cd ..

.PHONY: docker-dev
docker-dev:
	docker build -t sbit/janus:latest -t sbit/janus:dev -t sbit/janus:${GIT_TAG} -t sbit/janus:${GIT_REV} .
	
.PHONY: local-dev
local-dev: check-env install
	docker run --rm --name sbit_testchain -d -p 22002:22002 sbit/sbit sbitd -regtest -rpcbind=0.0.0.0:22002 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=sbit -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/SBit-Project/janus/docker/fill_user_account.sh sbit_testchain:.
	docker exec sbit_testchain /bin/sh -c ./fill_user_account.sh
	SBIT_RPC=http://sbit:testpasswd@localhost:22002 SBIT_NETWORK=auto $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev

.PHONY: local-dev-https
local-dev-https: check-env install
	docker run --rm --name sbit_testchain -d -p 22002:22002 sbit/sbit sbitd -regtest -rpcbind=0.0.0.0:22002 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=sbit -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/SBit-Project/janus/docker/fill_user_account.sh sbit_testchain:.
	docker exec sbit_testchain /bin/sh -c ./fill_user_account.sh > /dev/null&
	SBIT_RPC=http://sbit:testpasswd@localhost:22002 SBIT_NETWORK=auto $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev --https-key https/key.pem --https-cert https/cert.pem

.PHONY: local-dev-logs
local-dev-logs: check-env install
	docker run --rm --name sbit_testchain -d -p 22002:22002 sbit/sbit:dev sbitd -regtest -rpcbind=0.0.0.0:22002 -rpcallowip=0.0.0.0/0 -logevents=1 -rpcuser=sbit -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole | true
	sleep 3
	docker cp ${GOPATH}/src/github.com/SBit-Project/janus/docker/fill_user_account.sh sbit_testchain:.
	docker exec sbit_testchain /bin/sh -c ./fill_user_account.sh
	SBIT_RPC=http://sbit:testpasswd@localhost:22002 SBIT_NETWORK=auto $(GOBIN)/janus --port $(JANUS_PORT) --accounts ./docker/standalone/myaccounts.txt --dev > janus_dev_logs.txt

.PHONY: unit-tests
unit-tests: check-env
	go test -v ./... -timeout 30s

docker-build-unit-tests:
	docker build -t sbit/tests.janus -f ./docker/unittests.Dockerfile .

docker-unit-tests:
	docker run --rm -v `pwd`:/go/src/github.com/SBit-Project/janus sbit/tests.janus

docker-tests: docker-build-unit-tests docker-unit-tests openzeppelin-docker-compose

docker-configure-https: docker-configure-https-build
	docker/setup_self_signed_https.sh

docker-configure-https-build:
	docker build -t sbit/openssl.janus -f ./docker/openssl.Dockerfile ./docker

# -------------------------------------------------------------------------------------------------------------------
# NOTE:
# 	The following make rules are only for local test purposes
# 
# 	Both run-janus and run-sbit must be invoked. Invocation order may be independent, 
# 	however it's much simpler to do in the following order: 
# 		(1) make run-sbit 
# 			To stop Sbit node you should invoke: make stop-sbit
# 		(2) make run-janus
# 			To stop Janus service just press Ctrl + C in the running terminal

# Runs current Janus implementation
run-janus:
	@ printf "\nRunning Janus...\n\n"

	go run `pwd`/main.go \
		--sbit-rpc=http://${test_user}:${test_user_passwd}@0.0.0.0:22002 \
		--sbit-network=auto \
		--bind=0.0.0.0 \
		--port=22402 \
		--accounts=`pwd`/docker/standalone/myaccounts.txt \
		--log-file=janusLogs.txt \
		--dev

run-janus-https:
	@ printf "\nRunning Janus...\n\n"

	go run `pwd`/main.go \
		--sbit-rpc=http://${test_user}:${test_user_passwd}@0.0.0.0:22002 \
		--sbit-network=auto \
		--bind=0.0.0.0 \
		--port=22402 \
		--accounts=`pwd`/docker/standalone/myaccounts.txt \
		--log-file=janusLogs.txt \
		--dev \
		--https-key https/key.pem \
		--https-cert https/cert.pem

test_user = sbit
test_user_passwd = testpasswd

# Runs docker container of sbit locally and starts sbitd inside of it
run-sbit:
	@ printf "\nRunning sbit...\n\n"
	@ printf "\n(1) Starting container...\n\n"
	docker run ${sbit_container_flags} sbit/sbit sbitd ${sbitd_flags} > /dev/null

	@ printf "\n(2) Importing test accounts...\n\n"
	@ sleep 3
	docker cp ${shell pwd}/docker/fill_user_account.sh ${sbit_container_name}:.

	@ printf "\n(3) Filling test accounts wallets...\n\n"
	docker exec ${sbit_container_name} /bin/sh -c ./fill_user_account.sh > /dev/null
	@ printf "\n... Done\n\n"

seed-sbit:
	@ printf "\n(2) Importing test accounts...\n\n"
	docker cp ${shell pwd}/docker/fill_user_account.sh ${sbit_container_name}:.

	@ printf "\n(3) Filling test accounts wallets...\n\n"
	docker exec ${sbit_container_name} /bin/sh -c ./fill_user_account.sh
	@ printf "\n... Done\n\n"

sbit_container_name = test-chain

# TODO: Research -v
sbit_container_flags = \
	--rm -d \
	--name ${sbit_container_name} \
	-v ${shell pwd}/dapp \
	-p 22002:22002

# TODO: research flags
sbitd_flags = \
	-regtest \
	-rpcbind=0.0.0.0:22002 \
	-rpcallowip=0.0.0.0/0 \
	-logevents \
	-addrindex \
	-reindex \
	-txindex \
	-rpcuser=${test_user} \
	-rpcpassword=${test_user_passwd} \
	-deprecatedrpc=accounts \
	-printtoconsole

# Starts continuously printing Sbit container logs to the invoking terminal
follow-sbit-logs:
	@ printf "\nFollowing sbit logs...\n\n"
		docker logs -f ${sbit_container_name}

open-sbit-bash:
	@ printf "\nOpening sbit bash...\n\n"
		docker exec -it ${sbit_container_name} bash

# Stops docker container of sbit
stop-sbit:
	@ printf "\nStopping sbit...\n\n"
		docker kill `docker container ps | grep ${sbit_container_name} | cut -d ' ' -f1` > /dev/null
	@ printf "\n... Done\n\n"

restart-sbit: stop-sbit run-sbit

submodules:
	git submodules init

# Run openzeppelin tests, Janus/SBIT needs to already be running
openzeppelin:
	cd testing && make openzeppelin

# Run openzeppelin tests in docker
# Janus and SBIT need to already be running
openzeppelin-docker:
	cd testing && make openzeppelin-docker

# Run openzeppelin tests in docker-compose
openzeppelin-docker-compose:
	cd testing && make openzeppelin-docker-compose