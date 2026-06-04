package main

import (
	"context"
	"encoding/json"

	"g.tizu.dev/CCWSUI/components"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type HostMsg int

const (
	HostMsgUpdate HostMsg = iota + 1
	HostMsgWantSlug
	HostMsgFreeze
	HostMsgReady
	HostMsgHello
)

// Host to Server

type HostUpdatePayload struct {
	Client uuid.UUID           `json:"client"`
	Root   components.WireNode `json:"root"`
}

type HostWantSlugPayload struct {
	Slug string `json:"slug"`
}

// Server to Host

type HostReadyPayload struct {
	URL string `json:"url"`
}

type HostHelloPayload struct {
	Client uuid.UUID `json:"client"`
}

type HostEnvelope struct {
	Type HostMsg         `json:"t"`
	Data json.RawMessage `json:"d"`
}

func sendHostMsg(conn *websocket.Conn, typ HostMsg, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	env, err := json.Marshal(HostEnvelope{Type: typ, Data: data})
	if err != nil {
		return err
	}
	return conn.Write(context.TODO(), websocket.MessageText, env)
}


