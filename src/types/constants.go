package types

// Query types
const (
	TypeAuthenticate uint = 0
)

const (
	TypeCreateDatabase uint = 100
	TypeDropDatabase   uint = 101
	TypeOpenDatabase   uint = 102
	TypeListDatabases  uint = 103
)

const (
	TypeCreateCollection uint = 200
	TypeDropCollection   uint = 201
	TypeListCollections  uint = 202
)

const (
	TypeObjectGet         uint = 300
	TypeObjectWrite       uint = 301
	TypeObjectDelete      uint = 302
	TypeGetObjectVersions uint = 303
	TypeGetObjectDiff     uint = 304
	TypeFindObjects       uint = 305
)

const (
	TypeSubscribe uint = 400
)

// File names
const (
	DataFileName         = "objects.data"
	PrimaryIndexFileName = "primary.index"
	ObjectIndexFileName  = "objects.index"
)
