package ssidd_test

import (
	"fmt"
	"log"
	"testing"

	ssidd "github.com/hanifff/ssidd/server"
)

func BenchmarkAddfile(b *testing.B) {
	ipfs := ssidd.NewIpfsConn()
	data := ssidd.RandomString(10000)
	b.Run(fmt.Sprintf("data_size_%d kb", 10), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ipfs.Addfile(data)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkGetfileContent(b *testing.B) {
	ipfs := ssidd.NewIpfsConn()
	data := ssidd.RandomString(10000)
	cid, err := ipfs.Addfile(data)
	if err != nil {
		log.Fatal(err)
	}
	b.Run(fmt.Sprintf("data_size_%d kb", 10), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ipfs.GetfileContent(cid)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}
