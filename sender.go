package main

import "net"

func connectToDevice(address string) (net.Conn, error)  {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}