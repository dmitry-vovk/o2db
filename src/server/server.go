package server

import (
	"client"
	"config"
	"db"
	. "logger"
	"net/http"
)

type ServerType struct {
	Config  *config.ConfigType
	Clients []*client.Client
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
func (s *ServerType) Run() error {
	SystemLog.Print("Starting core processor")
	go s.Core.Processor()
	SystemLog.Printf("Starting socket listener on %s", config.Config.ListenTCP)
	err := s.runSocketListener()
	if err != nil {
		return err
	}
	http.HandleFunc("/", s.wsHandler)
	SystemLog.Printf("Starting HTTP listener on %s", config.Config.ListenHTTP)
	return http.ListenAndServe(config.Config.ListenHTTP, nil)
}
