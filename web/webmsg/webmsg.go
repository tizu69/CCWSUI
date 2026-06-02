package webmsg

import (
	"encoding/json"

	"g.tizu.dev/CCWSUI/components"
)

type Type int

// Server to Client
const (
	TypeUpdate Type = iota + 1
	TypeTexture
)

// Client to Server
const (
	TypeEvent Type = iota + 1
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
	Data string `json:"data"`
}

type Event struct {
	ID    string          `json:"id"`
	Event json.RawMessage `json:"event"`
}
