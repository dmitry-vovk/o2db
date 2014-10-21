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
)

const (
	flushDelay    = 100 * time.Millisecond
	FIELD_ID      = "id"
	FIELD_VERSION = "__version__"
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
	buf, err := c.encodeObject(&p.Data)
	if err != nil {
		return err
	}
	offset := c.getFreeSpaceOffset()
	err = c.DataFile.Write(buf.Bytes(), offset)
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
	wo.Id = getInt(wo.Data["id"])
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
	defer handler.Close()
	_, err = handler.Write(b.Bytes())
	return err
}

// Read index structure from disk
func (c *Collection) readObjectIndex() error {
	handler, err := os.Open(c.IndexPointerFile)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(handler)
	err = dec.Decode(&c.Objects)
	if err != nil {
		return err
	}
	// Find the end of storage space
	for _, obj := range c.Objects {
		for _, ver := range obj {
			end := ver.Offset + ver.Len
			if end > c.freeSlotOffset {
				c.freeSlotOffset = end
			}
		}
	}
	return handler.Close()
}

// Reads object from collection file
func (c *Collection) ReadObject(p ReadObject) (*ObjectFields, error) {
	// Special case when object selected by ID only
	var id int
	// Prettify ID
	if rawId, ok := p.Fields[FIELD_ID]; ok {
		id = getInt(rawId)
	}
	// Object with zero ID cannot exist
	if id == 0 {
		return nil, nil
	}
	// No objects there
	if len(c.Objects[id]) == 0 {
		return nil, nil
	}
	// Get by ID
	if len(p.Fields) == 1 {
		var version = len(c.Objects[id]) - 1
		return c.getObjectByIdAndVersion(id, version)
	}
	// Get by ID and version
	if rawVersion, ok := p.Fields[FIELD_VERSION]; ok && len(p.Fields) == 2 {
		version := getInt(rawVersion)
		return c.getObjectByIdAndVersion(id, version)
	}
	// TODO get by conditions
	return nil, nil
}

func (c *Collection) getObjectByIdAndVersion(id, version int) (*ObjectFields, error) {
	data, err := c.DataFile.Read(c.Objects[id][version].Offset, c.Objects[id][version].Len)
	if err != nil {
		return nil, err
	}
	// Decode bytes into object
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	obj := ObjectFields{}
	err = dec.Decode(&obj)
	obj[FIELD_VERSION] = version
	return &obj, err
}

// Returns the number of object versions
func (c *Collection) GetObjectVersions(p GetObjectVersions) (ObjectVersions, error) {
	return ObjectVersions(len(c.Objects[p.Id])), nil
}
