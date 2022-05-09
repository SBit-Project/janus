#!/bin/sh
while
     scli help &> /dev/null
     rc=$?; if [[ $rc == 0 ]]; then break; fi
do :;  done

balance=`scli getbalance`
if [ "${balance:0:1}" == "0" ]
then
    set -x
	scli generate 600 > /dev/null
	set -
fi

WALLETFILE=test-wallet
LOCKFILE=${SBIT_DATADIR}/import-test-wallet.lock

if [ ! -e $LOCKFILE ]; then
  while
       scli getaddressesbyaccount "" &> /dev/null
       rc=$?; if [[ $rc != 0 ]]; then continue; fi

       set -x

       scli importprivkey "cMbgxCJrTYUqgcmiC1berh5DFrtY1KeU4PXZ6NZxgenniF1mXCRk" # addr=sUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW hdkeypath=m/88'/0'/1'
       scli importprivkey "cRcG1jizfBzHxfwu68aMjhy78CpnzD9gJYZ5ggDbzfYD3EQfGUDZ" # addr=sLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf hdkeypath=m/88'/0'/2'

       solar prefund sUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW 500
       solar prefund sLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf 500
       touch $LOCKFILE

       set -
       break
  do :;  done
fi
