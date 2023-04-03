package main

import (
	"fmt"
	"net"
)

func main() {
	address, _ := net.ResolveUDPAddr("udp", ":5000")
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	store := make(map[string]string)
	store["version"] = "1.0.0"

	for {
		buf := make([]byte, 1000)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Size of message: %v", n)
		message := string(buf[:n])
		fmt.Printf("Read message from %v: %v", addr, message)

		isInsert := false

		for pos, char := range message {
			if char == '=' {
				isInsert = true
				key := message[:pos]
				value := message[pos+1:]
				fmt.Printf("storing key '%v' with value '%v'", key, value)
				if key != "version" {
					store[key] = value
				}
				break
			}
		}

		if !isInsert {
			_, err := conn.WriteToUDP([]byte(store[message]), addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Sent message to %v: %v", addr, store[message])
		} else {
			_, err := conn.WriteToUDP([]byte("OK"), addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Sent message to %v: %v", addr, "OK")
		}
	}
}
