// Collection routines for writing objects
package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"logger"
	"reflect"
	"strconv"
	"strings"
	. "types"
)

// Check object for validity (index types)
func (c *Collection) ObjectValid(p *WriteObject) (uint, error) {
	for indexName, index := range c.Indices {
		if val, ok := p.Data[indexName]; ok {
			// If type mismatch detected
			if reflect.TypeOf(val) != index.GetType() {
				// try to sanitize (convert type)
				switch val.(type) {
				case float64: // Convert float64 to ...
					switch index.GetType().String() {
					case "string":
						p.Data[indexName] = strings.TrimRight(strings.TrimRight(strconv.FormatFloat(val.(float64), 'f', 64, 64), "0"), ".")
					case "int":
						p.Data[indexName] = int(val.(float64))
					}
				case int: // Convert int to ...
					switch index.GetType().String() {
					case "string":
						p.Data[indexName] = strconv.Itoa(val.(int))
					case "float64":
						p.Data[indexName] = float64(val.(float64))
					}
				case string: // Convert string to ...
					switch index.GetType().String() {
					case "int":
						intVal, err := strconv.Atoi(val.(string))
						if err == nil {
							p.Data[indexName] = intVal
						} else {
							return RObjectInvalid, errors.New(fmt.Sprintf("Type of field %s invalid: expected %s, got %s", indexName, index.GetType(), reflect.TypeOf(val)))
						}
					case "float64":
						floatVal, err := strconv.ParseFloat(val.(string), 64)
						if err == nil {
							p.Data[indexName] = floatVal
						} else {
							return RObjectInvalid, errors.New(fmt.Sprintf("Type of field %s invalid: expected %s, got %s", indexName, index.GetType(), reflect.TypeOf(val)))
						}
					}
				default:
					return RObjectInvalid, errors.New(fmt.Sprintf("Type of field %s invalid: expected %s, got %s", indexName, index.GetType(), reflect.TypeOf(val)))
				}
			}
		}
	}
	return RNoError, nil
}

// Writes (inserts/updates) object instance into collection
func (c *Collection) WriteObject(p WriteObject) (uint, error) {
	if code, err := c.ObjectValid(&p); err != nil {
		return code, err
	}
	buf, err := c.encodeObject(&p.Data)
	if err != nil {
		return RObjectEncodeError, err
	}
	offset := c.getFreeSpaceOffset()
	err = c.DataFile.Write(buf.Bytes(), offset)
	if err != nil {
		return RDataWriteError, err
	}
	go c.AddObjectToIndices(&p, c.addObjectToIndex(&p, offset, buf.Len()))
	return RNoError, nil
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
