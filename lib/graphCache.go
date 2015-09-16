package carbon

import "fmt"
import (
	whisper "github.com/lomik/go-whisper"
)

type Subscription struct {
	closeChan chan bool
	dataChan chan []DataPoint
	Range TimeRange
}

type DataPoint whisper.TimeSeriesPoint

type DataPoints []*DataPoint

type TimeRange struct {
	From int
	Until int
}

type Kind int

const (
	AVERAGE Kind = iota
	MIN
	MAX
	NONE
)

type graphCache interface {
	Flush(flushLimit int)
	Insert(dataPoints DataPoints)
	Get(subscription Subscription)
	Suscribe(subscription Subscription)
}

// Get the cache for the graph identified by a name
func Open(name string, kind Kind) (*graphCache) {
	fmt.Printf("Getting GraphCache: %v\n", name)
	var gc graphCache
	var exists bool
	if gc, exists = caches[kind][name]; !exists { // The graphCache does not exist yet, create it
		fmt.Println("New Cache")
		if kind == AVERAGE {
			cache := NewAverageCache(name)
			go cache.run() // Start the eventloop for this graphCache
			gc = cache
		}
	} else {
		fmt.Println("Existing Cache")
	}
	return &gc
}

var caches [Kind][string]*graphCache = map[Kind][string]*graphCache{
	AVERAGE: make(map[string]*graphCache),
	MIN: make(map[string]*graphCache),
	MAX: make(map[string]*graphCache)}
