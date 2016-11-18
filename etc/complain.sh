#!/bin/bash

# This script runs all of the complain functions,
# like `go test`, `go vet`, `golint`, and `go fmt`.

set -e

echo "Testing go code..."
go test ./models
go test ./requesthandler
go test

echo "Vetting..."
go vet ./models

echo "Linting..."
golint --set_exit_status ./models

echo "Formatting..."
gofmt -w ./

echo "Checking flow types..."
./app/node_modules/.bin/flow check app/main.js
./app/node_modules/.bin/flow coverage app/main.js

echo "Checking for linter errors..."
./app/node_modules/.bin/eslint -c app/.eslint.json app/main.js

echo
echo
echo "         ------------"
echo "        |            |"
echo "        |    PASS    |"
echo "        |            |"
echo "         ------------"
echo
echo
