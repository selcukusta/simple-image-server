package helper

import (
	"os"
)

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
