package benchmark

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/green-lantern-id/mq-benchmarking/benchmark/distribution"
)

type MessageSender interface {
	Send([]byte)
}

type SendEndpoint struct {
	MessageSender MessageSender
}

func (endpoint SendEndpoint) sendMsg(msgSize int, fin int64) {
	message := make([]byte, msgSize)
	binary.PutVarint(message, time.Now().UnixNano())
	binary.PutVarint(message[9:], fin)
	endpoint.MessageSender.Send(message)
}

func (endpoint SendEndpoint) StartDuration(numMsg int, duration int, delayUs int, msgSize int, poissonRate float64, isPoisson bool, finEnabled bool) {
	started := time.Now().UnixNano()
	poisson := distribution.GeneratePoisson(poissonRate)
	msgSizeChan := make(chan int)
	done := make(chan bool)
	delay := delayUs
	doneSign := false

	//Start sampling
	go func() {
		for !doneSign {
			if isPoisson {
				delay = poisson.Sample()
			}
			<-time.After(time.Microsecond * time.Duration(delay))
			msgSizeChan <- msgSize
		}
	}()

	// Start ending clock
	go func() {
		<-time.After(time.Millisecond * time.Duration(duration))
		done <- true
	}()

	// Sending message
	msgCount := 0
	for !doneSign {
		select {
		case msgSize := <-msgSizeChan:
			go endpoint.sendMsg(msgSize, 0)
			msgCount++
		case d := <-done:
			doneSign = d
		}
	}

	// Send fin
	ended := time.Now().UnixNano()
	if finEnabled {
		log.Printf("Sending FIN messages")
		for i := 0; i < 1000; i++ {
			endpoint.sendMsg(1024, 0xff)
			<-time.After(time.Millisecond)
		}
	}

	ms := float32(ended-started) / 1000000
	log.Printf("Time: %f ms", ms)
	log.Printf("Message sent: %d", msgCount)
}

func (endpoint SendEndpoint) StartPoisson(numMsg int, duration int, delayUs int, msgSize int, poissonRate float64, isPoisson bool) {

	// Sample wait time.
	poisson := distribution.GeneratePoisson(poissonRate)
	started := time.Now().UnixNano()
	delay := delayUs
	var ended int64
	if numMsg != 0 { // number of messages mode
		for i := 0; i < numMsg; i++ {
			endpoint.sendMsg(msgSize, 0)
			fmt.Printf("\rMessage sent: %d", i)
			if isPoisson {
				delay = poisson.Sample()
			}
			<-time.After(time.Microsecond * time.Duration(delay))
		}
		// Sent FIN message
		endpoint.sendMsg(msgSize, 0xff)
		fmt.Printf("\nSent FIN\n")

	} else { // assume that duration is not zero
		done := false

		// Set timer
		go func() {
			ticker := time.NewTicker(time.Millisecond * time.Duration(duration)).C
			for {
				<-ticker
				done = true
			}
		}()

		i := 0
		xmsgSentChan := make(chan int)

		go func() {
			for {
				endpoint.sendMsg(msgSize, 0)
				xmsgSentChan <- 1
				if isPoisson {
					delay = poisson.Sample()
				}
				<-time.After(time.Microsecond * time.Duration(delay))
			}
		}()

		for !done {
			<-xmsgSentChan
			i++
		}

		// Sent FIN
		ended = time.Now().UnixNano()

		for j := 0; j < 1000; j++ {
			endpoint.sendMsg(1024, 0xff)
			<-time.After(time.Millisecond)
		}
		fmt.Printf("\nSend FIN (every 1ms for 1 sec)\n")

		log.Printf("Message sent: %d", i)
	}

	ms := float32(ended-started) / 1000000
	log.Printf("Time: %f ms", ms)
}

func (endpoint SendEndpoint) SendMsg(msgSize <-chan int, doneSignal <-chan bool) {
	done := false
	log.Printf("Start sender")
	i := 1
	started := time.Now().UnixNano()
	for done != true {
		select {
		case mSize := <-msgSize:
			message := make([]byte, mSize)
			binary.PutVarint(message, time.Now().UnixNano())
			endpoint.MessageSender.Send(message)
			fmt.Printf("\rMessage Sent: %d", i)
			i++
		case signal := <-doneSignal:
			done = signal
			fmt.Printf("\n")
		}
	}
	ended := time.Now().UnixNano()

	for j := 0; j < 1000; j++ {
		endpoint.sendMsg(1024, 0xff)
		<-time.After(time.Millisecond)
	}
	fmt.Printf("\nSend FIN (every 1ms for 1 sec)\n")
	log.Printf("Message sent: %d", i)
	ms := float32(ended-started) / 1000000
	log.Printf("Elapse time: %f", ms)
}

// Merge TestLatency and TestThroughput in one single test
func (endpoint SendEndpoint) TestAll(messageSize int, numberToSend int) {
	message := make([]byte, messageSize)
	start := time.Now().UnixNano()

	for i := 0; i < numberToSend; i++ {
		binary.PutVarint(message, time.Now().UnixNano())
		endpoint.MessageSender.Send(message)
	}

	stop := time.Now().UnixNano()
	ms := float32(stop-start) / 1000000
	log.Printf("Send %d messages in %f ms\n", numberToSend, ms)
}

func (endpoint SendEndpoint) TestThroughput(messageSize int, numberToSend int) {
	message := make([]byte, messageSize)
	start := time.Now().UnixNano()
	for i := 0; i < numberToSend; i++ {
		endpoint.MessageSender.Send(message)
	}

	stop := time.Now().UnixNano()
	ms := float32(stop-start) / 1000000
	log.Printf("Sent %d messages in %f ms\n", numberToSend, ms)
	log.Printf("Sent %f per second\n", 1000*float32(numberToSend)/ms)
}

func (endpoint SendEndpoint) TestLatency(messageSize int, numberToSend int) {
	start := time.Now().UnixNano()
	b := make([]byte, 9)
	for i := 0; i < numberToSend; i++ {
		binary.PutVarint(b, time.Now().UnixNano())
		endpoint.MessageSender.Send(b)
		log.Printf("Message Sent")
	}

	stop := time.Now().UnixNano()
	ms := float32(stop-start) / 1000000
	log.Printf("Sent %d messages in %f ms\n", numberToSend, ms)
	log.Printf("Sent %f per second\n", 1000*float32(numberToSend)/ms)
}
