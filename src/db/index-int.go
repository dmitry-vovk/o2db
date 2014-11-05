// Maintains index of integer values of versions and ids
// 	int(value) > []int(versions) > ids
package db

import (
	"bytes"
	"encoding/gob"
	"logger"
	"os"
	"time"
)

type IntIndex struct {
	Name          string         // Field name
	Map           map[int]idList // index to id map
	IndexFileName string         // Name of the index file storage
	Flush         chan bool      // chan to let index know to flush the index to file
}

// Create new empty string index
func NewIntIndex(fileName string) *IntIndex {
	intIndex := IntIndex{
		Map: make(map[int]idList),
	}
	intIndex.IndexFileName = fileName
	intIndex.Flush = make(chan bool, 100)
	go intIndex.indexFlusher()
	return &intIndex
}

// Read existing int index from file
func OpenIntIndex(fileName string) (*IntIndex, error) {
	handler, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer handler.Close()
	dec := gob.NewDecoder(handler)
	i := NewIntIndex(fileName)
	err = dec.Decode(&i.Map)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *IntIndex) encode() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(i.Map)
	return b.Bytes()
}

func (i *IntIndex) DoFlush() {
	i.Flush <- true
}

func (i *IntIndex) indexFlusher() {
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
func (i *IntIndex) FlushToFile() error {
	handler, err := os.OpenFile(i.IndexFileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer handler.Close()
	_, err = handler.Write(i.encode())
	return err
}

// Add value/id/version to index
func (i *IntIndex) Add(value interface{}, id, version int) {
	intVal := value.(int)
	if i.Map[intVal] == nil {
		i.Map[intVal] = idList{}
	}
	if i.Map[intVal][id] == nil {
		i.Map[intVal][id] = versionsList{}
	}
	i.Map[intVal][id] = append(i.Map[intVal][id], version)
}

// Remove value/id/version from the index
func (i *IntIndex) Delete(value interface{}, id, version int) {
	intVal := value.(int)
	if i.Map[intVal][id] != nil {
		versions := i.Map[intVal][id]
		for n, ver := range versions {
			if ver == version {
				i.Map[intVal][id] = append(versions[:n], versions[n+1:]...)
				if len(i.Map[intVal][id]) == 0 {
					delete(i.Map[intVal], id)
				}
				break
			}
		}
	}
}

// Find map["id"]"versions"
func (i *IntIndex) Find(value interface{}) map[int][]int {
	ids := make(map[int][]int)
	for k, v := range i.Map[value.(int)] {
		ids[k] = v
	}
	return ids
}
