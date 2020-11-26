package main

import (
	"Mongo/internal/pkg/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := server.NewServer()
	s.SetConfig()
	s.SetStorage()
	s.SetHandlers()
	s.Run()
}
