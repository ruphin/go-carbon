package app

import(
	carbon "github.com/ruphin/go-carbon/lib"
)
import "time"

func main() {
	var cache = carbon.GetGraphCache("test")
	cache.Insert(&carbon.DataPoint{1, time.Now().Unix()})
	time.Sleep(1)
	cache.Insert(&carbon.DataPoint{2, time.Now().Unix()})
	time.Sleep(1)
	cache.Insert(&carbon.DataPoint{3, time.Now().Unix()})
	time.Sleep(1)
	cache.Insert(&carbon.DataPoint{4, time.Now().Unix()})
	time.Sleep(1)
	cache.Insert(&carbon.DataPoint{5, time.Now().Unix()})
	time.Sleep(1)
	cache.Insert(&carbon.DataPoint{6, time.Now().Unix()})
	time.Sleep(1)
	cache.Insert(&carbon.DataPoint{7, time.Now().Unix()})
	time.Sleep(1)
	cache.Flush((time.Now().Unix() - 5))
}
