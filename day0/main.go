package main

import (
	"fmt"
	"io"
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

	// scanner := bufio.NewScanner(conn)
	// for scanner.Scan() {
	// 	data := scanner.Bytes()
	// 	_, err := conn.Write(data)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error: ", err.Error())
			}
			return
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
	}
}
