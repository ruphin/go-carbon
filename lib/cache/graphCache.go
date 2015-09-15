package cache

type Subscription struct {
	closeChan chan bool
	dataChan chan []DataPoint
	Range TimeRange
}

type DataPoint struct {
	Time  int
	Value float64
}

type DataPoints []*DataPoint

type TimeRange struct {
	From int
	Until int
}

type GraphCache interface {
	Flush(flushLimit int)
	Insert(dataPoints DataPoints)
	Get(subscription Subscription)
	Suscribe(subscription Subscription)
}
