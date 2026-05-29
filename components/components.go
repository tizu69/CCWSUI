package components

import (
	"html/template"
	"image/color"
	"io"
)

var rootTemplate, _ = template.New("root").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Title}}</title>
	<script src="/static/wasm_exec.js"></script>
	<link rel="stylesheet" href="/static/ccwsui.css">
	<script src="/static/ccwsui_next.js" defer></script>
	<script id="ccwsui-socketurl" type="text/plain">{{.SocketURL}}</script>
</head>
<body>
	<canvas id="ccwsui-root"></canvas>
</body>
</html>
`)

func Root(url, title string, w io.Writer) error {
	return rootTemplate.Execute(w, map[string]any{
		"SocketURL": url,
		"Title":     title,
	})
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
	MeasureContext

	GetDimensions() (w, h int)

	RenderText(x, y int, text string, color color.Color)
	RenderTex(x, y, w, h int, path string, sx, sy int)
	RenderTexPattern(x, y, w, h int, path string, sx, sy, sw, sh int)
	GetTexSize(path string) (x, y int)
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
