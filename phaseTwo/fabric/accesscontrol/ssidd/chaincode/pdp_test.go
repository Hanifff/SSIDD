package ssidd_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	ssidd "github.com/hanifff/ssiddSC/chaincode"
	mocks "github.com/hanifff/ssiddSC/chaincode/chaincodefakes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
)

type Cli map[string]interface{}

var (
	optsPdpUser      map[string]interface{}
	optsPdpResource  map[string]interface{}
	policiesPdpTest  []*ssidd.Policy
	optsPdpResTest   map[string]interface{}
	resPdpTest       []*ssidd.Resource
	clients          []Cli
	decideTest       = []string{"dec01", "dec02"}
	expectedDecRead  = []bool{false, false, true, true}
	expectedDecWrite = []bool{true, false, true, false}
)

func InitPolicyPdpTest() {
	optsPdpUser = make(map[string]interface{})
	optsPdpResource = make(map[string]interface{})

	status := "Active"
	expiration := "Date of expiration"
	NrOfBann := "Number of times uses is banned"
	usersOrgId := "Organisation ID"
	optsPdpUser["nftid"] = "NFT id of the origanisation of user"
	optsPdpResource["organisationid"] = "Resource owner organisation id"

	p0xx := &ssidd.Policy{PolicyID: "p0xx",
		User:     &ssidd.UserRequirement{Status: status, NrOfBann: NrOfBann, Expiration: expiration, OrganisationID: usersOrgId, OptionalUserAttr: optsPdpUser},
		Resource: &ssidd.ResourceRequirement{Status: status, Expiration: expiration, OptionalResAttr: optsPdpResource},
		Rules: []*ssidd.PolicyRules{
			{Type: "READ", Title: "status", ComparisionType: "bool", Description: "User must be an active user.", Expression: `status == true`},
			{Type: "READ", Title: "nrofbann", ComparisionType: "bool", Description: "User with more than 5 times being banned does not get access.", Expression: `nrofbann <= 5`},
			{Type: "READ", Title: "organisationid", ComparisionType: "bool", Description: "User must be employe in organisation spain01.", Expression: `organisationid == 'spain01'`},
			{Type: "WRITE", Title: "status", ComparisionType: "bool", Description: "User must be an active user.", Expression: `status == true`},
		},
	}
	policiesPdpTest = []*ssidd.Policy{p0xx}
	p0x1 := &ssidd.Policy{PolicyID: "p0x1",
		User:     &ssidd.UserRequirement{Status: status, NrOfBann: NrOfBann, Expiration: expiration, OrganisationID: usersOrgId, OptionalUserAttr: optsPdpUser},
		Resource: &ssidd.ResourceRequirement{Status: status, Expiration: expiration, OptionalResAttr: optsPdpResource},
		Rules: []*ssidd.PolicyRules{
			{Type: "READ", Title: "status", ComparisionType: "bool", Description: "User must be an active user.", Expression: `status == true`},
			{Type: "READ", Title: "nrofbann", ComparisionType: "bool", Description: "User with more than 5 times being banned does not get access.", Expression: `nrofbann <= 5`},
			{Type: "READ", Title: "organisationid", ComparisionType: "bool", Description: "User must be employe in organisation spain01.", Expression: `organisationid == 'spain01'`},
			{Type: "WRITE", Title: "status", ComparisionType: "bool", Description: "User must be an active user.", Expression: `status == true`},
			{Type: "WRITE", Title: "nrofbann", ComparisionType: "bool", Description: "User with more than 5 times being banned does not get access.", Expression: `nrofbann <= 5`},
			{Type: "WRITE", Title: "organisationid", ComparisionType: "bool", Description: "User must be employe in organisation spain01.", Expression: `organisationid == 'spain01'`},
			{Type: "WRITE", Title: "nftid", ComparisionType: "bool", Description: "User's organisation must hold the NFT.", Expression: `nftid == 'nfxx01'`},
		},
	}
	policiesPdpTest = append(policiesPdpTest, p0x1)
}

func InitResPdpTest() {
	optsPdpResTest = make(map[string]interface{})
	optsPdpResTest["nftid"] = "nfxs01"
	resPdpTest = []*ssidd.Resource{{ResouceID: "r001", Attributes: &ssidd.ResourceAttr{
		Status:             true,
		AssociatedPolicyId: "p0Xxx",
		OwnerOrgID:         "spain01",
		Expiration:         time.Now().AddDate(1, 1, 1).Format("02-Jan-2006"), // 1 year, 1 day, 1 month
		OptionalAttr:       optsPdpResTest,
	},
	},
	}
	optsPdpResTest["onlyadmin"] = true
	resPdpTest = append(resPdpTest, &ssidd.Resource{ResouceID: "r002", Attributes: &ssidd.ResourceAttr{
		Status:             false,
		AssociatedPolicyId: "p0x1",
		OwnerOrgID:         "spain01",
		Expiration:         time.Now().AddDate(2, 1, 1).Format("02-Jan-2006"), // 2 year, 1 day, 1 month
		OptionalAttr:       optsPdpResTest,
	},
	})
}

