package db

type Object struct {
	Class  Schema
	Id     uint64
	Fields []Field
}
