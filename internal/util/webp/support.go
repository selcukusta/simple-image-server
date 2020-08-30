package webp

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

func webpSupportIsEnabled() bool {
	value, defined := os.LookupEnv("WEBP_ENABLED")
	if !defined || value == "" {
		return false
	}
	return true
}

//ConvertToWebp is using to convert the image to WebP
func ConvertToWebp(source []byte) ([]byte, error) {
	if !webpSupportIsEnabled() {
		return nil, errors.New("WebP feature is not enabled. Please check the 'WEBP_ENABLED' env variable")
	}

	if source == nil {
		return nil, errors.New("Source could not be nil")
	}

	tempSource, err := ioutil.TempFile("", "source-*.png")
	if err != nil {
		return nil, err
	}

	defer os.Remove(tempSource.Name())

	tempDestination, err := ioutil.TempFile("", "destination-*.webp")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempDestination.Name())

	_, err = tempSource.Write(source)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("cwebp", tempSource.Name(), "-o", tempDestination.Name())

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(tempDestination.Name())
	if err != nil {
		return nil, err
	}
	return data, nil
}
