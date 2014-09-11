package auth

import (
	"server/client"
	"server/message"
	"config"
)

// Check client credentials
func Authenticate(c *client.ClientType, p message.Payload) bool {
	if p["name"] == config.Config.User.Name && p["password"] == config.Config.User.Password {
		c.Authenticated = true
	}
	return c.Authenticated
}
