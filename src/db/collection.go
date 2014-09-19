package db

import (
	. "types"
)

type Collection struct {
	Name    string
	Objects []interface{}
}

func (this *Collection) WriteObject(p WriteObject) error {
	return nil
}
