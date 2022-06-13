package ssidd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// InitLedger initalizes the ledger with an instance of Decision
func (pdp *PDPSC) InitLedger(ctx contractapi.TransactionContextInterface) error {
	decId := "decisionid"
	decision := Decision{DecisionID: decId,
		SubjectID: "client:DID:12345abcd", ResourceID: "resourceid",
		Decision: true, Description: "Got read access",
		Timestamp: time.Now().Format("02-Jan-2006"),
	}
	desJSON, err := json.Marshal(decision)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(decId, desJSON)
}

// DecideRead make decision for read request based on a policy and
// resource object requested from PAPSC and PIPSC.
// Returns an instance of decision to the SDK.
func (pdp *PDPSC) DecideRead(ctx contractapi.TransactionContextInterface, decisionid,
	clientDID, resourceid, cliAttr string) (*Decision, error) {
	readResource, readPolicy, auditator, db :=
		PIPSC{}, PAPSC{}, AuditSC{}, DBSC{}

	resouce, err := readResource.ReadResource(ctx, resourceid)
	if err != nil {
		return nil, fmt.Errorf("reading resource %s failed, %v", resourceid, err)
	}
	findPolicy := resouce.Attributes.AssociatedPolicyId
	policy, err := readPolicy.ReadPoicy(ctx, findPolicy)
	if err != nil {
		return nil, fmt.Errorf("reading policy %s failed, %v", findPolicy, err)
	}
	clientAttrs := make(map[string]interface{})
	err = json.Unmarshal([]byte(cliAttr), &clientAttrs)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling client attributes failed, %v", err)
	}
	// Set the number of times users have been denied access
	aud, err := auditator.ReadAudit(ctx, clientDID)
	if err != nil {
		return nil, fmt.Errorf("reading audit for %s failed, %v", clientDID, err)
	}
	clientAttrs["nrofbann"] = aud.NrOfBann
	ok := pdp.checkAttr(clientAttrs)
	if !ok {
		return pdp.decisionMaker(decisionid, clientDID, resourceid,
			"missing required user attribute(s)", "", false), nil
	}
	// validate access based on policy rules and provided client attributes
	for _, rule := range policy.Rules {
		if rule.Type == "READ" {
			attrValue, exists := clientAttrs[rule.Title]
			if !exists {
				return nil, fmt.Errorf("client does not fullfill the requriements for the %v attribute", rule.Title)
			}
			ok, err := Evalute(rule.Title, rule.Expression, attrValue)
			if err != nil {
				return nil, fmt.Errorf("client attribute for the %v is in invalid format, %v", rule.Title, err)
			}
			if !ok {
				return pdp.decisionMaker(decisionid, clientDID, resourceid,
					fmt.Sprintf("Client does not fulfill requirements for the %v attribute.", rule.Title), "", false), nil
			}
		}
	}
	// get the hash of data from db chaincode (invoke privetly)
	txn, err := db.readTxn(ctx, resouce.TxnId)
	if err != nil {
		return nil, fmt.Errorf("reading resource hash for resource %s failed, %v", resourceid, err)
	}
	// submitt the decision to the ledger
	decision := pdp.decisionMaker(decisionid, clientDID, resourceid,
		fmt.Sprintf("Client can get access to the resource %s.", resourceid), txn.TxnHash, true)
	decisionJSON, err := json.Marshal(decision)
	if err != nil {
		return nil, fmt.Errorf("marshaling decision failed, %v", err)
	}
	return decision, ctx.GetStub().PutState(decisionid, decisionJSON)
}

