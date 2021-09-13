#!/bin/bash -x

cd "$(dirname "$0")"/.. || exit

trap 'exit 130' INT
while true; do
	find . -name "*.go" | entr -dsr "go run ./cmd/aightreader"
done
