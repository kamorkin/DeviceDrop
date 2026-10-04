package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type MessageType byte

const (
messageTypeText MessageType = 1
messageTypeFile MessageType = 2
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Device connected")

	header := make([]byte, 8)

	messageType := make([]byte, 1)

	_, err = io.ReadFull(conn, messageType)
	if err != nil {
		log.Fatal(err)
	}

	receivedType := MessageType(messageType[0])

	_, err = io.ReadFull(conn, header)
	if err != nil {
		log.Fatal(err)
	}

	size := binary.BigEndian.Uint64(header)

	payload := make([]byte, size)

	_, err = io.ReadFull(conn, payload)
	if err != nil {
		log.Fatal(err)
	}

	switch receivedType {
	case messageTypeText:
		fmt.Println("text message")
		fmt.Println(string(payload))

	case messageTypeFile:
		fmt.Println("File message")

	default:
		fmt.Println("unknown message type")
	}

	fmt.Println(messageType[0])
	fmt.Println(header)
	fmt.Println(size)
}
