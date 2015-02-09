package types

// Query types
const (
	TypeAuthenticate uint = 0

	TypeCreateDatabase uint = 100
	TypeDropDatabase   uint = 101
	TypeOpenDatabase   uint = 102
	TypeListDatabases  uint = 103

	TypeCreateCollection uint = 200
	TypeDropCollection   uint = 201
	TypeListCollections  uint = 202

	TypeObjectGet         uint = 300
	TypeObjectWrite       uint = 301
	TypeObjectDelete      uint = 302
	TypeGetObjectVersions uint = 303
	TypeGetObjectDiff     uint = 304
	TypeSelectObjects     uint = 305

	TypeSubscribe          uint = 400
	TypeAddSubscription    uint = 401
	TypeCancelSubscription uint = 402
	TypeListSubscriptions  uint = 403
)

// Response codes
const (
	RNoError = iota

	RAuthenticated
	RNotAuthenticated
	RAuthenticationRequired // 3

	RDatabaseCreated
	RDatabaseDeleted
	RDatabaseSelected
	RDatabaseList
	RDatabaseAlreadyExist
	RDatabaseNotSelected
	RDatabaseDoesNotExist

	RCollectionCreated
	RCollectionDeleted
	RCollectionAlreadyExists
	RCollectionDoesNotExist
	RCollectionList

	RObject
	RObjectWritten
	RObjectInvalid
	RObjectEncodeError
	RObjectDecodeError
	RObjectDoesNotExist
	RObjectNotFound

	RDataWriteError
	RDataReadError

	RSubscribed
	RUnsubscribed
	RSubscriptionInvalidFormat
	RSubscriptionCreated
	RSubscriptionCancelled
	RSubscriptionAlreadyExists
	RSubscriptionDoesNotExist
	RSubscriptionList

	RUnknownQueryType
)

// File names
const (
	DataFileName         = "objects.data"
	PrimaryIndexFileName = "primary.index"
	ObjectIndexFileName  = "objects.index"
	CollectionSchema     = "schema"
)

const (
	MessageDelimiter byte = 0 // Message delimiter. Every message should end with this byte
)
