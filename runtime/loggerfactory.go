package runtime

import (
	"github.com/gclkaze/evamodulerepositoryserver/logger"
	stdlog "github.com/gclkaze/evamodulerepositoryserver/logger/stdlogger"
	"github.com/magiconair/properties"
)

func CreateLogger(props *properties.Properties) logger.ILogger {
	return stdlog.NewSTDLogger()
}
