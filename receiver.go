package main

import (
	"fmt"
	"log"
	"net"
)

func runReceiver() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		receivedType, payload, err := readMessage(conn)

		if err != nil {
			fmt.Println("Connection closed:", err)
			break
		}

		handleMessage(conn, receivedType, payload)
	}
}
