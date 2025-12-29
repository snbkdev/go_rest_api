package main

import (
    "context"
    "log"
    "net"

    pb "grpc_example/generated"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

const port = ":50051"

type server struct {
    pb.UnimplementedMoneyTransactionServer
}

func (s *server) MakeTransaction(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {
    log.Printf("Received: $%.2f from %s to %s", req.Amount, req.From, req.To)
    return &pb.TransactionResponse{Confirmation: true}, nil
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    srv := grpc.NewServer()
    pb.RegisterMoneyTransactionServer(srv, &server{})
    reflection.Register(srv)
    
    log.Printf("Server listening on %s", port)
    if err := srv.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}