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
		make(chan Message, 1),
	}

	go room.dispatch()
	room.AddClient(client)
	room.Cast(Message("ping"))
	select {
	case msg := <-client.c:
		if string(msg) != "ping" {
			t.Errorf("Got wrong message: %v", msg)
		}
	default:
		t.Errorf("Not received.")
	}
}

func TestClose(t *testing.T) {
	room := NewSimpleRoom()

	client := &TestClient{
		"test",
		make(chan Message, 1),
	}

	go room.dispatch()
	room.AddClient(client)
	room.Close()
	actual := room.closed
	expected := true
	if expected != actual {
		t.Errorf("room isn't closed: %v %v", expected, actual)
	}

	room.AddClient(client)
	room.Cast(Message("ping"))
	select {
	case <-client.c:
		t.Errorf("Not closed.")
	default:
	}
}
