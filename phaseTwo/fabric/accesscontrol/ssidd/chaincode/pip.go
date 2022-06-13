package ssidd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// InitPip initalizes the PIPSC smart contract with a sample resource.
func (pip *PIPSC) InitPip(ctx contractapi.TransactionContextInterface) error {
	opts := make(map[string]interface{})
	opts["nftid"] = "nf00x"
	res := &Resource{
		ResouceID: "r0x1",
		TxnId:     "someUniqueid",
		Attributes: &ResourceAttr{
			Status:             true,
			AssociatedPolicyId: "p001",
			OwnerOrgID:         "spain01",
			Expiration:         time.Now().AddDate(1, 1, 1).Format("02-Jan-2026"),
			OptionalAttr:       opts,
		},
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("r0x1", resJSON)
}

// CreateResource writes a new resource object to the world state with given details.
func (pip *PIPSC) CreateResource(ctx contractapi.TransactionContextInterface, resourceid,
	associatedPolicyId, ownerOrgId, expiration, txnid, txnHash, optAttr string) error {
	exists, err := pip.ResourceExists(ctx, resourceid)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the resouce %s already exists", resourceid)
	}
	resouce := Resource{
		ResouceID: resourceid,
		TxnId:     txnid,
		Attributes: &ResourceAttr{
			Status:             true,
			AssociatedPolicyId: associatedPolicyId,
			OwnerOrgID:         ownerOrgId,
			Expiration:         expiration,
		},
	}
	resouce.Attributes.OptionalAttr = map[string]interface{}{}
	optionalAttrs := make(map[string]interface{})
	err = json.Unmarshal([]byte(optAttr), &optionalAttrs)
	if err != nil {
		return err
	}
	for k, v := range optionalAttrs {
		resouce.Attributes.OptionalAttr[k] = v
	}
	resJSON, err := json.Marshal(resouce)
	if err != nil {
		return err
	}
	db := DBSC{}
	err = db.createTxn(ctx, txnid, txnHash, resourceid)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(resourceid, resJSON)
}

// ReadResource returns the resource object stored in the world state with given resourceid.
func (pip *PIPSC) ReadResource(ctx contractapi.TransactionContextInterface, resourceid string) (*Resource, error) {
	resourceJSON, err := ctx.GetStub().GetState(resourceid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if resourceJSON == nil {
		return nil, fmt.Errorf("the resource %s does not exist", resourceid)
	}
	//var resource Resource
	resAttrs := &ResourceAttr{}
	resAttrs.OptionalAttr = make(map[string]interface{})
	resource := &Resource{Attributes: resAttrs}

	err = json.Unmarshal(resourceJSON, resource)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// UpdateResource updates an existing resource in the world state with provided parameters.
func (pip *PIPSC) UpdateResource(ctx contractapi.TransactionContextInterface,
	resourceid, associatedPolicyId, ownerOrgId, expiration, txnid, optAttr string) error {
	exists, err := pip.ResourceExists(ctx, resourceid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the resource %s does not exist", resourceid)
	}
	resouce := Resource{
		ResouceID: resourceid,
		TxnId:     txnid,
		Attributes: &ResourceAttr{
			Status:             true,
			AssociatedPolicyId: associatedPolicyId,
			OwnerOrgID:         ownerOrgId,
			Expiration:         expiration,
		},
	}
	resouce.Attributes.OptionalAttr = map[string]interface{}{}
	optionalAttrs := make(map[string]interface{})
	err = json.Unmarshal([]byte(optAttr), &optionalAttrs)
	if err != nil {
		return err
	}
	for k, v := range optionalAttrs {
		resouce.Attributes.OptionalAttr[k] = v
	}
	resJSON, err := json.Marshal(resouce)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(resourceid, resJSON)
}

// DeleteResource deletes an given resource from the world state.
func (pip *PIPSC) DeleteResource(ctx contractapi.TransactionContextInterface, resourceid string) error {
	exists, err := pip.ResourceExists(ctx, resourceid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the resource %s does not exist", resourceid)
	}
	return ctx.GetStub().DelState(resourceid)
}

// ResourceExists returns true when resource with given ID exists in world state
func (pip *PIPSC) ResourceExists(ctx contractapi.TransactionContextInterface, resourceid string) (bool, error) {
	resJSON, err := ctx.GetStub().GetState(resourceid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return resJSON != nil, nil
}

// GetAllResources returns all resources found in world state
func (pip *PIPSC) GetAllResources(ctx contractapi.TransactionContextInterface) ([]*Resource, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var resources []*Resource
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var resource Resource
		err = json.Unmarshal(queryResponse.Value, &resource)
		if err != nil {
			return nil, err
		}
		resources = append(resources, &resource)
	}
	return resources, nil
}
