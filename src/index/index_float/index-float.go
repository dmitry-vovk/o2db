// Maintains index of integer values of versions and ids
// 	int(value) > []int(versions) > ids
package index_float

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
	MapV map[float64]idList // versioned index to id map
	Map  map[float64][]int  // index to id map
}

type FloatIndex struct {
	Name          string    // Field name
	maps          maps      // maps of value to id/version
	IndexFileName string    // Name of the index file storage
	Flush         chan bool // chan to let index know to flush the index to file
}

// Create new empty string index
func NewFloatIndex(fileName string) *FloatIndex {
	floatIndex := FloatIndex{
		maps: maps{
			MapV: make(map[float64]idList),
			Map:  make(map[float64][]int),
		},
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
	err = dec.Decode(&i.maps)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *FloatIndex) encode() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(i.maps)
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
					if os.IsNotExist(err) {
						return
					}
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
	index := value.(float64)
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
func (i *FloatIndex) Delete(value interface{}, id, version int) {
	index := value.(float64)
	i.deleteMostRecent(index, id)
	i.deleteVersioned(index, id, version)
}

func (i *FloatIndex) deleteVersioned(index float64, id, version int) {
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

func (i *FloatIndex) deleteMostRecent(index float64, id int) {
	if i.maps.Map[index] != nil {
		for n, eid := range i.maps.Map[index] {
			if id == eid {
				i.maps.Map[index] = append(i.maps.Map[index][:n], i.maps.Map[index][n+1:]...)
			}
		}
	}
}

// Find list of ids
func (i *FloatIndex) Find(value interface{}) []int {
	return i.maps.Map[value.(float64)]
}

func (i *FloatIndex) ConditionalFind(op string, value interface{}) []int {
	ids := []int{}
	index := value.(float64)
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

func (i *FloatIndex) GetType() reflect.Type {
	return reflect.TypeOf(0.1)
}
