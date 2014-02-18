package main

/*
これを参考にした習作
https://github.com/akrennmair/telnet-chat/blob/master/chat.go
*/

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type NetClient struct {
	con net.Conn
	c   chan<- Message
}

type Client interface {
	Id() string
	Channel() chan<- Message
}

type Message string

const MAX_MSG_BUF int = 64

func (nc *NetClient) Channel() chan<- Message {
	return nc.c
}

func (nc *NetClient) Id() string {
	return fmt.Sprint(nc.con.RemoteAddr())
}

func handle(con net.Conn, addclient chan<- Client, deleteclient chan<- Client, msgchan chan Message) {
	c := make(chan Message, MAX_MSG_BUF)
	client := NetClient{con, c}

	io.WriteString(client.con, "> ")
	go func() {
		defer client.con.Close()
		for s := range c {
			if _, err := io.WriteString(client.con, string(s)); err != nil {
				deleteclient <- &client
				return
			}
			io.WriteString(client.con, "> ")
		}
	}()

	addclient <- &client

	buf := bufio.NewReader(client.con)
	for {
		l, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		msgchan <- Message(string(l) + "\r\n")
	}
}

func distribute(addclient <-chan Client, deleteclient <-chan Client, msgchan <-chan Message) {
	clients := make(map[Client]bool)
	for {
		select {
		case client := <-addclient:
			fmt.Printf("new client: %v\n", client.Id())
			clients[client] = true
		case client := <-deleteclient:
			fmt.Printf("delete client: %v\n", client.Id())
			delete(clients, client)
		case msg := <-msgchan:
			for client, _ := range clients {
				select {
				case client.Channel() <- msg:
				default:
				}
			}
		}
	}
}

func main() {
	port := ":8080"
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	addclient := make(chan Client)
	deleteclient := make(chan Client)
	msgchan := make(chan Message)

	go distribute(addclient, deleteclient, msgchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handle(conn, addclient, deleteclient, msgchan)
	}
}
