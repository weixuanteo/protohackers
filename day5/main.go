package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"unicode"
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
			upstreamMsg = strings.TrimRight(upstreamMsg, "\n")
			fmt.Printf("upstream->client: '%v'", upstreamMsg)
			// if upstreamMsg != "" {
			// 	words := strings.Split(upstreamMsg, " ")
			// 	for i := 0; i < len(words); i++ {
			// 		words[i] = rewriteIfBogusAddr(words[i])
			// 	}
			// 	upstreamMsg = strings.Join(words[:], " ")
			// }
			words := strings.Split(upstreamMsg, " ")
			for i, word := range words {
				if word == "" {
					continue
				}
				words[i] = rewriteIfBogusAddr(word)
			}
			upstreamMsg = strings.Join(words[:], " ")
			upstreamMsg += "\n"
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
			clientMsg = strings.TrimRight(clientMsg, "\n")
			fmt.Printf("client->upstream: '%v'", clientMsg)
			// if clientMsg != "" {
			// 	words := strings.Split(clientMsg, " ")
			// 	for i := 0; i < len(words); i++ {
			// 		if i == len(words)-1 {
			// 			words[len(words)-1] = rewriteIfBogusAddr(strings.TrimRight(words[len(words)-1], "\n"))
			// 		}
			// 		words[i] = rewriteIfBogusAddr(words[i])
			// 	}
			// 	// fmt.Println("client->server:", words)
			// 	clientMsg = strings.Join(words[:], " ")
			// }
			words := strings.Split(clientMsg, " ")
			for i, word := range words {
				if word == "" {
					continue
				}
				words[i] = rewriteIfBogusAddr(word)
			}
			clientMsg = strings.Join(words[:], " ")
			clientMsg += "\n"
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

func rewriteIfBogusAddr(str string) string {
	if str[0] != '7' {
		// fmt.Println("not a bogus addr:", str)
		return str
	}

	if len(str) < 26 || len(str) > 35 {
		// fmt.Println("< 26 or > 35:", str)
		return str
	}

	for _, c := range str {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			// fmt.Printf("char: '%c'", c)
			// fmt.Println("not a letter or digit:", str)
			return str
		}
	}
	return bogus
}
