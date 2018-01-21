package mq

import (

	"github.com/pebbe/zmq4"
	"github.com/green-lantern-id/mq-benchmarking/benchmark"
)

type Zeromq struct {
	handler  benchmark.MessageHandler
	sender   *zmq4.Socket
	receiver *zmq4.Socket
}

func zeromqReceive(zeromq *Zeromq) {
	for {
		// TODO: Some messages come back empty. Is this a slow-consumer problem?
		// Should DONTWAIT be used?
		message, _ := zeromq.receiver.RecvBytes(zmq4.DONTWAIT)
		if zeromq.handler.ReceiveMessage(message) {
			break
		}
	}
}

func NewZeromq(numberOfMessages int, clientMode string) *Zeromq {
	ctx, _ := zmq4.NewContext()
	pub, _ := ctx.NewSocket(zmq4.PUB)
	sub, _ := ctx.NewSocket(zmq4.SUB)
	if clientMode == "consumer" {
		sub.Connect(getEnv("MQ_CONNECTION_STRING", "tcp://localhost:5555"))
	} else {
		pub.Bind(getEnv("MQ_CONNECTION_STRING", "tcp://*5555"))
	}

	var handler benchmark.MessageHandler

	handler = &benchmark.AllInOneMessageHandler{
		NumberOfMessages: numberOfMessages,
		Latencies: []float32{},
	}

	return &Zeromq{
		handler:  handler,
		sender:   pub,
		receiver: sub,
	}
}

func (zeromq *Zeromq) Setup() {
	// Sleep is needed to avoid race condition with receiving initial messages.
//	time.Sleep(30 * time.Second)
	go zeromqReceive(zeromq)
}

func (zeromq *Zeromq) Teardown() {
	zeromq.sender.Close()
	zeromq.receiver.Close()
}

func (zeromq *Zeromq) Send(message []byte) {
	// TODO: Should DONTWAIT be used? Possibly overloading consumer.
	zeromq.sender.SendBytes(message, zmq4.DONTWAIT)
}

func (zeromq *Zeromq) MessageHandler() *benchmark.MessageHandler {
	return &zeromq.handler
}
