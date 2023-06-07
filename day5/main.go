package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
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

func handleConn(downstream net.Conn) {
	addr := downstream.RemoteAddr()
	fmt.Printf("accepted connection: %v\n", addr)

	upstream, err := net.Dial("tcp", upstream)
	if err != nil {
		fmt.Println(err)
		return
	}

	once := sync.Once{}
	relay := func(dst io.WriteCloser, src io.ReadCloser) {
		defer once.Do(func() { src.Close(); dst.Close(); fmt.Printf("closed connection: %v", addr) })

		for r := bufio.NewReader(src); ; {
			msg, err := r.ReadString('\n')
			if err != nil {
				return
			}
			msg = strings.TrimRight(msg, "\n")
			fmt.Printf("upstream->client: '%v'\n", msg)
			words := strings.Split(msg, " ")
			for i, word := range words {
				if word == "" {
					continue
				}
				words[i] = rewriteIfBogusAddr(word)
			}
			msg = strings.Join(words, " ")
			msg += "\n"

			_, err = dst.Write([]byte(msg))
			if err != nil {
				return
			}

		}
	}

	go relay(upstream, downstream)
	relay(downstream, upstream)
}

func rewriteIfBogusAddr(str string) string {
	if len(str) < 26 || len(str) > 35 || str[0] != '7' {
		return str
	}

	for _, c := range str {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return str
		}
	}
	return bogus
}
