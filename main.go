package main

import (
	"numbers/mgmt"
	"log"
)

func main() {
	server := mgmt.NewServer()
	log.Fatal(server.Start())
}
