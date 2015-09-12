package carbon

import "time"
import "fmt"
// import (
// 	whisper "github.com/lomik/go-whisper"
// )


// Determines what value is stored in each individual timeslot
type AggregationMethod int
const (
	Average AggregationMethod = iota // Uses two backing whisper databases
	Min
	Max
	None // Store the last value written
)

// This is used for input values, timestamp is a unix timestamp
type DataPoint struct {
	Value int64
	Timestamp int64
}

// This represents a value that can be written to whisper. Timeslot can only be a multiple of X where X is the storage interval used in whisper
type timeSlot struct {
	value int64
	timeslot int64
}

// This represents the cached value for a single timeslot, with the last time it was updated
type cacheSlot struct {
	value int64
	mutationTime int64
}

// TODO: Make insert and flush work for averaged values (backed by two whisperfiles for, 'count' and 'value')
//       At this time it just overwrites the current value, which is kinda useless
func (cs *cacheSlot) insert(value int64) {
	cs.value = value
}

// This represents the cache for a single graph (identified uniquely by a string)
type graphCache struct {
	name string
	aggregation int
	caches map[int64]*cacheSlot
	inputChan chan *DataPoint
	flushChan chan int64
}

// The event loop for this graphCache
// It guarantees sequential inserts and flushes
func (gc *graphCache) run() {
	var dp *DataPoint
	var flushLimit int64
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
func (gc *graphCache) flush(flushLimit int64) {
	fmt.Println("GraphCache %v - INTERNAL FLUSH", gc.name)
	fmt.Println("Caches: ", gc.caches)
	var flushTargets = []*timeSlot{}
	for timeslot, cs := range gc.caches {
		if cs.mutationTime <= flushLimit {
			flushTargets = append(flushTargets, &timeSlot{cs.value, timeslot})
			delete(gc.caches, timeslot)
		}
	}
	fmt.Println("FlushTargets: ", flushTargets)
	fmt.Println("Caches: ", gc.caches)

	// TODO: Write flushTargets to whisper
	//       In another fiber perhaps to make this non-blocking?
}

// Internal insert function. This is called in the runloop for the graphCache to insert dataPoints recieved on the inputChannel
func (gc *graphCache) insert(dp *DataPoint) {
	timeslot := dp.Timestamp - (dp.Timestamp % 10)
	if cs, exists := gc.caches[timeslot]; exists {
		cs.insert(dp.Value)
		cs.mutationTime = time.Now().Unix()
	} else {
		gc.caches[timeslot] = &cacheSlot{dp.Value, time.Now().Unix()}
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

// Get the cache for the graph identified by a name
func GetGraphCache(name string, type int) *graphCache {
	fmt.Printf("Getting GraphCache: %v\n", name)
	var gc *graphCache
	var exists bool
	if gc, exists = graphCaches[name]; !exists { // The graphCache does not exist yet, create it
		fmt.Println("New Cache")
		// Maybe the inputChannel needs a longer queue size to avoid blocking
		gc = &graphCache{name, make(map[int64]*cacheSlot), make(chan *DataPoint), make(chan int64)}
		go gc.run() // Start the eventloop for this graphCache
	} else {
		fmt.Println("Existing Cache")
	}
	return gc
}