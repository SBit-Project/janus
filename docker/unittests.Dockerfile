FROM golang:1.14

WORKDIR $GOPATH/src/github.com/SBit-Project/janus
COPY . $GOPATH/src/github.com/SBit-Project/janus
RUN go get -d ./...

CMD [ "go", "test", "-v", "./..."]