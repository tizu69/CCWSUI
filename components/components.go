package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Root(ctx RenderContext, title string, root Native) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []Node{
			Script(Src("/static/htmx.min.js")),
			Script(Src("/static/hx-ws.js")),
			Link(Rel("stylesheet"), Href("/static/ccwsui.css")),
			Script(Src("/static/ccwsui.js"), Defer()),
		},
		Body: []Node{
			Div(Attr("hx-ws:connect", ctx.SocketURL()), root.Render(ctx),
				Attr("hx-swap:inherited", "none"), ID("ccwsui-root")),
		},
	})
}

func HtmxSwap(to string, children ...Node) Node {
	return Div(Attr("hx-swap-oob", "true"), ID(to), Group(children))
}

// Native components are components that are not composed of other components.
// They are the core building blocks of the UI, and only extensible in Go.
type Native interface {
	Render(ctx RenderContext) Node
}

type RenderContext interface {
	SocketURL() string
	EventURL(ev string) string
}

func ifelse[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
