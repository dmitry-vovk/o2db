package db

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"os"
)

type hashIndex [20]byte

type versionsList []int

type idList map[int]versionsList

type StringIndex struct {
	Name string               // Field name
	Map  map[hashIndex]idList // index to id map
}

// Create new empty string index
func NewStringIndex() *StringIndex {
	return &StringIndex{
		Map: make(map[hashIndex]idList),
	}
}

// Read existing string index from file
func OpenStringIndex(fileName string) (*StringIndex, error) {
	handler, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer handler.Close()
	dec := gob.NewDecoder(handler)
	i := NewStringIndex()
	err = dec.Decode(&i.Map)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// Return list of IDs matching value
func (i *StringIndex) Find(value interface{}) map[int][]int {
	index := i.getHash(value.(string))
	ids := make(map[int][]int)
	for k, v := range i.Map[index] {
		ids[k] = v
	}
	return ids
}

// Add value to index
func (i *StringIndex) Add(value interface{}, id, version int) {
	index := i.getHash(value.(string))
	if i.Map[index] == nil {
		i.Map[index] = idList{}
	}
	if i.Map[index][id] == nil {
		i.Map[index][id] = versionsList{}
	}
	i.Map[index][id] = append(i.Map[index][id], version)
}

// Remove id associated with value
func (i *StringIndex) Delete(value interface{}, id, version int) {
	index := i.getHash(value.(string))
	if i.Map[index][id] != nil {
		versions := i.Map[index][id]
		for n, ver := range versions {
			if ver == version {
				i.Map[index][id] = append(versions[:n], versions[n+1:]...)
				if len(i.Map[index][id]) == 0 {
					delete(i.Map[index], id)
				}
				break
			}
		}
	}
}

// Flush the index to file
func (i *StringIndex) FlushToFile(fileName string) error {
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

func (i *StringIndex) getHash(value string) hashIndex {
	sh := sha1.New()
	sh.Write([]byte(value))
	var s hashIndex
	copy(s[:], sh.Sum(nil)[0:20])
	return s
}
