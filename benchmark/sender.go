package benchmark

import (
	"encoding/binary"
	"log"
	"time"
)

type MessageSender interface {
	Send([]byte)
}

type SendEndpoint struct {
	MessageSender MessageSender
}


// Merge TestLatency and TestThroughput in one single test
func (endpoint SendEndpoint) TestAll(messageSize int, numberToSend int){
	message := make([]byte, messageSize)
	start := time.Now().UnixNano()

	for i := 0; i < numberToSend; i++ {
		binary.PutVarint(message, time.Now().UnixNano())
		endpoint.MessageSender.Send(message)
	}

	stop := time.Now().UnixNano()
	ms := float32(stop-start)/ 1000000
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
	}

	stop := time.Now().UnixNano()
	ms := float32(stop-start) / 1000000
	log.Printf("Sent %d messages in %f ms\n", numberToSend, ms)
	log.Printf("Sent %f per second\n", 1000*float32(numberToSend)/ms)
}
