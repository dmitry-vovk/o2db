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

type FloatIndex struct {
	Name          string             // Field name
	Map           map[float64]idList // index to id map
	IndexFileName string             // Name of the index file storage
	Flush         chan bool          // chan to let index know to flush the index to file
}

// Create new empty string index
func NewFloatIndex(fileName string) *FloatIndex {
	floatIndex := FloatIndex{
		Map: make(map[float64]idList),
	}
	floatIndex.IndexFileName = fileName
	floatIndex.Flush = make(chan bool, 100)
	go floatIndex.indexFlusher()
	return &floatIndex
}

// Read existing int index from file
func OpenFloatIndex(fileName string) (*FloatIndex, error) {
	handler, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer handler.Close()
	dec := gob.NewDecoder(handler)
	i := NewFloatIndex(fileName)
	err = dec.Decode(&i.Map)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *FloatIndex) encode() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(i.Map)
	return b.Bytes()
}

func (i *FloatIndex) DoFlush() {
	i.Flush <- true
}

func (i *FloatIndex) indexFlusher() {
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
func (i *FloatIndex) FlushToFile() error {
	handler, err := os.OpenFile(i.IndexFileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer handler.Close()
	_, err = handler.Write(i.encode())
	return err
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
