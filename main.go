package main

import (
	"log"
	"net"
)

func handle(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4*1024)
	for {
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			break
		}
		n, err = conn.Write(buf[0:n])
		if err != nil {
			break
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}
