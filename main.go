package main

import (
	"fmt"
	"os"
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
	if len(os.Args) < 2 {
		fmt.Println("use: go run . send | receive")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "receive":
		runReceiver()

	case "send":
		runSender()

	default:
		fmt.Println("unknown mode")
	}
}