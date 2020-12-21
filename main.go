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
		log.Fatalf("Unable to start TCP server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Server Started on :8888")

	for {
		cnn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %s", err.Error())
		}

		go s.newClient(cnn)
	}
}
