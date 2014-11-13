package db

import (
	. "types"
)

// Compares two object versions and returns list of differentiating fields
func (c *Collection) GetObjectDiff(p GetObjectDiff) (ObjectDiff, error) {
	obj1, err := c.getObjectByIdAndVersion(p.Id, p.From)
	if err != nil {
		return ObjectDiff{}, err
	}
	obj2, err := c.getObjectByIdAndVersion(p.Id, p.To)
	if err != nil {
		return ObjectDiff{}, err
	}
	var diff ObjectDiff = make(map[string]interface{})
	o1 := *obj1
	o2 := *obj2
	for k, v := range o1 {
		if o1[k] != o2[k] {
			diff[k] = v
		}
	}
	return diff, nil
}

// Returns the number of object versions
func (c *Collection) GetObjectVersions(p GetObjectVersions) (ObjectVersions, error) {
	return ObjectVersions(len(c.Objects[p.Id])), nil
}
