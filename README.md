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

## Testing Broker

To create a message in the right topic:

1. ``` $ kafka-console-producer --broker-list kafka:9092 --topic jenkins.requests```

NOTE: add `kafka` hostname in `/etc/host` and point to the ip address of the
machine where kafka-server is hosted. if running locally, add the next line to
the hosts file `kafka  127.0.0.1`

Below are some methods implemented, write a message with the structure defined:

### Create a project
```json
{
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
  "action":"delete",
  "body":{
    "name": string
  }
}
```

### List projects
```json
{
  "action": "list"
}
```

To see all messages produced as a response use next command:

2. ``` $ kafka-console-consumer --bootstrap-server kafka:9092 --topic jenkins.responses --from-beginning```

Optional you can also see the logs using:

3. ``` $ docker-compose logs -f broker```
