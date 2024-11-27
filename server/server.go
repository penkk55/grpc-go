package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	pb "github.com/penkk55/grpc-go/generated"
	helper "github.com/penkk55/grpc-go/helper"
	"google.golang.org/grpc"
)

type BaconServer struct {
	pb.UnimplementedBaconServiceServer
	pb.UnimplementedBeefSummaryServiceServer
}

// Implement the GetBacon RPC
func (s *BaconServer) GetBacon(ctx context.Context, req *pb.BaconRequest) (*pb.BaconResponse, error) {
	fmt.Println("baaacon")
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Return the response as a gRPC response
	return &pb.BaconResponse{Text: string(body)}, nil
}

// Implement GetBeefSummary RPC
func (s *BaconServer) GetBeefSummary(ctx context.Context, req *pb.Empty) (*pb.BeefSummaryResponse, error) {

	//call api

	url := "https://baconipsum.com/api/?type=meat-and-filler&paras=99&format=text"
	// url = fmt.Sprintf("https://baconipsum.com/api/?type=%s&paras=%d&format=text", req.Type, req.Paras)

	// Make the HTTP call
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bacon ipsum: %v", err)
	}
	defer resp.Body.Close()
	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Sample meat data string (can be replaced by real input)
	// meatData := "Fatback t-bone t-bone, pastrami .. t-bone. pork, meatloaf jowl enim. Bresaola t-bone."
	meatData := string(body)

	fmt.Println("ffffff")
	// Clean and process the meat names

	meats := helper.CleanMeats(meatData)

	// Count occurrences of each meat
	meatCount := make(map[string]int32)
	for _, meat := range meats {
		// fmt.Println("meee", meat)
		meatCount[meat]++
	}

	// Prepare and return the response
	return &pb.BeefSummaryResponse{Beef: meatCount}, nil
}

func main() {
	// Create a listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register both BaconService and BeefSummaryService
	pb.RegisterBaconServiceServer(grpcServer, &BaconServer{})
	pb.RegisterBeefSummaryServiceServer(grpcServer, &BaconServer{})

	// Start the server
	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
