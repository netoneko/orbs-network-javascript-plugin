#!/bin/bash -xe

docker run -d -p 8080:8080 orbs:gamma-js /opt/orbs/gamma-server -override-config '{"experimental-external-processor-plugin-path": "/opt/orbs/plugins/orbs-javascript-plugin"}'

sleep 5

export E2E=true

export GO111MODULE=on

go get -d -v github.com/orbs-network/orbs-client-sdk-go/codec
go get -d -v github.com/stretchr/testify/require

go test ./test/e2e/... -v
