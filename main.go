package main

import (
	"log"
	"numbers/server"
)

func main() {
	srv := server.New()
	log.Fatal(srv.Start())
}
