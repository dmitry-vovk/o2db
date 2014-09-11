package message

const (
	TypeAuth uint = iota
	TypeCreateDatabase
	TypeCreateCollection
)

type Container struct {
	Type    uint    `json:"type"`
	Payload Payload `json:"payload"`
}

type Payload map[string]string

type CreateDatabase struct {
	Name string `json:"name"`
}
