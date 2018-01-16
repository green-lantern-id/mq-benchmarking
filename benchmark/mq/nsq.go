package mq

import (
	"log"
	"os"

	"github.com/bitly/go-nsq"
	"github.com/green-lantern-id/mq-benchmarking/benchmark"
)

type Nsq struct {
	handler benchmark.MessageHandler
	pub     *nsq.Producer
	sub     *nsq.Consumer
	topic   string
	channel string
}

func NewNsq(numberOfMessages int, testLatency bool) *Nsq {
	topic := getEnv("TOPIC_NAME", "default")
	channel := "test"
	conn := getEnv("MQ_CONNECTION_STRING", "localhost:4150")
	log.Printf("[NSQClient] Connect to %s", conn)
	pub, _ := nsq.NewProducer(conn, nsq.NewConfig())
	sub, _ := nsq.NewConsumer(topic, channel, nsq.NewConfig())

	var handler benchmark.MessageHandler

	handler = &benchmark.AllInOneMessageHandler{
		NumberOfMessages: numberOfMessages,
		Latencies: []float32{},
	}
	/*
	if testLatency {
		handler = &benchmark.LatencyMessageHandler{
			NumberOfMessages: numberOfMessages,
			Latencies:        []float32{},
		}
	} else {
		handler = &benchmark.ThroughputMessageHandler{NumberOfMessages: numberOfMessages}
	}
	*/

	return &Nsq{
		handler: handler,
		pub:     pub,
		sub:     sub,
		topic:   topic,
		channel: channel,
	}
}

func (n *Nsq) Setup() {
	n.sub.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		n.handler.ReceiveMessage(message.Body)
		return nil
	}))
	conn := getEnv("MQ_CONNECTION_STRING", "localhost:4150")
	log.Printf("[NSQClient] Subscribe to %s", conn)
	n.sub.ConnectToNSQD(conn)
}

func (n *Nsq) Teardown() {
	n.sub.Stop()
	n.pub.Stop()
}

func (n *Nsq) Send(message []byte) {
	//n.pub.PublishAsync(n.topic, message, nil)
	n.pub.Publish(n.topic, message)
}

func (n *Nsq) MessageHandler() *benchmark.MessageHandler {
	return &n.handler
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
