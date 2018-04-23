package utils

// Check if the error is not nil printing the error and
// panic after that.
func FatalCheck(err error, errorMsg string, v ...interface{}) {
	if err != nil {
		if errorMsg == "" {
			Log.Errorf(errorMsg, v...)
		}
		Log.Panic(err)
	}
}
