#!/bin/bash
docker run -it --rm -e TEST='nsq' \
-e CLIENT_MODE='consumer' \
-e MQ_CONNECTION_STRING='xx.xx.xx.xx:4150' \
-e TOPIC_NAME='testname' \
-v `pwd`:/var/log \
green-lantern/mq-benchmarking:0.1