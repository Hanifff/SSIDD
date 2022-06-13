package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	pb "github.com/hanifff/ssidd/protos"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		fmt.Println("err : ", err)
	}

	client := pb.NewSsiddClient(conn)
	for {
		reader := bufio.NewReader(os.Stdout)
		fmt.Print("Message to send: ")
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", -1)

		if text == "exit" || text == "Exit" {
			break
		} else if text == "" {
			continue
		}

		message := &pb.ReadRequest{ClientDID: text}

		returnMessage, err := client.Read(context.Background(), message)
		if err != nil {
			fmt.Println("err : ", err)
		} else {
			fmt.Println("repsone  :", returnMessage)
		}
	}

}
