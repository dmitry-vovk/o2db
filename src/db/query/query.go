package query

import (
	"server/message"
	"server/client"
)

// This is the main entry for processing queries
func ProcessQuery(c *client.ClientType, q *message.Container) interface {} {
	if q.Type == message.TypeAuth {
		c.Authenticated = true
		return "client authenticated"
	}
	if c.Authenticated {
		return "client IS authenticated"
	} else {
		return "client IS NOT authenticated"
	}
}
