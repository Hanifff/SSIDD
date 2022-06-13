package ssidd

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PDPSC is an instance of an access decision
type PDPSC struct {
	contractapi.Contract
}

// Decision is an instance of decision made by PDPSC
type Decision struct {
	DecisionID  string `json:"decisionID"`
	SubjectID   string `json:"subjectID"` // can be clientdid
	ResourceID  string `json:"resourceID"`
	Decision    bool   `json:"decision"`
	Description string `json:"description"`
	TxnHash     string `json:"txnHash"`
	Timestamp   string `json:"timestamp"`
}
