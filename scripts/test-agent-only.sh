#!/bin/bash

# Script untuk test agent saja dengan output langsung

set -e

cd "$(dirname "$0")/.."

# Load env
export $(grep -v '^#' .env | xargs)

# Build agent
echo "Building agent..."
cd agent
go build -o ../tmp/agent/main .
cd ..

# Run agent dengan output langsung
echo "Starting agent..."
echo "=================="
./tmp/agent/main
