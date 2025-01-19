# Define the binary name
BINARY_NAME = timetable_api

# Define the Go flags
GO_FLAGS = -v

# Build the Go application
build:
	go build $(GO_FLAGS) -o $(BINARY_NAME) main.go

# Run the Go application
run:
	go run $(GO_FLAGS) main.go

# Format the Go code
fmt:
	go fmt ./...

# Tidy up dependencies
tidy:
	go mod tidy

# Test the application
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)

# Build and run the application
all: build run