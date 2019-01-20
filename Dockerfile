FROM golang

COPY . /go/src/github.intel.com/kubernetes/ci-cd-broker
WORKDIR /go/src/github.intel.com/kubernetes/ci-cd-broker
RUN go get ./...
