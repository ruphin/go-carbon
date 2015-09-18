package main

import (
	carbon "github.com/ruphin/go-carbon/lib"
	whisper "github.com/lomik/go-whisper"
)

import "time"

func main() {
	var cache = carbon.Open("test", carbon.AVERAGE)
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 1}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 2}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 3}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 4}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 5}})
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 1}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 2}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 3}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 4}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 5}})
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 1}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 2}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 3}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 4}})
	time.Sleep(1 * time.Second)
	cache.Insert([]*whisper.TimeSeriesPoint{&whisper.TimeSeriesPoint{int(time.Now().Unix()), 5}})
	time.Sleep(10 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
	cache.Flush(int(time.Now().Unix() - 30))
	time.Sleep(1 * time.Second)
}
