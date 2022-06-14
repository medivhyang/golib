package orm

import (
	"fmt"
	"log"
	"os"
)

var (
	ErrNotFoundDialect    = errorf("not found dialect")
	ErrInvalidDialect     = errorf("invalid dialect")
	ErrRequirePointerType = errorf("require pointer type")
	ErrRequireSliceType   = errorf("require slice type")
	ErrRequireStructType  = errorf("require struct type")
	ErrCannotSetValue     = errorf("can not set value")
)

var TagKey = "orm"

var errorPrefix = "orm: "

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(errorPrefix+format, args...)
}

var (
	debug              = false
	debugLogger Logger = log.New(os.Stdout, "", log.LstdFlags)
	debugPrefix        = "orm: "
)

type Logger interface {
	Printf(format string, args ...interface{})
}

func EnableDebug(b bool) {
	debug = b
}

func SetLogger(l Logger) {
	debugLogger = l
}

func debugf(format string, args ...interface{}) {
	if debug && debugLogger != nil {
		debugLogger.Printf(debugPrefix+format, args...)
	}
}
