package server

import (
	"config"
	"net"
	"bufio"
	"strings"
	"log"
	"io"
	"bytes"
	"server/client"
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
		go s.handler(conn)
	}
}

func (s *ServerType) handler(c net.Conn) {
	defer c.Close()
	for {
		line, err := bufio.NewReader(c).ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Client diconnected")
			} else {
				log.Printf("Read error: %s", err)
			}
			return
		}
		s := strings.TrimRight(string(line), "\n")
		log.Printf("In: %s", s)
		io.Copy(c, bytes.NewBufferString(s+"\n"))
	}
}
