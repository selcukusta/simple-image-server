package helper

import (
	"log"
	"time"
)

//TimeTrack will be used to calculate elapsed time of execution.
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
