package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

func connectToDevice(address string) (net.Conn, error) {
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
	defer conn.Close()

	fmt.Println("Device connected")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Command (text/file/exit): ")

		if !scanner.Scan() {
			return
		}

		command := scanner.Text()

		switch command {
		case "text":
			fmt.Print("Text: ")

			if !scanner.Scan() {
				return
			}

			text := scanner.Text()

			err = writeMessage(conn, messageTypeText, []byte(text))
			if err != nil {
				fmt.Println("Failed to send message:", err)
				return
			}

		case "file":
			fmt.Print("File name: ")

			if !scanner.Scan() {
				return
			}

			fileName := scanner.Text()
			fileNameData := []byte(fileName)
			fileNameSize := len(fileNameData)

			fileNameHeader := make([]byte, 2)
			binary.BigEndian.PutUint16(
				fileNameHeader,
				uint16(fileNameSize),
			)

			fileData, err := os.ReadFile(fileName)
			if err != nil {
				fmt.Println("Failed to read file:", err)
				continue
			}

			filePayload := append(fileNameHeader, fileNameData...)
			filePayload = append(filePayload, fileData...)

			err = writeMessage(conn, messageTypeFile, filePayload)
			if err != nil {
				fmt.Println("Failed to send file:", err)
				return
			}

		case "exit":
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}
