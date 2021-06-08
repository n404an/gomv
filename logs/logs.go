package log

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	err   *log.Logger
	fatal *log.Logger
	f     *os.File
}

func NewLogger() *Logger {
	l := &Logger{}

	f, err := os.OpenFile("gomv.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	l.info = log.New(f, "INFO\t", log.Ldate|log.Ltime)
	l.err = log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	l.fatal = log.New(f, "FATAL\t", log.Ldate|log.Ltime|log.Lshortfile)
	l.f = f
	return l
}

func (l *Logger) Close() {
	l.f.Close()
}

func (l *Logger) Info(items ...interface{}) {
	l.info.Println(items...)

}

func (l *Logger) Err(items ...interface{}) {
	l.err.Println(items...)
}

func (l *Logger) Fatal(items ...interface{}) {
	l.fatal.Fatalln(items...)
}
