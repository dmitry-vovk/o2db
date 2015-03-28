package client

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	. "logger"
	"net"
	. "types"
)

type Client struct {
	Conn          net.Conn
	WsConn        *websocket.Conn
	Authenticated bool
	Db            string
	QueryCounter  uint
}

func (c *Client) Respond(resp Response) {
	out, err := json.Marshal(resp)
	if err != nil {
		ErrorLog.Printf("Error encoding response: %s", err)
		return
	}
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
