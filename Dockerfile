FROM orbsnetwork/v8worker2:latest

ADD ./orbs-network-javascript-plugin /src-plugin

ADD ./orbs-network-go /src/

WORKDIR /src-plugin

RUN cp go.mod.v8worker2 /go/src/github.com/ry/v8worker2/go.mod && \
    cp go.mod.docker go.mod && \
    ./build-binaries.sh