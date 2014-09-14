package auth

import (
	"server/client"
	"config"
	"server/types"
)

// Check client credentials
func Authenticate(c *client.ClientType, p types.Authenticate) bool {
	if p.Name == config.Config.User.Name && p.Password == config.Config.User.Password {
		c.Authenticated = true
	}
	return c.Authenticated
}
