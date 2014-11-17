// Collection routines for object reading
package db

import (
	"bytes"
	"encoding/gob"
	. "types"
)

// Reads object from collection file
func (c *Collection) ReadObject(p ReadObject) (*ObjectFields, error) {
	// Special case when object selected by ID only
	var id int
	// Prettify ID
	if rawId, ok := p.Fields[FIELD_ID]; ok {
		id = getInt(rawId)
	}
	// Object with zero ID cannot exist
	if id == 0 {
		return nil, nil
	}
	// No objects there
	if len(c.Objects[id]) == 0 {
		return nil, nil
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
	return nil, nil
}

func (c *Collection) getObjectByIdAndVersion(id, version int) (*ObjectFields, error) {
	data, err := c.DataFile.Read(c.Objects[id][version].Offset, c.Objects[id][version].Len)
	if err != nil {
		return nil, err
	}
	// Decode bytes into object
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	obj := ObjectFields{}
	err = dec.Decode(&obj)
	obj[FIELD_VERSION] = version
	return &obj, err
}
