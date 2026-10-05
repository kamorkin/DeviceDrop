package main

import (
	"encoding/binary"
	"io"
	"net"
)

type MessageType byte

const (
	messageTypeText MessageType = 1
	messageTypeFile MessageType = 2
)

func writeMessage(conn net.Conn, messageType MessageType, payload []byte) error {
	size := uint64(len(payload))

	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, size)

	typeData := []byte{byte(messageType)}

	message := append(typeData, header...)
	message = append(message, payload...)

	_, err := conn.Write(message)
	return err
}

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