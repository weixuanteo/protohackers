package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"unicode"
)

const (
	address = "0.0.0.0:5000"
)

var users = make(map[net.Addr]string)
var connections = make(map[net.Conn]bool)

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

	defer func() {
		conn.Close()
		fmt.Printf("closed connection: %v\n", addr)
		if _, ok := users[addr]; ok {
			delete(connections, conn)
			for c := range connections {
				_, err := c.Write([]byte("* " + users[addr] + " has left the room\n"))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		delete(users, addr)
	}()

	scanner := bufio.NewScanner(conn)

	if _, ok := users[addr]; !ok {
		_, err := conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
		for scanner.Scan() {
			in := scanner.Bytes()
			name := string(in)
			fmt.Printf("received name request: %v\n", string(name))
			if isValid(name) {
				var connectedUsers []string
				for _, u := range users {
					connectedUsers = append(connectedUsers, u)
				}
				_, err := conn.Write([]byte("* The room contains: " + strings.Join(connectedUsers, ", ") + "\n"))
				if err != nil {
					fmt.Println(err)
					break
				}
				users[addr] = name
				connections[conn] = true
				for c := range connections {
					if c != conn {
						_, err := c.Write([]byte("* " + name + " has entered the room\n"))
						if err != nil {
							fmt.Println(err)
							break
						}
					}
				}
			} else {
				_, err := conn.Write([]byte("Invalid name, try again\n"))
				if err != nil {
					fmt.Println(err)
					break
				}
			}
			break
		}

	}

	if _, ok := users[addr]; ok {
		for scanner.Scan() {
			in := scanner.Bytes()
			fmt.Printf("received message: %v\n", string(in))
			for c := range connections {
				if c != conn {
					_, err := c.Write([]byte("[" + users[addr] + "] " + string(in) + "\n"))
					if err != nil {
						fmt.Println(err)
						break
					}
				}
			}
		}
	}
}

func isValid(str string) bool {
	for _, ch := range str {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}
