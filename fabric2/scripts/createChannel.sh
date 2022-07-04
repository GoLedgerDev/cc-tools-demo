#!/bin/bash

# imports  
. scripts/envVar.sh
. scripts/utils.sh

CHANNEL_NAME="$1"
DELAY="$2"
MAX_RETRY="$3"
VERBOSE="$4"
ORG_QNTY=$5
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${MAX_RETRY:="5"}
: ${VERBOSE:="false"}
: ${ORG_QNTY:=3}

if [ ! -d "channel-artifacts" ]; then
	mkdir channel-artifacts
fi

createChannelTx() {
	set -x
	if [ $ORG_QNTY -gt 1 ]
	then
		configtxgen -configPath "./configtx" -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME
	else
		configtxgen -configPath "./configtx-1org" -profile OneOrgChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME
	fi
	res=$?
	{ set +x; } 2>/dev/null
  	verifyResult $res "Failed to generate channel configuration transaction..."
}

createChannel() {
	if [ $ORG_QNTY -gt 1 ]
	then
		setGlobals 1
	else
		setGlobals 0
	fi
	# Poll in case the raft leader is not set yet
	local rc=1
	local COUNTER=1
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
		sleep $DELAY
		set -x
		peer channel create -o localhost:7050 -c $CHANNEL_NAME --ordererTLSHostnameOverride orderer.example.com -f ./channel-artifacts/${CHANNEL_NAME}.tx --outputBlock $BLOCKFILE --tls --cafile $ORDERER_CA >&log.txt
		res=$?
		{ set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	verifyResult $res "Channel creation failed"
}

# joinChannel ORG
joinChannel() {
  FABRIC_CFG_PATH=$PWD/./config/
  ORG=$1
  setGlobals $ORG
	local rc=1
	local COUNTER=1
	## Sometimes Join takes time, hence retry
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
    peer channel join -b $BLOCKFILE >&log.txt
    res=$?
    { set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	verifyResult $res "After $MAX_RETRY attempts, peer0.org${ORG} has failed to join channel '$CHANNEL_NAME' "
}

setAnchorPeer() {
  ORG=$1
  docker exec cli ./scripts/setAnchorPeer.sh $ORG $CHANNEL_NAME 
}

if [ $ORG_QNTY -gt 1 ]
then
	FABRIC_CFG_PATH=${PWD}/configtx

	## Create channeltx
	infoln "Generating channel create transaction '${CHANNEL_NAME}.tx'"
	createChannelTx

	FABRIC_CFG_PATH=$PWD/./config/
	BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"

	## Create channel
	infoln "Creating channel ${CHANNEL_NAME}"
	createChannel
	successln "Channel '$CHANNEL_NAME' created"

	## Join all the peers to the channel
	infoln "Joining org1 peer to the channel..."
	joinChannel 1
	infoln "Joining org2 peer to the channel..."
	joinChannel 2
	infoln "Joining org3 peer to the channel..."
	joinChannel 3

	## Set the anchor peers for each org in the channel
	infoln "Setting anchor peer for org1..."
	setAnchorPeer 1
	infoln "Setting anchor peer for org2..."
	setAnchorPeer 2
	infoln "Setting anchor peer for org2..."
	setAnchorPeer 3
else
	FABRIC_CFG_PATH=${PWD}/configtx-1org

	## Create channeltx
	infoln "Generating channel create transaction '${CHANNEL_NAME}.tx'"
	createChannelTx

	FABRIC_CFG_PATH=$PWD/./config/
	BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"

	## Create channel
	infoln "Creating channel ${CHANNEL_NAME}"
	createChannel
	successln "Channel '$CHANNEL_NAME' created"

	## Join all the peers to the channel
	infoln "Joining org peer to the channel..."
	joinChannel 0

	## Set the anchor peers for each org in the channel
	infoln "Setting anchor peer for org..."
	setAnchorPeer 0
fi

successln "Channel '$CHANNEL_NAME' joined"
