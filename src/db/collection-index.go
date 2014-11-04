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
		indexFileName := hash(k) + ".index"
		switch v.Type {
		case "string":
			c.Indices[k] = NewStringIndex(indexFileName)
		case "int":
			c.Indices[k] = NewIntIndex(indexFileName)
		case "float":
			c.Indices[k] = NewFloatIndex(indexFileName)
		default:
			logger.ErrorLog.Printf("Index handler of type %s not implemented", v.Type)
		}
	}
}

func (c *Collection) AddObjectToIndices(o *WriteObject, version int) {
	// TODO add object to indices
	for field, index := range c.Indices {
		index.Add(o.Data[field], o.Id, version)

	}
	logger.ErrorLog.Printf("Indices: %# v", pretty.Formatter(c.Indices))
}
