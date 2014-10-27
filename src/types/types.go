package types

import "net"

// Query message container
type Container struct {
	Type    uint        `json:"type"`    // one of the Type... constant values
	Payload interface{} `json:"payload"` // any type of payload
}

// Response container
type Response struct {
	Result   bool        `json:"result"`
	Response interface{} `json:"response"`
}

// Authentication request payload
type Authentication struct {
	Name     string `json:"name"`     // username
	Password string `json:"password"` // and password
}

// Create database request
type CreateDatabase struct {
	Name string `json:"name"` // Name for the new database. Must be correct file name.
}

// Drop database request
type DropDatabase struct {
	Name string `json:"name"` // Name for the database
}

// List databases according to mask
type ListDatabases struct {
	Mask string `json:"mask"` // Mask is glob expression
}

// Open (and set as default for connected client) database
type OpenDatabase struct {
	Name string `json:"name"` // Name of the database to open
}

// Create collection with Name and Fields
type CreateCollection struct {
	Name   string           `json:"class"`  // Collection name (class in terms of OOP)
	Fields map[string]Field `json:"fields"` // Collection of named fields (indices)
}

// Collection field description
type Field struct {
	Type string `json:"type"` // type
}

type DropCollection struct {
	Name string `json:"class"` // Collection name (class in terms of OOP)
}

type Index struct {
	Name string
}

type Schema struct {
	ClassName string
	Fields    []Field
	Indices   []Index
}

type Object struct {
	Class  Schema
	Id     uint64
	Fields []Field
}

type Client struct {
	Conn          net.Conn
	Authenticated bool
	Db            string
}

type ObjectFields map[string]interface{}

// Query of type TypeObjectWrite
type WriteObject struct {
	Collection string       `json:"class"`
	Id         int          `json:"-"`
	Data       ObjectFields `json:"data"`
}

// Query of type TypeObjectGet
type ReadObject struct {
	Collection string       `json:"class"`
	Fields     ObjectFields `json:"data"`
}

// Query of type TypeGetObjectVersions
type GetObjectVersions struct {
	Collection string `json:"class"`
	Id         int    `json:"id"`
}

// Response
type ObjectVersions int

// Get diff between two object versions
type GetObjectDiff struct {
	Collection string `json:"class"`
	Id         int    `json:"id"`
	From       int    `json:"from"` // Version one
	To         int    `json:"to"`   // Version two
}

// Difference between two object versions
type ObjectDiff ObjectFields
