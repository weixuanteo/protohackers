package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
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
			// upstreamMsg = strings.TrimRight(upstreamMsg, "\n")
			// fmt.Printf("upstream->client: '%v'", upstreamMsg)
			// words := strings.Split(upstreamMsg, " ")
			// for i, word := range words {
			// 	if word == "" {
			// 		continue
			// 	}
			// 	words[i] = rewriteIfBogusAddr(word)
			// }
			// upstreamMsg = strings.Join(words[:], " ")
			// upstreamMsg += "\n"
			boguscoinRegex := regexp.MustCompile(`(?m)(^|\s)7[\w\d]{25,34}(\s|$)`)
			upstreamMsg = boguscoinRegex.ReplaceAllString(upstreamMsg, " "+bogus)
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
			// clientMsg = strings.TrimRight(clientMsg, "\n")
			// fmt.Printf("client->upstream: '%v'", clientMsg)
			// words := strings.Split(clientMsg, " ")
			// for i, word := range words {
			// 	if word == "" {
			// 		continue
			// 	}
			// 	words[i] = rewriteIfBogusAddr(word)
			// }
			// clientMsg = strings.Join(words[:], " ")
			// clientMsg += "\n"
			boguscoinRegex := regexp.MustCompile(`(?m)(^|\s)7[\w\d]{25,34}(\s|$)`)
			clientMsg = boguscoinRegex.ReplaceAllString(clientMsg, " "+bogus)
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
		return str
	}

	if len(str) < 26 || len(str) > 35 {
		return str
	}

	for _, c := range str {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return str
		}
	}
	return bogus
}
