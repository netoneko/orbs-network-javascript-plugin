#!/bin/bash -x

rm -rf _bin

cp go.mod.v8worker2 ./vendor/github.com/ry/v8worker2/go.mod

# Pack arguments

time go test ./pack/... -v -run TestPackArguments

# Run tests

time go test ./worker/... -v
time go test ./test/... -v

# Build binary

time go build -o _bin/orbs-javascript-plugin -buildmode=plugin -a main.go
