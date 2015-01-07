package index_string

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"logger"
	"os"
	"reflect"
	"time"
)

const flushDelay = 100 * time.Millisecond

type hashIndex [20]byte

type versionsList []int

type idList map[int]versionsList

type maps struct {
	MapV map[hashIndex]idList // versioned index to id map
	Map  map[hashIndex][]int  // index to id map
}
type StringIndex struct {
	Name          string    // Field name
	maps          maps      // index to id map
	IndexFileName string    // Name of the index file storage
	Flush         chan bool // chan to let index know to flush the index to file
}

// Create new empty string index
func NewStringIndex(fileName string) *StringIndex {
	stringIndex := StringIndex{
		maps: maps{
			MapV: make(map[hashIndex]idList),
			Map:  make(map[hashIndex][]int),
		},
	}
	stringIndex.IndexFileName = fileName
	stringIndex.Flush = make(chan bool, 100)
	go stringIndex.indexFlusher()
	return &stringIndex
}

func (i *StringIndex) encode() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(i.maps)
	return b.Bytes()
}

func (i *StringIndex) DoFlush() {
	i.Flush <- true
}

func (i *StringIndex) indexFlusher() {
	var flag bool = false
	for {
		select {
		case <-i.Flush:
			flag = true
		default:
			if flag {
				if err := i.FlushToFile(); err != nil {
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
func (i *StringIndex) FlushToFile() error {
	handler, err := os.OpenFile(i.IndexFileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
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
	err = dec.Decode(&i.maps)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *StringIndex) getEncodedData() *bytes.Buffer {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(i.maps)
	return &b
}

// Return list of IDs matching value
func (i *StringIndex) Find(value interface{}) []int {
	return i.maps.Map[i.getHash(value.(string))]
}

// Does not implement
func (i *StringIndex) ConditionalFind(op string, value interface{}) []int {
	ids := []int{}
	stringVal := i.getHash(value.(string))
	switch op {
	case "!=", "<>", "ne": // not equal
		for v, n := range i.maps.Map {
			if v != stringVal {
				ids = append(ids, n...)
			}
		}
	}
	return ids
}

// Add value to index
func (i *StringIndex) Add(value interface{}, id, version int) {
	index := i.getHash(value.(string))
	if i.maps.MapV[index] == nil {
		i.maps.MapV[index] = idList{}
	}
	if i.maps.MapV[index][id] == nil {
		i.maps.MapV[index][id] = versionsList{}
	}
	i.maps.MapV[index][id] = append(i.maps.MapV[index][id], version)
	i.deleteMostRecent(index, id)
	i.maps.Map[index] = append(i.maps.Map[index], id)
}

// Remove id associated with value
func (i *StringIndex) Delete(value interface{}, id, version int) {
	index := i.getHash(value.(string))
	i.deleteMostRecent(index, id)
	i.deleteVersioned(index, id, version)
}

func (i *StringIndex) deleteVersioned(index hashIndex, id, version int) {
	if i.maps.MapV[index][id] != nil {
		versions := i.maps.MapV[index][id]
		for n, ver := range versions {
			if ver == version {
				i.maps.MapV[index][id] = append(versions[:n], versions[n+1:]...)
				if len(i.maps.MapV[index][id]) == 0 {
					delete(i.maps.MapV[index], id)
				}
				break
			}
		}
	}
}

func (i *StringIndex) deleteMostRecent(index hashIndex, id int) {
	if i.maps.Map[index] != nil {
		for n, eid := range i.maps.Map[index] {
			if id == eid {
				i.maps.Map[index] = append(i.maps.Map[index][:n], i.maps.Map[index][n+1:]...)
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

func (i *StringIndex) GetType() reflect.Type {
	return reflect.TypeOf("string")
}
