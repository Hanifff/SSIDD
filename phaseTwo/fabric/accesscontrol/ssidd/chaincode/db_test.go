package ssidd

import (
	"encoding/json"
	"fmt"
	"testing"

	mocks "github.com/hanifff/ssiddSC/chaincode/chaincodefakes"
	"github.com/stretchr/testify/require"
)

var txns = []*Txn{&Txn{TxnId: "someid10", TxnHash: "somenonreadablehash", ResourceId: "r001"},
	&Txn{TxnId: "someid11", TxnHash: "somenonreada1234567890", ResourceId: "r002"},
	&Txn{TxnId: "someid22", TxnHash: "123456789readablehash", ResourceId: "r003"}}

func TestCreateTxn(t *testing.T) {
	auditStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(auditStub)
	db := DBSC{}

	for k, v := range txns {
		testName := fmt.Sprintf("%d, %s", k, v.TxnId)
		t.Run(testName, func(t *testing.T) {
			err := db.createTxn(transactionContext, v.TxnId, v.TxnHash, v.ResourceId)
			require.NoError(t, err)
		})
	}
}

func TestReadTxn(t *testing.T) {
	dbStub := &mocks.FakeSsiddStub{}
	db := DBSC{}
	ctx := &mocks.FakeTransactionContext{}
	ctx.GetStubReturns(dbStub)
	expectedTxn := txns[0]
	bytes, err := json.Marshal(expectedTxn)
	require.NoError(t, err)

	dbStub.GetStateReturns(bytes, nil)
	txn, err := db.readTxn(ctx, "")
	require.NoError(t, err)
	require.Equal(t, expectedTxn, txn)

	dbStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve transaction"))
	_, err = db.readTxn(ctx, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve transaction")
}

func TestDeleteTxn(t *testing.T) {
	dbStub := &mocks.FakeSsiddStub{}
	db := DBSC{}
	ctx := &mocks.FakeTransactionContext{}
	ctx.GetStubReturns(dbStub)
	expectedTxn := txns[0]
	bytes, err := json.Marshal(expectedTxn)
	require.NoError(t, err)

	dbStub.GetStateReturns(bytes, nil)
	dbStub.DelStateReturns(nil)
	err = db.deleteTxn(ctx, "")
	require.NoError(t, err)

	dbStub.GetStateReturns(nil, nil)
	err = db.deleteTxn(ctx, expectedTxn.TxnId)
	require.EqualError(t, err, fmt.Sprintf("the transaction %s does not exist", expectedTxn.TxnId))

	dbStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve transaction"))
	err = db.deleteTxn(ctx, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve transaction")
}
