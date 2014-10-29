// Maintains index of integer values of versions and ids
// 	int(value) > []int(versions) > ids
package db

import (
	"bytes"
	"encoding/gob"
	"os"
)

type FloatIndex struct {
	Name string             // Field name
	Map  map[float64]idList // index to id map
}

// Create new empty string index
func NewFloatIndex() *FloatIndex {
	return &FloatIndex{
		Map: make(map[float64]idList),
	}
}

// Read existing int index from file
func OpenFloatIndex(fileName string) (*FloatIndex, error) {
	handler, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer handler.Close()
	dec := gob.NewDecoder(handler)
	i := NewFloatIndex()
	err = dec.Decode(&i.Map)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// Add value/id/version to index
func (i *FloatIndex) Add(value interface{}, id, version int) {
	floatVal := value.(float64)
	if i.Map[floatVal] == nil {
		i.Map[floatVal] = idList{}
	}
	if i.Map[floatVal][id] == nil {
		i.Map[floatVal][id] = versionsList{}
	}
	i.Map[floatVal][id] = append(i.Map[floatVal][id], version)
}

// Remove value/id/version from the index
func (i *FloatIndex) Delete(value interface{}, id, version int) {
	floatVal := value.(float64)
	if i.Map[floatVal][id] != nil {
		versions := i.Map[floatVal][id]
		for n, ver := range versions {
			if ver == version {
				i.Map[floatVal][id] = append(versions[:n], versions[n+1:]...)
				if len(i.Map[floatVal][id]) == 0 {
					delete(i.Map[floatVal], id)
				}
				break
			}
		}
	}
}

// Find map["id"]"versions"
func (i *FloatIndex) Find(value interface{}) map[int][]int {
	ids := make(map[int][]int)
	for k, v := range i.Map[value.(float64)] {
		ids[k] = v
	}
	return ids
}

// Write index data to file
func (i *FloatIndex) FlushToFile(fileName string) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(i.Map)
	if err != nil {
		return err
	}
	handler, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer handler.Close()
	_, err = handler.Write(b.Bytes())
	return err
}
