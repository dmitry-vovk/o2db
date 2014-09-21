package main

import (
	"config"
	"flag"
	"log"
	. "logger"
	"runtime"
	"server"
)

func main() {
	log.Fatal(
		server.CreateNew(config.Config).Run(),
	)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var configFile = flag.String("config", "o2db.json", "Path to config.json")
	flag.Parse()
	if err := config.Read(*configFile); err != nil {
		log.Fatal(err)
	}
	SetupLogs()
}
