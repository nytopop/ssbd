// Package logs provides logging functionality for ssbd.
package logs

import (
	"log"
	"os"

	"github.com/nytopop/ssbd/config"
)

// Logger contains two pointers to log.Logger, E represents the error log,
// and A represents the access log.
type Logger struct {
	E *log.Logger
	A *log.Logger
}

var logger Logger

// InitLoggers initializes logger with parameters from configuration.
func InitLoggers() error {
	errlog, err := os.Create(config.CFG.Srv.ErrorLog)
	if err != nil {
		return err
	}

	acclog, err := os.Create(config.CFG.Srv.AccessLog)
	if err != nil {
		return err
	}

	logger = Logger{
		E: log.New(errlog, "", log.Ldate|log.Ltime),
		A: log.New(acclog, "", log.Ldate|log.Ltime),
	}

	return nil
}

// Error logs one or more entries to the error log.
func Error(err ...interface{}) {
	logger.E.Println(err...)
}

// Access logs one or more entries to the access log.
func Access(err ...interface{}) {
	logger.A.Println(err...)
}

// App wide Err type.
type Err string

func (e Err) Error() string { return string(e) }
