package main

import (
	"fmt"
	"log"
	"net"
)

func handleMessage(receivedType MessageType, payload []byte)  {
	switch receivedType {
	case messageTypeText:
		fmt.Println("Text message:")
		fmt.Println(string(payload))

	case messageTypeFile:
		fmt.Println("File message")

	default:
		fmt.Println("Unknown message type")
	}
}

func main() {
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

	fmt.Println("Device connected")

	receivedType, payload, err := readMessage(conn)
		if err != nil {
		log.Fatal(err)
	}

	handleMessage(receivedType, payload)
}