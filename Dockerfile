FROM golang

COPY . /go/src/github.com/juliocri/ci-cd-broker
WORKDIR /go/src/github.com/juliocri/ci-cd-broker
RUN go get ./...
