package server

import (
	"bufio"
	"config"
	"db"
	"fmt"
	"io"
	. "logger"
	"net"
	"server/message"
	"time"
	. "types"
)

func (s *ServerType) runSocketListener() error {
	socket, err := net.Listen("tcp4", config.Config.ListenTCP)
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
		DebugLog.Print("---------------------------------------------")
		msg, err := bufio.NewReader(client.Conn).ReadBytes(MessageDelimiter)
		if err != nil {
			if err == io.EOF {
				DebugLog.Printf("Client diconnected")
			} else {
				ErrorLog.Printf("Read error: %s", err)
			}
			return
		}
		query, err := message.Parse(msg[:len(msg)-1]) // cut out delimiter
		if err != nil {
			ErrorLog.Printf("Parse error: %s", err)
			ErrorLog.Printf("Message was: %s", msg)
			client.Respond(
				Response{
					Result:   false,
					Response: fmt.Sprintf("%s", err),
				},
			)
		} else {
			DebugLog.Printf("Message: %v", query)
			pkg := &db.Package{
				Container: query,
				Client:    client,
				RespChan:  make(chan Response),
			}
			/* FIXME temporarily disabled for sequential logs
			go func(in *db.Package) {
				go client.Respond(<-in.RespChan)
			}(pkg)
			*/
			s.Core.Input <- pkg
			client.Respond(<-pkg.RespChan)
		}
	}
}
