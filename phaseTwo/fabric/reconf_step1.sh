#!/bin/bash

# Exit on first error
set -e

pushd ../test-network


# fecth with ordere as admin
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config/

export CORE_PEER_LOCALMSPID="OrdererMSP"
export CH_NAME="mychannel"
#export TLS_ROOT_CA=${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/tls/ca.crt 
export TLS_ROOT_CA=${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/cacerts/localhost-9054-ca-orderer.pem
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp
export ORDERER_CONTAINER=localhost:7050



# Step 1: pull and translate config
# fetch the channel configuration in protobuf format
peer channel fetch config config_block.pb -o $ORDERER_CONTAINER -c $CH_NAME --tls --cafile $TLS_ROOT_CA
# convert the protobuf of the channel config into JSON 
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
# trim config
jq .data.data[0].payload.data.config config_block.json > config.json
# copy config
cp config.json modified_config.json
popd


cat <<EOF

Now: modify modified_config.json (copy of config file)
defualt values for batch:
         "BatchSize": {
            "mod_policy": "Admins",
            "value": {
              "absolute_max_bytes": 103809024,
              "max_message_count": 10,
              "preferred_max_bytes": 524288

use either jq or open config file in vscode
Or you can use: -jq -s '.[0] * {"channel_group":{"groups":{"Orderer": {"values": {"BatchSize": { absolute_max_bytes: "$a" max_message_count: "$i" , preferred_max_bytes: "$j" }}}}}}' config.json ./capabilities.json > modified_config.json

Next, run the reconf_step2.sh

EOF
