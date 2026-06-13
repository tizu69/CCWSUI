package components

import (
	"encoding/json"
	"fmt"
)

type WireNode struct {
	ID       string          `json:"id,omitempty"`
	Kind     string          `json:"kind"`
	Props    json.RawMessage `json:"props,omitempty"`
	Children []WireNode      `json:"children,omitempty"`
}

type WireCodec interface {
	ToWire() (WireNode, error)
}

type WireDecoder func(n WireNode) (Native, error)

var wireRegistry = map[string]WireDecoder{}

func RegisterWire(kind string, dec WireDecoder) {
	if _, ok := wireRegistry[kind]; ok {
		panic(fmt.Sprintf("kind %s already registered", kind))
	}
	wireRegistry[kind] = dec
}

func FromWire(n WireNode) (Native, error) {
	dec, ok := wireRegistry[n.Kind]
	if !ok {
		return nil, fmt.Errorf("unknown kind: %s", n.Kind)
	}
	return dec(n)
}
