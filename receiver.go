package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func runReceiver() {
	config, err := createTLSConfig()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := tls.Listen("tcp", ":8080", config)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		fmt.Println("Expected TLS connection")
		return
	}

	if err := tlsConn.Handshake(); err != nil {
		fmt.Println("TLS handshake failed:", err)
		return
	}

	fmt.Println("TLS handshake successful; authentication not implemented")
}

// for {
// 	receivedType, payload, err := readMessage(conn)

// 	if err != nil {
// 		fmt.Println("Connection closed:", err)
// 		break
// 	}

// 	handleMessage(conn, receivedType, payload)
// }
