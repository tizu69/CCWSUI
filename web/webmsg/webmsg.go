package webmsg

import (
	"context"
	"encoding/json"

	"g.tizu.dev/CCWSUI/components"
	"github.com/coder/websocket"
)

type Type int

const (
	TypeUpdate Type = iota + 1
	TypeTexture
	TypeEvent
	TypeMetadata
	TypeRedirect
)

type Envelope struct {
	Type Type            `json:"t"`
	Data json.RawMessage `json:"d"`
}

type Update struct {
	Root components.WireNode `json:"root"`
}

type Texture struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

type Event struct {
	ID    string          `json:"id"`
	Event json.RawMessage `json:"event"`
}

type Metadata struct {
	Title *string `json:"title"`
}

type Redirect struct {
	URL string `json:"url"`
}

func SendMsg(conn *websocket.Conn, typ Type, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	env, err := json.Marshal(Envelope{Type: typ, Data: data})
	if err != nil {
		return err
	}
	return conn.Write(context.Background(), websocket.MessageText, env)
}
