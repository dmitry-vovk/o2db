package server

import (
	"bufio"
	"bytes"
	"config"
	"db"
	"encoding/json"
	"fmt"
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
	Clients []*Client
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
		c := &Client{
			Conn: conn,
		}
		go this.handler(c)
	}
}

// Handle single client connection
func (this *ServerType) handler(client *Client) {
	defer client.Conn.Close()
	for {
		log.Print("---------------------------------------------")
		msg, err := bufio.NewReader(client.Conn).ReadBytes(messageDelimiter)
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
			log.Printf("Parser returned error: %s", err)
			handle(
				Response{
					Result:   false,
					Response: fmt.Sprintf("%s", err),
				},
				client,
			)
		} else {
			log.Printf("Message: %v", query)
			pkg := &db.Package{
				Container: query,
				Client:    client,
				RespChan:  make(chan Response),
			}
			go func(in *db.Package) {
				go handle(<-in.RespChan, client)
			}(pkg)
			this.Core.Input <- pkg
		}
	}
}

// Send response to client
func handle(resp Response, client *Client) {
	out, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		return
	}
	log.Printf("Response: %s", out)
	_, err = io.Copy(client.Conn, bytes.NewBuffer(append(out, messageDelimiter)))
	if err != nil {
		log.Printf("Error sending response to client: %s", err)
		return
	}
}
