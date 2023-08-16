package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	address = "0.0.0.0:5001"

	ErrorType         = 10
	NumberType        = 20
	TicketType        = 21
	WantHeartbeatType = 40
	HeartbeatType     = 41
	IAmCameraType     = 80
	IAmDispatcherType = 81
)

// Message Types
type Error struct {
	msg string
}

type Plate struct {
	plate     string
	timestamp uint32
}

type Ticket struct {
	plate      string
	road       uint16
	mile1      uint16
	timestamp1 uint32
	mile2      uint16
	timestamp2 uint32
	speed      uint16
}

type WantHeartbeat struct {
	interval uint32
}

type IAmCamera struct {
	road  uint16
	mile  uint16
	limit uint16
}

type IAmDispatcher struct {
	numroads uint8
	roads    []uint16
}

func main() {
	fmt.Println("Starting server...")

	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("an error occured during tcp listen, err: %v\n", err)
		return
	}

	defer func() {
		ln.Close()
		fmt.Println("Closing tcp listener")
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("an error occured when accepting conn, err: %v\n", err)
			continue
		}

		go handleConn(conn)

	}

}

func handleConn(conn net.Conn) {
	addr := conn.RemoteAddr()
	fmt.Printf("accepted connection from %s\n", addr.String())

	defer func() {
		conn.Close()
		fmt.Printf("closed connection from %v\n", addr)
	}()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		in := scanner.Bytes()
		fmt.Printf("received message: %v, string(message): %s\n", in, string(in))
	}

}
