//go:build wasm && js

package main

import (
	"fmt"
	"reflect"
)

func (j *jsapi) UseContext(id string, v any) {
	pv := reflect.ValueOf(v)
	inner := pv.Elem()

	known, ok := j.context[id]
	if !ok {
		j.context[id] = inner.Interface()
		return
	}

	kv := reflect.ValueOf(known)
	if inner.Elem().Type() != kv.Elem().Type() {
		panic(fmt.Sprintf("context %q already exists with different type %T", id, known))
	}
	inner.Set(kv)
}
