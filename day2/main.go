package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	address = "0.0.0.0:5001"
)

type PriceEntry struct {
	Timestamp int32
	Price     int32
}

func main() {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
	}
	defer ln.Close()

	priceMap := make(map[string][]PriceEntry)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn, priceMap)
	}
}

func handleConn(conn net.Conn, priceMap map[string][]PriceEntry) {
	addr := conn.RemoteAddr()
	fmt.Printf("accepted connection: %v\n", addr)

	defer func() {
		conn.Close()
		fmt.Printf("closed connection: %v\n", addr)
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		in := scanner.Bytes()

		fmt.Printf("received request: %v\n", string(in))

		op := string(in[0])
		num1 := int32(binary.BigEndian.Uint32(in[1:5]))
		num2 := int32(binary.BigEndian.Uint32(in[5:9]))

		if op == "I" {
			priceMap[addr.String()] = append(priceMap[addr.String()], PriceEntry{num1, num2})
		}

		if op == "Q" {
			minTime := &num1
			maxTime := &num2
			var out []byte

			count := int32(0)
			total := int32(0)
			for _, entry := range priceMap[addr.String()] {
				timestamp := entry.Timestamp
				price := entry.Price
				if timestamp >= *minTime && timestamp <= *maxTime {
					count++
					total += price
				}
			}
			out = make([]byte, 4)
			if count > 0 {
				binary.BigEndian.PutUint32(out, uint32(total/count))
			} else {
				binary.BigEndian.PutUint32(out, 0)
			}

			_, err := conn.Write(out)
			if err != nil {
				fmt.Println(err)
				return
			} else {
				fmt.Printf("sent response: %v\n", out)
			}
		}

	}
}
