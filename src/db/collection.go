// Collection definition and methods to work with collection objects
package db

import (
	"bytes"
	"encoding/gob"
	"logger"
	"time"
	. "types"
)

const (
	flushDelay = 100 * time.Millisecond
)

type Hash [20]byte // SHA1 hash

type ObjectIndex map[Hash][]int

type ObjectPointer struct {
	Offset int
	Len    int
}

type Collection struct {
	Name             string                 // Collection/class name
	Objects          map[int]ObjectPointer  // Objects
	Indices          map[string]ObjectIndex // collection of indices
	DataFile         *DbFile                // Objects storage
	IndexFile        map[string]*DbFile     // List of indices
	freeSlotOffset   int
	IndexPointerFile *DbFile
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
	logger.ErrorLog.Printf("Object written at offset %d", offset)
	c.freeSlotOffset += length
	logger.ErrorLog.Printf("Next offset is %d", c.freeSlotOffset)
	wo.Id = len(c.Objects)
	c.Objects[wo.Id] = ObjectPointer{offset, length}
	logger.ErrorLog.Printf("Object index: %v", c.Objects)
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
				c.flushObjectIndex()
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
	enc := gob.NewEncoder(&b)
	err := enc.Encode(c.Objects)
	if err != nil {
		return err
	}
	// TODO use simple file writes, not DbFile
	return c.IndexPointerFile.Dump(&b)
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
