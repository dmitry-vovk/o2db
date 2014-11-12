// Database definition and methods to work with database collections
package db

type Database struct {
	DataDir     string
	Collections map[string]*Collection
}
