package main

import (
	"testing"
)

type TestClient struct {
	id string
	c  chan Message
}

func (client *TestClient) Id() string {
	return client.id
}

func (client *TestClient) Channel() chan<- Message {
	return client.c
}

func (client *TestClient) Join(Room) {
	panic("Not implemented")
}

func TestSimpleRoomType(t *testing.T) {
	room := NewSimpleRoom()
	// Type check
	_ = Room(room)
}

func TestDispatch(t *testing.T) {
	room := NewSimpleRoom()

	client := &TestClient{
		"test",
		make(chan Message),
	}

	go room.dispatch()
	room.AddClient(client)
	room.Cast(Message("ping"))
	select {
	case msg := <-client.c:
		if string(msg) != "ping" {
			t.Errorf("Got wrong message: %v", msg)
		}
	}
}
