package logger

import (
	"log"
	"os"
)

var (
	DebugLog  *log.Logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
	ErrorLog  *log.Logger = log.New(os.Stderr, "Error in ", log.Lshortfile)
	SystemLog *log.Logger = log.New(os.Stdout, "", log.LstdFlags)
)
