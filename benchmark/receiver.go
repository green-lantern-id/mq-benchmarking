package benchmark

import (
	"encoding/binary"
	"log"
	"sync"
	"time"
	"strings"
	"fmt"
	"os"
)

type MessageReceiver interface {
	MessageHandler() *MessageHandler
	Setup()
	Teardown()
}

type ReceiveEndpoint struct {
	MessageReceiver  MessageReceiver
	NumberOfMessages int
	Handler          *MessageHandler
}

func NewReceiveEndpoint(receiver MessageReceiver, numberOfMessages int) *ReceiveEndpoint {
	return &ReceiveEndpoint{
		MessageReceiver:  receiver,
		NumberOfMessages: numberOfMessages,
		Handler:          receiver.MessageHandler(),
	}
}

type MessageHandler interface {
	// Process a received message. Return true if it's the last message, otherwise
	// return false.
	ReceiveMessage([]byte) bool

	// Indicate whether the handler has been marked complete, meaning all messages
	// have been received.
	HasCompleted() bool

}

type AllInOneMessageHandler struct {
	NumberOfMessages int
	Timeout int
	Latencies []float32
	messageCounter int
	hasStarted bool
	hasCompleted bool
	started int64
	stopped int64
	completionLock sync.Mutex
}

func (handler *AllInOneMessageHandler) HasCompleted() bool {
	handler.completionLock.Lock()
	defer handler.completionLock.Unlock()
	return handler.hasCompleted
}

// Merge Latency and Throughput to a single handler + write report to file
func (handler *AllInOneMessageHandler) ReceiveMessage(message []byte) bool {
	now := time.Now().UnixNano()
	if !handler.hasStarted{
		handler.hasStarted = true
		handler.started = time.Now().UnixNano()
		if handler.Timeout != 0 {
			handler.SetTimer()
		}
	}
	// Update message counter
	handler.messageCounter++

	// Record latency
	then, _ := binary.Varint(message[0:9])	// First 8 bytes is sending time
	fin, _ := binary.Varint(message[9:18])	// FIN

	if then != 0 {
		handler.Latencies = append(handler.Latencies, (float32(now-then))/1000000.0)
	}

	if fin != 0 {
		handler.stopped = time.Now().UnixNano()
		handler.WriteReport()
		handler.completionLock.Lock()
		handler.hasCompleted = true
		handler.completionLock.Unlock()
		return true
	}
	return false
}

func (endpoint ReceiveEndpoint) WaitForCompletion() {
	for {
		if (*endpoint.Handler).HasCompleted() {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Set timer to stop (milliseconds)
func (handler *AllInOneMessageHandler) SetTimer() {
	log.Printf("Set consumer timeout: %d ms", handler.Timeout)
	go func(){
		<-time.After(time.Duration(handler.Timeout) * time.Millisecond)
		handler.stopped = time.Now().UnixNano()
		handler.WriteReport()
		handler.completionLock.Lock()
		handler.hasCompleted = true
		handler.completionLock.Unlock()
	}()
}

func (handler *AllInOneMessageHandler) WriteReport(){
	ms := float32(handler.stopped-handler.started)/1000000.0
	fmt.Printf("\n\n")
	log.Printf("Received %d messages in %f ms\n", handler.messageCounter, ms)
	log.Printf("Throughput %f msg per second\n", float32(handler.messageCounter*1000)/ms)


	sum := float32(0)
	for _, latency := range handler.Latencies {
		sum += latency
	}
	avgLatency := float32(sum) / float32(len(handler.Latencies))

	// Write report.csv
	file, err := os.Create("/var/log/mq_latency.csv")
	if err != nil {
		log.Fatal("Cannot create file")
	}
	defer file.Close()

	latencies := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(handler.Latencies)), ","), "[]")

	file.WriteString(latencies)

	log.Printf("Mean latency for %d messages: %f ms\n", handler.messageCounter,
		avgLatency)
}
