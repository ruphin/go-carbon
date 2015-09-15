package cache

type averageCache struct {
	name String
	subscriptions []Subscription
	inputChan chan []DataPoint
	flushChan chan int64
	closeChan chan bool
	cache map[int]*cacheSlot
	avgBackend WhisperBackend
	countBackend WhisperBackend
}

type cacheSlot struct {
	Average int64
	Count int64
	LastUpdated int64
}

// The event loop for this averageCache
// It guarantees sequential inserts and flushes
func (self *averageCache) run() {
	var dataPoints []DataPoint
	var flushLimit int64
	for {
		select {
		case stop = <- closeChan: // The cache is ordered to close
			self.close()
		case flushLimit = <- self.flushChan: // A flush is queued
			self.flush(flushLimit)
		case dataPoints = <- self.inputChan: // An insert is queued
			self.insert(dataPoints)
		}
	}
}

// Flush all cacheSlots with a mutationTime smaller than flushLimit
func (self *averageCache) flush(flushLimit int) {
	fmt.Println("averageCache %v - INTERNAL FLUSH", self.name)
	var avgFlushTargets DataPoints = make(DataPoints)
	var countFlushTargets DataPoints = make(DataPoints)
	for timeSlot, cacheSlot := range self.cache {
		// TODO: 
		// Write all changes to subscriptions
		if cacheSlot.LastUpdated <= flushLimit {
			avgFlushTargets = append(avgFlushTargets, DataPoint{timeSlot, cacheSlot.Average})
			countFlushTargets = append(countFlushTargets, DataPoint{timeSlot, cacheSlot.Count})
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
func (self *averageCache) insert(dataPoints []DataPoint) {
	var timeSlot int64
	var cacheSlot cacheSlot
	var exists bool
	for _, dataPoint := range dataPoints {
		timeSlot = dataPoint.Timestamp - (dataPoint.Timestamp % 10) // TODO: Custom interval
		if cacheSlot, exists = self.cache[timeSlot]; exists {
			cacheSlot.Value = ((cacheSlot.Count * cacheSlot.Value) + dataPoint.Value) / (cacheSlot.Count + 1)
			cacheSlot.Count += 1
			cacheSlot.LastUpdated = time.Now().Unix()
		} else {
			self.cache[timeSlot] = &cacheSlot{dataPoint.Value, 1, time.Now().Unix()}
		}
	}
}

// Flush the cache points in the graphCache to whisper
func (self *averageCache) Flush(flushLimit int64) {
	// TODO: Remove a graphCache if it has no cacheSlots and nothing queued in inputChan
	//       Only flush otherwise
	fmt.Printf("GraphCache %v - FLUSH\n", self.name)
	self.flushChan <- flushLimit
}

// Insert a DataPoint into the graphCache
func (self *averageCache) Insert(dataPoints []DataPoint) {
	fmt.Printf("GraphCache %v - INSERT\n", self.name)
	self.inputChan <- dataPoints
}

func (self *averageCache) Get(subscription Subscription) {
	subscription.dataChan <- self.avgBackend.Read(subscription.Range)
}

func (self *averageCache) Suscribe(subscription Subscription) {
	self.subscriptions = append(self.subscriptions. subscription)
	self.Get(subscription)
}
