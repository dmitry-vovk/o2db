package message

const (
	TypeAuth uint = iota

	TypeCreateDatabase
	TypeDropDatabase
	TypeCreateCollection
	TypeDropCollection
	TypeOpenDatabase
	TypeListDatabases
	TypeListCollections

	TypeObjectInsert
	TypeObjectUpdate
	TypeObjectDelete
	TypeObjectSelect

	TypeTransactionStart
	TypeTransactionCommit
	TypeTransactionAbort
)

type Container struct {
	Type    uint    `json:"type"`
	Payload Payload `json:"payload"`
}

type Payload map[string]string

type CreateDatabase struct {
	Name string `json:"name"`
}

type Field struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Length uint   `json:"length"`
}

type Index struct {
	Field string `json:"field"`
	Type  uint   `json:"type"`
}

type CreateCollection struct {
	Name    string  `json:"name"`
	Fields  []Field `json:"fields"`
	Indices []Index `json:"indices"`
}
