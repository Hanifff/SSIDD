package ssidd

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// RecordAudit records or updates a banned or non-banned operation with given details.
func (a *AuditSC) RecordAudit(ctx contractapi.TransactionContextInterface, clientDID string) error {
	exists, err := a.AuditExists(ctx, clientDID)
	if err != nil {
		return err
	}
	if exists {
		// get the audit
		auditJSON, err := ctx.GetStub().GetState(clientDID)
		if err != nil {
			return fmt.Errorf("failed to read from world state: %v", err)
		}
		if auditJSON == nil {
			return fmt.Errorf("there us no records for %s", clientDID)
		}
		var action Action
		err = json.Unmarshal(auditJSON, &action)
		if err != nil {
			return err
		}
		// update
		err = updateAudit(ctx, clientDID, action.NrOfBann+1)
		if err != nil {
			return err
		}
		return nil
	}
	auditJSON, err := json.Marshal(Action{ClientDID: clientDID, NrOfBann: 1})
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(clientDID, auditJSON)
}

// ReadAudit returns the audit object stored in the world state with given clientDID.
func (a *AuditSC) ReadAudit(ctx contractapi.TransactionContextInterface, clientDID string) (*Action, error) {
	exists, err := a.AuditExists(ctx, clientDID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &Action{ClientDID: clientDID, NrOfBann: 0}, nil
	}
	auditJSON, err := ctx.GetStub().GetState(clientDID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if auditJSON == nil {
		return &Action{}, fmt.Errorf("there is no records for %s", clientDID)
	}
	var action Action
	err = json.Unmarshal(auditJSON, &action)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// DeleteAudit deletes an given audit from the world state.
func (a *AuditSC) DeleteAudit(ctx contractapi.TransactionContextInterface, clientDID string) (bool, error) {
	exists, err := a.AuditExists(ctx, clientDID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("the audit %s does not exist", clientDID)
	}
	return true, ctx.GetStub().DelState(clientDID)
}

// AuditExists returns true when audit with given ID exists in world state
func (a *AuditSC) AuditExists(ctx contractapi.TransactionContextInterface, clientDID string) (bool, error) {
	auditJSON, err := ctx.GetStub().GetState(clientDID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return auditJSON != nil, nil
}

// access returns all audits found in world state
func (a *AuditSC) GetAllAudits(ctx contractapi.TransactionContextInterface) ([]*Action, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var actions []*Action
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var action Action
		err = json.Unmarshal(queryResponse.Value, &action)
		if err != nil {
			return nil, err
		}
		actions = append(actions, &action)
	}
	return actions, nil
}

// updateAudit updates an existing audit with the new number of banned.
func updateAudit(ctx contractapi.TransactionContextInterface, clientDID string, nrOfBann int) error {
	auditJSON, err := json.Marshal(Action{clientDID, nrOfBann})
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(clientDID, auditJSON)
}
