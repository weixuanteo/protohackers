package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
)

const (
	address = "0.0.0.0:5000"
)

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

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

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		in := scanner.Bytes()

		fmt.Printf("received request: %v\n", string(in))

		var out []byte
		var req Request
		err := json.Unmarshal(in, &req)
		if err != nil || req.Method != "isPrime" || req.Number == nil {
			out = []byte("invalid request")
		} else {
			out, _ = json.Marshal(Response{"isPrime", isPrime(int(*req.Number))})
		}
		out = append(out, byte('\n'))

		_, err = conn.Write(out)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("sent response: %v\n", string(out))
		}
	}
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}

	if n == 2 {
		return true
	}

	if n%2 == 0 {
		return false
	}

	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}
