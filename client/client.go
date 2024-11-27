package main

import (
	"context"
	"log"
	"time"

	pb "github.com/penkk55/grpc-go/generated"

	"google.golang.org/grpc"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBaconServiceClient(conn)

	// Make a request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.BaconRequest{
		Type:  "meat-and-filler",
		Paras: 2,
	}

	resp, err := client.GetBacon(ctx, req)
	if err != nil {
		log.Fatalf("failed to call GetBacon: %v", err)
	}

	log.Printf("Response: %s", resp.Text)
}
