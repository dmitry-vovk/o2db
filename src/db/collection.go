// Collection definition and methods to work with collection objects
package db

import (
	"bytes"
	"encoding/gob"
	"github.com/kr/pretty"
	"logger"
	"os"
	"time"
	. "types"
	"strconv"
)

const (
	flushDelay = 100 * time.Millisecond
)

type Hash [20]byte // SHA1 hash

type ObjectIndex map[Hash][]int

// Object instance
type ObjectVersion struct {
	Offset int // Offset from the beginning of object data file
	Len    int // Number of bytes to read
}

// List of object instances. index is version number.
type ObjectPointer map[int]ObjectVersion

type Collection struct {
	Name             string                 // Collection/class name
	Objects          map[int]ObjectPointer  // Objects. map index is object ID
	Indices          map[string]ObjectIndex // collection of indices
	DataFile         *DbFile                // Objects storage
	IndexFile        map[string]*DbFile     // List of indices
	freeSlotOffset   int
	IndexPointerFile string
	ObjectIndexFlush chan (bool)
}

// Writes (inserts/updates) object instance into collection
func (c *Collection) WriteObject(p WriteObject) error {
	// TODO
	// 1. Encode object
	// 2. Write to file
	// 3. Add field indices
	// 4. Add starting position and length to data index
	buf, err := c.encodeObject(&p.Data)
	f, err := OpenFile(c.DataFile.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	//len := buf.Len()
	offset := c.getFreeSpaceOffset()
	err = f.Write(buf.Bytes(), offset)
	if err != nil {
		return err
	}
	c.addObjectToIndex(&p, offset, buf.Len())
	return nil
}

// GOB encodes object
func (c *Collection) encodeObject(data *ObjectFields) (*bytes.Buffer, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(data)
	if err != nil {
		logger.ErrorLog.Printf("%s", err)
		return nil, err
	}
	return &b, nil
}

// Returns pointer to the start of unallocated file space
func (c *Collection) getFreeSpaceOffset() int {
	return c.freeSlotOffset
}

// Adds object to indices
func (c *Collection) addObjectToIndex(wo *WriteObject, offset, length int) {
	c.freeSlotOffset += length
	switch wo.Data["id"].(type) {
	case string:
		id, _ := strconv.Atoi(wo.Data["id"].(string))
		wo.Id = id
	case float64: wo.Id = int(wo.Data["id"].(float64))
	default:
		wo.Id = len(c.Objects)
	}
	next := len(c.Objects[wo.Id])
	if len(c.Objects[wo.Id]) == 0 {
		c.Objects[wo.Id] = ObjectPointer{}
	}
	c.Objects[wo.Id][next] = ObjectVersion{
		Len:    length,
		Offset: offset,
	}
	c.ObjectIndexFlush <- true
}

// Goroutine to trigger object index flushing to disk
func (c *Collection) objectIndexFlusher() {
	var flag bool = false
	for {
		select {
		case <-c.ObjectIndexFlush:
			flag = true
		default:
			if flag {
				if err := c.flushObjectIndex(); err != nil {
					logger.ErrorLog.Printf("Error flushing objects index: %s", err)
				}
				flag = false
			} else {
				time.Sleep(flushDelay)
			}
		}
	}
}

// Dump index structure to disk
func (c *Collection) flushObjectIndex() error {
	var b bytes.Buffer
	logger.ErrorLog.Printf("%# v", pretty.Formatter(c.Objects))
	enc := gob.NewEncoder(&b)
	err := enc.Encode(c.Objects)
	if err != nil {
		return err
	}
	handler, err := os.OpenFile(c.IndexPointerFile, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	_, err = handler.Write(b.Bytes())
	return err
}

// Reads object from collection file
func (c *Collection) ReadObject(p ReadObject) (*ObjectFields, error) {
	// TODO
	// 1. Get conditions from p
	// 2. Figure corresponding indices
	// 3. Find object id according to indices
	// 4. Figure out starting position and length of GOB record
	// 5. Read and decode it
	f, err := OpenFile(c.DataFile.FileName)
	if err != nil {
		logger.ErrorLog.Printf("Open file: %s", err)
		return nil, err
	}
	defer f.Close()
	dec := gob.NewDecoder(f.Handler)
	var fields ObjectFields
	err = dec.Decode(&fields)
	if err != nil {
		logger.ErrorLog.Printf("Decoding: %s", err)
		return nil, err
	}
	logger.DebugLog.Printf("Read fields: %v", fields)
	return &fields, nil
}
