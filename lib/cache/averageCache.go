
type averageCache struct {
	name String
	subscriptions []*Subscription
	inputChan chan *DataPoint
	flushChan chan int64
	closeChan chan bool
	cache map[int64]*cacheSlot
	avgBackend WhisperBackend
	countBackend WhisperBackend
}

type cacheSlot struct {
	average int64
	count int64
	lastUpdated int64
}
// The event loop for this averageCache
// It guarantees sequential inserts and flushes
func (self *averageCache) run() {
	var dataPoint *DataPoint
	var flushLimit int64
	for {
		select {
		case stop = <- closeChan: // The cache is ordered to close
			self.close()
		case flushLimit = <- self.flushChan: // A flush is queued
			self.flush(flushLimit)
		case dataPoint = <- self.inputChan: // An insert is queued
			self.insert(dataPoint)
		}
	}
}

// Flush all cacheSlots with a mutationTime smaller than flushLimit
func (self *averageCache) flush(flushLimit int64) {
	fmt.Println("averageCache %v - INTERNAL FLUSH", self.name)
	var avgFlushTargets = []*DataPoint
	var countFlushTargets = []*DataPoint
	for timeslot, cacheSlot := range self.cache {
		if cacheSlot.lastUpdated <= flushLimit { // If the
			avgFlushTargets = append(avgFlushTargets, &DataPoint{cacheSlot.average, timeslot})
			countFlushTargets = append(countFlushTargets, &DataPoint{cacheSlot.count, timeslot})

			delete(self.cache, timeslot)
		}
	}
	fmt.Println("FlushAverages: ", avgFlushTargets)
	fmt.Println("FlushCounts: ", countFlushTargets)

	// TODO: Write flush targets to whisper
	//       In another fiber perhaps to make this non-blocking?
}

// Internal insert function. This is called in the runloop for the graphCache to insert dataPoints recieved on the inputChannel
func (self *averageCache) insert(dataPoint *DataPoint) {
	timeslot := dataPoint.Timestamp - (dataPoint.Timestamp % 10) // TODO: Custom interval
	if avgCache, exists := self.avgCaches[timeslot]; exists {
		countCache := self.countCaches[timeslot]
		avgCache.Value = (countCache.Value * avgCache.Value + dataPoint.Value) / countCache.value + 1
		countCache.Value += 1
		avgCache.Timestamp = time.Now().Unix() // We don't need the countCache timestamp, since this is always the same
	} else {
		gc.avgCache[timeslot] = &DataPoint{dataPoint.Value, time.Now().Unix()}
		gc.countCache[timeslot] = &DataPoint{1, 0} // We don't need to track the countCache timestamp
	}
}

// Flush the cache points in the graphCache to whisper
func (gc *graphCache) Flush(flushLimit int64) {
	// TODO: Remove a graphCache if it has no cacheSlots and nothing queued in inputChan
	//       Only flush otherwise
	fmt.Printf("GraphCache %v - FLUSH\n", gc.name)
	gc.flushChan <- flushLimit
}

// Insert a DataPoint into the graphCache
func (gc *graphCache) Insert(dp *DataPoint) {
	fmt.Printf("GraphCache %v - INSERT\n", gc.name)
	gc.inputChan <- dp
}

// The map of graphCaches in the cache package
var graphCaches map[string]*graphCache