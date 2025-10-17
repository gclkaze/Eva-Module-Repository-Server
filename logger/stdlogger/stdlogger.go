package stdlog

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/logger"
	log "github.com/sirupsen/logrus"
)

type STDLogger struct {
	Errors []string
}

func NewSTDLogger() *STDLogger {
	return &STDLogger{}
}
func (rd *STDLogger) GetType() logger.LoggerType {
	return logger.STD
}
func (logger *STDLogger) Printf(id string, format string, args ...interface{}) {
	log.Printf(format, args...)
}
func (logger *STDLogger) Errorf(id string, format string, args ...interface{}) {
	log.Printf(format, args...)
	logger.Errors = append(logger.Errors, fmt.Sprintf(format, args...))
}
func (logger *STDLogger) Init(id string) {

}
func (logger *STDLogger) PutSuccessMessage(id string, result bool, message string) {
	log.Printf("id:%s,result:%t,message:%s", id, result, message)
}

func (logger STDLogger) IsOnError() bool {
	return false
}
func (logger *STDLogger) Clear() {
}

func (logger STDLogger) PrintMessage(id string, message string) {
	log.Printf("Execution Id : %s, message : %s", id, message)
}

func (logger STDLogger) PrintErrorMessage(id string, message string) {
	log.Printf("Execution Id : %s, message : %s", id, message)
}
