package types

// Common interface for field index
// Typically index handler will have more methods, these are common for all handlers
type FieldIndex interface {
	Add(value interface{}, id int)     // Add value to the index
	Find(value interface{}) []int      // Find ids by value (direct match)
	Delete(value interface{}, id int)  // Delete value - id association
	FlushToFile(fileName string) error // Write index to file
}
