package logger

//Level is using to categorize level
type Level int

const (
	//TRACE is using to categorize the log as Trace
	TRACE Level = iota
	//DEBUG is using to categorize the log as Debug
	DEBUG
	//INFO is using to categorize the log as Info
	INFO
	//WARN is using to categorize the log as Warn
	WARN
	//ERROR is using to categorize the log as Error
	ERROR
	//FATAL is using to categorize the log as Fatal
	FATAL
)

func (l Level) String() string {
	return [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}
