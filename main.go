package main

import (
	"fmt"
	"log"
	"net"
)




func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
}

	fmt.Println("Device connected")

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
}

fmt.Println(n)
fmt.Println(string(buffer[:n]))

}