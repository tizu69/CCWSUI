package main

import (
	"g.tizu.dev/CCWSUI/components"
	"github.com/coder/websocket"
)

type Room struct {
	Title string
	Root  components.Native

	listeners []*websocket.Conn
}

func NewRoom(title string, root components.Native) *Room {
	return &Room{
		Title: title,
		Root:  root,

		listeners: make([]*websocket.Conn, 0),
	}
}

func (r *Room) Add(conn *websocket.Conn) {
	r.listeners = append(r.listeners, conn)
}

func (r *Room) Remove(conn *websocket.Conn) {
	for i, c := range r.listeners {
		if c == conn {
			r.listeners = append(r.listeners[:i], r.listeners[i+1:]...)
			break
		}
	}
}

func getCoreRooms() map[string]*Room {
	return map[string]*Room{
		"home": getCoreRoomHome(),
	}
}

func getCoreRoomHome() *Room {
	return NewRoom("CCWSUI!",
		components.Padded(8, 8, 8, 8,
			components.LiteralOf("Hey there!")))
}
