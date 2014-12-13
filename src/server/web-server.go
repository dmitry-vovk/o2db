package server

import (
	"db"
	"github.com/gorilla/websocket"
	. "logger"
	"net/http"
	"server/message"
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
			ErrorLog.Printf("Got message type %d: %s", msg)
			continue
		}
		client := &Client{
			WsConn: conn,
		}
		query, err := message.Parse(msg) // cut out delimiter
		if err != nil {
			ErrorLog.Printf("Parse error: %s", err)
			client.Respond(
				Response{
					Result:   false,
					Response: "message parse error",
				},
			)
		} else {
			DebugLog.Printf("Message: %v", query)
			pkg := &db.Package{
				Container: query,
				Client:    client,
				RespChan:  make(chan Response),
			}
			go func(in *db.Package) {
				go client.Respond(<-in.RespChan)
			}(pkg)
			s.Core.Input <- pkg
		}
	}
}
