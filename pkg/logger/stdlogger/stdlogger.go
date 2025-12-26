// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package stdlog

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
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
