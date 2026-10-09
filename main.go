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
		deviceID, err := loadOrCreateDeviceID()
		if err != nil {
			log.Fatal(err)
		}

		deviceName, err := getDeviceName()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(deviceName)
		fmt.Println(deviceID)

	case "pairing-test":
		testID := "123e4567-e89b-42d3-a456-426614174000"

		err := addTrustedDevice(testID)
		if err != nil {
			log.Fatal("Failed to add device:", err)
		}

		loadedDevices, err := loadTrustedDevices()
		if err != nil {
			log.Fatal("Failed to load devices:", err)
		}
		fmt.Printf("Trusted devices: %+v\n", loadedDevices)

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
