package benchmark

import (
	"log"
	"os"
	"strconv"
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
	if tester.Mode == "consumer" {
		tester.Setup()
	}
	defer tester.Teardown()
	testDuration, _ := strconv.Atoi(getEnv("TEST_DURATION", "0"))
	msgSize, _ := strconv.Atoi(getEnv("MSG_UNIFORM_SIZE", "1024"))
	fin, _ := strconv.ParseBool(getEnv("FIN_ENABLED", "false"))

	if tester.Mode == "producer" {
		log.Printf("Running producer mode")
		sender := &SendEndpoint{MessageSender: tester}

		msgRateGenerator := getEnv("MSG_RATE_GENERATOR", "uniform") // uniform|poisson

		if msgRateGenerator == "uniform" {
			uniformRate, _ := strconv.Atoi(getEnv("MSG_UNIFORM_DELAY_US", "1000"))
			log.Printf("======= Test configuation ======")
			log.Printf("Distribution: Uniform")
			log.Printf("Delay Time: %d micro-seconds", uniformRate)
			log.Printf("Message Size: %d bytes", msgSize)
			log.Printf("================================")
			sender.StartDuration(tester.MessageCount, testDuration, uniformRate, msgSize, 0, false, fin)
		} else { // poisson
			poissonAvgRate, _ := strconv.ParseFloat(getEnv("MSG_POISSON_AVG_DELAY", "500.0"), 64)
			log.Printf("======= Test configuation ======")
			log.Printf("Distribution: Poisson")
			log.Printf("Average Rate: %f micro-seconds", poissonAvgRate)
			log.Printf("Message Size: %d bytes", msgSize)
			log.Printf("================================")
			sender.StartDuration(tester.MessageCount, testDuration, 0, msgSize, poissonAvgRate, true, fin)
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
