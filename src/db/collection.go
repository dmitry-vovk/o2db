// Collection definition and methods to work with collection objects
package db

import (
	"bytes"
	"encoding/gob"
	"logger"
	. "types"
)

type Hash [20]byte // SHA1 hash

type ObjectIndex map[Hash][]uint64

type Collection struct {
	Name           string                 // Collection/class name
	Objects        map[uint64]interface{} // Objects
	Indices        map[string]ObjectIndex // collection of indices
	DataFile       *DbFile                // Objects storage
	IndexFile      map[string]*DbFile     // List of indices
	freeSlotOffset int
}

// Writes (inserts/updates) object instance into collection
func (this *Collection) WriteObject(p WriteObject) error {
	// TODO
	// 1. Encode object
	// 2. Write to file
	// 3. Add field indices
	// 4. Add starting position and length to data index
	buf, err := this.encodeObject(&p.Data)
	f, err := OpenFile(this.DataFile.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	//len := buf.Len()
	offset := this.getFreeSpaceOffset()
	err = f.Write(buf.Bytes(), offset)
	if err != nil {
		return err
	}
	this.addObjectToIndex(&p.Data, offset, buf.Len())
	return nil
}

// GOB encodes object
func (this *Collection) encodeObject(data *ObjectFields) (*bytes.Buffer, error) {
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
func (this *Collection) getFreeSpaceOffset() int {
	return this.freeSlotOffset
}

// Adds object to indices
func (this *Collection) addObjectToIndex(data *ObjectFields, offset, len int) {
	this.freeSlotOffset += len

}

// Reads object from collection file
func (this *Collection) ReadObject(p ReadObject) (*ObjectFields, error) {
	// TODO
	// 1. Get conditions from p
	// 2. Figure corresponding indices
	// 3. Find object id according to indices
	// 4. Figure out starting position and length of GOB record
	// 5. Read and decode it
	f, err := OpenFile(this.DataFile.FileName)
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
