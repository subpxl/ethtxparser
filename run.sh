#!/bin/bash

# verify connection to blcokchain node 
curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://127.0.0.1:7545
echo "\n"

echo "\nStarting blockchain parser..."

# Start the main application
cd tx-parser-main
go run cmd/main.go &
APP_PID=$!

# Wait for server to start
sleep 3

curl -X POST "http://localhost:8000/subscribe?address=0xdD93e92dc32d0B2F51430b0e6dA29BDd01AF68D6"



curl -X POST "http://localhos   t:8000/subscribe?address=0xC22c7f8bA7dE381A299ee4EB3a11E1316525ce45"

curl "http://localhost:8000/transactions?address=0xdD93e92dc32d0B2F51430b0e6dA29BDd01AF68D6"

curl "http://localhost:8000/subscribers"



# Cleanup on exit
trap "kill $APP_PID" EXIT

# # Keep script running
# wait $APP_PID

# etherium nodes
# curl -X POST "http://localhost:8000/subscribe?address=0x6d2e03b7EfFEae98BD302A9F836D0d6Ab0002766"
