package benchmark

import (
	"log"
	"strconv"
	"os"
	"github.com/green-lantern-id/mq-benchmarking/benchmark/clock"
	"github.com/green-lantern-id/mq-benchmarking/benchmark/generator"
)

type Tester struct {
	Name         string
	MessageSize  int
	MessageCount int
	TestLatency  bool
	MessageSender
	MessageReceiver
	Mode string
}

func (tester Tester) Test() {
	log.Printf("Begin %s test", tester.Name)
	if tester.Mode == "consumer"{
		tester.Setup()
	}
	defer tester.Teardown()

	if tester.Mode == "producer"{
		log.Printf("Running producer mode")
		sender := &SendEndpoint{MessageSender: tester}
		log.Printf("Start uniform clock rate")

		msgSizeGenerator := getEnv("MSG_SIZE_GENERATOR", "uniform")	// uniform|poisson
		msgRateGenerator := getEnv("MSG_RATE_GENERATOR", "uniform")	// uniform|poisson
		testDuration,_ := strconv.Atoi(getEnv("TEST_DURATION", "0"))

		var msgGenerator generator.MessageGenerator

		if msgSizeGenerator == "uniform" {
			msgSize, _ := strconv.Atoi(getEnv("MSG_UNIFORM_SIZE", "1024"))
			msgGenerator = generator.NewUniformGenerator(msgSize)
		}


	//	var msgSizeChan chan int
	//	var end chan bool
		if msgRateGenerator == "uniform" {
			uniformRate, _ := strconv.ParseFloat(getEnv("MSG_UNIFORM_TPS_RATE", "1000"), 64)
			msgSizeChan, end := clock.UniformRate(msgGenerator, uniformRate, tester.MessageCount, testDuration)
			sender.Start(msgSizeChan, end)
		} else {	// poisson
			poissonAvgRate, _ := strconv.ParseFloat(getEnv("MSG_POISSON_AVG_DELAY", "500"), 64)
			msgSizeChan, end := clock.PoissonRate(msgGenerator, poissonAvgRate, tester.MessageCount, testDuration)
			sender.Start(msgSizeChan, end)
		}

		//sender.Start(msgSizeChan, end)
	} else {
		log.Printf("Running consumer mode")
		receiver := NewReceiveEndpoint(tester, tester.MessageCount)
		receiver.WaitForCompletion()
	}

	log.Printf("End %s test", tester.Name)
}

func (tester Tester) testAll(){
	if tester.Mode == "consumer" {
		log.Printf("[Consumer] Running test throughput in consumer mode")
		receiver := NewReceiveEndpoint(tester, tester.MessageCount)
		receiver.WaitForCompletion()
	} else {
		log.Printf("[Producer] Running test throughput in producer mode")
		sender := &SendEndpoint{MessageSender: tester}
		sender.TestAll(tester.MessageSize, tester.MessageCount)
	}
}

func (tester Tester) testThroughput() {
	if tester.Mode == "consumer" {
		log.Printf("[Consumer] Running test throughput in consumer mode")
		receiver := NewReceiveEndpoint(tester, tester.MessageCount)
		receiver.WaitForCompletion()
	} else {
		log.Printf("[Producer] Running test throughput in producer mode")
		sender := &SendEndpoint{MessageSender: tester}
		sender.TestThroughput(tester.MessageSize, tester.MessageCount)
	}
}

func (tester Tester) testLatency() {
	if tester.Mode == "consumer" {
		log.Printf("[Consumer] Running test latency in consumer mode")
		receiver := NewReceiveEndpoint(tester, tester.MessageCount)
		receiver.WaitForCompletion()
	} else { // producer mode
		log.Printf("[Producer] Running test latency in producer mode")
		sender := &SendEndpoint{MessageSender: tester}
		sender.TestLatency(tester.MessageSize, tester.MessageCount)
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
