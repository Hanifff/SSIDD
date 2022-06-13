package ssidd

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var GW *gateway.Network
var IpfsInstance *IPFSModel

func ConfigNet() (*gateway.Network, error) {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		return nil, fmt.Errorf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		return nil, fmt.Errorf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			return nil, fmt.Errorf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"phaseTwo",
		"fabric",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	// docker
	/* gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("/var/keys/connection-org1.yaml")),
		gateway.WithIdentity(wallet, "appUser"),
	)
	*/
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return nil, fmt.Errorf("Failed to get network: %v", err)
	}
	GW = network
	return network, nil
}

// populateWallet populates a wallet in case it is not configuered already.
func populateWallet(wallet *gateway.Wallet) error {
	// local
	credPath := filepath.Join(
		"..",
		"..",
		"phaseTwo",
		"fabric",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}
	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}
	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	return wallet.Put("appUser", identity)
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// TODO:  test: In case clients will be enrolled with their DID to MSP..
/*
var ADMIN = "did:SPAIN01:" + RandomString(10)
var CLIENT = "did:SPAIN01:" + RandomString(10)

func RegisterUsers(ctx context.ClientProvider) {
	cli, err := msp.New(ctx,
		msp.WithOrg("org1.example.com"))
	if err != nil {
		log.Fatalf("an error occured while creating a msp client, %v", err)
	}
	secr, err := cli.Register(&msp.RegistrationRequest{Name: CLIENT})
	if err != nil {
		log.Fatalf("an error occured while registering user, %v", err)
	}
	err = cli.Enroll(CLIENT, msp.WithSecret(secr))
	if err != nil {
		log.Fatalf("an error occured while enrolling user, %v", err)
	}
} */
