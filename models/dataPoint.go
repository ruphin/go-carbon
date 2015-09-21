package models

type DataPoint struct {
	Value float64
	Time int64 // Unix time
}

type DataPoints []*DataPoint