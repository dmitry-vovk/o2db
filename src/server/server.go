package server

import (
	"bufio"
	"bytes"
	"config"
	dbQuery "db/query"
	"io"
	"log"
	"net"
	"server/client"
	"server/message"
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
		log.Print("---------------------------------------------")
		msg, err := bufio.NewReader(c.Conn).ReadBytes(messageDelimiter)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client diconnected")
			} else {
				log.Printf("Read error: %s", err)
			}
			return
		}
		query, err := message.Parse(msg[:len(msg)-1]) // cut out delimiter
		if err != nil {
			log.Printf("Parse error: %s", err)
			// TODO add proper handling
			// err = s.respond(c, fmt.Sprintf("%s", err))
		} else {
			log.Printf("Message: %v", query)
			out := dbQuery.ProcessQuery(c, query)
			_, err = io.Copy(c.Conn, bytes.NewBuffer(append(out, messageDelimiter)))
			if err != nil {
				log.Printf("Error sending response to client: %s", err)
				return
			}
		}
	}
}
