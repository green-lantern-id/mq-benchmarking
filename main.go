package main

import (
	"log"
	"os"
	"strconv"

	"github.com/green-lantern-id/mq-benchmarking/benchmark"
	"github.com/green-lantern-id/mq-benchmarking/benchmark/mq"
)

func newTester(subject string, testLatency bool, msgCount, msgSize int, mode string) *benchmark.Tester {
	var messageSender benchmark.MessageSender
	var messageReceiver benchmark.MessageReceiver

	log.Printf("Testing %s", subject)

	switch subject {
	case "nsq":
		nsq := mq.NewNsq(msgCount, mode)
		messageSender = nsq
		messageReceiver = nsq
	case "zeromq":
		zeromq := mq.NewZeromq(msgCount, mode)
		messageSender = zeromq
		messageReceiver = zeromq
	default:
		return nil
	}

	return &benchmark.Tester{
		subject,
		msgSize,
		msgCount,
		testLatency,
		messageSender,
		messageReceiver,
		mode,
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

func parseEnv() (string, bool, int, int, string) {
	test := getEnv("TEST", "nsq")
	messageCount, err := strconv.Atoi(getEnv("MESSAGE_COUNT", "0"))
	messageSize, err := strconv.Atoi(getEnv("MESSAGE_SIZE", "1024"))
	mode := getEnv("CLIENT_MODE", "consumer") // consumer vs producer
	testLatency, err := strconv.ParseBool(getEnv("TEST_LATENCY", "false"))

	if err != nil {
		log.Printf("[ERROR] Cannot get environment variables %s", err)
	}

	return test, testLatency, messageCount, messageSize, mode
}

func main() {

	tester := newTester(parseEnv())
	if tester == nil {
		os.Exit(1)
	}

	tester.Test()
}
