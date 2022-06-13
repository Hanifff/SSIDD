package ssidd_test

import (
	"fmt"
	"log"
	"testing"

	ssidd "github.com/hanifff/ssiddSC/chaincode"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6  . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ssiddStub
type ssiddStub interface {
	shim.ChaincodeStubInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . stateIterator
type stateIterator interface {
	shim.StateQueryIteratorInterface
}

var (
	attrValuesBool = make([]interface{}, 0)
	attrValuesStr  = make([]interface{}, 0)
	attrValuesStr2 = make([]interface{}, 0)
)
var evalTest = []struct {
	title    string
	expr     string
	attrVals []interface{}
	want     []bool
}{
	{"status", `status == true`, append(attrValuesBool, true, false, true, false),
		[]bool{true, false, true, false}},
	{"organisationid", `organisationid == 'spain01'`, append(attrValuesStr, "Spain01", "spain01", "ENGLAND01", "SPAIN01"),
		[]bool{false, true, false, false}}, // CASE SENSTIVIE
	{"nftid", `nftid == 'nfx001'`, append(attrValuesStr2, "NFXR01", "nfx001", "nfxx01", "GFGX01"),
		[]bool{false, true, false, false}},
}

func TestEvaluate(t *testing.T) {
	for k, v := range evalTest {
		testName := fmt.Sprintf("%d, %s", k, v.title)
		t.Run(testName, func(t *testing.T) {
			for j, attr := range v.attrVals {
				verify, err := ssidd.Evalute(v.title, v.expr, attr)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
				if verify != v.want[j] {
					t.Errorf("Rule evaluation failed: got %v, want %v", verify, v.want[j])
				}
			}
		})
	}
}
