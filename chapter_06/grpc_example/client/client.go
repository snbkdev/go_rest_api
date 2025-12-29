package main

import (
	"context"
	"log"

	pb "grpc_example/generated"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewMoneyTransactionClient(conn)

	from := "1234"
	to := "5678"
	amount := float32(1250.75)

	r, err := c.MakeTransaction(context.Background(), &pb.TransactionRequest{From: from, To: to, Amount: amount})
	if err != nil {
		log.Fatalf("Could not transact: %v", err)
	}

	log.Printf("Transaction confirmed: %t", r.Confirmation)
}