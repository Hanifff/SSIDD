package ssidd

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// InitPap initlaizes the PAPSC contract
func (pap *PAPSC) InitPap(ctx contractapi.TransactionContextInterface) error {
	optsUser := make(map[string]interface{})
	optsResource := make(map[string]interface{})

	status := "Active"
	expiration := "Date of expiration"
	NrOfBann := "Number of times uses is banned"
	usersOrgId := "Organisation ID"
	optsUser["nftid"] = "NFT id of the origanisation of user"
	optsResource["ownerorgid"] = "Resource owner organisation id"

	p001 := &Policy{PolicyID: "p001",
		User:     &UserRequirement{Status: status, NrOfBann: NrOfBann, Expiration: expiration, OrganisationID: usersOrgId, OptionalUserAttr: optsUser},
		Resource: &ResourceRequirement{Status: status, Expiration: expiration, OptionalResAttr: optsResource},
		Rules: []*PolicyRules{
			{Type: "READ", Title: "status", ComparisionType: "bool", Description: "User must be an active user.", Expression: `status == true`},
			{Type: "READ", Title: "nrofbann", ComparisionType: "bool", Description: "User with more than 5 times being banned does not get access.", Expression: `nrofbann <= 5`},
			{Type: "WRITE", Title: "organisationid", ComparisionType: "bool", Description: "User must be employe in organisation spain01.", Expression: `organisationid == 'spain01'`},
		},
	}
	policyJSON, err := json.Marshal(p001)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(p001.PolicyID, policyJSON)
}

// CreatePolicy issues a new policy object to the world state with given details.
func (pap *PAPSC) CreatePolicy(ctx contractapi.TransactionContextInterface,
	policyid, userStatus, resStatuc, nrOfBann, organisationID, userExpiration, resExpiration,
	optUserAttr, optResAttr, policyRules string) error {

	optionalUserAttr := make(map[string]interface{})
	err := json.Unmarshal([]byte(optUserAttr), &optionalUserAttr)
	if err != nil {
		return err
	}

	optionalResAttr := make(map[string]interface{})
	err = json.Unmarshal([]byte(optResAttr), &optionalResAttr)
	if err != nil {
		return err
	}

	var rules []*PolicyRules
	err = json.Unmarshal([]byte(policyRules), &rules)
	if err != nil {
		return err
	}
	exists, err := pap.PolicyExists(ctx, policyid)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the policy %s already exists", policyid)
	}
	policy := Policy{
		PolicyID: policyid,
		User: &UserRequirement{
			Status:           userStatus,
			NrOfBann:         nrOfBann,
			OrganisationID:   organisationID,
			Expiration:       userExpiration,
			OptionalUserAttr: optionalUserAttr,
		},
		Resource: &ResourceRequirement{
			Status:          resStatuc,
			Expiration:      resExpiration,
			OptionalResAttr: optionalResAttr,
		},
		Rules: rules,
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(policyid, policyJSON)
}

// ReadPoicy returns the policy object stored in the world state with given policyid.
func (pap *PAPSC) ReadPoicy(ctx contractapi.TransactionContextInterface, policyid string) (*Policy, error) {
	policyJSON, err := ctx.GetStub().GetState(policyid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if policyJSON == nil {
		return nil, fmt.Errorf("the policy %s does not exist", policyid)
	}
	user := &UserRequirement{}
	user.OptionalUserAttr = make(map[string]interface{})
	resource := &ResourceRequirement{}
	resource.OptionalResAttr = make(map[string]interface{})
	policy := &Policy{User: user, Resource: resource, Rules: []*PolicyRules{}}
	err = json.Unmarshal(policyJSON, &policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

// UpdatePolicy updates an existing policy in the world state with provided parameters.
func (pap *PAPSC) UpdatePolicy(ctx contractapi.TransactionContextInterface,
	policyid, userStatus, resStatuc, nrOfBann, organisationID, userExpiration, resExpiration,
	optUserAttr, optResAttr, policyRules string) error {

	optionalUserAttr := make(map[string]interface{})
	err := json.Unmarshal([]byte(optUserAttr), &optionalUserAttr)
	if err != nil {
		return err
	}
	optionalResAttr := make(map[string]interface{})
	err = json.Unmarshal([]byte(optResAttr), &optionalResAttr)
	if err != nil {
		return err
	}
	var rules []*PolicyRules
	err = json.Unmarshal([]byte(policyRules), &rules)
	if err != nil {
		return err
	}
	exists, err := pap.PolicyExists(ctx, policyid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the policy %s does not exist", policyid)
	}
	// overwriting original asset with new asset
	policy := Policy{
		PolicyID: policyid,
		User: &UserRequirement{
			Status:           userStatus,
			NrOfBann:         nrOfBann,
			Expiration:       userExpiration,
			OptionalUserAttr: optionalUserAttr,
		},
		Resource: &ResourceRequirement{
			Status:          resStatuc,
			Expiration:      resExpiration,
			OptionalResAttr: optionalResAttr,
		},
		Rules: rules,
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(policyid, policyJSON)
}

// DeletePolicy deletes an given policy from the world state.
func (pap *PAPSC) DeletePolicy(ctx contractapi.TransactionContextInterface, policyid string) error {
	exists, err := pap.PolicyExists(ctx, policyid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the policy %s does not exist", policyid)
	}
	return ctx.GetStub().DelState(policyid)
}

// PolicyExists returns true when policy with given ID exists in world state
func (pap *PAPSC) PolicyExists(ctx contractapi.TransactionContextInterface, policyid string) (bool, error) {
	policyJSON, err := ctx.GetStub().GetState(policyid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return policyJSON != nil, nil
}

// GetAllPoicies returns all policy found in world state
func (pap *PAPSC) GetAllPoicies(ctx contractapi.TransactionContextInterface) ([]*Policy, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var policies []*Policy
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var policy Policy
		err = json.Unmarshal(queryResponse.Value, &policy)
		if err != nil {
			return nil, err
		}
		policies = append(policies, &policy)
	}

	return policies, nil
}
