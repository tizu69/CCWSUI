package webmsg

import (
	"encoding/json"

	"g.tizu.dev/CCWSUI/components"
)

type Type int

const (
	TypeUpdate Type = iota + 1
	TypeTexture
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
