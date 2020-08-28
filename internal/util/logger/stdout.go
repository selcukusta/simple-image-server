package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	realip "github.com/Ferluci/fast-realip"
	"github.com/valyala/fasthttp"
)

type entity struct {
	Version         int       `json:"@version"`
	LogLineNumber   int       `json:"LogLineNumber"`
	CallStack       string    `json:"CallStack"`
	MethodName      string    `json:"MethodName"`
	ApplicationName string    `json:"ApplicationName"`
	Date            time.Time `json:"Date"`
	MachineName     string    `json:"MachineName"`
	Level           string    `json:"Level"`
	Message         string    `json:"Message"`
	Logger          string    `json:"Logger"`
	URL             string    `json:"Url"`
	ClientIP        string    `json:"Client-IP"`
	Referer         string    `json:"Referer"`
	QueryString     string    `json:"QueryString"`
}

//Log is using to create log object
type Log struct {
	Level   Level
	Message string
	Rq      *fasthttp.RequestCtx
}

func createLog(logObject Log) entity {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	hostname, _ := os.Hostname()
	appName, exists := os.LookupEnv("APPLICATION_NAME")
	if !exists {
		appName = "UNKNOWN APPLICATION NAME"
	}
	logger, exists := os.LookupEnv("LOGGER_NAME")
	if !exists {
		logger = "UNKNOWN LOGGER NAME"
	}

	output := entity{Version: 1, LogLineNumber: line, CallStack: file, MethodName: fn.Name(), ApplicationName: appName, Date: time.Now(), MachineName: hostname, Level: logObject.Level.String(), Message: logObject.Message, Logger: logger}
	if logObject.Rq != nil {
		output.Referer = string(logObject.Rq.Request.Header.Peek("Referer"))
		uri := logObject.Rq.Request.URI()
		output.URL = fmt.Sprintf("%s://%s%s", uri.Scheme(), uri.Host(), uri.Path())
		output.QueryString = string(uri.QueryString())
		output.ClientIP = realip.FromRequest(logObject.Rq)
	}
	return output
}

//WriteLog is using to write log to stdout
func WriteLog(logObject Log) {
	output := createLog(logObject)
	var jsonData []byte
	jsonData, err := json.Marshal(output)
	if err != nil {
		return
	}
	fmt.Println(string(jsonData))
}
