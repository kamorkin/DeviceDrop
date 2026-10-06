package main

import (
	"fmt"
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

	for {
	receivedType, payload, err := readMessage(conn)
		if err != nil {
		fmt.Println("Connection closed:", err)
		break
	}

	handleMessage(receivedType, payload)
	}
}