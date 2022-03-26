package logger

import (
	"io"
	"log"
)

type BLogger struct {
	infoLogger *log.Logger
	errLogger  *log.Logger
}

func New(out io.Writer) *BLogger {
	return &BLogger{
		log.New(out, "[INFO]  ", log.LstdFlags),
		log.New(out, "[ERROR] ", log.LstdFlags),
	}
}

func (bl *BLogger) Info(v ...interface{}) {
	bl.infoLogger.Println(v...)
}

func (bl *BLogger) Error(v ...interface{}) {
	bl.errLogger.Println(v...)
}

func (bl *BLogger) Fatal(v ...interface{}) {
	bl.errLogger.Fatal(v...)
}

func (bl *BLogger) Infof(format string, v ...interface{}) {
	bl.infoLogger.Printf(format, v...)
}

func (bl *BLogger) Errorf(format string, v ...interface{}) {
	bl.errLogger.Printf(format, v...)
}
