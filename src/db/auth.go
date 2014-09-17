package db

import (
	"config"
	"types"
)

// Check client credentials
func Authenticate(c *ClientType, p types.TAuthenticate) bool {
	if p.Name == config.Config.User.Name && p.Password == config.Config.User.Password {
		c.Authenticated = true
	}
	return c.Authenticated
}
