package app

import(
	carbon "github.com/ruphin/go-carbon/lib"
)
import "time"

cache = carbon.GetGraphCache('test')
cache.Insert(Datapoint{1, time.Now().Unix()}))
time.Sleep(1)
cache.Insert(Datapoint{2, time.Now().Unix()}))
time.Sleep(1)
cache.Insert(Datapoint{3, time.Now().Unix()}))
time.Sleep(1)
cache.Insert(Datapoint{4, time.Now().Unix()}))
time.Sleep(1)
cache.Insert(Datapoint{5, time.Now().Unix()}))
time.Sleep(1)
cache.Insert(Datapoint{6, time.Now().Unix()}))
time.Sleep(1)
cache.Insert(Datapoint{7, time.Now().Unix()}))
time.Sleep(1)
cache.Flush((time.Now().Unix() - 5))