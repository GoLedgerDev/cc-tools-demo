#!/usr/bin/env bash

function trap_ctrlc ()
{
    if [[ ! -z "$CUSTOMASSETSFILE" ]]
    then
        printf "%s\n" "$CUSTOMASSETSFILE" > chaincode/assettypes/customAssets.go
    fi

    if [[ ! -z "$HEADERFILE" ]]
    then
        printf "%s\n" "$HEADERFILE" > chaincode/header/header.go
    fi

    exit 2
}

# Make sure go mod is up to date
cd chaincode && go mod vendor && cd ..

# Copy customAssets.go content
CUSTOMASSETSFILE=$(cat chaincode/assettypes/customAssets.go)
HEADERFILE=$(cat chaincode/header/header.go)

# Prevent loss of the customAssets.go file
trap "trap_ctrlc" 2

# Delete customAssets.go from tree before compressing
rm chaincode/assettypes/customAssets.go
rm chaincode/header/header.go

# Compress file
tar -czf cc-tools-demo.tar.gz chaincode

# Restore customAssets.go file
printf "%s\n" "$CUSTOMASSETSFILE" > chaincode/assettypes/customAssets.go
printf "%s\n" "$HEADERFILE" > chaincode/header/header.go

# Clear trap
trap - 2

# Upload to S3 bucket
aws s3 cp cc-tools-demo.tar.gz s3://gofabric-templates/trial/