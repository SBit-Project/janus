#!/bin/sh
docker-compose -f ${GOPATH}/src/github.com/SBit-Project/janus/docker/quick_start/docker-compose.regtest.yml up -d
sleep 3 #executing too fast causes some errors
docker cp ${GOPATH}/src/github.com/SBit-Project/janus/docker/fill_user_account.sh sbit_regtest:.
docker exec sbit_regtest /bin/sh -c ./fill_user_account.sh