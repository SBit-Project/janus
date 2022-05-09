#!/bin/sh
docker-compose -f ${GOPATH}/src/github.com/SBit-Project/janus/docker/quick_start/docker-compose.testnet.yml up -d
# sleep 3 #executing too fast causes some errors
# docker cp ${GOPATH}/src/github.com/SBit-Project/janus/docker/fill_user_account.sh sbit_testchain:.
# docker exec sbit_testnet /bin/sh -c ./fill_user_account.sh