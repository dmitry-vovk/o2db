package db

import (
	. "logger"
	. "types"
)

type Collection struct {
	Name    string
	Objects []interface{}
}

func (this *Collection) WriteObject(p WriteObject) error {
	DebugLog.Printf("Writing object %v", p)
	return nil
}
