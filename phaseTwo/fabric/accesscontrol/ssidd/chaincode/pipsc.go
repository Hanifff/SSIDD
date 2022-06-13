package ssidd

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
Example of a Resource in json:
R0XX = {
	ResouceID : "R0XX",
	Attributes : {
		Status: true,
		AssociatedPolicyId: "P0XX"
		OwnerOrgID: "SPAIN-01",
		Expiration: "YYYY-MM-DD",
		OptionalAttr: {
			NFTID: "NX01",
			etc...
		}
	}
}
*/

// PIPSC is an instance of PIP smart contract
type PIPSC struct {
	contractapi.Contract
}

// Resource is an instance of objects to access
type Resource struct {
	ResouceID  string        `json:"resourceid"`
	TxnId      string        `json:"txnid"`
	Attributes *ResourceAttr `json:"attributes"`
}

// ResourceAttr is an instance of all attributes associated with the resource
type ResourceAttr struct {
	Status             bool                   `json:"status"`
	AssociatedPolicyId string                 `json:"associatedpolicyid"`
	OwnerOrgID         string                 `json:"ownerorgid"`
	Expiration         string                 `json:"expiration"`
	OptionalAttr       map[string]interface{} `json:"optionalattr"`
}
