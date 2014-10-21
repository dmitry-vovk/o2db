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

const flushDelay = 100 * time.Millisecond

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
