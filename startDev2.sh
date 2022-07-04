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

# Clear unused images and volumes
docker rmi $(docker images --quiet --filter "dangling=true")
docker volume rm $(docker volume ls -qf dangling=true)

# Script used to start the development environment.
if [ ! -d "chaincode/vendor" ]; then
    cd ./chaincode; go mod vendor; cd ..
fi
cd ./chaincode; go fmt ./...; cd ..
cd ./fabric2; ./startDev.sh -n $ORG_QNTY; cd ..
cd ./rest-server; ./startDev2.sh -n $ORG_QNTY; cd ..
