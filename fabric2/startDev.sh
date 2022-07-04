#!/usr/bin/env bash

ORG_QNTY=3

while getopts n: opt; do
    case $opt in
        n)  ORG_QNTY=${OPTARG}
            ;;
    esac
done

if [ $ORG_QNTY != 3 -a $ORG_QNTY != 1 ]
then
  echo 'WARNING: The number of organizations allowed is either 3 or 1.'
  echo 'Defaulting to 3 organizations.'
  ORG_QNTY=3
fi

if [ $ORG_QNTY == 1 ]
then
  CCCG_PATH="../chaincode/collections2-org.json"
else
  CCCG_PATH="../chaincode/collections2.json"
fi

./network.sh down -n $ORG_QNTY
rm -rf organizations/peerOrganizations
rm -rf organizations/ordererOrganizations
rm -rf organizations/rest-certs
docker network create cc-tools-demo-net
./network.sh up createChannel -n $ORG_QNTY
./network.sh deployCC -ccn cc-tools-demo -ccp ../chaincode -ccl go -n $ORG_QNTY -cccg $CCCG_PATH