package backend


import (
	log "github.com/Sirupsen/logrus"
)
import (
	whisper "github.com/lomik/go-whisper"
)


type WhisperBackend struct {}

func (self *WhisperBackend) Write(dataPoints []*whisper.TimeSeriesPoint) {
	if len(dataPoints) == 0 {
		return
	}
	path := "/tmp/test.wsp"
	log.WithFields(log.Fields{
		"dataPoints": dataPoints,
	}).Info("Writing to Backend")
	wsp, err := whisper.Open(path)
	if err != nil {
		retentions, _ := whisper.ParseRetentionDefs("10s:14d,1m:90d")
		log.WithFields(log.Fields{
			"path": path,
			"retentions": "10s:14d,1m:90d",
			"aggregation": whisper.Sum,
		}).Info("New Whisper Database")
		wsp, err = whisper.Create(path, retentions, whisper.Sum, 0.5)
		if err != nil {
			log.WithFields(log.Fields{
				"path": path,
			}).Error("Failed to create file")
			return
		}
	}
	log.WithFields(log.Fields{
		"wsp": wsp,
	}).Info("Opened Whisper")

	defer wsp.Close()

	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"path": path,
				"values": dataPoints,
			}).Error("Whisper Database write failed")
		}
	}()

	wsp.UpdateMany(dataPoints)
	log.Info("Updated Whisper")
	return
}


func (self *WhisperBackend) Read(timeRange TimeRange) ([]*whisper.TimeSeriesPoint) {
	log.Info("Reading from Backend")
	return nil
}