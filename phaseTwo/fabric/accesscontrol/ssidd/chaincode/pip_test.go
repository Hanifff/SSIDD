package ssidd_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	ssidd "github.com/hanifff/ssiddSC/chaincode"
	mocks "github.com/hanifff/ssiddSC/chaincode/chaincodefakes"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
)

var (
	opts          map[string]interface{}
	resourcesTest []*ssidd.Resource
	txnHashes     = []string{"hash1", "hash2"}
)

func InitPipTest() {
	opts = make(map[string]interface{})
	opts["nftid"] = "nf00x"
	resourcesTest = []*ssidd.Resource{{ResouceID: "r001", TxnId: "aha",
		Attributes: &ssidd.ResourceAttr{
			Status:             true,
			AssociatedPolicyId: "p0xx",
			OwnerOrgID:         "spain01",
			Expiration:         time.Now().AddDate(1, 1, 1).Format("02-Jan-2006"), // 1 year, 1 day, 1 month
			OptionalAttr:       opts,
		},
	},
	}
	opts["onlyadmin"] = true
	resourcesTest = append(resourcesTest, &ssidd.Resource{ResouceID: "r002", TxnId: "chs",
		Attributes: &ssidd.ResourceAttr{
			Status:             false,
			AssociatedPolicyId: "p0x2",
			OwnerOrgID:         "spain01",
			Expiration:         time.Now().AddDate(2, 1, 1).Format("02-Jan-2006"), // 1 year, 1 day, 1 month
			OptionalAttr:       opts,
		},
	})
}

func TestInitPip(t *testing.T) {
	pipStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(pipStub)
	pip := ssidd.PIPSC{}
	err := pip.InitPip(transactionContext)
	require.NoError(t, err)
}

func TestCreateResource(t *testing.T) {
	InitPipTest()
	pipStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(pipStub)
	pip := ssidd.PIPSC{}
	for k, v := range resourcesTest {
		testName := fmt.Sprintf("%d, %s", k, v.ResouceID)
		t.Run(testName, func(t *testing.T) {
			optResJson, err := json.Marshal(v.Attributes.OptionalAttr)
			require.NoError(t, err)
			err = pip.CreateResource(transactionContext, v.ResouceID, v.Attributes.AssociatedPolicyId, v.Attributes.OwnerOrgID,
				v.Attributes.Expiration, v.TxnId, txnHashes[k], string(optResJson))
			require.NoError(t, err)
		})
	}
}

func TestReadResource(t *testing.T) {
	pipStub := &mocks.FakeSsiddStub{}
	pip := ssidd.PIPSC{}
	ctx := &mocks.FakeTransactionContext{}
	ctx.GetStubReturns(pipStub)
	expectedResource := resourcesTest[0]
	bytes, err := json.Marshal(expectedResource)
	require.NoError(t, err)

	pipStub.GetStateReturns(bytes, nil)
	res, err := pip.ReadResource(ctx, "")
	require.NoError(t, err)
	require.Equal(t, expectedResource, res)

	pipStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve resource"))
	_, err = pip.ReadResource(ctx, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve resource")

	pipStub.GetStateReturns(nil, nil)
	res, err = pip.ReadResource(ctx, "r001")
	require.EqualError(t, err, "the resource r001 does not exist")
	require.Nil(t, res)
}

func TestUpdateResource(t *testing.T) {
	InitPipTest()
	pip := ssidd.PIPSC{}
	pipStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(pipStub)

	expectedResource := resourcesTest[0]
	bytes, err := json.Marshal(expectedResource)
	require.NoError(t, err)

	pipStub.GetStateReturns(bytes, nil)
	optResJson, err := json.Marshal(expectedResource.Attributes.OptionalAttr)
	require.NoError(t, err)
	err = pip.UpdateResource(transactionContext, "", "", "", "", "", string(optResJson))
	require.NoError(t, err)

	pipStub.GetStateReturns(nil, nil)
	err = pip.UpdateResource(transactionContext, "r001", "", "", "", "", string(optResJson))
	require.EqualError(t, err, "the resource r001 does not exist")

	pipStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve resource"))
	err = pip.UpdateResource(transactionContext, "r001", "", "", "", "", string(optResJson))
	require.EqualError(t, err, "failed to read from world state: unable to retrieve resource")
}

func TestDeleteResource(t *testing.T) {
	InitPipTest()
	pip := ssidd.PIPSC{}
	pipStub := &mocks.FakeSsiddStub{}
	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(pipStub)

	res := resourcesTest[0]
	bytes, err := json.Marshal(res)
	require.NoError(t, err)

	pipStub.GetStateReturns(bytes, nil)
	pipStub.DelStateReturns(nil)
	err = pip.DeleteResource(transactionContext, "")
	require.NoError(t, err)

	pipStub.GetStateReturns(nil, nil)
	err = pip.DeleteResource(transactionContext, res.ResouceID)
	require.EqualError(t, err, "the resource r001 does not exist")

	pipStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve resource"))
	err = pip.DeleteResource(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve resource")
}

func TestGetAllResources(t *testing.T) {
	InitPipTest()
	pip := ssidd.PIPSC{}
	pipStub := &mocks.FakeSsiddStub{}
	res := resourcesTest[0]
	bytes, err := json.Marshal(res)
	require.NoError(t, err)

	iterator := &mocks.FakeStateIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	transactionContext := &mocks.FakeTransactionContext{}
	transactionContext.GetStubReturns(pipStub)

	pipStub.GetStateByRangeReturns(iterator, nil)
	resources, err := pip.GetAllResources(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*ssidd.Resource{res}, resources)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	resources, err = pip.GetAllResources(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, resources)

	pipStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all resources"))
	resources, err = pip.GetAllResources(transactionContext)
	require.EqualError(t, err, "failed retrieving all resources")
	require.Nil(t, resources)
}
