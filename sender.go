package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
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
	/*
	text := "hello"
	payload := []byte(text)

	err = writeMessage(conn, messageTypeText, payload)
	if err != nil {
		log.Fatal(err)
	}
	*/


	fileName := "test.txt"
	fileNameData := []byte(fileName)
	fileNameSize := len(fileNameData)
	fileNameHeader := make([]byte, 2)
	binary.BigEndian.PutUint16(fileNameHeader, uint16(fileNameSize))

	fileData, err := os.ReadFile(fileName)
		if err != nil {
		log.Fatal(err)
	}

	filePayload := append(fileNameHeader, fileNameData...)
	filePayload = append(filePayload, fileData...)

	err = writeMessage(conn, messageTypeFile, filePayload)
		if err != nil {
		log.Fatal(err)

	}
}