// DecideWrite make decision for write request based on a policy requested from PAPSC.
// Returns an instance of decision to the SDK.
func (pdp *PDPSC) DecideWrite(ctx contractapi.TransactionContextInterface, decisionid,
	clientDID, policyId, cliAttr string) (*Decision, error) {
	readPolicy, auditator := PAPSC{}, AuditSC{}
	policy, err := readPolicy.ReadPoicy(ctx, policyId)
	if err != nil {
		return nil, fmt.Errorf("reading policy %s failed, %v", policyId, err)
	}

	clientAttrs := make(map[string]interface{})
	err = json.Unmarshal([]byte(cliAttr), &clientAttrs)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling client attributes failed, %v", err)
	}
	// Set the number of times users have been denied access
	aud, err := auditator.ReadAudit(ctx, clientDID)
	if err != nil {
		return nil, fmt.Errorf("reading audit for %s failed, %v", clientDID, err)
	}
	clientAttrs["nrofbann"] = aud.NrOfBann
	ok := pdp.checkAttr(clientAttrs)
	if !ok {
		return pdp.decisionMaker(decisionid, clientDID, "",
			"missing required user attribute(s)", "", false), nil
	}
	// validate access based on policy rules and provided client attributes
	for _, rule := range policy.Rules {
		if rule.Type == "WRITE" {
			attrValue, exists := clientAttrs[rule.Title]
			if !exists {
				return nil, fmt.Errorf("client does not fullfill the requriements for the %v attribute", rule.Title)
			}
			ok, err := Evalute(rule.Title, rule.Expression, attrValue)
			if err != nil {
				return nil, fmt.Errorf("client attribute for the %v is in invalid format, %v", rule.Title, err)
			}
			if !ok {
				return pdp.decisionMaker(decisionid, clientDID, "",
					fmt.Sprintf("Client does not fulfill requirements for the %v attribute.", rule.Title), "", false), nil
			}
		}
	}
	// submitt the decision to the ledger
	decision := pdp.decisionMaker(decisionid, clientDID, "",
		fmt.Sprintf("Client can get share her data, decision id: %s.", decisionid), "", true)
	decisionJSON, err := json.Marshal(decision)
	if err != nil {
		return nil, fmt.Errorf("marshaling decision failed, %v", err)
	}
	return decision, ctx.GetStub().PutState(decisionid, decisionJSON)
}

// DeletePolicy deletes an given policy from the world state.
func (pap *PDPSC) DeleteDecision(ctx contractapi.TransactionContextInterface, decisionid string) error {
	desJSON, err := ctx.GetStub().GetState(decisionid)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if desJSON == nil {
		return fmt.Errorf("the decision %s does not exist", decisionid)
	}
	return ctx.GetStub().DelState(decisionid)
}

// GetAllDecisions returns all decisions found in world state
func (pdp *PDPSC) GetAllDecisions(ctx contractapi.TransactionContextInterface) ([]*Decision, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var decisions []*Decision
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var decision Decision
		err = json.Unmarshal(queryResponse.Value, &decision)
		if err != nil {
			return nil, err
		}
		decisions = append(decisions, &decision)
	}
	return decisions, nil
}

// decisionMaker produces an instance of Decision based on the given details.
func (pdp *PDPSC) decisionMaker(decisionid, clientDID, resourceid, description, txnid string, decision bool) *Decision {
	return &Decision{DecisionID: decisionid,
		SubjectID: clientDID, ResourceID: resourceid,
		Decision: decision, Description: description,
		TxnHash: txnid, Timestamp: time.Now().Format("02-Jan-2006"),
	}
}

// checkAttr is a private method that is used for checking required clients attributes are provided
func (pdp *PDPSC) checkAttr(clientAttrs map[string]interface{}) bool {
	// check if all required user attributes are provide by user
	if clientAttrs["status"] == nil || clientAttrs["nrofbann"] == nil ||
		clientAttrs["organisationid"] == nil || clientAttrs["expiration"] == nil {
		return false
	}
	return true
}

// decisionExists is a private method that returns true when descision with given ID exists in world state.
func (pdp *PDPSC) decisionExists(ctx contractapi.TransactionContextInterface, decisionid,
	ClientDID, ResourceID string) (bool, error) {
	desJSON, err := ctx.GetStub().GetState(decisionid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	var decision Decision
	err = json.Unmarshal(desJSON, &decision)
	if err != nil {
		return false, err
	}
	if decision.SubjectID == ClientDID && decision.ResourceID == ResourceID {
		return decision.Decision, nil
	}
	return false, fmt.Errorf("could not evaluate the decision from world state")
}
