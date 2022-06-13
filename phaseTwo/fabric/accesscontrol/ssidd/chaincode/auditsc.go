package ssidd

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AuditSC struct {
	contractapi.Contract
}

// Action describe an illegal activity
type Action struct {
	ClientDID string `json:"clientdid"`
	NrOfBann  int    `json:"nrofbann"`
	/* Type        string    `json:"type"`
	ValidAction bool      `json:"validAction"`
	Timestamp   time.Time `json:"timestamp"` */
}
