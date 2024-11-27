GO = go
# Define the modules you want to update
GRPC_MODULE=google.golang.org/grpc
PROTOBUF_MODULE=google.golang.org/protobuf
# Test flags
TEST_FLAGS = -v -cover -coverprofile=coverage.out

# HTML coverage report
COVERAGE_HTML = coverage.html


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




# Run tests with coverage
test:
	$(GO) test $(TEST_FLAGS) ./...

# Generate HTML report
coverage: test
	go tool cover -html=coverage.out -o $(COVERAGE_HTML)

# Show coverage in the terminal
coverage-terminal: test
	go tool cover -func=coverage.out

# Clean up generated files
clean:
	rm -f coverage.out $(COVERAGE_HTML)	