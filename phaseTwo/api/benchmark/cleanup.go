//usr/bin/env go run $0 "$@"; exit
package main

// This script cleans up the HLF ledger.
import (
	"fmt"
	"os"
	"strconv"

	ssidd "github.com/hanifff/ssidd/server"
)

func main() {
	gw, err := ssidd.ConfigNet()
	if err != nil {
		fmt.Printf("could not configuere getway to getaway, %v", err)
		panic(err)
	}

	pdpContract := gw.GetContractWithName("ssidd", "PDPSC")
	auditContract := gw.GetContractWithName("ssidd", "AuditSC")
	start, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	reqs, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	fmt.Println("deleting decisions... ", start, reqs)
	for i := start; i <= reqs; i++ {
		fmt.Printf("Dec id: %d\n", i)
		_, err = pdpContract.SubmitTransaction("DeleteDecision", fmt.Sprintf("%d", i))
		if err != nil {
			fmt.Printf("pdp removing decision: %v\n", err)
			continue
		}
		_, err = auditContract.SubmitTransaction("DeleteAudit", fmt.Sprintf("client_%d", i))
		if err != nil {
			fmt.Printf("auidt removing decision: %v\n", err)
			continue
		}
	}
}
