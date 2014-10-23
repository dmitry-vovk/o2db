package db

import (
	"github.com/kr/pretty"
	"logger"
	. "types"
)

func (c *Collection) AddObjectToIndices(o *WriteObject) {
	// TODO add object to indices
	logger.ErrorLog.Printf("Indices: %# v", pretty.Formatter(c.Indices))
}
