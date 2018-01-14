package benchmark

import "log"

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
	tester.Setup()
	defer tester.Teardown()

	if tester.TestLatency {
		tester.testLatency()
	} else {
		tester.testThroughput()
	}

	log.Printf("End %s test", tester.Name)
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
