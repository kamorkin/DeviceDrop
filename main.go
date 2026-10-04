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

func readMessage(conn net.Conn) (MessageType, []byte, error) {
	messageType := make([]byte, 1)

	_, err := io.ReadFull(conn, messageType)
	if err != nil {
		return 0, nil, err
	}

	header := make([]byte, 8)

	_, err = io.ReadFull(conn, header)
	if err != nil {
		return 0, nil, err
	}

	size := binary.BigEndian.Uint64(header)

	payload := make([]byte, size)

	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return 0, nil, err
	}

	receivedType := MessageType(messageType[0])

	return receivedType, payload, nil
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