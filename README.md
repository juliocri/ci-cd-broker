# CI/CD Broker

## Requisites
* [Docker] (https://docs.docker.com/install/)
* [Docker-compose] (https://docs.docker.com/compose/install/)

### For development and Testing
* GoLang version 1.11+
* [Kafka] (https://kafka.apache.org/quickstart)

## Installation

1. ``` $ git clone https://gitlab.devtools.intel.com/kubernetes/ci-cd-broker```
1. ``` $ cd ci-cd-broker ```
1. ``` $ docker-compose build && docker-compose up ```

NOTE: If it is installed behind a proxy please update the .env file with the corresponding proxy values.

## Testing Broker

To create a message in the right topic:

1. ``` $ kafka-console-producer --broker-list kafka:9092 --topic jenkins-requests```

NOTE: add `kafka` hostname in `/etc/host` and point to the ip address of the
machine where kafka-server is hosted. if running locally, add the next line to
the hosts file `kafka  127.0.0.1`

Below are some methods implemented, write a message with the structure defined:

### Create a project
```json
{
  "id":uuid,
  "action":"create",
  "body":{
    "name": string,
    "description": string
  }
}
```

### Update project
```json
{
  "id":uuid,
  "action":"update",
  "body":{
    "name": string,
    "description": string
  }
}
```

### Delete project
```json
{
  "id": uuid,
  "action":"delete",
  "body":{
    "name": string
  }
}
```

### List projects
```json
{
  "id": uuid,
  "action": "list"
}
```

NOTE: all request must to sent a id with a valid uuid formatting, this id is
required to match a response and identify the response to a related request.

To see all messages produced as a response use next command:

2. ``` $ kafka-console-consumer --bootstrap-server kafka:9092 --topic jenkins-responses --from-beginning```

Optional you can also see the logs using:

3. ``` $ docker-compose logs -f broker```
