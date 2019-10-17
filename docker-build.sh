#!/bin/bash

cp .dockerignore ..
docker build --no-cache -f Dockerfile \
    -t orbs:build ..

[ "$(docker ps -a | grep orbs_build)" ] && docker rm -f orbs_build

docker run --name orbs_build orbs:build sleep 1

rm -rf _bin && mkdir -p _bin
docker cp orbs_build:/src-plugin/_bin/orbs-javascript-plugin _bin