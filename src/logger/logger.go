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

func SetLogPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}
	debugLogFile, err := os.OpenFile(path+"/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, (os.FileMode)(0644))
	if err != nil {
		return err
	}
	DebugLog = log.New(debugLogFile, "", log.Lmicroseconds|log.Lshortfile)
	errorLogFile, err := os.OpenFile(path+"/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, (os.FileMode)(0644))
	if err != nil {
		return err
	}
	ErrorLog = log.New(errorLogFile, "Error in ", log.Lshortfile)
	systemLogFile, err := os.OpenFile(path+"/system.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, (os.FileMode)(0644))
	if err != nil {
		return err
	}
	SystemLog = log.New(systemLogFile, "", log.LstdFlags)
	return nil
}
