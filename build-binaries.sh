#!/bin/bash -x

rm -rf _bin

cp go.mod.v8worker2 ./vendor/github.com/ry/v8worker2/go.mod
time go build -o _bin/orbs-javascript-plugin -buildmode=plugin -a main.go