func InitClientPdpTest() {
	client01 := make(map[string]interface{})
	client02 := make(map[string]interface{})

	client01["status"] = true
	client01["expiration"] = time.Now().AddDate(1, 1, 1).Format("02-Jan-2006")
	client01["organisationid"] = "usa01"
	client01["nftid"] = "nfxx01"
	// client02
	client02["status"] = true
	client02["expiration"] = time.Now().AddDate(1, 1, 1).Format("02-Jan-2006")
	client02["organisationid"] = "spain01"
	client02["nftid"] = "nf0xs1" // differnt with client01

	client03 := make(map[string]interface{})
	data := `{"status":"true","expiration": "02-Jan-2026", "organisationid": "spain01","nftid": "nfxx01"}`
	err := json.Unmarshal([]byte(data), &client03)
	if err != nil {
		panic(err)
	}
	clients = append(clients, client01, client02, client03)
}

// decideTester is used for testing purpose
func decideTester(ctx contractapi.TransactionContextInterface, decisionid,
	requestType, clientDID, resourceid string, resouce *ssidd.Resource, policy *ssidd.Policy,
	clientAttrs map[string]interface{}) (*ssidd.Decision, error) {
	decision := &ssidd.Decision{DecisionID: decisionid,
		SubjectID: clientDID, ResourceID: resourceid,
		Timestamp: time.Now().Format("02-Jan-2006"),
	}
	clientAttrs["nrofbann"] = 0
	// check if all required user attributes are provide by user
	if clientAttrs["status"] == nil || clientAttrs["nrofbann"] == nil ||
		clientAttrs["organisationid"] == nil || clientAttrs["expiration"] == nil {
		decision = decisionMakerTester(decisionid, clientDID, resourceid,
			"pdp: missing one or more user attributes required user attribute", "", false)
		return decision, nil
	}

	// validate access based on policy rules and provided client attributes
	for _, rule := range policy.Rules {
		if rule.Type == requestType {
			attrValue, exists := clientAttrs[rule.Title]
			if !exists {
				return nil, fmt.Errorf("pdp: client does not fullfill the requriements for the %v attribute", rule.Title)
			}
			ok, err := ssidd.Evalute(rule.Title, rule.Expression, attrValue)
			if err != nil {
				return decision, fmt.Errorf("pdp: client attribute for the %v is in invalid format\nerr: %v", rule.Title, err)
			}
			if !ok {
				decision = decisionMakerTester(decisionid, clientDID, resourceid,
					fmt.Sprintf("Client does not fulfill requirements for the %v attribute.", rule.Title), "", false)
				return decision, nil
			}
		}
	}
	decisionJSON, err := json.Marshal(decision)
	if err != nil {
		return nil, err
	}
	decision = decisionMakerTester(decisionid, clientDID, resourceid,
		fmt.Sprintf("Client does not fulfill requirements for the %v attribute.", " rule.Title"), "", true)
	return decision, ctx.GetStub().PutState(clientDID, decisionJSON)
}

func decisionMakerTester(decisionid, clientDID, resourceid, description, txnid string, decision bool) *ssidd.Decision {
	return &ssidd.Decision{DecisionID: decisionid,
		SubjectID: clientDID, ResourceID: resourceid,
		Decision: decision, Description: description,
		TxnHash: txnid, Timestamp: time.Now().Format("02-Jan-2006"),
	}
}
func TestDecide(t *testing.T) {
	InitPolicyPdpTest()
	InitResPdpTest()
	InitClientPdpTest()
	clientDID := "did:SPAIN01:123456789abcdefghigklmn"
	pdpStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(pdpStub)
	currentReadDec := 0
	currentWriteDec := 0
	for k, v := range decideTest {
		testName := fmt.Sprintf("%d, %s", k+1, v)
		// READ access
		t.Run(testName, func(t *testing.T) {
			for j, p := range policiesPdpTest {
				decision, err := decideTester(transactionContext, v, "READ", clientDID,
					resPdpTest[k].ResouceID, resPdpTest[j], p, clients[k])
				require.NoError(t, err)
				require.EqualValues(t, expectedDecRead[currentReadDec], decision.Decision)
				currentReadDec++
			}
		})
		// WRITE ACCESS
		t.Run(testName, func(t *testing.T) {
			for j, p := range policiesPdpTest {
				decision, err := decideTester(transactionContext, v, "WRITE", clientDID,
					resPdpTest[k].ResouceID, resPdpTest[j], p, clients[k])
				require.NoError(t, err)
				require.EqualValues(t, expectedDecWrite[currentWriteDec], decision.Decision)
				currentWriteDec++
			}
		})
	}
}

func TestGetAllDecisions(t *testing.T) {
	InitPolicyPdpTest()
	pdpStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(pdpStub)
	pdp := ssidd.PDPSC{}
	des := ssidd.Decision{DecisionID: "someId"}
	bytes, err := json.Marshal(des)
	require.NoError(t, err)

	iterator := &mocks.FakeStateIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	transactionContext.GetStubReturns(pdpStub)

	pdpStub.GetStateByRangeReturns(iterator, nil)
	decisions, err := pdp.GetAllDecisions(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*ssidd.Decision{&des}, decisions)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	decisions, err = pdp.GetAllDecisions(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, decisions)

	pdpStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all decisions"))
	decisions, err = pdp.GetAllDecisions(transactionContext)
	require.EqualError(t, err, "failed retrieving all decisions")
	require.Nil(t, decisions)
}
