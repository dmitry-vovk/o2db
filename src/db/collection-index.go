package db

import (
	"github.com/kr/pretty"
	"logger"
	. "types"
)

func (c *Collection) CreateIndices(fields map[string]Field) {
	if c.Indices == nil {
		c.Indices = make(map[string]FieldIndex)
	}
	for k, v := range fields {
		switch v.Type {
		case "string":
			c.Indices[k] = NewStringIndex()
		case "int":
			c.Indices[k] = NewIntIndex()
		case "float":
			c.Indices[k] = NewFloatIndex()
		default:
			logger.ErrorLog.Printf("Index handler of type %s not implemented", v.Type)
		}
	}
}

func (c *Collection) AddObjectToIndices(o *WriteObject) {
	// TODO add object to indices
	logger.ErrorLog.Printf("Indices: %# v", pretty.Formatter(c.Indices))
}
