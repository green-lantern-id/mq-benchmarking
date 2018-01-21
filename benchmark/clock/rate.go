package clock

import (
	"time"
	"log"
	"github.com/green-lantern-id/mq-benchmarking/benchmark/distribution"
	"github.com/green-lantern-id/mq-benchmarking/benchmark/generator"
)

func UniformRate(g generator.MessageGenerator, tps float64, messageCount int, duration int) (<-chan int,<-chan bool) {
	msgSizeChan := make(chan int)
	ticker := time.NewTicker(time.Millisecond * time.Duration(1000.0/tps))
	endSignal := make(chan bool)

	if messageCount > 0 {
		go func(){
			for i:=0; i<messageCount; i++ {
				t:=<-ticker.C
				msgSizeChan <-g.GetMessageSize()
				log.Printf(t.String())
			}

			endSignal<-true

			log.Printf("Stop uniform rate clock")
		}()
	} else {
		done:= make(chan bool)
		stop:=false

		go func(){
			<-time.After(time.Millisecond * time.Duration(duration))
			done<-true
		}()

		go func(){
			for !stop {
				select {
					case <- ticker.C:
						msgSizeChan <- g.GetMessageSize()
					case d:= <- done:
						stop=d
				}
			}

			endSignal<-true
			log.Printf("Stop uniform rate clock")
		}()
	}
	

	return msgSizeChan, endSignal
}

func PoissonRate(g generator.MessageGenerator, avgSize float64, messageCount int, duration int)(<-chan int,<-chan bool){
	waitTimeGenerator := distribution.GeneratePoisson(avgSize)
	msgSizeChan := make(chan int)
	endSignal := make(chan bool)

	if messageCount > 0{
		go func(){
			for i:=0; i<messageCount; i++ {
				wait := waitTimeGenerator.Sample()
				<-time.After(time.Microsecond * time.Duration(wait))
				msgSizeChan<-g.GetMessageSize()
			}
			endSignal<-true
			log.Printf("Stop poisson rate clock")
		}()
	} else {
		done:= make(chan bool)
		stop:=false
		go func(){
			<-time.After(time.Duration(duration) * time.Millisecond)
			done<-true
		}()

		go func() {
			for !stop{
				wait := waitTimeGenerator.Sample()
				select{
					case <-time.After(time.Microsecond * time.Duration(wait)):
						msgSizeChan<-g.GetMessageSize()
					case d:= <- done:
						stop=d
				}
			}
			endSignal<-true
			log.Printf("Stop poisson rate clock")
		}()
	}

	

	return msgSizeChan, endSignal
}