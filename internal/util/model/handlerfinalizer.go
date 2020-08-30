package model

import (
	"fmt"

	"github.com/selcukusta/simple-image-server/internal/processor"
	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/logger"
	"github.com/selcukusta/simple-image-server/internal/util/webp"
	"github.com/valyala/fasthttp"
)

//SucceededFinalizer is using to create a model for succeeded finalizing requests
type SucceededFinalizer struct {
	ResponseWriter *fasthttp.RequestCtx
	ContentType    string
	Headers        map[string]string
}

//FailedFinalizer is using to create a model for failed finalizing requests
type FailedFinalizer struct {
	ResponseWriter *fasthttp.RequestCtx
	StdOut         *CustomError
}

//CustomError is using to create a model for custom exception
type CustomError struct {
	Message string
	Detail  error
}

//Finalize is using to finalize the request unsuccessfully
func (hf FailedFinalizer) Finalize() {
	if hf.StdOut != nil {
		if hf.StdOut.Detail != nil {
			logger.WriteLog(logger.Log{Level: logger.ERROR, Message: fmt.Sprintf(constant.LogErrorFormat, hf.StdOut.Message, hf.StdOut.Detail.Error()), Rq: hf.ResponseWriter})
		} else {
			logger.WriteLog(logger.Log{Level: logger.ERROR, Message: hf.StdOut.Message, Rq: hf.ResponseWriter})
		}
	}

	hf.ResponseWriter.Response.Header.Set("Content-Type", "text/html")
	hf.ResponseWriter.SetStatusCode(fasthttp.StatusInternalServerError)
	_, err := hf.ResponseWriter.WriteString(constant.ErrorMessage)
	if err != nil {
		logger.WriteLog(logger.Log{Level: logger.ERROR, Message: fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()), Rq: hf.ResponseWriter})
	}
}

//Finalize is using to finalize the request successfully
func (hf SucceededFinalizer) Finalize(params map[string]string, imageAsByte []byte) {
	result, errMessage, err := processor.ImageProcess(params, imageAsByte, hf.ContentType)
	if err != nil {
		customError := CustomError{Message: errMessage, Detail: err}
		FailedFinalizer{ResponseWriter: hf.ResponseWriter, StdOut: &customError}.Finalize()
		return
	}

	if result == nil {
		customError := CustomError{Message: errMessage}
		FailedFinalizer{ResponseWriter: hf.ResponseWriter, StdOut: &customError}.Finalize()
		return
	}

	if constant.CacheControlMaxAge != -1 {
		maxAge := constant.CacheControlMaxAge * 24 * 60 * 60
		hf.ResponseWriter.Response.Header.Add("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	}

	if hf.Headers != nil && len(hf.Headers) > 0 {
		for key, value := range hf.Headers {
			hf.ResponseWriter.Response.Header.Add(key, value)
		}
	}

	contentType := hf.ContentType

	if params["webp"] != "" {
		converted, err := webp.ConvertToWebp(result)
		if err != nil {
			logger.WriteLog(logger.Log{Level: logger.WARN, Message: err.Error(), Rq: hf.ResponseWriter})
		} else {
			result = converted
			contentType = "image/webp"
		}
	}

	hf.ResponseWriter.Response.Header.Set("Content-Type", contentType)
	_, err = hf.ResponseWriter.Write(result)
	if err != nil {
		logger.WriteLog(logger.Log{Level: logger.ERROR, Message: fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()), Rq: hf.ResponseWriter})
	}
}
