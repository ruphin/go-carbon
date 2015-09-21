package cache

import (
	"time"
	log "github.com/Sirupsen/logrus"
	models "github.com/ruphin/go-carbon/models"
	models "github.com/ruphin/go-carbon/backend"
)

type averageCache struct {
	name string
	subscriptions models.Subscriptions
	inputChan chan models.DataPoints
	flushChan chan int
	closeChan chan bool
	cache map[int]*cacheSlot
	valueBackend *backend.WhisperBackend
	countBackend *backend.WhisperBackend
}

type cacheSlot struct {
	Value float64
	Count float64
	LastUpdated int
}

// The event loop for this averageCache
// It guarantees sequential inserts and flushes
func (self *averageCache) run() {
	var flushLimit int
	var dataPoints []*whisper.TimeSeriesPoint
	for {
		select {
		case <- self.closeChan: // The cache is ordered to close
			log.Debug("Close signal")
			self.close()
		case flushLimit = <- self.flushChan: // A flush is queued
			log.Debug("Flush Signal")
			self.flush(flushLimit)
		case dataPoints = <- self.inputChan: // An insert is queued
			log.Debug("Data Signal")
			self.insert(dataPoints)
		}
	}
}

func (self *averageCache) close() {

}

// Flush all cacheSlots with a mutationTime smaller than flushLimit
func (self *averageCache) flush(flushLimit int) {
	log.WithFields(log.Fields{
		"cache": self.name,
	}).Debug("Internal Flush")
	var valueFlushTargets []*whisper.TimeSeriesPoint = make([]*whisper.TimeSeriesPoint, 0)
	var countFlushTargets []*whisper.TimeSeriesPoint = make([]*whisper.TimeSeriesPoint, 0)
	for timeSlot, cacheSlot := range self.cache {
		// TODO: 
		// Write all changes to subscriptions
		if cacheSlot.LastUpdated <= flushLimit {
			valueFlushTargets = append(valueFlushTargets, &whisper.TimeSeriesPoint{timeSlot, cacheSlot.Value})
			countFlushTargets = append(countFlushTargets, &whisper.TimeSeriesPoint{timeSlot, cacheSlot.Count})
			delete(self.cache, timeSlot)
		}
	}
	log.Debug("FlushAverages: ", valueFlushTargets)
	log.Debug("FlushCounts: ", countFlushTargets)

	// TODO: Write flush targets to whisper
	//       In another fiber perhaps to make this non-blocking?
	self.valueBackend.Write(valueFlushTargets)
	self.countBackend.Write(countFlushTargets)
}

// Internal insert function. This is called in the runloop for the graphCache to insert dataPoints recieved on the inputChannel
func (self *averageCache) insert(dataPoints []*whisper.TimeSeriesPoint) {
	log.WithFields(log.Fields{
		"cache": self.name,
	}).Debug("Internal Insert")
	var timeSlot int
	var cs *cacheSlot
	var exists bool
	for _, dataPoint := range dataPoints {
		timeSlot = dataPoint.Time - (dataPoint.Time % 10) // TODO: Custom interval
		if cs, exists = self.cache[timeSlot]; exists {
			cs.Value += dataPoint.Value
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
	log.WithFields(log.Fields{
		"cache": self.name,
	}).Debug("Flush")
	self.flushChan <- flushLimit
}

// Insert a DataPoint into the graphCache
func (self *averageCache) Insert(dataPoints []*whisper.TimeSeriesPoint) {
	log.WithFields(log.Fields{
		"cache": self.name,
	}).Debug("Insert")
	self.inputChan <- dataPoints
}

func (self *averageCache) Get(subscription Subscription) {
	values := self.valueBackend.Read(subscription.Range)
	counts := self.countBackend.Read(subscription.Range)
	var result = make([]*whisper.TimeSeriesPoint, len(values))
	for i, value := range values {
		result[i] = &whisper.TimeSeriesPoint{value.Time, value.Value / counts[i].Value}
	}
	subscription.dataChan <- result
}

func (self *averageCache) Suscribe(subscription Subscription) {
	self.subscriptions = append(self.subscriptions, subscription)
	self.Get(subscription)
}

func NewAverageCache(name string) graphCache {
	gc := &averageCache{
		name,
		make([]Subscription, 10),
		make(chan []*whisper.TimeSeriesPoint),
		make(chan int),
		make(chan bool),
		make(map[int]*cacheSlot),
		nil,
		nil}
	go gc.run()
	return graphCache(gc)
}