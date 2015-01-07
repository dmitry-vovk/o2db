// Maintains index of integer values of versions and ids
// 	int(value) > []int(versions) > ids
package index_int

import (
	"bytes"
	"encoding/gob"
	"logger"
	"os"
	"reflect"
	"time"
)

const flushDelay = 100 * time.Millisecond

type versionsList []int

type idList map[int]versionsList

type maps struct {
	MapV map[int]idList // versioned index to id map
	Map  map[int][]int  // index to id map
}

type IntIndex struct {
	Name          string    // Field name
	maps          maps      // index to id map
	IndexFileName string    // Name of the index file storage
	Flush         chan bool // chan to let index know to flush the index to file
}

// Create new empty string index
func NewIntIndex(fileName string) *IntIndex {
	intIndex := IntIndex{
		maps: maps{
			MapV: make(map[int]idList),
			Map:  make(map[int][]int),
		},
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
	err = dec.Decode(&i.maps)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *IntIndex) encode() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(i.maps)
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
	var index int
	if ind, ok := value.(float64); ok {
		index = int(ind)
	} else if ind, ok := value.(int); ok {
		index = int(ind)
	}
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

// Remove value/id/version from the index
func (i *IntIndex) Delete(value interface{}, id, version int) {
	var index int
	if ind, ok := value.(float64); ok {
		index = int(ind)
	} else if ind, ok := value.(int); ok {
		index = int(ind)
	}
	i.deleteMostRecent(index, id)
	i.deleteVersioned(index, id, version)
}

func (i *IntIndex) deleteVersioned(index, id, version int) {
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

func (i *IntIndex) deleteMostRecent(index, id int) {
	if i.maps.Map[index] != nil {
		for n, eid := range i.maps.Map[index] {
			if id == eid {
				i.maps.Map[index] = append(i.maps.Map[index][:n], i.maps.Map[index][n+1:]...)
			}
		}
	}
}

// Find map["id"]"versions"
func (i *IntIndex) Find(value interface{}) []int {
	var index int
	if ind, ok := value.(float64); ok {
		index = int(ind)
	} else if ind, ok := value.(int); ok {
		index = int(ind)
	}
	return i.maps.Map[index]
}

func (i *IntIndex) ConditionalFind(op string, value interface{}) []int {
	ids := []int{}
	var index int
	if ind, ok := value.(float64); ok {
		index = int(ind)
	} else if ind, ok := value.(int); ok {
		index = int(ind)
	}
	switch op {
	case "<", "lt": // less than
		for v, n := range i.maps.Map {
			if v < index {
				ids = append(ids, n...)
			}
		}
	case ">", "gt": // greater than
		for v, n := range i.maps.Map {
			if v > index {
				ids = append(ids, n...)
			}
		}
	case "<=", "=<", "le": // less or equal
		for v, n := range i.maps.Map {
			if v <= index {
				ids = append(ids, n...)
			}
		}
	case ">=", "=>", "ge": // greater or equal
		for v, n := range i.maps.Map {
			if v >= index {
				ids = append(ids, n...)
			}
		}
	case "!=", "<>", "ne": // not equal
		for v, n := range i.maps.Map {
			if v != index {
				ids = append(ids, n...)
			}
		}
	}
	return ids
}

func (i *IntIndex) GetType() reflect.Type {
	return reflect.TypeOf(1)
}
