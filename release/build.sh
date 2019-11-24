#!/bin/bash -e

export PLUGIN_PATH=_bin/orbs-javascript-plugin

if [ ! -f $PLUGIN_PATH ]; then
    echo "Plugin not found: $PLUGIN_PATH does not exist"
fi

rm -rf ./release/_bin
cp -rf _bin ./release/_bin

docker build -f ./release/Dockerfile.node.javascript -t orbs:export-js ./release

docker build -f ./release/Dockerfile.gamma.javascript -t orbs:gamma-js ./release
