package main

import (
	"log"

	ssidd "github.com/hanifff/ssiddSC/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	ssiddSCCC, err := contractapi.NewChaincode(&ssidd.PDPSC{}, &ssidd.PAPSC{},
		&ssidd.PIPSC{}, &ssidd.AuditSC{}, &ssidd.DBSC{})
	if err != nil {
		log.Panicf("Error creating ssiddSCCC chaincode: %v", err)
	}

	if err := ssiddSCCC.Start(); err != nil {
		log.Panicf("Error starting ssiddSCCC chaincode: %v", err)
	}
}
