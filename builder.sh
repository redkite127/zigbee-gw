#!/bin/bash

SERVICE=zigbee-gw
BUILD_CONTAINER=build_$SERVICE
GO_VERSION=1.9

docker pull golang:$GO_VERSION
if [[ `docker ps -a | grep $BUILD_CONTAINER` ]]; then
    docker rm -f -v $BUILD_CONTAINER
fi

docker run --rm -i --name $BUILD_CONTAINER -v $PWD:/go/src/$SERVICE  golang:$GO_VERSION /bin/sh -c "cd /go/src/$SERVICE/src ; go get -insecure -d -v ; CGO_ENABLED=0 go build -a -installsuffix cgo -o ../$SERVICE -v"
