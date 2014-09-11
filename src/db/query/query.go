package query

import (
	"db"
	"db/auth"
	"fmt"
	"server/client"
	"server/message"
)

// This is the main entry for processing queries
func ProcessQuery(c *client.ClientType, q *message.Container) interface{} {
	if q.Type == message.TypeAuth {
		if auth.Authenticate(c, q.Payload) {
			return "Authenticated"
		} else {
			return "Authentication failed"
		}
	}
	if c.Authenticated {
		switch q.Type {
		case message.TypeCreateDatabase:
			err := db.CreateDatabase(q.Payload)
			if err == nil {
				return "Database created"
			}
			return string(fmt.Sprintf("%s", err))
			break
		case message.TypeOpenDatabase:
			err := db.OpenDatabase(q.Payload)
			if err == nil {
				return "Database opened"
			}
			return string(fmt.Sprintf("%s", err))
			break
		default:
			return "Unknown query type"
		}
	}
	return "Authentication required"
}
