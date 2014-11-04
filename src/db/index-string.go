package db

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"logger"
	"os"
	"time"
)

type hashIndex [20]byte

type versionsList []int

type idList map[int]versionsList

type StringIndex struct {
	Name          string               // Field name
	Map           map[hashIndex]idList // index to id map
	IndexFileName string               // Name of the index file storage
	Flush         chan bool            // chan to let index know to flush the index to file
}

// Create new empty string index
func NewStringIndex(fileName string) *StringIndex {
	stringIndex := StringIndex{
		Map: make(map[hashIndex]idList),
	}
	stringIndex.IndexFileName = fileName
	stringIndex.Flush = make(chan bool, 100)
	go stringIndex.indexFlusher()
	return &stringIndex
}

func (i *StringIndex) encode() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(i.Map)
	return b.Bytes()
}

func (i *StringIndex) indexFlusher() {
	var flag bool = false
	for {
		select {
		case <-i.Flush:
			flag = true
		default:
			if flag {
				if err := i.FlushToFile(i.IndexFileName); err != nil {
					logger.ErrorLog.Printf("Error flushing index: %s", err)
				}
				flag = false
			} else {
				time.Sleep(flushDelay)
			}
		}
	}
}

// Flush the index to file
func (i *StringIndex) FlushToFile(fileName string) error {
	handler, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer handler.Close()
	_, err = handler.Write(i.encode())
	return err
}

// Read existing string index from file
func OpenStringIndex(fileName string) (*StringIndex, error) {
	handler, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer handler.Close()
	dec := gob.NewDecoder(handler)
	i := NewStringIndex(fileName)
	err = dec.Decode(&i.Map)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *StringIndex) getEncodedData() *bytes.Buffer {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(i.Map)
	return &b
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

func (i *StringIndex) getHash(value string) hashIndex {
	sh := sha1.New()
	sh.Write([]byte(value))
	var s hashIndex
	copy(s[:], sh.Sum(nil)[0:20])
	return s
}
