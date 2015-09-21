package cache

import (
	log "github.com/Sirupsen/logrus"
	models "github.com/ruphin/go-carbon/models"
)


type Kind int

const (
	AVERAGE Kind = iota
	MIN
	MAX
	NONE
)

type graphCache interface {
	Flush(flushLimit int64)
	Insert(dataPoints model.DataPoints)
	Get(subscription Subscription)
	Suscribe(subscription Subscription)
}

// Get the cache for the graph identified by a name
func Open(name string, kind Kind) (graphCache) {
	log.WithFields(log.Fields{
		"cache": name,
	}).Info("Opening Cache")
	var gc graphCache
	var exists bool
	if gc, exists = caches[kind][name]; !exists { // The graphCache does not exist yet, create it
		log.Info("New Cache")
		if kind == AVERAGE {
			log.Info("Average Mode")
			gc = NewAverageCache(name)
		}
	} else {
		log.Info("Existing Cache")
	}
	return gc
}

var caches = map[Kind]map[string]graphCache{
	AVERAGE: make(map[string]graphCache),
	MIN: make(map[string]graphCache),
	MAX: make(map[string]graphCache)}
