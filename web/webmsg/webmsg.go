package webmsg

import (
	"encoding/json"

	"g.tizu.dev/CCWSUI/components"
)

type Type int

const (
	TypeUpdate Type = iota + 1
)

type Envelope struct {
	Type Type            `json:"t"`
	Data json.RawMessage `json:"d"`
}

type Update struct {
	Root components.WireNode `json:"root"`
}
