package distribution

import (
  "time"

  "math"
  "math/rand"

  //"fmt"
)

type PoissonSample struct {
  Lambda float64
  seeded bool
  cdf []float64
}

type PoissonResult struct {
  messageCount int
  sumDataSize int
  timeUsed int64
}

type sendFunc func(packetSizeInKB int)
type conditionFunc func(startUnixNanoTime int64, result PoissonResult) bool

func GeneratePoisson(lambda float64) PoissonSample {
  instance := PoissonSample {
    Lambda: lambda,
    seeded: false,
    cdf: []float64{math.Pow(math.E,-lambda)}};
  instance.CDF(int(math.Ceil(float64(2)*lambda))); //force pre-compute CDF from 0 to 2*lambda
  return instance;
}

func (p PoissonSample) PMF(k int) float64{
  //cannot use Dynamic Programming memoization,
  //first few PMF may be store as zero (cause by e^-lambda) and corrupt further calculation 
  lambda := p.Lambda;
  result := float64(1);
  //formula: (lambda^k) / ((e^lambda)(k!))
  for i := 1; i <= k ; i++ {
    result *= lambda/float64(i)
    if(float64(i) <= lambda) {
      result /= math.E;
    }
  }
  if(lambda > float64(k)) {
    result /= math.Pow(math.E,lambda - float64(k));
  }
  return result;
  //simple code like following may reach "Infinity" between calculation 
  //return (math.Pow(lambda, kf)*math.Pow(math.E, -lambda))/float64(Fact(k))
}

func (p *PoissonSample) CDF(k int) float64 {
  //memoization
  start := len(p.cdf) - 1;
  if(start < k) {
    for i := start + 1; i <= k; i++ {
      p.cdf = append(p.cdf,p.cdf[i-1] + p.PMF(i));
    }
  }
  return p.cdf[k];
}

func (p *PoissonSample) Sample() int {
  if !p.seeded {
    p.seeded = !p.seeded
    rand.Seed(time.Now().UTC().UnixNano())
  }
  sample := rand.Float64()
  i := 0
  for {
    if sample < p.CDF(i) {
      return i
    }
    i++
  }    
}

func poissonBenchmark(avgSizeInKB float64, avgWaitTimeInMicrosec float64,
  send sendFunc, stopCondition conditionFunc) PoissonResult {

  packetSizeGenerator := GeneratePoisson(avgSizeInKB);
  waitTimeGenerator := GeneratePoisson(avgWaitTimeInMicrosec);
  //fmt.Println("Instances initialized");
  var waitTime,packetSize int;
  result := PoissonResult {
    messageCount: 0,
    sumDataSize: 0,
    timeUsed: 0 }

  start := time.Now().UnixNano();
  for {
    waitTime = waitTimeGenerator.Sample();
    //fmt.Printf("Wait %d microseconds\n",waitTime)
    packetSize = packetSizeGenerator.Sample();
    //fmt.Printf("Packet size %d KB\n",packetSize)

    time.Sleep(time.Duration(waitTime) * time.Microsecond);
    if(stopCondition(start,result)) {
      result.timeUsed = time.Now().UnixNano() - start;
      return result;
    }

    send(packetSize);
    result.messageCount++;
    result.sumDataSize += packetSize;
    if(stopCondition(start,result)) {
      result.timeUsed = time.Now().UnixNano();
      return result;
    }

  }

}
