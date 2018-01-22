mq-benchmarking
==========================

### Changelog
- Remove environment variable `LATENCY_TEST`. All test case will produce both latency and throughput results

### Environment Variables
- TEST: "nsq"(default)|"zmq"
- CLIENT_MODE:    "consumer"(default), "producer"
- MQ_CONNECTION_STRING: connection string to message queue endpoint
- MESSAGE_COUNT: number of message for perform testing (set to `0` when want to specify duration)
- TEST_DURATION: string of int (milliseconds) for testing (set to `0` when want to specify message count)
- MSG_SIZE_GENERATOR: `uniform`(default), `possion`
- MSG_RATE_GENERATOR: `uniform`(default), `possion`
- MSG_UNIFORM_SIZE: string of int(default 1024), message size (in byte), available only `MSG_SIZE_GENERATOR` is `uniform`
- MSG_UNIFORM_TPS_RATE" string of float(default 1000.0), rate of sending message, available only when `MSG_RATE_GENERATOR` is `uniform`
- MSG_POISSON_AVG_DELAY: string of float(default 500.0) Average delay (between sending message). Available only when `MSG_RATE_GENERATOR`=`poisson`


### TODO
- Configuration topic name (probably use topic name for each test case)
- End signal (by time), currently available only number of messages
- Mount volumn for test result


### Docker build

`git clone https://github.com/green-lantern-id/mq-benchmarking`
`docker build -t green-lantern/mq-benchmarking:0.1 .`

#### Run Message Broker 
- NSQ
    `docker-compose -f docker-compose-nsq.yml up`

#### Run Consumer
`docker run -it --rm -e CLIENT_MODE='consumer' -e MQ_CONNECTION_STRING='somewhere:someport' green-lantern/mq-benchmarking:0.1`

##### Get Latency report
with mounted volume to /var/log
example: `-v /var/log:/var/log`
Report filename: mq_latency.csv


#### Run Producer
`docker run -it --rm -e CLIENT_MODE='producer' -e MQ_CONNECTION_STRING='somewhere:someport' green-lantern/mq-benchmarking:0.1`

NOTE: When run all three component without docker-compose, do not use 'localhost' as a connection string
use your machine real IP instead
