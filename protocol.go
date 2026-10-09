package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type MessageType byte

const (
	maxMessageSize                     = 10 * 1024 * 1024
	messageTypeText        MessageType = 1
	messageTypeFile        MessageType = 2
	messageTypePairRequest MessageType = 3
	messageTypePairAccept  MessageType = 4
	messageTypePairReject  MessageType = 5
)

func parseFilePayload(payload []byte) (string, []byte, error) {
	if len(payload) < 2 {
		return "", nil, errors.New("invalid file payload")
	}
	fileNameHeader := payload[0:2]
	fileNameSize := binary.BigEndian.Uint16(fileNameHeader)
	fileNameEnd := 2 + int(fileNameSize)
	if fileNameEnd > len(payload) {
		return "", nil, errors.New("invalid file payload")
	}
	fileNameData := payload[2:fileNameEnd]
	fileName := string(fileNameData)
	fileData := payload[fileNameEnd:]
	saveName := "received_" + fileName
	return saveName, fileData, nil
}

func writeMessage(conn net.Conn, messageType MessageType, payload []byte) error {
	size := uint64(len(payload))

	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, size)

	typeData := []byte{byte(messageType)}

	message := append(typeData, header...)
	message = append(message, payload...)

	reader := bytes.NewReader(message)
	_, err := io.Copy(conn, reader)
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
	if size > maxMessageSize {
		return 0, nil, errors.New("message too large")
	}
	payload := make([]byte, size)

	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return 0, nil, err
	}

	receivedType := MessageType(messageType[0])

	switch receivedType {
	case messageTypeText,
		messageTypeFile,
		messageTypePairRequest,
		messageTypePairAccept,
		messageTypePairReject:
	default:
		return 0, nil, errors.New("unknown message type")
	}

	return receivedType, payload, nil
}
