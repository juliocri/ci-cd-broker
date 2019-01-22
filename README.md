# CI/CD Broker

## Requisites
* Docker
* Docker-compose

### For development and Testing
* GoLang version 1.11+
* Kafka

## Installation

1. ``` $ git clone https://github.intel.com/kubernetes/ci-cd-broker```
1. ``` $ cd ci-cd-broker ```
1. ``` $ docker-compose up ```

## Testing Producer

1. ``` $ kafka-console-producer --broker-list kafka:9092 --topic jenkins.requests```

NOTE: add `kafka` hostname in `/etc/host` and point to the ip address of the
machine where kafka-server is hosted.
