package helper

import (
	"log"
	"os"
	"time"
)

//TimeTrack will be used to calculate elapsed time of execution.
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

//GoogleCredentialIsAvailable will be check credential file exists or not.
func GoogleCredentialIsAvailable() bool {
	value, defined := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !defined || value == "" {
		return false
	}

	if _, err := os.Stat(value); os.IsNotExist(err) {
		return false
	}

	return true
}
