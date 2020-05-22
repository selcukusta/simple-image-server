package helper

import (
	"log"
	"time"
)

//TraceObject is using to store information about logged object
type TraceObject struct {
	HandlerName string
	Parameter   string
	Took        time.Duration
}

//TimeTrack will be used to calculate elapsed time of execution.
func (o TraceObject) TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	o.Took = elapsed
	log.Printf("%+v", o)
}
