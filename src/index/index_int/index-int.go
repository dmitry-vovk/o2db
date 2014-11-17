// Maintains index of integer values of versions and ids
// 	int(value) > []int(versions) > ids
package index_int

import (
	"bytes"
	"encoding/gob"
	"logger"
	"os"
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
	intVal := int(value.(int))
	if i.maps.MapV[intVal] == nil {
		i.maps.MapV[intVal] = idList{}
	}
	if i.maps.MapV[intVal][id] == nil {
		i.maps.MapV[intVal][id] = versionsList{}
	}
	i.maps.MapV[intVal][id] = append(i.maps.MapV[intVal][id], version)
	i.maps.Map[intVal] = append(i.maps.Map[intVal], id)
}

// Remove value/id/version from the index
func (i *IntIndex) Delete(value interface{}, id, version int) {
	intVal := value.(int)
	if i.maps.MapV[intVal][id] != nil {
		versions := i.maps.MapV[intVal][id]
		for n, ver := range versions {
			if ver == version {
				i.maps.MapV[intVal][id] = append(versions[:n], versions[n+1:]...)
				if len(i.maps.MapV[intVal][id]) == 0 {
					delete(i.maps.MapV[intVal], id)
				}
				break
			}
		}
	}
	if i.maps.Map[intVal] != nil {
		for n, eid := range i.maps.Map[intVal] {
			if id == eid {
				i.maps.Map[intVal] = append(i.maps.Map[intVal][:n], i.maps.Map[intVal][n+1:]...)
			}
		}
	}
}

// Find map["id"]"versions"
func (i *IntIndex) Find(value interface{}) []int {
	return i.maps.Map[int(value.(int))]
}

func (i *IntIndex) ConditionalFind(op string, value interface{}) []int {
	ids := []int{}
	intVal := int(value.(float64))
	switch op {
	case "<", "lt": // less than
		for v, n := range i.maps.Map {
			if v < intVal {
				ids = append(ids, n...)
			}
		}
	case ">", "gt": // greater than
		for v, n := range i.maps.Map {
			if v > intVal {
				ids = append(ids, n...)
			}
		}
	case "<=", "=<", "le": // less or equal
		for v, n := range i.maps.Map {
			if v <= intVal {
				ids = append(ids, n...)
			}
		}
	case ">=", "=>", "ge": // greater or equal
		for v, n := range i.maps.Map {
			if v >= intVal {
				ids = append(ids, n...)
			}
		}
	case "!=", "<>", "ne": // not equal
		for v, n := range i.maps.Map {
			if v != intVal {
				ids = append(ids, n...)
			}
		}
	}
	return ids
}
