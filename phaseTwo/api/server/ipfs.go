package ssidd

import (
	"fmt"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// IPFSModel is an instance of go-ipfs-api
type IPFSModel struct {
	shell *shell.Shell
}

// NewIpfsConn makes a new connection with IPFS cluster.
func NewIpfsConn() *IPFSModel {
	ipfsApi := shell.NewShell("localhost:5001")
	return &IPFSModel{
		shell: ipfsApi,
	}
}

//	Addfile addes content of file to private ipfs cluster and returns the transaction hash.
// For now: data are not retrived from file, but directly inserted as string.
func (ipfs *IPFSModel) Addfile(fileContent string) (string, error) {
	cid, err := ipfs.shell.Add(strings.NewReader(fileContent))
	if err != nil {
		return "", fmt.Errorf("ipfs: %s", err)
	}
	return cid, nil
}

// Getfile gets the content of file from the ipfs cluster given a transaciton hash.
// For now: It writes the resposne in the current directory
func (ipfs *IPFSModel) GetfileContent(cid string) ([]byte, error) {
	err := ipfs.shell.Get(cid, "./ipfs/readResults")
	if err != nil {
		return nil, fmt.Errorf("ipfs: %s", err)
	}
	// TODO: shold send the file through https connection
	data, err := os.ReadFile("/home/ubuntu/DID/api/ipfs/readResults/" + cid)
	if err != nil {
		return nil, fmt.Errorf("could not find the file localy, %v", err)
	}
	// clean up the directory
	err = os.Remove("/home/ubuntu/DID/api/ipfs/readResults/" + cid)
	if err != nil {
		return nil, fmt.Errorf("could not cleanup the directoy locally, %v", err)
	}
	return data, nil
}

func (ipfs *IPFSModel) UnpintFile(cid string) error {
	return ipfs.shell.Unpin(cid)
}
