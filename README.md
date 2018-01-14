mq-benchmarking
==========================

### Environment Variables
- CLIENT_MODE:    "consumer"(default), "producer"
- MQ_CONNECTION_STRING: connection string to message queue endpoint
- MESSAGE_COUNT: number of message for perform testing
- MESSAGE_SIZE: size of message (byte)
- LATENCY_TEST: "true", "false"(default) -- string of boolean

### TODO
- add ZMQ
- Configuration topic name (probably use topic name for each test case)


### Docker build

`git clone https://github.com/green-lantern-id/mq-benchmarking`
`docker build -t green-lantern/mq-benchmarking:0.1 .`

#### Run Message Broker 
- NSQ
    `docker-compose -f docker-compose-nsq.yml up`

#### Run Consumer
`docker run -it --rm -e CLIENT_MODE='consumer' -e MQ_CONNECTION_STRING='somewhere:someport' green-lantern/mq-benchmarking:0.1`

#### Run Producer
`docker run -it --rm -e CLIENT_MODE='producer' -e MQ_CONNECTION_STRING='somewhere:someport' green-lantern/mq-benchmarking:0.1`

NOTE: When run all three component without docker-compose, do not use 'localhost' as a connection string
use your machine real IP instead
