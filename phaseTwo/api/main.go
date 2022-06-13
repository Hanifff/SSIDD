package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"time"

	pb "github.com/hanifff/ssidd/protos"
	ssidd "github.com/hanifff/ssidd/server"
	"google.golang.org/grpc"
)

func main() {
	var (
		/* address = flag.String("address", "10.0.0.14", "address to listen on") */
		address = flag.String("address", "0.0.0.0", "address to listen on")
		port    = flag.String("port", "8081", "port to listen on")
	)
	// Connect to DB
	/* db, err := ssidd.NewOffchainDB("root", "password", "tcp", "127.0.0.1:3306", "offchain")
	if err != nil {
		log.Printf("Something went wrong while connectig to the database.\n%v", err)
		return
	}
	defer db.Conn.Close() */

	//  connect to fabric gateaway
	_, err := ssidd.ConfigNet()
	if err != nil {
		log.Fatalf("could not configuere getway to getaway, %v", err)
	}

	ssidd.IpfsInstance = ssidd.NewIpfsConn()
	// invoke chain code function for request read access
	initSmartContracts()
	go func() {
		if err := serve(*address, *port); err != nil {
			log.Fatal(err)
		}
	}()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
	log.Println("Shutting down..")
	os.Exit(0)
}

func serve(addr, port string) error {
	ssidd := ssidd.NewSsiddServer(addr, port)
	l, err := net.Listen("tcp", ssidd.Address+":"+ssidd.Port)
	if err != nil {
		return err
	}
	defer l.Close()
	// max 1gb for now, change if more needede
	var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxRecvMsgSize(1024*2024*1024), grpc.MaxSendMsgSize(1024*1024*1024))
	runtime.GOMAXPROCS(runtime.NumCPU())
	gRPCserver := grpc.NewServer(opts...)
	pb.RegisterSsiddServer(gRPCserver, ssidd)
	log.Printf("server %s running", l.Addr())
	if err := gRPCserver.Serve(l); err != nil {
		return err
	}
	//gRPCserver.GracefulStop()
	log.Println("Server have shut down.")
	return nil
}

// initSmartContracts initializes the smart contracts with some sample inputs.
func initSmartContracts() {
	start := time.Now()
	gw, err := ssidd.ConfigNet()
	if err != nil {
		fmt.Printf("could not configuere getway to getaway, %v", err)
		panic(err)
	}
	fmt.Printf("PIP initialization...\n")
	contract := gw.GetContractWithName("ssidd", "PIPSC")
	_, err = contract.SubmitTransaction("InitPip")
	if err != nil {
		fmt.Printf("error initializing pip contract, %v\n", err)
		panic(err)
	}
	// submit data to ipfs and hash to the resoruce
	ipfs := ssidd.NewIpfsConn()
	data := ssidd.RandomString(10000) // 10kb DATA
	cid, err := ipfs.Addfile(data)
	if err != nil {
		log.Fatalf("ipfs: %s", err)
	}
	// if initpip not working
	opts := make(map[string]interface{})
	opts["nftid"] = "nf00x"
	optsJSON, err := json.Marshal(opts)
	if err != nil {
		fmt.Printf("errr marsharling , %v\n", err)
		panic(err)
	}
	contract = gw.GetContractWithName("ssidd", "PIPSC")
	_, err = contract.SubmitTransaction("CreateResource", "r004", "p001", "spain01",
		time.Now().AddDate(1, 1, 1).Format("02-Jan-2006"), "txn02", cid, string(optsJSON))
	if err != nil {
		fmt.Printf("errr calling pip create contract, %v\n", err)
		panic(err)
	}
	res, err := contract.SubmitTransaction("ReadResource", "r004")
	if err != nil {
		fmt.Printf("errr calling contract, %v\n", err)
		panic(err)
	}
	fmt.Printf("Test init pip-remove :%s\n", res)

	fmt.Printf("PAP initialization...\n")
	contract = gw.GetContractWithName("ssidd", "PAPSC")
	_, err = contract.SubmitTransaction("InitPap")
	if err != nil {
		fmt.Printf("error initializing pap contract, %v\n", err)
		panic(err)
	}

	fmt.Printf("DB (smart contract) initialization...\n")
	contract = gw.GetContractWithName("ssidd", "DBSC")
	_, err = contract.SubmitTransaction("InitDB")
	if err != nil {
		fmt.Printf("error initializing db contract, %v\n", err)
		panic(err)
	}

	fmt.Printf("PDP initialization...\n")
	contract = gw.GetContractWithName("ssidd", "PDPSC")
	_, err = contract.SubmitTransaction("InitLedger")
	if err != nil {
		fmt.Printf("error initializing pdp contract, %v\n", err)
		panic(err)
	}

	// Test decide integrity
	contract = gw.GetContractWithName("ssidd", "PDPSC")
	cliattr := make(map[string]interface{})
	cliattr["status"] = true
	cliattr["expiration"] = time.Now().AddDate(1, 1, 1).Format("02-Jan-2006")
	cliattr["organisationid"] = "spain01"
	cliattr["nftid"] = "nfxx01"
	cliAttrsJSON, err := json.Marshal(cliattr)
	if err != nil {
		fmt.Printf("errr marsharling , %v\n", err)
		panic(err)
	}
	hasAcces, err := contract.SubmitTransaction("DecideRead", ssidd.RandomString(3),
		"SOMEDID", "r004", string(cliAttrsJSON))
	if err != nil {
		fmt.Printf("errr calling contract pdp, %v\n", err)
		panic(err)
	}
	fmt.Printf("Testing PIP response:%s\n", hasAcces)

	auditContract := gw.GetContractWithName("ssidd", "AuditSC")
	_, err = auditContract.SubmitTransaction("RecordAudit", "acew:123132asdsad")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Tested Audit\n")
	elapsed := time.Since(start)
	fmt.Printf("SDK initalisations took %s\n", elapsed)
}
