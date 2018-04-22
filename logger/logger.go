package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/shiena/ansicolor"
)

// Text colors
const (
	Black   string = "\x1b[30m"
	Red     string = "\x1b[31m"
	Green   string = "\x1b[32m"
	Yellow  string = "\x1b[33m"
	Blue    string = "\x1b[34m"
	Magenta string = "\x1b[35m"
	Cyan    string = "\x1b[36m"
	White   string = "\x1b[37m"
	Reset   string = "\x1b[0m"
)

type Logger struct {
	out              io.Writer // destination for output
	infoPrefix       string
	logPrefix        string
	warnPrefix       string
	errorPrefix      string
	resetColorSuffix string
}

// Returns a new console logger
func New(out io.Writer) *Logger {
	return &Logger{
		out:              out,
		infoPrefix:       fmt.Sprintf("%s[*] ", Blue),
		logPrefix:        fmt.Sprint("[+] "),
		warnPrefix:       fmt.Sprintf("%s[!] ", Yellow),
		errorPrefix:      fmt.Sprintf("%s[-] ", Red),
		resetColorSuffix: fmt.Sprintf("%s\n", Reset),
	}
}

// Print output to logger output writer
func (l *Logger) Output(s string) error {
	writer := ansicolor.NewAnsiColorWriter(l.out)
	_, err := fmt.Fprintf(writer, s+l.resetColorSuffix)
	return err
}

// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
}

// Logs log
func (l *Logger) Log(v ...interface{}) {
	l.Output(l.logPrefix + fmt.Sprint(v...))
}

func (l *Logger) Logf(format string, v ...interface{}) {
	l.Output(l.logPrefix + fmt.Sprintf(format, v...))
}

// Logs info
func (l *Logger) Info(v ...interface{}) {
	l.Output(l.infoPrefix + fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(l.infoPrefix + fmt.Sprintf(format, v...))
}

// Logs warning
func (l *Logger) Warn(v ...interface{}) {
	l.Output(l.warnPrefix + fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(l.warnPrefix + fmt.Sprintf(format, v...))
}

// Logs error
func (l *Logger) Error(v ...interface{}) {
	l.Output(l.errorPrefix + fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(l.errorPrefix + fmt.Sprintf(format, v...))
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(l.errorPrefix + fmt.Sprint(v...))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(s)
	panic(s)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(s)
	panic(s)
}

var Log *Logger = New(os.Stdout)
