#!/bin/bash
source creds.bashrc

go run admin.go &
ADMIN_PID=$!

go run main.go
SERVER_PID=$!

LOGS_DIR="./logs"
mkdir -p $LOGS_DIR

function cleanup() {
	kill -9 $ADMIN_PID &> "${LOGS_DIR}/admin_logs"
	kill -9 $SERVER_PID &> "${LOGS_DIR}/server_logs" 
}

trap cleanup EXIT

wait
