package carbon

import "time"
import "fmt"
import "encoding/json"
// import (
// 	whisper "github.com/lomik/go-whisper"
// )

// A unix timestamp
type timestamp int64

// This is used for input values, timestamp can be anything in the form of a unix timestamp
type DataPoint struct {
	value int64
	timestamp timestamp
}

// This represents a value that can be written to whisper. Timeslot can only be a multiple of X where X is the storage interval used in whisper
type timeSlot struct {
	value int64
	timeslot timestamp
}

// This represents the cached value for a single timeslot, with the last time it was updated
type cacheSlot struct {
	value int64
	mutationTime timestamp
}

// TODO: Make insert and flush work for averaged values (backed by two whisperfiles for, 'count' and 'value')
//       At this time it just overwrites the current value, which is kinda useless
func (cs *cacheSlot) insert(value int64) {
	cs.value = value
}

// This represents the cache for a single graph (identified uniquely by a string)
type graphCache struct {
	name string
	caches map[timestamp]*cacheSlot
	inputChan chan *DataPoint
	flushChan chan timestamp
}

// The event loop for this graphCache
// It guarantees sequential inserts and flushes
func (gc *graphCache) run() {
	fmt.Println("GraphCache %v - INTERNAL RUN", gc.name)
	var dp *DataPoint
	var flushLimit timestamp
	for {
		select {
		case flushLimit = <- gc.flushChan: // A flush is queued
			gc.flush(flushLimit)
		case dp = <- gc.inputChan: // An insert is queued
			gc.insert(dp)
		}
	}
}

// Flush all cacheSlots with a mutationTime smaller than flushLimit
func (gc *graphCache) flush(flushLimit timestamp) {
	fmt.Println("GraphCache %v - INTERNAL FLUSH", gc.name)
	var flushTargets = []*timeSlot{}
	for timeslot, cs := range gc.caches {
		if cs.mutationTime <= flushLimit {
			flushTargets = append(flushTargets, &timeSlot{cs.value, timeslot})
			delete(gc.caches, timeslot)
		}
	}

	// TODO: Write flushTargets to whisper
	//       In another fiber perhaps to make this non-blocking?
}

// Internal insert function. This is called in the runloop for the graphCache to insert dataPoints recieved on the inputChannel
func (gc *graphCache) insert(dp *DataPoint) {
	fmt.Println("GraphCache %v - INTERNAL INSERT", gc.name)
	fmt.Println("VALUE %v", string(json.Marshal(dp)))
	timeslot := dp.timestamp - (dp.timestamp % 10)
	if cs, exists := gc.caches[timeslot]; exists {
		cs.insert(dp.value)
		cs.mutationTime = timestamp(time.Now().Unix())
	} else {
		gc.caches[timeslot] = &cacheSlot{cs.value, timestamp(time.Now().Unix())}
	}
}

// Flush the cache points in the graphCache to whisper
func (gc *graphCache) Flush(flushLimit timestamp) {
	// TODO: Remove a graphCache if it has no cacheSlots and nothing queued in inputChan
	//       Only flush otherwise

	fmt.Println("GraphCache %v - FLUSH", gc.name)
	gc.flushChan <- flushLimit
}

// Insert a DataPoint into the graphCache
func (gc *graphCache) Insert(dp *DataPoint) {
	fmt.Println("GraphCache %v - INSERT", gc.name)
	gc.inputChan <- dp
}

// The map of graphCaches in the cache package
var graphCaches map[string]*graphCache

// Get the cache for the graph identified by a name
func GetGraphCache(name string) *graphCache {
	var gc *graphCache
	if gc, exists := graphCaches[name]; !exists { // The graphCache does not exist yet, create it
		// Maybe the inputChannel needs a longer queue size to avoid blocking
		gc = &graphCache{name, make(map[timestamp]*cacheSlot), make(chan *DataPoint), make(chan timestamp)} 
		go gc.run() // Start the eventloop for this graphCache
	}
	return gc
}