package googledrive

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/selcukusta/simple-image-server/internal/processor"
	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/helper"
	drive "google.golang.org/api/drive/v2"
)

//GoogleDriveHandler is using connect to Google Drive subscription and get the image
func GoogleDriveHandler(w http.ResponseWriter, r *http.Request) {
	if !helper.GoogleCredentialIsAvailable() {
		log.Println(`Google credential file cannot be found! Please create the file and set the "GOOGLE_APPLICATION_CREDENTIALS" environment variable.`)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	path := mux.Vars(r)["path"]
	defer helper.TimeTrack(time.Now(), path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service, err := drive.NewService(ctx)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to create Drive service", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	res, err := service.Files.Get(path).Download()
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to download file", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(res.Body)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Downloaded data has been corrupted", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	result, errMessage, err := processor.ImageProcess(mux.Vars(r), downloadedData.Bytes(), res.Header.Get("Content-Type"))
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, errMessage, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	if constant.CacheControlMaxAge != -1 {
		maxAge := constant.CacheControlMaxAge * 24 * 60 * 60
		w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	}
	w.Header().Add("ETag", string(res.Header.Get("Etag")))
	_, err = w.Write(result)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
	}
}
