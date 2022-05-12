FROM golang:1.14-alpine

RUN apk add --no-cache make gcc musl-dev git

WORKDIR $GOPATH/src/github.com/SBit-Project/janus
COPY ./ $GOPATH/src/github.com/SBit-Project/janus

ARG GIT_SHA
ENV GIT_SHA=$GIT_SHA

RUN go build \
        -ldflags \
            "-X 'github.com/qtumproject/janus/pkg/params.GitSha=`./sha.sh`'" \
        -o $GOPATH/bin $GOPATH/src/github.com/SBit-Project/janus/... && \
    rm -fr $GOPATH/src/github.com/SBit-Project/janus/.git

ENV SBIT_RPC=http://sbit:testpasswd@localhost:22002
ENV SBIT_NETWORK=auto

EXPOSE 22402
EXPOSE 23890

ENTRYPOINT [ "janus" ]