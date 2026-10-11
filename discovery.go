package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func runDiscoveryListener() {
	publicKey, _, err := loadOrCreateDeviceIdentity()
	if err != nil {
		log.Fatal(err)
	}

	deviceID := getDeviceID(publicKey)

	addr := &net.UDPAddr{
		Port: 9999,
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println("Failed to start discovery listener:", err)
		return
	}
	defer conn.Close()

	for {

		buffer := make([]byte, 1024)

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Failed to read packet", err)
			return
		}

		packetData := buffer[:n]

		parts := strings.Split(string(packetData), "|")
		if len(parts) != 3 {
			continue
		}

		if parts[0] != "DEVICEDROP" {
			continue
		}

		remoteDeviceID := parts[1]

		decodedID, err := hex.DecodeString(remoteDeviceID)
		if err != nil || len(decodedID) != 32 {
			continue
		}

		if remoteDeviceID == deviceID {
			continue
		}

		tcpPort := parts[2]

		tcpAddress := net.JoinHostPort(addr.IP.String(), tcpPort)

		fmt.Println(tcpAddress)
	}

}

func sendDiscoveryAnnouncement() {
	publicKey, _, err := loadOrCreateDeviceIdentity()
	if err != nil {
		log.Fatal(err)
	}

	deviceID := getDeviceID(publicKey)

	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 9999,
	}

	conn, err := net.DialUDP("udp4", nil, broadcastAddr)
	if err != nil {
		fmt.Println("Failed to create UDP connect:", err)
		return
	}
	defer conn.Close()

	announcement := "DEVICEDROP|" + deviceID + "|8080"
	packetData := []byte(announcement)

	for {
		_, err = conn.Write(packetData)
		if err != nil {
			fmt.Println("Failed to send UDP packet:", err)
			return
		}

		time.Sleep(5 * time.Second)
	}
}
