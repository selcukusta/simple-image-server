package googledrive

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/helper"
	"github.com/selcukusta/simple-image-server/internal/util/model"
	"github.com/valyala/fasthttp"
	"google.golang.org/api/drive/v2"
)

//Handler is using connect to Google Drive subscription and get the image
func Handler(ctx *fasthttp.RequestCtx, vars map[string]string) {
	if !helper.GoogleCredentialIsAvailable() {
		log.Println(`Google credential file cannot be found! Please create the file and set the "GOOGLE_APPLICATION_CREDENTIALS" environment variable.`)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, err := ctx.WriteString(constant.ErrorMessage)
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	path := vars["path"]
	defer helper.TraceObject{HandlerName: "Google Drive", Parameter: path}.TimeTrack(time.Now())

	_, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	service, err := drive.NewService(ctx)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to create Drive service", err.Error()))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, err = ctx.WriteString(constant.ErrorMessage)
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	res, err := service.Files.Get(path).Download()
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to download file", err.Error()))
		ctx.SetStatusCode(http.StatusInternalServerError)
		_, err = ctx.WriteString(constant.ErrorMessage)
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(res.Body)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Downloaded data has been corrupted", err.Error()))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, err = ctx.WriteString(constant.ErrorMessage)
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	headers := make(map[string]string)
	headers["ETag"] = string(res.Header.Get("Etag"))
	model.HandlerFinalizer{ResponseWriter: ctx, Headers: headers}.Finalize(vars, downloadedData.Bytes(), res.Header.Get("Content-Type"))
}
