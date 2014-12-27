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
func (s *Subscription) Match(object ObjectFields) bool {
	// go over all queries
	for k, v := range s.Query {
		// see if the object has the field
		if field, ok := object[k]; ok {
			// TODO compare?
			if v != field {
				return false
			}
		}
	}
	return true
}
