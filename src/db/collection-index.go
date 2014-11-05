package db

import (
	"logger"
	. "types"
)

func (c *Collection) CreateIndices(fields map[string]Field) {
	if c.Indices == nil {
		c.Indices = make(map[string]FieldIndex)
	}
	for k, v := range fields {
		indexFileName := c.BaseDir + hash(k) + ".index"
		switch v.Type {
		case "string":
			c.Indices[k] = NewStringIndex(indexFileName)
			c.Indices[k].DoFlush()
		case "int":
			c.Indices[k] = NewIntIndex(indexFileName)
			c.Indices[k].DoFlush()
		case "float":
			c.Indices[k] = NewFloatIndex(indexFileName)
			c.Indices[k].DoFlush()
		default:
			logger.ErrorLog.Printf("Index handler of type %s not implemented", v.Type)
		}
	}
}

// TODO fix panic when field value it of wrong type
func (c *Collection) AddObjectToIndices(o *WriteObject, version int) {
	for field, index := range c.Indices {
		if index != nil {
			index.Add(o.Data[field], o.Id, version)
			index.DoFlush()
		}
	}
}
