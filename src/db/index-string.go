package db

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"os"
)

type hashIndex [20]byte

type idList []int

type StringIndex struct {
	Map map[hashIndex]idList
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
func (i *StringIndex) Find(value string) []int {
	return i.Map[i.getHash(value)]
}

// Add value to index
func (i *StringIndex) Add(value string, id int) {
	index := i.getHash(value)
	i.Map[index] = append(i.Map[index], id)
}

// Remove id associated with value
func (i *StringIndex) Delete(value string, id int) {
	index := i.getHash(value)
	if ids, ok := i.Map[index]; ok {
		for n, item := range ids {
			if item == id {
				i.Map[index] = append(ids[:n], ids[n+1:]...)
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
