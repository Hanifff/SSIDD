package ssidd

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pb "github.com/hanifff/ssidd/protos"
)

// SSIDD is an instance of a gRPC server for the access control
type SsiddServer struct {
	mu      sync.Mutex
	Address string
	Port    string
	pb.UnimplementedSsiddServer
}

// AccessDecision is an instance of decision made by PDPSC
type AccessDecision struct {
	DecisionID  string `json:"decisionID"`
	SubjectID   string `json:"subjectID"` // can be clientdid
	ResourceID  string `json:"resourceID"`
	Decision    bool   `json:"decision"`
	Description string `json:"description"`
	TxnHash     string `json:"txnHash"`
	Timestamp   string `json:"timestamp"`
}

// NewSsiddServer intilizes a new instance of go gRPC server.
func NewSsiddServer(address string, port string) *SsiddServer {
	return &SsiddServer{
		Address: address,
		Port:    port,
	}
}

// Read invokes the DecideRead method of the PDPSC smart contract.
// If an access permission was given by PDPSC, this method reads
// data from IPFS.
// Read retunes the result to the gatekeeper.
func (s *SsiddServer) Read(ctx context.Context, readRequest *pb.ReadRequest) (*pb.ReadResponse, error) {
	if readRequest.IsPolicy {
		fmt.Println("got a read  policy request: ", readRequest)
		return nil, nil
	}
	// invoke chain code function for request read access
	pdpContract := GW.GetContractWithName("ssidd", "PDPSC")
	cliAttrsJSON, err := json.Marshal(readRequest.ClientAttributes)
	if err != nil {
		return nil, err
	}
	verify, err := pdpContract.SubmitTransaction("DecideRead", RandomString(4),
		readRequest.ClientDID, readRequest.ResourceID, string(cliAttrsJSON))
	if err != nil {
		return nil, fmt.Errorf("could not invoke chaincode, %v", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	var decision AccessDecision
	err = json.Unmarshal(verify, &decision)
	if err != nil {
		return nil, err
	}
	if !decision.Decision {
		// audit client request
		auditContract := GW.GetContractWithName("ssidd", "AuditSC")
		_, err = auditContract.SubmitTransaction("RecordAudit", decision.SubjectID)
		if err != nil {
			return nil, err
		}
		return &pb.ReadResponse{
			ClientDID:  decision.SubjectID,
			ResourceID: decision.ResourceID,
			Data:       nil,
			Accept:     false,
			Message:    "Access denied",
		}, nil
	}
	// read the data from ipfs
	data, err := IpfsInstance.GetfileContent(decision.TxnHash)
	if err != nil {
		return nil, fmt.Errorf("could not find data in ipfs, %v", err)
	}
	return &pb.ReadResponse{
		ClientDID:  decision.SubjectID,
		ResourceID: decision.ResourceID,
		Data:       data,
		Accept:     true,
		Message:    "you got access to the resource",
	}, nil
}

// Write invokes the DecideWrite method of the PDPSC smart contract.
// If an access permission was given by PDPSC, this method stores
// shared data to IPFS and invokes CreateResource method from the PIPSC smart contract.
// Write retunes the result to the gatekeeper.
func (s *SsiddServer) Write(ctx context.Context, writeRequest *pb.WriteRequest) (*pb.WriteResponse, error) {
	if writeRequest.IsPolicy {
		fmt.Println("got a write policy request: ", writeRequest)
		return nil, nil
	}
	// invoke chain code function for request write access, can be invoked concurrently
	pdpContract := GW.GetContractWithName("ssidd", "PDPSC")
	cliAttrsJSON, err := json.Marshal(writeRequest.ClientAttributes)
	if err != nil {
		return nil, err
	}
	verify, err := pdpContract.SubmitTransaction("DecideWrite", RandomString(4),
		writeRequest.ClientDID, writeRequest.PolicyID, string(cliAttrsJSON))
	if err != nil {
		return nil, fmt.Errorf("could not invoke chaincode, %v", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var decision AccessDecision
	err = json.Unmarshal(verify, &decision)
	if err != nil {
		return nil, err
	}
	if !decision.Decision {
		// audit client request
		auditContract := GW.GetContractWithName("ssidd", "AuditSC")
		_, err = auditContract.SubmitTransaction("RecordAudit", decision.SubjectID)
		if err != nil {
			return nil, err
		}
		return &pb.WriteResponse{
			ClientDID:  decision.SubjectID,
			ResourceID: "",
			Accept:     false,
			Message:    "access denied",
		}, nil
	}
	resAttrsJSON, err := json.Marshal(writeRequest.OptionalResAttributes)
	if err != nil {
		return nil, err
	}
	// write the data to db/ ipfs - more time consuming
	cid, err := IpfsInstance.Addfile(string(writeRequest.Data))
	if err != nil {
		return nil, err
	}
	pipContract := GW.GetContractWithName("ssidd", "PIPSC")
	_, err = pipContract.SubmitTransaction("CreateResource", writeRequest.ResourceID, writeRequest.PolicyID,
		writeRequest.OwnerOrgID, time.Now().AddDate(1, 1, 1).Format("02-Jan-2006"),
		RandomString(4), cid, string(resAttrsJSON))
	if err != nil {
		return nil, err
	}
	return &pb.WriteResponse{
		ClientDID:  decision.SubjectID,
		ResourceID: writeRequest.ResourceID,
		Accept:     true,
		Message:    "thanks for sharing",
	}, nil
}

func (s *SsiddServer) Delete(ctx context.Context, deleteRequest *pb.DeleteRequest) (*pb.DeleteResponse, error) {

	fmt.Println("got a delete request: ", deleteRequest)
	return nil, nil
}

func (s *SsiddServer) Update(ctx context.Context, updateRequest *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	fmt.Println("got a updaterequest: ", updateRequest)
	return nil, nil
}

func (s *SsiddServer) AddOrganisation(ctx context.Context, orgAddRequest *pb.OrgAddRequest) (*pb.OrgAddResponse, error) {

	fmt.Println("got a orgaddrequest: ", orgAddRequest)
	return nil, nil
}

func (s *SsiddServer) RemoveOrganisation(ctx context.Context, orgRemoveRequest *pb.OrgRemoveRequest) (*pb.OrgRemoveResponse, error) {

	fmt.Println("got a orgRemoveRequest: ", orgRemoveRequest)
	return nil, nil
}
