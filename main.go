package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Unable to Start Server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Server Started on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to Connect %s", err.Error())
			continue
		}

		c := s.newUser(conn)
		go c.readInput()
	}
}
