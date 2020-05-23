package googledrive

import (
	"bytes"
	"context"
	"time"

	"github.com/selcukusta/simple-image-server/internal/util/helper"
	"github.com/selcukusta/simple-image-server/internal/util/model"
	"github.com/valyala/fasthttp"
	"google.golang.org/api/drive/v2"
)

//Handler is using connect to Google Drive subscription and get the image
func Handler(ctx *fasthttp.RequestCtx, vars map[string]string) {
	if !helper.GoogleCredentialIsAvailable() {
		customError := model.CustomError{Message: `Google credential file cannot be found! Please create the file and set the "GOOGLE_APPLICATION_CREDENTIALS" environment variable.`}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	path := vars["path"]
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	service, err := drive.NewService(ctx)
	if err != nil {
		customError := model.CustomError{Message: "Unable to create Drive service", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	res, err := service.Files.Get(path).Download()
	if err != nil {
		customError := model.CustomError{Message: "Unable to download file", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(res.Body)
	if err != nil {
		customError := model.CustomError{Message: "Downloaded data has been corrupted", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	headers := make(map[string]string)
	headers["ETag"] = string(res.Header.Get("Etag"))
	model.SucceededFinalizer{ResponseWriter: ctx, ContentType: res.Header.Get("Content-Type"), Headers: headers}.Finalize(vars, downloadedData.Bytes())
}
