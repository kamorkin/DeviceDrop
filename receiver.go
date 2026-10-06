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

		for {
			receivedType, payload, err := readMessage(conn)
			if err != nil {
				fmt.Println("Connection closed:", err)
				break
			}

			handleMessage(receivedType, payload)
		}

		conn.Close()
	}
}