package ssidd

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// initDB initalizes the db smartcontract.
func (t *DBSC) InitDB(ctx contractapi.TransactionContextInterface) error {
	txnid, txnHash, resourceid := "somenonreadablehash", "somehashee", "r0x1"
	db := Txn{TxnId: txnid, TxnHash: txnHash, ResourceId: resourceid}
	dbJSON, err := json.Marshal(db)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(resourceid, dbJSON)
}

// createTxn submits an instance of TXN to the world state.
func (t *DBSC) createTxn(ctx contractapi.TransactionContextInterface, txnid, txnHash, resourceid string) error {
	db := Txn{TxnId: txnid, TxnHash: txnHash, ResourceId: resourceid}
	dbJSON, err := json.Marshal(db)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(txnid, dbJSON)
}

// readTxn reads a Txn from world state with the given resourceid.
func (t *DBSC) readTxn(ctx contractapi.TransactionContextInterface, txnid string) (*Txn, error) {
	txnJSON, err := ctx.GetStub().GetState(txnid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if txnJSON == nil {
		return nil, fmt.Errorf("there is no transaciton with the given id")
	}
	txn := Txn{}
	err = json.Unmarshal(txnJSON, &txn)
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

// deleteTxn deletes an Txn from the world state.
func (t *DBSC) deleteTxn(ctx contractapi.TransactionContextInterface, txnid string) error {
	exists, err := t.txnExist(ctx, txnid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the transaction %s does not exist", txnid)
	}
	return ctx.GetStub().DelState(txnid)
}

// txnExists chaecks if a Txn exists in the world state.
func (t *DBSC) txnExist(ctx contractapi.TransactionContextInterface, txnid string) (bool, error) {
	txnJSON, err := ctx.GetStub().GetState(txnid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return txnJSON != nil, nil
}
