FROM golang

COPY . /go/src/gitlab.devtools.intel.com/kubernetes/ci-cd-broker
WORKDIR /go/src/gitlab.devtools.intel.com/kubernetes/ci-cd-broker
RUN go get ./...
