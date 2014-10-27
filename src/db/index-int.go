package db

import (
	"bytes"
	"encoding/gob"
	"os"
)

type IntIndex struct {
	Map map[int]idList
}

// Create new empty string index
func NewIntIndex() *IntIndex {
	return &IntIndex{
		Map: make(map[int]idList),
	}
}

// Read existing int index from file
func OpenIntIndex(fileName string) (*IntIndex, error) {
	handler, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer handler.Close()
	dec := gob.NewDecoder(handler)
	i := NewIntIndex()
	err = dec.Decode(&i.Map)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *IntIndex) Add(value interface{}, id int) {
	i.Map[value.(int)] = append(i.Map[value.(int)], id)
}

func (i *IntIndex) Delete(value interface{}, id int) {
	index := value.(int)
	if ids, ok := i.Map[index]; ok {
		for n, item := range ids {
			if item == id {
				i.Map[index] = append(ids[:n], ids[n+1:]...)
				break
			}
		}
	}
}

func (i *IntIndex) Find(value interface{}) []int {
	return i.Map[value.(int)]
}

func (i *IntIndex) FlushToFile(fileName string) error {
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
