package main

import (
	"fmt"
	"log"
	"net"
)

func connectToDevice(address string) (net.Conn, error)  {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func runSender() {
	conn, err := connectToDevice("127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Device connected")

	text := "hello"
	payload := []byte(text)

	err = writeMessage(conn, messageTypeText, payload)
	if err != nil {
		log.Fatal(err)
	}
}
