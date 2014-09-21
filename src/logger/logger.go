package logger

import (
	"log"
	"os"
)

var (
	DebugLog *log.Logger
	ErrorLog *log.Logger
)

func SetupLogs() {
	ErrorLog = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLog = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
}
