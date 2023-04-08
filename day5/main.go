package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	address  = "0.0.0.0:5001"
	upstream = "chat.protohackers.com:16963"
	bogus    = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
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
	addr := conn.RemoteAddr()
	fmt.Printf("accepted connection: %v\n", addr)

	upstreamConn, err := net.Dial("tcp", upstream)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		conn.Close()
		upstreamConn.Close()
		fmt.Printf("closed connection: %v\n", addr)
	}()

	upstreamReader := bufio.NewReader(upstreamConn)
	clientReader := bufio.NewReader(conn)

	upstreamMsgs := make(chan string)
	clientMsgs := make(chan string)

	go func() {
		for {
			upstreamMsg, err := upstreamReader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			upstreamMsgs <- upstreamMsg
		}
	}()

	go func() {
		for {
			clientMsg, err := clientReader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			clientMsgs <- clientMsg
		}
	}()

	for {
		select {
		case upstreamMsg := <-upstreamMsgs:
			_, err = conn.Write([]byte(upstreamMsg))
			if err != nil {
				fmt.Println(err)
				return
			}

		case clientMsg := <-clientMsgs:
			_, err = upstreamConn.Write([]byte(clientMsg))
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
		}
	}

}
