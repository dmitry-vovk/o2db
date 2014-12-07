package server

import (
	"bufio"
	"bytes"
	"config"
	"db"
	"encoding/json"
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
			handle(
				Response{
					Result:   false,
					Response: fmt.Sprintf("%s", err),
				},
				client,
			)
		} else {
			DebugLog.Printf("Message: %v", query)
			pkg := &db.Package{
				Container: query,
				Client:    client,
				RespChan:  make(chan Response),
			}
			go func(in *db.Package) {
				go handle(<-in.RespChan, client)
			}(pkg)
			s.Core.Input <- pkg
		}
	}
}

// Send response to client
func handle(resp Response, client *Client) {
	out, err := json.Marshal(resp)
	if err != nil {
		ErrorLog.Printf("Error encoding response: %s", err)
		return
	}
	DebugLog.Printf("Response: %s", out)
	_, err = io.Copy(client.Conn, bytes.NewBuffer(append(out, MessageDelimiter)))
	if err != nil {
		ErrorLog.Printf("Error sending response to client: %s", err)
		return
	}
}
