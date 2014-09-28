// Collection definition and methods to work with collection objects
package db

import (
	"bytes"
	"encoding/gob"
	"logger"
	. "logger"
	. "types"
)

type Hash [20]byte // SHA1 hash

type ObjectIndex map[Hash][]uint64

type Collection struct {
	Name      string                 // Collection/class name
	Objects   map[uint64]interface{} // Objects
	Indices   map[string]ObjectIndex // collection of indices
	DataFile  *DbFile                // Objects storage
	IndexFile map[string]*DbFile     // List of indices
}

// Writes (inserts/updates) object instance into collection
func (this *Collection) WriteObject(p WriteObject) error {
	// TODO
	// 1. Encode object
	// 2. Write to file
	// 3. Add field indices
	// 4. Add starting position and length to data index
	DebugLog.Printf("Writing object data %v", p.Data)
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(&p.Data)
	if err != nil {
		logger.ErrorLog.Printf("%s", err)
		return err
	}
	f, err := OpenFile("file.db")
	if err != nil {
		logger.ErrorLog.Printf("Open file: %s", err)
		return err
	}
	defer f.Close()
	len := b.Len()
	logger.DebugLog.Printf("Writing %d bytes", len)
	err = f.Write(b.Bytes(), 0)
	if err != nil {
		logger.ErrorLog.Printf("Writing: %s", err)
		return err
	}
	return nil
}

// Reads object from collection file
func (this *Collection) ReadObject(p ReadObject) (*ObjectFields, error) {
	// TODO
	// 1. Get conditions from p
	// 2. Figure corresponding indices
	// 3. Find object id according to indices
	// 4. Figure out starting position and length of GOB record
	// 5. Read and decode it
	f, err := OpenFile("file.db")
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
