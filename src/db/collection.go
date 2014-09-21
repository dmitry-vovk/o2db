// Collection definition and methods to work with collection objects
package db

import (
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
	return nil
}
