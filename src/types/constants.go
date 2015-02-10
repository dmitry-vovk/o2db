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
	RNoError = iota // 0

	RAuthenticated          // 1
	RNotAuthenticated       // 2
	RAuthenticationRequired // 3

	RDatabaseCreated      // 4
	RDatabaseDeleted      // 5
	RDatabaseSelected     // 6
	RDatabaseList         // 7
	RDatabaseAlreadyExist // 8
	RDatabaseNotSelected  // 9
	RDatabaseDoesNotExist // 10

	RCollectionCreated       // 11
	RCollectionDeleted       // 12
	RCollectionAlreadyExists // 13
	RCollectionDoesNotExist  // 14
	RCollectionList          // 15

	RObject             // 16
	RObjectWritten      // 17
	RObjectInvalid      // 18
	RObjectEncodeError  // 19
	RObjectDecodeError  // 20
	RObjectDoesNotExist // 21
	RObjectNotFound     // 22

	RDataWriteError // 23
	RDataReadError  // 24

	RSubscribed                // 25
	RUnsubscribed              // 26
	RSubscriptionInvalidFormat // 27
	RSubscriptionCreated       // 28
	RSubscriptionCancelled     // 29
	RSubscriptionAlreadyExists // 30
	RSubscriptionDoesNotExist  // 31
	RSubscriptionList          // 32

	RUnknownQueryType // 33
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
