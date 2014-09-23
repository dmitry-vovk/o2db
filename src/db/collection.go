// Collection definition and methods to work with collection objects
package db

import (
	"logger"
	. "logger"
	. "types"
)

type Collection struct {
	Name    string
	Objects []interface{}
}

// Writes (inserts/updates) object instance into collection
func (this *Collection) WriteObject(p WriteObject) error {
	DebugLog.Printf("Writing object %v", p)
	f, err := OpenFile("file.txt")
	if err != nil {
		logger.ErrorLog.Printf("Open file: %s", err)
		return err
	}
	err = f.Write([]byte("Hello there"), 100)
	if err != nil {
		logger.ErrorLog.Printf("Writing: %s", err)
		return err
	}
	out, err := f.Read(100, 5)
	if err != nil {
		return err
	}
	logger.ErrorLog.Printf("Read from file: %s", out)
	return nil
}
