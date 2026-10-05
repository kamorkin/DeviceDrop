package main

import (
	"log"
	"net"
)

func runReceiver()  {
		listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	receivedType, payload, err := readMessage(conn)
		if err != nil {
		log.Fatal(err)
	}

	handleMessage(receivedType, payload)
}