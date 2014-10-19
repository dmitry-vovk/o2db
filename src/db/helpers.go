// Helper functions
package db

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
)

// Shorthand to get SHA1 string
func hash(s string) string {
	sh := sha1.New()
	sh.Write([]byte(s))
	return hex.EncodeToString(sh.Sum(nil))
}

// Convert value into int whatever it is
// JSON encoded number can come as string or float
func getInt(in interface{}) int {
	var id int
	switch in.(type) {
	case int:
		id = in.(int)
	case uint:
		id = in.(int)
	case string:
		id, _ = strconv.Atoi(in.(string))
	case float64:
		id = int(in.(float64))
	}
	return id
}
