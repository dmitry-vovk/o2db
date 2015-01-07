package types

import "reflect"

// Common interface for field index
// Typically index handler will have more methods, these are common for all handlers
type FieldIndex interface {
	Add(value interface{}, id, version int)           // Add value to the index
	Find(value interface{}) []int                     // Find ids by value (direct match)
	ConditionalFind(op string, val interface{}) []int // Find ids by special condition
	Delete(value interface{}, id, version int)        // Delete value - id association
	FlushToFile() error                               // Write index to file
	DoFlush()                                         // trigger flushing
	GetType() reflect.Type                            // Returns type of index
	// TODO add versioned methods to find an object
}
