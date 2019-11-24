#!/bin/bash -e

if [[ $CIRCLE_TAG == v* ]] ;
then
  VERSION=$CIRCLE_TAG
else
  VERSION=experimental
fi

VERSION="${VERSION}-js"

docker login -u $DOCKER_HUB_LOGIN -p $DOCKER_HUB_PASSWORD

docker tag orbs:export-js orbsnetwork/node:$VERSION
docker push orbsnetwork/node:$VERSION

docker tag orbs:gamma-js orbsnetwork/gamma:$VERSION
docker push orbsnetwork/gamma:$VERSION
