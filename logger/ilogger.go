package logger

type LoggerType int

const (
	STD LoggerType = iota
)

type ILogger interface {
	PutSuccessMessage(id string, result bool, message string)
	Printf(id string, format string, args ...interface{})
	Errorf(id string, format string, args ...interface{})
	GetType() LoggerType
	Init(id string)
	IsOnError() bool
	Clear()
}
