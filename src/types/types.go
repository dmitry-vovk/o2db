package types

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	. "logger"
	"net"
)

// Query message container
type Container struct {
	Type    uint        `json:"type"`    // one of the Type... constant values
	Payload interface{} `json:"payload"` // any type of payload
}

// Response container
type Response struct {
	Result   bool        `json:"result"`
	Code     uint        `json:"code"`
	Response interface{} `json:"response"`
}

// Authentication request payload
type Authentication struct {
	Name     string `json:"name"`     // username
	Password string `json:"password"` // and password
}

// Create database request
type CreateDatabase struct {
	Name string `json:"name"` // Name for the new database. Must be correct file name.
}

// Drop database request
type DropDatabase struct {
	Name string `json:"name"` // Name for the database
}

// List databases according to mask
type ListDatabases struct {
	Mask string `json:"mask"` // Mask is glob expression
}

// Open (and set as default for connected client) database
type OpenDatabase struct {
	Name string `json:"name"` // Name of the database to open
}

// Create collection with Name and Fields
type CreateCollection struct {
	Name   string           `json:"class"`  // Collection name (class in terms of OOP)
	Fields map[string]Field `json:"fields"` // Collection of named fields (indices)
}

// Collection field description
type Field struct {
	Type string `json:"type"` // type
}

type DropCollection struct {
	Name string `json:"class"` // Collection name (class in terms of OOP)
}

type Index struct {
	Name string
}

type Schema struct {
	ClassName string
	Fields    []Field
	Indices   []Index
}

type Object struct {
	Class  Schema
	Id     uint64
	Fields []Field
}

type Client struct {
	Conn          net.Conn
	WsConn        *websocket.Conn
	Authenticated bool
	Db            string
}

func (c *Client) Respond(resp Response) {
	out, err := json.Marshal(resp)
	if err != nil {
		ErrorLog.Printf("Error encoding response: %s", err)
		return
	}
	DebugLog.Printf("Response: %s", out)
	if c.WsConn != nil {
		c.WsConn.WriteMessage(websocket.TextMessage, out)
	} else {
		_, err = io.Copy(c.Conn, bytes.NewBuffer(append(out, MessageDelimiter)))
		if err != nil {
			ErrorLog.Printf("Error sending response to client: %s", err)
			return
		}
	}
}

type ObjectFields map[string]interface{}

// Query of type TypeObjectWrite
type WriteObject struct {
	Collection string       `json:"class"`
	Id         int          `json:"-"`
	Data       ObjectFields `json:"data"`
}

// Query of type TypeObjectGet
type ReadObject struct {
	Collection string       `json:"class"`
	Fields     ObjectFields `json:"data"`
}

// Query of type TypeGetObjectVersions
type GetObjectVersions struct {
	Collection string `json:"class"`
	Id         int    `json:"id"`
}

// Response
type ObjectVersions int

// Get diff between two object versions
type GetObjectDiff struct {
	Collection string `json:"class"`
	Id         int    `json:"id"`
	From       int    `json:"from"` // Version one
	To         int    `json:"to"`   // Version two
}

type SelectObjects struct {
	Collection string       `json:"class"`
	Query      ObjectFields `json:"query"`
}

// Difference between two object versions
type ObjectDiff ObjectFields

// Create a subscription for updates when new objects are written
type AddSubscription struct {
	Collection string       `json:"class"` // Collection name
	Key        string       `json:"key"`   // Authorisation key
	Query      ObjectFields `json:"query"` // Set of conditions for events
}

// Remove subscription
type CancelSubscription struct {
	Collection string `json:"class"` // Collection name
	Key        string `json:"key"`   // Authorisation key
}

// Client call to receive object updates by specified key
type Subscribe struct {
	Database   string `json:"database"`
	Collection string `json:"class"`
	Key        string `json:"key"`
}

type SubscriptionItem struct {
	Collection string       `json:"collection"`
	Key        string       `json:"key"`
	Query      ObjectFields `json:"query"`
}

type ListSubscriptions struct {
	Collections []string `json:"classes"`
}
