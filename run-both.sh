#!/bin/bash

echo "Starting chat server..."
go run ./cmd/server/main.go &
SERVER_PID=$!

cleanup() {
    echo
    echo "Stopping chat server..."
    kill "$SERVER_PID" 2>/dev/null
}

trap cleanup EXIT INT TERM

sleep 1

echo "Starting chat client..."
go run ./cmd/client/client.go