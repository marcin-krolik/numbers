package main

import (
	"log"
	"numbers/mgmt"
)

func main() {
	server := mgmt.NewServer()
	log.Fatal(server.Start())
}
