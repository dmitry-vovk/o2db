package db

import (
	. "types"
)

// Compares two object versions and returns list of differentiating fields
func (c *Collection) GetObjectDiff(p GetObjectDiff) (ObjectDiff, uint, error) {
	obj1, code, err := c.getObjectByIdAndVersion(p.Id, p.From)
	if err != nil {
		return ObjectDiff{}, code, err
	}
	obj2, code, err := c.getObjectByIdAndVersion(p.Id, p.To)
	if err != nil {
		return ObjectDiff{}, code, err
	}
	var diff ObjectDiff = make(map[string]interface{})
	o1 := *obj1
	o2 := *obj2
	for k, v := range o1 {
		if o1[k] != o2[k] {
			diff[k] = v
		}
	}
	return diff, RNoError, nil
}

// Returns the number of object versions
func (c *Collection) GetObjectVersions(p GetObjectVersions) (ObjectVersions, uint, error) {
	return ObjectVersions(len(c.Objects[p.Id])), RNoError, nil
}
