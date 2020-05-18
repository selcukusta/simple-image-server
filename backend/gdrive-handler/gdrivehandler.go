package gdrivehandler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/selcukusta/simple-image-server/backend/constant"
	helper "github.com/selcukusta/simple-image-server/backend/helper"
	processor "github.com/selcukusta/simple-image-server/backend/image-processor"
	drive "google.golang.org/api/drive/v2"
)

var (
	errorMessage = `
		<h1>Oops! Something went wrong...</h1>
		<p>We seem to be having some technical difficulties. Hang tight.</p>
	`
)

//GoogleDriveHandler is using connect to Google Drive subscription and get the image
func GoogleDriveHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	defer helper.TimeTrack(time.Now(), path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service, err := drive.NewService(ctx)
	if err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Unable to create Drive service", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	res, err := service.Files.Get(path).Download()
	if err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Unable to download file", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(res.Body)
	if err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Downloaded data has been corrupted", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	result, errMessage, err := processor.ImageProcess(mux.Vars(r), downloadedData.Bytes(), res.Header.Get("Content-Type"))
	if err != nil {
		log.Println(fmt.Sprintf("%s: %s", errMessage, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	if constant.CacheControlMaxAge != -1 {
		maxAge := constant.CacheControlMaxAge * 24 * 60 * 60
		w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	}
	w.Header().Add("ETag", string(res.Header.Get("Etag")))
	w.Write(result)
}
