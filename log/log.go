package log

import (
	"io"
	"log"
	"os"
)

type dummy struct{}

func (d *dummy) Write(_ []byte) (n int, err error) {
	return 0, nil
}

var logger = log.New(&dummy{}, "", 0)

func SetVerboseLogger(out io.Writer) {
	logger = log.New(out, "[DEBUG] ", 0)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func Println(v ...interface{}) {
	logger.Println(v...)
}

func Fatal(v ...interface{}) {
	logger.Println(v...)
	os.Exit(1)
}
