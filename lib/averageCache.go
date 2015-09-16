package carbon

import "fmt"
import "time"

type averageCache struct {
	name string
	subscriptions []Subscription
	inputChan chan DataPoints
	flushChan chan int
	closeChan chan bool
	cache map[int]*cacheSlot
	avgBackend *WhisperBackend
	countBackend *WhisperBackend
}

type cacheSlot struct {
	Average float64
	Count float64
	LastUpdated int
}

// The event loop for this averageCache
// It guarantees sequential inserts and flushes
func (self *averageCache) run() {
	var stop bool
	var flushLimit int
	var dataPoints DataPoints
	for {
		select {
		case stop = <- self.closeChan: // The cache is ordered to close
			self.close()
		case flushLimit = <- self.flushChan: // A flush is queued
			self.flush(flushLimit)
		case dataPoints = <- self.inputChan: // An insert is queued
			self.insert(dataPoints)
		}
	}
}

func (self *averageCache) close() {

}

// Flush all cacheSlots with a mutationTime smaller than flushLimit
func (self *averageCache) flush(flushLimit int) {
	fmt.Println("averageCache %v - INTERNAL FLUSH", self.name)
	var avgFlushTargets DataPoints = make(DataPoints, 2)
	var countFlushTargets DataPoints = make(DataPoints, 2)
	for timeSlot, cacheSlot := range self.cache {
		// TODO: 
		// Write all changes to subscriptions
		if cacheSlot.LastUpdated <= flushLimit {
			avgFlushTargets = append(avgFlushTargets, &DataPoint{timeSlot, cacheSlot.Average})
			countFlushTargets = append(countFlushTargets, &DataPoint{timeSlot, cacheSlot.Count})
			delete(self.cache, timeSlot)
		}
	}
	fmt.Println("FlushAverages: ", avgFlushTargets)
	fmt.Println("FlushCounts: ", countFlushTargets)

	// TODO: Write flush targets to whisper
	//       In another fiber perhaps to make this non-blocking?
	self.avgBackend.Write(avgFlushTargets)
	self.countBackend.Write(countFlushTargets)
}

// Internal insert function. This is called in the runloop for the graphCache to insert dataPoints recieved on the inputChannel
func (self *averageCache) insert(dataPoints DataPoints) {
	var timeSlot int
	var cs *cacheSlot
	var exists bool
	for _, dataPoint := range dataPoints {
		timeSlot = dataPoint.Time - (dataPoint.Time % 10) // TODO: Custom interval
		if cs, exists = self.cache[timeSlot]; exists {
			cs.Average = ((cs.Count * cs.Average) + dataPoint.Value) / (cs.Count + 1)
			cs.Count += 1
			cs.LastUpdated = int(time.Now().Unix())
		} else {
			self.cache[timeSlot] = &cacheSlot{dataPoint.Value, 1, int(time.Now().Unix())}
		}
	}
}

// Flush the cache points in the graphCache to whisper
func (self *averageCache) Flush(flushLimit int) {
	// TODO: Remove a graphCache if it has no cacheSlots and nothing queued in inputChan
	//       Only flush otherwise
	fmt.Printf("GraphCache %v - FLUSH\n", self.name)
	self.flushChan <- flushLimit
}

// Insert a DataPoint into the graphCache
func (self *averageCache) Insert(dataPoints DataPoints) {
	fmt.Printf("GraphCache %v - INSERT\n", self.name)
	self.inputChan <- dataPoints
}

func (self *averageCache) Get(subscription Subscription) {
	subscription.dataChan <- self.avgBackend.Read(subscription.Range)
}

func (self *averageCache) Suscribe(subscription Subscription) {
	self.subscriptions = append(self.subscriptions, subscription)
	self.Get(subscription)
}

func NewAverageCache(name string) averageCache {
	return averageCache{
		name,
		make([]Subscription, 10),
		make(chan DataPoints),
		make(chan int),
		make(chan bool),
		make(map[int]*cacheSlot),
		nil,
		nil}
}