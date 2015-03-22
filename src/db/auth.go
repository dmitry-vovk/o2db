package db

import (
	"client"
	"config"
	. "types"
)

// Check client credentials
func (this *DbCore) Authenticate(c *client.Client, p Authentication) bool {
	if p.Name == config.Config.User.Name && p.Password == config.Config.User.Password {
		c.Authenticated = true
	}
	return c.Authenticated
}
