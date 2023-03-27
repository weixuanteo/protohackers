package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
)

const (
	address = "0.0.0.0:5000"
)

type Request struct {
	Method string `json:"method"`
	Number int    `json:"number"`
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
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error: ", err.Error())
			}
			return
		}

		req := &Request{}
		err = json.Unmarshal(buf[:n], req)
		if err != nil {
			_, err := conn.Write(buf[:n])
			if err != nil {
				fmt.Println("Error: ", err.Error())
				return
			}
			fmt.Println("Connection closed because it does not match the Request format")
			return
		}

		if req.Method == "" || req.Number == 0 || req.Method != "isPrime" {
			_, err := conn.Write(buf[:n])
			if err != nil {
				fmt.Println("Error: ", err.Error())
				return
			}
			fmt.Println("Connection closed because method is not isPrime or method/number is empty")
			return
		}

		resp := &Response{
			Method: req.Method,
			Prime:  isPrime(req.Number),
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = conn.Write(respBytes)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
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
