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
	HostMsgLeave
	HostMsgEvent
	HostMsgMetadata
	HostMsgTexture
)

// Host to Server

type HostUpdatePayload struct {
	Client uuid.UUID           `json:"client"`
	Root   components.WireNode `json:"root"`
}

type HostMetadataPayload struct {
	Client uuid.UUID `json:"client"`
	Title  *string   `json:"title,omitempty"`
}

type HostWantSlugPayload struct {
	Slug string `json:"slug"`
}

type HostTexturePayload struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

// Server to Host

type HostReadyPayload struct {
	URL string `json:"url"`
}

type HostHelloPayload struct {
	Client uuid.UUID `json:"client"`
	User   uuid.UUID `json:"user"`
}

type HostLeavePayload struct {
	Client uuid.UUID `json:"client"`
}

type HostEventPayload struct {
	Client uuid.UUID       `json:"client"`
	Event  string          `json:"event"`
	Data   json.RawMessage `json:"data"`
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
