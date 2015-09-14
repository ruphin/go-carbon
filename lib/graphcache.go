package carbon


type Subscription struct {
	closeChan chan bool
	dataChan chan []*DataPoint
	Range *TimeRange
}

type DataPoint struct{
	Value int64
	Timestamp int64
}

type TimeRange struct {
	From int64
	Until int64
}

type GraphCache interface {
	Flush(flushLimit int64)
	Insert(dp *DataPoint)
	Get(timeRange *TimeRange, receiver *ClientConnection)
	Suscribe(timeRange *TimeRange, subscription *Subscription)
}
