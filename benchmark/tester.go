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
	testDuration,_ := strconv.Atoi(getEnv("TEST_DURATION", "0"))

	if tester.Mode == "producer"{
		log.Printf("Running producer mode")
		sender := &SendEndpoint{MessageSender: tester}
		log.Printf("Start uniform clock rate")

		msgSizeGenerator := getEnv("MSG_SIZE_GENERATOR", "uniform")	// uniform|poisson
		msgRateGenerator := getEnv("MSG_RATE_GENERATOR", "uniform")	// uniform|poisson


		var msgGenerator generator.MessageGenerator

		if msgSizeGenerator == "uniform" {
			msgSize, _ := strconv.Atoi(getEnv("MSG_UNIFORM_SIZE", "1024"))
			msgGenerator = generator.NewUniformGenerator(msgSize)
		}

		if msgRateGenerator == "uniform" {
			uniformRate, _ := strconv.ParseFloat(getEnv("MSG_UNIFORM_TPS_RATE", "1000"), 64)
			msgSizeChan, end := clock.UniformRate(msgGenerator, uniformRate, tester.MessageCount, testDuration)
			sender.Start(msgSizeChan, end)
		} else {	// poisson
			poissonAvgRate, _ := strconv.ParseFloat(getEnv("MSG_POISSON_AVG_DELAY", "500"), 64)
			msgSizeChan, end := clock.PoissonRate(msgGenerator, poissonAvgRate, tester.MessageCount, testDuration)
			sender.Start(msgSizeChan, end)
		}
	} else {
		log.Printf("Running consumer mode")
		receiver := NewReceiveEndpoint(tester, tester.MessageCount)
		receiver.WaitForCompletion()
	}

	log.Printf("End %s test", tester.Name)
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
