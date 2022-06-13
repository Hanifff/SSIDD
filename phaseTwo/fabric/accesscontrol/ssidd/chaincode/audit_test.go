package ssidd_test

import (
	"encoding/json"
	"fmt"
	"testing"

	ssidd "github.com/hanifff/ssiddSC/chaincode"
	mocks "github.com/hanifff/ssiddSC/chaincode/chaincodefakes"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
)

var actionTest = []*ssidd.Action{
	{ClientDID: "did:SPAIN01:123456789abcdefghigklmn", NrOfBann: 3},
	{ClientDID: "did:USA01:123456799abcdefghigklmo", NrOfBann: 1},
	{ClientDID: "did:JAPAN01:123466797abcdefghigklms", NrOfBann: 20},
}
var banned = []bool{false, true, false}

func TestRecordAudit(t *testing.T) {
	auditStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(auditStub)

	aud := ssidd.AuditSC{}
	for k, v := range actionTest {
		testName := fmt.Sprintf("%d, %s", k, v.ClientDID)
		t.Run(testName, func(t *testing.T) {
			err := aud.RecordAudit(transactionContext, v.ClientDID, banned[k])
			require.NoError(t, err)
		})
	}
}

func TestReadAudit(t *testing.T) {
	auditStub := &mocks.FakeSsiddStub{}
	aud := ssidd.AuditSC{}
	ctx := &mocks.FakeTransactionContext{}
	ctx.GetStubReturns(auditStub)
	expectedAudit := actionTest[0]
	bytes, err := json.Marshal(expectedAudit)
	require.NoError(t, err)

	auditStub.GetStateReturns(bytes, nil)
	act, err := aud.ReadAudit(ctx, "")
	require.NoError(t, err)
	require.Equal(t, expectedAudit, act)

	auditStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve audit"))
	_, err = aud.ReadAudit(ctx, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve audit")

	auditStub.GetStateReturns(nil, nil)
	act, err = aud.ReadAudit(ctx, "did:SPAIN01:123456789abcdefghigklmn")
	//require.EqualError(t, err, "the audit did:SPAIN01:123456789abcdefghigklmn does not exist")
	require.NoError(t, err)
	//require.Nil(t, act)
	require.Equal(t, act, &ssidd.Action{ClientDID: expectedAudit.ClientDID, NrOfBann: 0})
}

func TestDeleteAudit(t *testing.T) {
	aud := ssidd.AuditSC{}
	auditStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(auditStub)

	act := actionTest[0]
	bytes, err := json.Marshal(act)
	require.NoError(t, err)

	auditStub.GetStateReturns(bytes, nil)
	auditStub.DelStateReturns(nil)
	result, err := aud.DeleteAudit(transactionContext, "")
	require.Equal(t, true, result)
	require.NoError(t, err)

	auditStub.GetStateReturns(nil, nil)
	_, err = aud.DeleteAudit(transactionContext, act.ClientDID)
	require.EqualError(t, err, "the audit did:SPAIN01:123456789abcdefghigklmn does not exist")

	auditStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve audit"))
	_, err = aud.DeleteAudit(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve audit")
}

func TestGetAllAudits(t *testing.T) {
	aud := ssidd.AuditSC{}
	auditStub := &mocks.FakeSsiddStub{}
	act := actionTest[0]
	bytes, err := json.Marshal(act)
	require.NoError(t, err)

	iterator := &mocks.FakeStateIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(auditStub)

	auditStub.GetStateByRangeReturns(iterator, nil)
	actions, err := aud.GetAllAudits(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*ssidd.Action{act}, actions)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	actions, err = aud.GetAllAudits(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, actions)

	auditStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all audits"))
	actions, err = aud.GetAllAudits(transactionContext)
	require.EqualError(t, err, "failed retrieving all audits")
	require.Nil(t, actions)
}
