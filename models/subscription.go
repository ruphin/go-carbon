package models

type Subscription struct {
	CloseChan chan bool
	DataChan chan DataPoints
	Range TimeRange
}

type Subscriptions []*Subscription

func (self *Subscription) Close() {
	self.CloseChan <- true
}