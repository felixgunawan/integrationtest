FROM golang:latest 

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/felixgunawan/integrationtest/examples/standalone
WORKDIR /go/src/github.com/felixgunawan/integrationtest/examples/standalone

ADD . /go/src/github.com/felixgunawan/integrationtest/examples/standalone

RUN dep ensure -v