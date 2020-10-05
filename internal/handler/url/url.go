package s3

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/model"
	"github.com/valyala/fasthttp"
)

//Handler is using connect to public host url and get the image
func Handler(ctx *fasthttp.RequestCtx, vars map[string]string) {
	path := vars["path"]

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	origin := fmt.Sprintf("%s", connection.URL)
	req, err := http.NewRequest("GET", origin, nil)
	if err != nil {
		customError := model.CustomError{Message: "An error has occurred while the request was prepared", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	req = req.WithContext(context)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		customError := model.CustomError{Message: "An error has occurred while the request was sent", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	defer res.Body.Close()
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(res.Body)
	if err != nil {
		customError := model.CustomError{Message: "An error has occurred while the body was read", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	headers := make(map[string]string)
	headers["ETag"] = string(res.Header.Get("Etag"))
	headers["LastModified"] = string(res.Header.Get("LastModified"))
	model.SucceededFinalizer{ResponseWriter: ctx, ContentType: res.Header.Get("Content-Type"), Headers: headers}.Finalize(vars, downloadedData.Bytes())
}
