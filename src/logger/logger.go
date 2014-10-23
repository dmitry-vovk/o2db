package logger

import (
	"log"
	"os"
)

var (
	DebugLog *log.Logger
	ErrorLog *log.Logger
)

func init() {
	SetupLogs()
}

func SetupLogs() {
	ErrorLog = log.New(os.Stderr, "Error in ", log.Lshortfile)
	DebugLog = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
}
