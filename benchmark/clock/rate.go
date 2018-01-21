package clock

import (
	"time"
	"log"
	"github.com/green-lantern-id/mq-benchmarking/benchmark/distribution"
	"github.com/green-lantern-id/mq-benchmarking/benchmark/generator"
)

func UniformRate(generator generator.MessageGenerator, tps float64, messageCount int) (<-chan int,<-chan bool) {
	msgSizeChan := make(chan int)
	ticker := time.NewTicker(time.Millisecond * 500)
	endSignal := make(chan bool)
	log.Printf("Create uniform clock")
	log.Printf("Messagecount: %d", messageCount)

	go func(){
		defer ticker.Stop()
		for t:= range ticker.C {
			msgSizeChan <- generator.GetMessageSize()
			log.Printf("Time ticking", t)
		}
	}()

	return msgSizeChan, endSignal
}

func PoissonRate(generator generator.MessageGenerator, avgSize float64, messageCount int)(<-chan int,<-chan bool){
	waitTimeGenerator := distribution.GeneratePoisson(avgSize)
	msgSizeChan := make(chan int)
	endSignal := make(chan bool)

	go func(){
		for {
			wait := waitTimeGenerator.Sample()
			<-time.After(time.Microsecond * time.Duration(wait))
			msgSizeChan<-generator.GetMessageSize()
		}
	}()

	return msgSizeChan, endSignal
}