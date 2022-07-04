#!/usr/bin/env bash

# Get env file
source .env

# Default organizations quantity
ORG_QTY=3

while getopts g:n: opt; do
    case $opt in
        g)
            GENERATE_CERT=true
            ;;
        n)  ORG_QTY=${OPTARG}
            ;;
    esac
done

if [ $ORG_QTY -gt 3 -o $ORG_QTY -lt 1 ]
then
  echo 'WARNING: The maximum number of organizations allowed is 3 and the minimum is 1.'
  echo 'Defaulting to 3 organizations.'
  ORG_QTY=3
fi

# Stop and remove containers
docker ps -a | grep "ccapi" | awk '{print $1}' | xargs docker stop
docker ps -a | grep "ccapi" | awk '{print $1}' | xargs docker rm

# Create docker-compose file
SED_PORT=80
if [ $HTTPS = "true" ]
then
  SED_PORT=443
  # Generate self-signed cert
fi

if [ $ORG_QTY -gt 1 ]
then
  sed "s/PORT/${SED_PORT}/g" template/docker-compose-org1-template.yaml > docker-compose-org1.yaml
  sed "s/PORT/${SED_PORT}/g" template/docker-compose-org2-template.yaml > docker-compose-org2.yaml
  docker-compose -f docker-compose-org1.yaml -f docker-compose-org2.yaml -f docker-compose-temp.yaml down --volumes
  if [ $ORG_QTY -eq 3 ]
  then
    sed "s/PORT/${SED_PORT}/g" template/docker-compose-org3-template.yaml > docker-compose-org3.yaml
    docker-compose docker-compose-org3.yaml down --volumes
  fi
else
  sed "s/PORT/${SED_PORT}/g" template/docker-compose-org-template.yaml > docker-compose-org.yaml
  docker-compose -f docker-compose-org.yaml docker-compose-temp.yaml down --volumes
fi

cd scripts
if [ "$HTTPS" == true ]; then
  if [ "$GENERATE_CERT" == true ]; then
    ./generate-dummy-cert.sh -d $DOMAIN
    if [ "$LETS_ENCRYPT" == true ];then
      ./letsencrypt-init.sh -g
    fi
  else
    if [ "$LETS_ENCRYPT" == true ];then
        ./letsencrypt-init.sh
    fi
  fi
fi
cd ..

# Start API
if [ "$GENERATE_CERT" != true ]; then
  docker-compose -f docker-compose-temp.yaml -p intermediate-container up >> /dev/null
  if [ $ORG_QTY -gt 1 ]
  then
    docker-compose -f docker-compose-org1.yaml -p ccapi.org1.example.com up -d
    docker-compose -f docker-compose-org2.yaml -p ccapi.org2.example.com up -d
    if [ $ORG_QTY -eq 3 ]
    then
      docker-compose -f docker-compose-org3.yaml -p ccapi.org3.example.com up -d
    fi
  else 
    docker-compose -f docker-compose-org.yaml -p ccapi.org.example.com up -d
  fi
fi
