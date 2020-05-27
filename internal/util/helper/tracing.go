package helper

import (
	"fmt"
	"time"

	"github.com/selcukusta/simple-image-server/internal/util/logger"
)

//TraceObject is using to store information about logged object
type TraceObject struct {
	HandlerName string
	Parameter   string
}

//TimeTrack will be used to calculate elapsed time of execution.
func (o TraceObject) TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	logger.WriteLog(logger.INFO, fmt.Sprintf("(%s) %s took %s", o.HandlerName, o.Parameter, elapsed))
}
