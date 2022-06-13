#!/bin/bash

# Exit on first error
set -e

# clean out any old identites in the wallets
rm -rf application-go/wallet/*

# launch network; create channel, join peers with the channel, and add third org
pushd ../test-network
./network.sh down
docker system prune --volumes -f
./network.sh up createChannel -ca -s couchdb
popd

pushd ../test-network/addOrg3
./addOrg3.sh up -c mychannel -s couchdb
popd

pushd ../test-network
./network.sh deployCC -ccn ssidd -ccv 1 -ccp ../accesscontrol/ssidd -ccl go -c mychannel

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org3MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
export CORE_PEER_ADDRESS=localhost:11051

# package the chaincode on org3
peer lifecycle chaincode package ssidd.tar.gz --path ../accesscontrol/ssidd --lang golang --label ssidd_1
# install the chain code on org3
peer lifecycle chaincode install ssidd.tar.gz

# query packageID
peer lifecycle chaincode queryinstalled

read -t 600 -p "please enter the packageID: " packageID
export CC_PACKAGE_ID=$packageID

# approve the  a definition of the the chaincode for Org3
# use the --package-id flag to provide the package identifier
# use the --init-required flag to request the ``Init`` function be invoked to initialize the chaincode
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --channelID mychannel --name ssidd --version 1 --package-id $CC_PACKAGE_ID --sequence 1 #--init-required

# verify chaincode definition
# use the --name flag to select the chaincode whose definition you want to query
peer lifecycle chaincode querycommitted --channelID mychannel --name ssidd --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# get endoresment from two orgs
#peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n ssidd --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt" -c '{ "Args": ["PAPSC:InitPap", "{}"] }'
# veirfy
#peer chaincode query -C mychannel -n ssidd -c '{"Args":["GetAllDecisions"]}'
#peer chaincode invoke -C mychannel -n ssidd -c '{ "Args": ["PAPSC:InitPap", "{}"] }'

popd

cat <<EOF
    Succefully added the third org and deployed chaincode on it.  
EOF
