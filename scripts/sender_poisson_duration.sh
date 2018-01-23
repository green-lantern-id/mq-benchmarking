#!/bin/bash
docker run -it --rm -e TEST='nsq' \
-e CLIENT_MODE='producer' \
-e MQ_CONNECTION_STRING='xx.xx.xx.xx:4150' \
-e MSG_UNIFORM_SIZE='1024' \
-e TEST_DURATION='300000' \
-e MSG_RATE_GENERATOR='poisson' \
-e MSG_POISSON_AVG_DELAY='100000' \
-e TOPIC_NAME='testname' \
green-lantern/mq-benchmarking:0.1