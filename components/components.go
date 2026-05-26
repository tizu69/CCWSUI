package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func RootV2(url, title string, root Native) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []Node{
			Script(Src("/static/wasm_exec.js")),
			Link(Rel("stylesheet"), Href("/static/ccwsui.css")),
			Script(Src("/static/ccwsui_next.js"), Defer()),
			Script(ID("ccwsui-socketurl"), Type("text/plain"), Text(url)),
		},
		Body: []Node{
			Canvas(ID("ccwsui-root")),
		},
	})
}

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
			Div(Attr("hx-ws:connect", "" /* ctx.SocketURL() */), /* root.Render(ctx) */
				Attr("hx-swap:inherited", "none"), ID("ccwsui-root")),
		},
	})
}

func HtmxSwap(to string, children ...Node) Node {
	return Div(Attr("hx-swap-oob", "true"), ID(to), Group(children))
}

// Native components are components that are not composed of other components.
// They are the core building blocks of the UI, and only extensible in Go.
// Don't forget to register them with [RegisterWire]!
type Native interface {
	// Measure returns its desired size, given a constraint (maximum allowed).
	Measure(ctx MeasureContext, constraint Size) Size
	// Layout returns a layout node, given a context and a rectangle to layout
	// inside of (which is partially the result of Measure).
	Layout(ctx LayoutContext, rect Rect) LayoutNode
	Render(ctx RenderContext, layout LayoutNode)
	ToWire() (WireNode, error)
}

type MeasureContext interface {
	GuessTextWidth(text string) int
	GetLineHeight() int
}

type LayoutContext interface {
	MeasureContext
}

type RenderContext interface {
	RenderText(x, y int, text string)
	RenderTex(x, y, w, h int, path string, sx, sy int)
}

type Size struct {
	W, H int
}

type Rect struct {
	X, Y, W, H int
}

type LayoutNode struct {
	Rect     Rect
	Children []LayoutNode
	Title    string
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ifelse[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
