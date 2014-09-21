package db

import (
	"log"
	. "types"
)

type Collection struct {
	Name    string
	Objects []interface{}
}

func (this *Collection) WriteObject(p WriteObject) error {
	log.Printf("Writing object %v", p)
	return nil
}
