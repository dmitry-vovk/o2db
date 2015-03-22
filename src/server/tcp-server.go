package server

import (
	"bufio"
	"db"
	"io"
	. "logger"
	"net"
	"time"
	. "types"
)

func (s *ServerType) runSocketListener() error {
	socket, err := net.Listen("tcp4", s.Config.ListenTCP)
	if err != nil {
		return err
	}
	go s.socketListener(socket)
	return nil
}

func (s *ServerType) socketListener(socket net.Listener) {
	for {
		conn, err := socket.Accept()
		if err != nil {
			ErrorLog.Printf("Accept error: %s", err)
			continue
		}
		DebugLog.Printf("TCP client connected")
		c := &Client{
			Conn: conn,
		}
		go s.handler(c)
	}
}

// Handle single client connection
func (s *ServerType) handler(client *Client) {
	start := time.Now()
	defer func() {
		DebugLog.Printf("Completed in %.4f seconds", time.Now().Sub(start).Seconds())
	}()
	defer client.Conn.Close()
	for {
		msg, err := bufio.NewReader(client.Conn).ReadBytes(MessageDelimiter)
		if err != nil {
			if err == io.EOF {
				DebugLog.Printf("Client diconnected")
			} else {
				ErrorLog.Printf("Read error: %s", err)
			}
			return
		}
		pkg := &db.Package{
			RawInput: msg[:len(msg)-1],
			Client:   client,
			RespChan: make(chan Response),
		}
		s.Core.Input <- pkg
		client.Respond(<-pkg.RespChan)
	}
}
