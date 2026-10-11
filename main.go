package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func handleMessage(conn net.Conn, receivedType MessageType, payload []byte) {
	switch receivedType {
	case messageTypeText:
		fmt.Println("Text message:")
		fmt.Println(string(payload))

	case messageTypeFile:
		saveName, fileData, err := parseFilePayload(payload)
		if err != nil {
			fmt.Println("Failed to parse file:", err)
			return
		}
		err = os.WriteFile(saveName, fileData, 0644)
		if err != nil {
			fmt.Println("Failed to save file:", err)
			return
		}
		fmt.Println("File saved:")

	case messageTypePairRequest:
		fmt.Println("Received pairing request")

		err := writeMessage(conn, messageTypePairReject, nil)
		if err != nil {
			fmt.Println("Failed to reject pairing:", err)
		}

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

	case "identity-test":
		publicKey, _, err := loadOrCreateDeviceIdentity()
		if err != nil {
			log.Fatal(err)
		}

		deviceID := getDeviceID(publicKey)

		deviceName, err := getDeviceName()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(deviceName)
		fmt.Println(deviceID)

	case "pair":
		if len(os.Args) < 3 {
			fmt.Println("Usage: DeviceDrop pair <address>")
			return
		}

		if err := sendPairRequest(os.Args[2]); err != nil {
			fmt.Println("Pairing error:", err)
		}
	default:
		fmt.Println("unknown mode")
	}
}
