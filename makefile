# Define the modules you want to update
GRPC_MODULE=google.golang.org/grpc
PROTOBUF_MODULE=google.golang.org/protobuf



# Install dependencies
install:
	go mod tidy


# Fetch and update the grpc and protobuf modules
update-modules:
	go get -u $(GRPC_MODULE)
	go get -u $(PROTOBUF_MODULE)

proto-gen:
	protoc --go_out=. --go-grpc_out=. bacon.proto

dev: 
	go run main.go

run-server:
	go run server/server.go	

run-client: 
	go run client/client.go