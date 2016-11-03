#!/bin/bash

# This script runs all of the complain functions,
# like `go test`, `go vet`, `golint`, and `go fmt`.

set -e

echo "TESTING..."
go test ./models

echo "VETTING..."
go vet ./models
echo "OK"

echo "LINTING..."
$GOPATH/bin/golint --set_exit_status ./models

echo "FORMATTING..."
gofmt -w ./
echo "DONE"
