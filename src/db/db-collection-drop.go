package db

import (
	"errors"
	"os"
	. "types"
)

// Deletes collection
func (d *Database) DropCollection(p DropCollection) error {
	var hashedName = hash(p.Name)
	var collectionPath = d.DataDir + string(os.PathSeparator) + hashedName
	// Check if collection exists
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return errors.New("Collection does not exist")
	}
	// Close all related files
	d.Collections[hashedName].DataFile.Close()
	for _, f := range d.Collections[hashedName].IndexFile {
		if f != nil {
			f.Close()
		}
	}
	// Delete collection reference from database
	delete(d.Collections, hashedName)
	// Delete all related files
	return os.RemoveAll(collectionPath)
}
