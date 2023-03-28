package main

import (
	"encoding/binary"
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
	addr := conn.RemoteAddr()
	fmt.Printf("accepted connection: %v\n", addr)

	defer func() {
		conn.Close()
		fmt.Printf("closed connection: %v\n", addr)
	}()

	priceMap := make(map[int32]int32)
	buf := make([]byte, 9)

	for {
		_, err := io.ReadFull(conn, buf)
		if err != nil || err == io.EOF {
			break
		}

		num1 := int32(binary.BigEndian.Uint32(buf[1:5]))
		num2 := int32(binary.BigEndian.Uint32(buf[5:]))

		if buf[0] == 'I' {
			priceMap[num1] = num2
			fmt.Printf("Insert operation: %v %v for address - %v\n", num1, num2, addr)
		}
		if buf[0] == 'Q' {
			n := 0
			total := 0
			mean := 0
			for timestamp, price := range priceMap {
				if timestamp >= num1 && timestamp <= num2 {
					n++
					total += int(price)
				}
			}

			if n > 0 {
				mean = total / n
			}

			out := make([]byte, 4)
			binary.BigEndian.PutUint32(out, uint32(mean))

			_, err := conn.Write(out)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Printf("Query operation: %v %v for address - %v\n", num1, num2, addr)
		}
	}

}
