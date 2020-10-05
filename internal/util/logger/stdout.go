package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	realip "github.com/Ferluci/fast-realip"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

//Init will be used to initialize the common fields of the log.
func Init() *log.Entry {
	pcs := make([]uintptr, 10)
	n := runtime.Callers(1, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)
	var stack []string
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		stack = append(stack, fmt.Sprintf("<%s> %s:%d", frame.File, frame.Function, frame.Line))
	}

	hostname, _ := os.Hostname()

	appName, exists := os.LookupEnv("APPLICATION_NAME")
	if !exists {
		appName = "UNKNOWN APPLICATION NAME"
	}

	logger, exists := os.LookupEnv("LOGGER_NAME")
	if !exists {
		logger = "UNKNOWN LOGGER NAME"
	}

	entry := log.WithFields(log.Fields{
		"@version":        1,
		"ApplicationName": appName,
		"Logger":          logger,
		"CallStack":       strings.Join(stack[:], "\n"),
		"Date":            time.Now(),
		"MachineName":     hostname,
	})

	return entry
}

//InitException will be used to initialize the common fields of the log and the exception detail.
func InitException(err error) *log.Entry {
	initial := Init()
	return initial.WithFields(log.Fields{
		"Exception": fmt.Sprintf("%+v", err),
	})
}

//InitWithRequest will be used to initialize the common fields of the log and the request detail.
func InitWithRequest(rq *fasthttp.RequestCtx) *log.Entry {
	initial := Init()
	uri := rq.Request.URI()
	return initial.WithFields(log.Fields{
		"Referer":     string(rq.Request.Header.Peek("Referer")),
		"Url":         fmt.Sprintf("%s://%s%s", uri.Scheme(), uri.Host(), uri.Path()),
		"QueryString": string(uri.QueryString()),
		"UserAgent":   string(rq.Request.Header.Peek("User-Agent")),
		"Client-IP":   realip.FromRequest(rq),
	})
}

//InitExceptionWithRequest will be used to initialize the common fields of the log, request and the exception detail.
func InitExceptionWithRequest(rq *fasthttp.RequestCtx, err error) *log.Entry {
	initial := InitWithRequest(rq)
	return initial.WithFields(log.Fields{
		"Exception": fmt.Sprintf("%+v", err),
	})
}
