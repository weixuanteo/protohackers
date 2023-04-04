package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	address, _ := net.ResolveUDPAddr("udp", ":5000")
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	store := make(map[string]string)
	store["version"] = "1.0.0"

	for {
		buf := make([]byte, 1000)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading UDP message:", err)
			return
		}

		message := string(buf[:n])
		fmt.Printf("Received message from %v: %v\n", addr, message)
		key, value, insert := strings.Cut(message, "=")
		if insert && key != "version" {
			store[key] = value
		} else {
			key := strings.TrimSpace(message)
			value, ok := store[key]
			if ok {
				data := fmt.Sprintf("%v=%v", key, value)
				_, err := conn.WriteToUDP([]byte(data), addr)
				if err != nil {
					fmt.Println("Error writing UDP message:", err)
					return
				}
			} else {
				_, err := conn.WriteToUDP([]byte("NOT FOUND"), addr)
				if err != nil {
					fmt.Println("Error writing UDP message:", err)
					return
				}
			}
		}
	}
}
