package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/green-lantern-id/mq-benchmarking/benchmark"
	"github.com/green-lantern-id/mq-benchmarking/benchmark/mq"
)

func newTester(subject string, testLatency bool, msgCount, msgSize int, mode string) *benchmark.Tester {
	var messageSender benchmark.MessageSender
	var messageReceiver benchmark.MessageReceiver

	switch subject {
	case "nsq":
		nsq := mq.NewNsq(msgCount, testLatency)
		messageSender = nsq
		messageReceiver = nsq
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
	test := getEnv("TEST_NAME", "nsq")
	messageCount, err := strconv.Atoi(getEnv("MESSAGE_COUNT", "10000"))
	messageSize, err := strconv.Atoi(getEnv("MESSAGE_SIZE", "1024"))
	mode := getEnv("CLIENT_MODE", "consumer") // consumer vs producer
	testLatency, err := strconv.ParseBool(getEnv("TEST_LATENCY", "false"))

	if err != nil {
		log.Printf("[ERROR] Cannot get environment variables %s", err)
	}

	return test, testLatency, messageCount, messageSize, mode
}

func parseArgs(usage string) (string, bool, int, int) {

	if len(os.Args) < 2 {
		log.Print(usage)
		os.Exit(1)
	}

	test := os.Args[1]
	messageCount := 1000000
	messageSize := 1000
	testLatency := false

	if len(os.Args) > 2 {
		latency, err := strconv.ParseBool(os.Args[2])
		if err != nil {
			log.Print(usage)
			os.Exit(1)
		}
		testLatency = latency
	}

	if len(os.Args) > 3 {
		count, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Print(usage)
			os.Exit(1)
		}
		messageCount = count
	}

	if len(os.Args) > 4 {
		size, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Print(usage)
			os.Exit(1)
		}
		messageSize = size
	}

	return test, testLatency, messageCount, messageSize
}

func main() {
	usage := fmt.Sprintf(
		"usage: %s "+
			"{"+
			"zeromq|"+
			"nsq|"+
			"} "+
			"[test_latency] [num_messages] [message_size]",
		os.Args[0])

	tester := newTester(parseEnv())
	if tester == nil {
		log.Println(usage)
		os.Exit(1)
	}

	tester.Test()
}
