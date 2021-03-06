package server

import (
	"client"
	"db"
	"github.com/gorilla/websocket"
	. "logger"
	"net/http"
	. "types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(_ *http.Request) bool { return true },
}

func (s *ServerType) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ErrorLog.Print(err)
		return
	}
	SystemLog.Printf("WebSocket client connect from %s", conn.RemoteAddr())
	go s.wsHandle(conn)
}

func (s *ServerType) wsHandle(conn *websocket.Conn) {
	defer conn.Close()
	client := &client.Client{
		WsConn: conn,
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			// EOF is caught here
			if err.Error() != "EOF" {
				ErrorLog.Printf("Error reading json message: %s", err)
			}
			break
		}
		if t != websocket.TextMessage {
			ErrorLog.Printf("Got message type %d: %s", t, msg)
			continue
		}
		pkg := &db.Package{
			RawInput: msg,
			Client:   client,
			RespChan: make(chan Response),
		}
		go func(in *db.Package) {
			go client.Respond(<-in.RespChan)
		}(pkg)
		s.Core.Input <- pkg
	}
}
