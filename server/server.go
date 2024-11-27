package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	pb "github.com/penkk55/grpc-go/generated"

	"google.golang.org/grpc"
)

type BaconServer struct {
	pb.UnimplementedBaconServiceServer
}

// Implement the GetBacon RPC
func (s *BaconServer) GetBacon(ctx context.Context, req *pb.BaconRequest) (*pb.BaconResponse, error) {
	// Build the URL
	url := "https://baconipsum.com/api/?type=meat-and-filler&paras=99&format=text"
	if req.Paras > 0 {
		url = fmt.Sprintf("https://baconipsum.com/api/?type=%s&paras=%d&format=text", req.Type, req.Paras)
	}

	// Make the HTTP call
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bacon ipsum: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Return the response as a gRPC response
	return &pb.BaconResponse{Text: string(body)}, nil
}

func main() {
	// Create a listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterBaconServiceServer(grpcServer, &BaconServer{})

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
