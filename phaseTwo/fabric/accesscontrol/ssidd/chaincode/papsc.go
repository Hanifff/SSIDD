package ssidd

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
Example of a Policy in json:
P0XX = {
	PolicyID :"P0XX", // REMOVED: 	Attributes: {}
	UserRequirement: {
		Status: "Active",
		Expiration: "Expiration date",
		OptionalUserAttr: {
			OrganisationID: "User's employer organisation", // an example
			// one of OPTIONAL fields
		}
	},
	ResourceRequirement: {
		Status: "Active",
		Expiration: "Expiration date",
		OptionalUserAttr: {
			OwnerOrgID: "Resource Owner org Id", // an example
			// one of OPTIONAL fields
		}
	},
	Rules: []PolicyRules{
		&PolicyRules{
			TYPE: "READ",
			Title: "NrOfBann", // IMPORTANT: MUST MATCH THE OptionalUserAttr.key
			ComparisionType: "bool",
			Description: "User must not have been banned more than X times.",
			Expression: "`nrOfBann =< 10`",
		},
		&PolicyRules{
			TYPE: "READ",
			Title: "OrganisationID", // IMPORTANT: MUST MATCH THE OptionalUserAttr.key
			ComparisionType: "bool",
			Description: "User must be employe in organisation SPAIN-01.",
			Expression: "`OrganisationID == SPAIN01`",
		},
		&PolicyRules{
			TYPE: "READ",
			Title: "Status"
			ComparisionType: "bool",
			Description: "User must be an active user.",
			Expression: "`Status == true`",
		},
		&PolicyRules{
			TYPE: "WRITE",
			Title: "NFTID"
			ComparisionType: "bool",
			Description: "User's organisation must hold the NFT.",
			Expression: "`NFTID == NF0XX1`",
		},
	}
}
*/

type PAPSC struct {
	contractapi.Contract
}

// Policy is an instance of policy stored onchain
type Policy struct {
	PolicyID string               `json:"policyid"`
	User     *UserRequirement     `json:"user"`
	Resource *ResourceRequirement `json:"resource"`
	Rules    []*PolicyRules       `json:"rules"`
}

type UserRequirement struct {
	Status           string                 `json:"status"`
	NrOfBann         string                 `json:"nrofbann"`
	Expiration       string                 `json:"expiration"`
	OrganisationID   string                 `json:"organisationid"`
	OptionalUserAttr map[string]interface{} `json:"optionaluserattr"`
}
type ResourceRequirement struct {
	Status          string                 `json:"status"`
	Expiration      string                 `json:"expiration"`
	OptionalResAttr map[string]interface{} `json:"optionalresattr"` // can for example be resource ID.
}

type PolicyRules struct {
	Type  string `json:"type"`  // what type of access request, i.e. READ, WRITE, DELETE, UPDATE
	Title string `json:"title"` // example: OrganisationID
	// IMPORTANT: MUST MATCH THE OptionalUserAttr.key
	ComparisionType interface{} `json:"comparisiontype"`
	Description     string      `json:"description"`
	Expression      string      `json:"expression"`
}
