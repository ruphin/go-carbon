package carbon

type GraphCache interface {
	Flush(flushLimit int64)
	Insert(dp *DataPoint)
}