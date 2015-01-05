// Collection routines for object reading
package db

import (
	"bytes"
	"encoding/gob"
	. "types"
)

// Reads object from collection file
func (c *Collection) ReadObject(p ReadObject) (*ObjectFields, uint, error) {
	// Special case when object selected by ID only
	var id int
	// Prettify ID
	if rawId, ok := p.Fields[FIELD_ID]; ok {
		id = getInt(rawId)
	}
	// Object with zero ID cannot exist
	if id == 0 {
		return nil, RObjectDoesNotExist, nil
	}
	// No objects there
	if len(c.Objects[id]) == 0 {
		return nil, RObjectNotFound, nil
	}
	// Get by ID
	if len(p.Fields) == 1 {
		var version = len(c.Objects[id]) - 1
		return c.getObjectByIdAndVersion(id, version)
	}
	// Get by ID and version
	if rawVersion, ok := p.Fields[FIELD_VERSION]; ok && len(p.Fields) == 2 {
		version := getInt(rawVersion)
		return c.getObjectByIdAndVersion(id, version)
	}
	return nil, RObjectNotFound, nil
}

func (c *Collection) getObjectByIdAndVersion(id, version int) (*ObjectFields, uint, error) {
	data, err := c.DataFile.Read(c.Objects[id][version].Offset, c.Objects[id][version].Len)
	if err != nil {
		return nil, RDataReadError, err
	}
	// Decode bytes into object
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	obj := ObjectFields{}
	if err = dec.Decode(&obj); err == nil {
		obj[FIELD_VERSION] = version
		return &obj, RNoError, err
	} else {
		return nil, RObjectDecodeError, nil
	}
}
