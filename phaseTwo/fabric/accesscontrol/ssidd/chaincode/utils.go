package ssidd

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Evalute evaluates and verifies rules of a Policy.
//
// Parameters:
// 	title: declares what should be evaluated
// 	expr: an expression that should be evaluated
// 	attrValue: a CASE SENSTIVIE client attribute
func Evalute(title, expr string, attrValue interface{}) (bool, error) {
	env, err := cel.NewEnv(
		cel.Declarations(decls.NewVar(title, decls.Dyn)))
	if err != nil {
		return false, err
	}
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return false, fmt.Errorf("type-check error: %s", issues.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		return false, fmt.Errorf("program construction error: %s", err)
	}
	result, _, err := prg.Eval(map[string]interface{}{title: attrValue})
	if err != nil {
		return false, fmt.Errorf("eval error: %s", err)
	}
	verify, ok := result.Value().(bool)
	if !ok {
		return false, fmt.Errorf("failed to convert %+v to bool", result)
	}
	return verify, nil
}

// verifyClientOrgMatchesPeerOrg is an internal function used verify client org id and matches peer org id.
func verifyClientOrgMatchesPeerOrg(ctx contractapi.TransactionContextInterface) (bool, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the peer's MSPID: %v", err)
	}

	if clientMSPID != peerMSPID {
		return false, nil
	}

	return true, nil
}
