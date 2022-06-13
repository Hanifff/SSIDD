package ssidd

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

type DBSC struct {
	contractapi.Contract
}

// Txn stores a hash of data on the ledger.
type Txn struct {
	TxnId      string `json:"txnid"`
	TxnHash    string `json:"txnhash"`
	ResourceId string `json:"resourceid"`
}
