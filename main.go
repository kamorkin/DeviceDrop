package main

import (
	"fmt"
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
	runReceive()
}