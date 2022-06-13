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

var (
	optsUser     map[string]interface{}
	optsResource map[string]interface{}
	policiesTest []*ssidd.Policy
)

func InitPapTest() {
	optsUser = make(map[string]interface{})
	optsResource = make(map[string]interface{})

	status := "Active"
	expiration := "Date of expiration"
	NrOfBann := "Number of times uses is banned"
	usersOrgId := "Organisation ID"
	optsUser["nftid"] = "NFT id of the origanisation of user"
	optsResource["ownerorgid"] = "Resource owner organisation id"

	p0xx := &ssidd.Policy{PolicyID: "p0xx",
		User:     &ssidd.UserRequirement{Status: status, NrOfBann: NrOfBann, Expiration: expiration, OrganisationID: usersOrgId, OptionalUserAttr: optsUser},
		Resource: &ssidd.ResourceRequirement{Status: status, Expiration: expiration, OptionalResAttr: optsResource},
		Rules: []*ssidd.PolicyRules{
			{Type: "READ", Title: "status", ComparisionType: "bool", Description: "User must be an active user.", Expression: `status == true`},
			{Type: "READ", Title: "nrofbann", ComparisionType: "bool", Description: "User with more than 5 times being banned does not get access.", Expression: `nrofbann <= 5`},
			{Type: "READ", Title: "organisationid", ComparisionType: "bool", Description: "User must be employe in organisation spain01.", Expression: `organisationid == 'spain01'`},
		},
	}
	policiesTest = []*ssidd.Policy{p0xx}
	p0x1 := &ssidd.Policy{PolicyID: "p0x1",
		User:     &ssidd.UserRequirement{Status: status, NrOfBann: NrOfBann, Expiration: expiration, OrganisationID: usersOrgId, OptionalUserAttr: optsUser},
		Resource: &ssidd.ResourceRequirement{Status: status, Expiration: expiration, OptionalResAttr: optsResource},
		Rules: []*ssidd.PolicyRules{
			{Type: "READ", Title: "status", ComparisionType: "bool", Description: "User must be an active user.", Expression: `status == true`},
			{Type: "READ", Title: "nrofbann", ComparisionType: "bool", Description: "User with more than 5 times being banned does not get access.", Expression: `nrofbann <= 5`},
			{Type: "READ", Title: "organisationid", ComparisionType: "bool", Description: "User must be employe in organisation spain01.", Expression: `organisationid == 'spain01'`},
			{Type: "WRITE", Title: "nftid", ComparisionType: "bool", Description: "User's organisation must hold the NFT.", Expression: `nftid == 'nf0xx1'`},
		},
	}
	policiesTest = append(policiesTest, p0x1)
}

func TestCreatePolicy(t *testing.T) {
	InitPapTest()
	papStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(papStub)

	pap := ssidd.PAPSC{}
	for k, v := range policiesTest {
		testName := fmt.Sprintf("%d, %s", k, v.PolicyID)
		t.Run(testName, func(t *testing.T) {
			optUserJson, err := json.Marshal(v.User.OptionalUserAttr)
			require.NoError(t, err)
			optResJson, err := json.Marshal(v.Resource.OptionalResAttr)
			require.NoError(t, err)
			rules, err := json.Marshal(v.Rules)
			require.NoError(t, err)
			err = pap.CreatePolicy(transactionContext, v.PolicyID, v.User.Status, v.Resource.Status,
				v.User.NrOfBann, v.User.OrganisationID, v.User.Expiration, v.Resource.Expiration, string(optUserJson),
				string(optResJson), string(rules))
			require.NoError(t, err)
		})
	}
}

func TestReadPoicy(t *testing.T) {
	papStub := &mocks.FakeSsiddStub{}
	pap := ssidd.PAPSC{}
	ctx := &mocks.FakeTransactionContext{}
	ctx.GetStubReturns(papStub)
	expectedPolicy := policiesTest[0]
	bytes, err := json.Marshal(expectedPolicy)
	require.NoError(t, err)

	papStub.GetStateReturns(bytes, nil)
	p, err := pap.ReadPoicy(ctx, "")
	require.NoError(t, err)
	require.Equal(t, expectedPolicy, p)

	papStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve policy"))
	_, err = pap.ReadPoicy(ctx, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve policy")

	papStub.GetStateReturns(nil, nil)
	p, err = pap.ReadPoicy(ctx, "p0xx")
	require.EqualError(t, err, "the policy p0xx does not exist")
	require.Nil(t, p)
}

func TestUpdatePolicy(t *testing.T) {
	InitPapTest()
	pap := ssidd.PAPSC{}
	papStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(papStub)

	expectedPolicy := policiesTest[0]
	bytes, err := json.Marshal(expectedPolicy)
	require.NoError(t, err)

	papStub.GetStateReturns(bytes, nil)
	optUserJson, err := json.Marshal(expectedPolicy.User.OptionalUserAttr)
	require.NoError(t, err)
	optResJson, err := json.Marshal(expectedPolicy.Resource.OptionalResAttr)
	require.NoError(t, err)
	rules, err := json.Marshal(expectedPolicy.Rules)
	require.NoError(t, err)
	err = pap.UpdatePolicy(transactionContext, "", "", "", "", "", "", "", string(optUserJson), string(optResJson), string(rules))
	require.NoError(t, err)

	papStub.GetStateReturns(nil, nil)
	err = pap.UpdatePolicy(transactionContext, "p0xx", "", "", "", "", "", "", string(optUserJson), string(optResJson), string(rules))
	require.EqualError(t, err, "the policy p0xx does not exist")

	papStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve policy"))
	err = pap.UpdatePolicy(transactionContext, "p0xx", "", "", "", "", "", "", string(optUserJson), string(optResJson), string(rules))
	require.EqualError(t, err, "failed to read from world state: unable to retrieve policy")
}

func TestDeletePoicy(t *testing.T) {
	InitPapTest()
	pap := ssidd.PAPSC{}
	papStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(papStub)

	p := policiesTest[0]
	bytes, err := json.Marshal(p)
	require.NoError(t, err)

	papStub.GetStateReturns(bytes, nil)
	papStub.DelStateReturns(nil)
	err = pap.DeletePolicy(transactionContext, "")
	require.NoError(t, err)

	papStub.GetStateReturns(nil, nil)
	err = pap.DeletePolicy(transactionContext, p.PolicyID)
	require.EqualError(t, err, "the policy p0xx does not exist")

	papStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve policy"))
	err = pap.DeletePolicy(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve policy")
}

func TestGetAllPolicies(t *testing.T) {
	InitPapTest()
	pap := ssidd.PAPSC{}
	pipStub := &mocks.FakeSsiddStub{}
	p := policiesTest[0]
	bytes, err := json.Marshal(p)
	require.NoError(t, err)

	iterator := &mocks.FakeStateIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(pipStub)

	pipStub.GetStateByRangeReturns(iterator, nil)
	policies, err := pap.GetAllPoicies(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*ssidd.Policy{p}, policies)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	policies, err = pap.GetAllPoicies(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, policies)

	pipStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all policies"))
	policies, err = pap.GetAllPoicies(transactionContext)
	require.EqualError(t, err, "failed retrieving all policies")
	require.Nil(t, policies)
}
