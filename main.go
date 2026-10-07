package main

import (
	"fmt"
	"os"
)

func handleMessage(receivedType MessageType, payload []byte) {
	switch receivedType {
	case messageTypeText:
		fmt.Println("Text message:")
		fmt.Println(string(payload))

	case messageTypeFile:
		saveName, fileData, err := parseFilePayload(payload)
		if err != nil {
			fmt.Println("Failed to parse file:" /*, err*/)
			return
		}
		err = os.WriteFile(saveName, fileData, 0644)
		if err != nil {
			fmt.Println("Failed to save file:", err)
			return
		}
		fmt.Println("File saved:")

	default:
		fmt.Println("Unknown message type")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("use: go run . send | receive | discovery-listen | discovery-send")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "receive":
		runReceiver()

	case "send":
		runSender()

	case "discovery-listen":
		runDiscoveryListener()

	case "discovery-send":
		sendDiscoveryAnnouncement()

	default:
		fmt.Println("unknown mode")
	}
}
