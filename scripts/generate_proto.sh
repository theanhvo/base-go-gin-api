#!/bin/bash

# Script to generate gRPC protobuf code

echo "Generating gRPC protobuf code..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install Protocol Buffers compiler."
    echo "Installation guide: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Generate protobuf code
echo "Generating Go code from protobuf..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/user_service.proto

echo "✅ Protobuf code generated successfully!"
echo "Generated files:"
echo "  - grpc/user_service.pb.go"
echo "  - grpc/user_service_grpc.pb.go"