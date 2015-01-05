// Collection routines for writing objects
package db

import (
	"bytes"
	"encoding/gob"
	"logger"
	. "types"
)

// Writes (inserts/updates) object instance into collection
func (c *Collection) WriteObject(p WriteObject) (uint, error) {
	buf, err := c.encodeObject(&p.Data)
	if err != nil {
		return RObjectEncodeError, err
	}
	offset := c.getFreeSpaceOffset()
	err = c.DataFile.Write(buf.Bytes(), offset)
	if err != nil {
		return RDataWriteError, err
	}
	version := c.addObjectToIndex(&p, offset, buf.Len())
	c.AddObjectToIndices(&p, version)
	return RNoError, nil
}

// GOB encodes object
func (c *Collection) encodeObject(data *ObjectFields) (*bytes.Buffer, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(data)
	if err != nil {
		logger.ErrorLog.Printf("%s", err)
		return nil, err
	}
	return &b, nil
}
