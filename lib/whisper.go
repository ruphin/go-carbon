package carbon

import (
	whisper "github.com/lomik/go-whisper"
)


type WhisperBackend struct {}

func (self *WhisperBackend) Write(dataPoints DataPoints) {
	wsp, err := whisper.Open("/tmp/bla.wsp")
	wsp.UpdateMany(dataPoints)
}


func (self *WhisperBackend) Read(timeRange TimeRange) ([]DataPoint) {
	return nil
}