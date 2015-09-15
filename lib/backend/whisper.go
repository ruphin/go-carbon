package backend

import (
	whisper "github.com/lomik/go-whisper"
)


type WhisperBackend struct {}

func (self *WhisperBackend) Write(dataPoints []DataPoint) {
	wsp, err := whisper.Open("/tmp/bla.wsp")
	for _, dataPoint := range dataPoints {

	}
	wsp.UpdateMany(&whisper.TimeSeriesPoint{})
}


func (self *WhisperBackend) Read(timeRange TimeRange) ([]DataPoint) {
	return nil
}