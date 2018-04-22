package utils

import (
	"github.com/fristonio/gvm/logger"
)

// Check if the error is not nil printing the error and
// panic after that.
func FatalCheck(err error, errorMsg string, v ...interface{}) {
	if err != nil {
		if errorMsg == nil {
			log.Errorf(errorMsg, v...)
		}
		log.Panic(err)
	}
}
