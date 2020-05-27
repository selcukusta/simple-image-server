package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

//Entity is using to create log message entity
type Entity struct {
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
}

//WriteLog is using to write log to stdout
func WriteLog(level Level, message string) {
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

	output := Entity{Version: 1, LogLineNumber: line, CallStack: file, MethodName: fn.Name(), ApplicationName: appName, Date: time.Now(), MachineName: hostname, Level: level.String(), Message: message, Logger: logger}
	var jsonData []byte
	jsonData, err := json.Marshal(output)
	if err != nil {
		return
	}
	fmt.Println(string(jsonData))
}
