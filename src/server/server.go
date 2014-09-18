package server

import (
	"bufio"
	"bytes"
	"config"
	"db"
	"encoding/json"
	"io"
	"log"
	"net"
	"server/message"
	. "types"
)

const (
	messageDelimiter byte = 0 // Message delimiter. Every message should end with this byte
)

type ServerType struct {
	Config  *config.ConfigType
	Clients []*ClientType
	Core    *db.DbCore
}

// Create and initialise new server instance
func CreateNew(config *config.ConfigType) *ServerType {
	return &ServerType{
		Config: config,
		Core: &db.DbCore{
			Input: make(chan *db.Package),
		},
	}
}

// Run processing
func (this *ServerType) Run() error {
	go this.Core.Processor()
	socket, err := net.Listen("tcp4", this.Config.ListenTCP)
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
		c := &ClientType{
			Conn: conn,
		}
		go this.handler(c)
	}
}

// Handle single client connection
func (this *ServerType) handler(c *ClientType) {
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
			pkg := &db.Package{
				Container: query,
				Client:    c,
				RespChan:  make(chan Response),
			}
			go handle(pkg)
			this.Core.Input <- pkg
		}
	}
}

// Wait for response from DbCore and send response to client
func handle(in *db.Package) {
	resp := <-in.RespChan
	out, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		return
	}
	log.Printf("Response: %s", out)
	_, err = io.Copy(in.Client.Conn, bytes.NewBuffer(append(out, messageDelimiter)))
	if err != nil {
		log.Printf("Error sending response to client: %s", err)
		return
	}
}
