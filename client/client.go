package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	pb "github.com/penkk55/grpc-go/generated"
	"google.golang.org/grpc"
)

type BaconServer struct {
	pb.UnimplementedBaconServiceServer
	pb.UnimplementedBeefSummaryServiceServer
}

// Implement the GetBacon RPC
func (s *BaconServer) GetBacon(ctx context.Context, req *pb.BaconRequest) (*pb.BaconResponse, error) {
	url := "https://baconipsum.com/api/?type=meat-and-filler&paras=99&format=text"
	if req.Paras > 0 {
		url = fmt.Sprintf("https://baconipsum.com/api/?type=%s&paras=%d&format=text", req.Type, req.Paras)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bacon ipsum: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return &pb.BaconResponse{Text: string(body)}, nil
}

// Clean and count meats in the meat string
func CleanMeats(input string) []string {
	reg := regexp.MustCompile(`[^\w\s]`)
	cleaned := reg.ReplaceAllString(input, "")

	meats := strings.Fields(cleaned)
	return meats
}

// Implement GetBeefSummary RPC
func (s *BaconServer) GetBeefSummary(ctx context.Context, req *pb.Empty) (*pb.BeefSummaryResponse, error) {
	meatData := "Fatback t-bone t-bone, pastrami .. t-bone. pork, meatloaf jowl enim. Bresaola t-bone."
	meats := CleanMeats(meatData)

	meatCount := make(map[string]int32)
	for _, meat := range meats {
		meatCount[meat]++
	}

	return &pb.BeefSummaryResponse{Beef: meatCount}, nil
}

// HTTP handler for the /beef/summary endpoint
func BeefSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// Setup the gRPC client to call the server on port 50051
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewBeefSummaryServiceClient(conn)

	// Call the GetBeefSummary RPC
	// resp, err := client.GetBeefSummary(context.Background(), &pb.Empty{})
	resp, err := client.GetBeefSummary(context.Background(), &pb.Empty{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get beef summary: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare the response (beef summary)
	response := map[string]interface{}{
		"beef": resp.Beef,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert response to JSON and send it
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

// HTTP handler for the /beef/summary endpoint
func BaconHandler(w http.ResponseWriter, r *http.Request) {
	// Setup the gRPC client to call the server on port 50051
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewBaconServiceClient(conn)

	resp, err := client.GetBacon(context.Background(), &pb.BaconRequest{Type: "meat-and-filler", Paras: 99})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get beef summary: %v", err), http.StatusInternalServerError)
		return
	}

	// Call the GetBeefSummary RPC

	// // Prepare the response (beef summary)
	// response := map[string]interface{}{
	// 	"beef": resp.Beef,
	// }

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert response to JSON and send it
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func main() {
	// Create a listener on port 50051 for gRPC
	// listener, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register both BaconService and BeefSummaryService
	pb.RegisterBaconServiceServer(grpcServer, &BaconServer{})
	pb.RegisterBeefSummaryServiceServer(grpcServer, &BaconServer{})

	// Start the gRPC server in a goroutine
	go func() {
		// log.Println("Starting gRPC server on port 50051...")
		// Uncomment the listener to start gRPC server
		// if err := grpcServer.Serve(listener); err != nil {
		// 		log.Fatalf("failed to serve gRPC: %v", err)
		// }
	}()

	// Set up the HTTP server for the RESTful API
	http.HandleFunc("/beef/summary", BeefSummaryHandler)
	http.HandleFunc("/bacon", BaconHandler)

	// Start the HTTP server on port 8080 for REST API
	log.Println("Starting HTTP server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
