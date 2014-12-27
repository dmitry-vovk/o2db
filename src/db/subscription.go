package db

import (
	. "types"
)

type Subscription struct {
	Key     string
	Query   ObjectFields
	Clients []*Client
}

// Match the object against subscription's Query
func (s *Subscription) Match(object *ObjectFields) bool {
	return true
}
