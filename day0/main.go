package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	address = "0.0.0.0:5000"
)

func main() {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		data := scanner.Bytes()
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Closing connection")
}
