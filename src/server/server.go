package server

import (
	"config"
	"net"
	"bufio"
	"log"
	"io"
	"server/client"
	"server/message"
	"bytes"
	"encoding/json"
)

const (
	messageDelimiter byte = 0 // Message delimiter. Every message should end with this byte
)

type ServerType struct {
	Config  *config.ConfigType
	Clients []*client.ClientType
}

func CreateNew(config *config.ConfigType) *ServerType {
	c := &ServerType{
		Config: config,
	}
	return c
}

func (s *ServerType) Run() error {
	socket, err := net.Listen("tcp4", s.Config.ListenTCP)
	if err != nil {
		return err
	}
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("Accept error: %s", err)
			continue
		}
		log.Printf("Client connected")
		c := &client.ClientType{
			Conn: conn,
		}
		go s.handler(c)
	}
}

// Handle single client connection
func (s *ServerType) handler(c *client.ClientType) {
	defer c.Conn.Close()
	for {
		msg, err := bufio.NewReader(c.Conn).ReadBytes(messageDelimiter)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client diconnected")
			} else {
				log.Printf("Read error: %s", err)
			}
			return
		}
		query, err := message.Parse(msg[:len(msg) - 1])
		if err != nil {
			log.Printf("%s", err)
		} else {
			log.Printf("Message: %v", query)
			// TODO Process message here and write response to client
			s.respond(c, query)
		}
	}
}

func (s *ServerType) respond(c *client.ClientType, r interface{}) error {
	out, err := json.Marshal(r)
	if err == nil {
		io.Copy(c.Conn, bytes.NewBuffer(append(out, messageDelimiter)))
		return nil
	} else {
		return err
	}
}
