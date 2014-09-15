package types

const (
	TypeAuthenticate uint = 0

	TypeCreateDatabase uint = 100
	TypeDropDatabase   uint = 101
	TypeOpenDatabase   uint = 102
	TypeListDatabases  uint = 103

	TypeCreateCollection uint = 200
	TypeDropCollection   uint = 201
	TypeListCollections  uint = 202

	TypeObjectGet    uint = 300
	TypeObjectWrite  uint = 301
	TypeObjectDelete uint = 302

	TypeTransactionStart  uint = 400
	TypeTransactionCommit uint = 401
	TypeTransactionAbort  uint = 402
)

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
type Authenticate struct {
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
	Fields map[string]Field `json:"fields"` // Collection of named fields
}

// Collection field description
type Field struct {
	Type  string `json:"type"`   // type
	Index string `json:"length"` // index or not. Index can be 'primary' or 'secondary'.
}

type DropCollection struct {
	Name string `json:"class"` // Collection name (class in terms of OOP)
}
