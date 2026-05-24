package components

import (
	"fmt"
	"net/url"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Root(roomid string, title string, root Native) Node {
	wsurl := fmt.Sprintf("/r/%s/service", url.PathEscape(roomid))
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []Node{
			Link(Rel("stylesheet"), Href("/static/ccwsui.css")),
			Script(Src("/static/htmx.min.js")),
			Script(Src("/static/hx-ws.js")),
		},
		Body: []Node{
			Div(Attr("hx-ws:connect", wsurl), root.Render()),
		},
	})
}

// Native components are components that are not composed of other components.
// They are the core building blocks of the UI, and only extensible in Go.
type Native interface {
	Render() Node
}

var nativeRegistry = make(map[string]Native)

func RegisterNative(name string, component Native) {
	if _, ok := nativeRegistry[name]; ok {
		panic(fmt.Sprintf("component %q already registered", name))
	}
	nativeRegistry[name] = component
}

func GetNative(name string) (v Native, ok bool) {
	v, ok = nativeRegistry[name]
	return
}
