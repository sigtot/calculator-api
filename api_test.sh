#!/usr/bin/env bash

go build

echo "Starting API..."
./rest-calculator & # Start the server process
export PID=$! # Save the process id for termination

echo "Sending expression '-1 * (2 * 6 / 3)'..."
echo "Expecting result -4"
result="$(curl -s -X POST -H 'Content-Type: application/json' -d '{"expression": "-1 * (2 * 6 / 3)"}' http://127.0.0.1:8000/api/calc/)" # Test the api
echo "Received response ${result}"

echo "Retrieving last expression evaluated on endpoint /api/history..."
result="$(curl -s http://127.0.0.1:8000/api/history/1)"
echo "Received response ${result}"

echo "Killing server and deleting binary..."
kill -15 ${PID} # Kill the server process
rm rest-calculator # Delete binary

echo "Test complete"
