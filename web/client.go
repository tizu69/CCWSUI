//go:build wasm && js

package main

import (
	"syscall/js"

	"g.tizu.dev/CCWSUI/components"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type jsapi struct {
	wrap             js.Value
	Root             components.Native
	wireRoot         *components.WireNode
	DevTools         bool
	requiredTextures []string
	overlays         []func()
	textWidthCache   map[string]int
	textWidthScale   int
	context          map[string]any
	texBordersCache  map[string][4]int
	ws               *websocket.Conn
	useruuid         uuid.UUID
}

func NewJSAPI() *jsapi {
	j := &jsapi{
		wrap:            js.Global().Get("ccwsui"),
		context:         make(map[string]any),
		texBordersCache: make(map[string][4]int),
	}

	j.wrap.Set("totalRerender", js.FuncOf(j.totalRerender))
	j.wrap.Set("openDevtools", js.FuncOf(j.openDevtools))

	wait := make(chan any)
	j.wrap.Call("prepare").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		close(wait)
		return nil
	}))
	<-wait

	var err error
	if j.useruuid, err = uuid.Parse(j.wrap.Get("user").String()); err != nil {
		panic(err)
	}

	return j
}

func (j *jsapi) SocketURL() string {
	return j.wrap.Get("socketURL").String() + "?user=" + j.useruuid.String()
}

func (j *jsapi) ValidateSocketURL() string {
	return j.wrap.Get("validateSocketURL").String()
}

func (j *jsapi) totalRerender(this js.Value, args []js.Value) any {
	j.TotalRerender(args[0].String())
	return nil
}

func (j *jsapi) openDevtools(this js.Value, args []js.Value) any {
	j.DevTools = !j.DevTools
	j.TotalRerender("DevTools toggled")
	if !j.DevTools {
		j.wrap.Call("closeDevtools")
	}
	return nil
}
