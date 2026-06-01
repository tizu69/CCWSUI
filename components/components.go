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

type ClippedRenderer interface {
	RenderClipped(ctx RenderContext, layout LayoutNode, clip Rect)
}

type MeasureContext interface {
	GuessTextWidth(text string) int
	GetLineHeight() int
	GetTexSize(path string) (x, y int)
	GetTexBorders(path string) (t, l, b, r int)
}

type interactivityContext interface {
	GetDimensions() (w, h int)
	GetMousePos() (x, y int)
	GetMouseScroll() (dx, dy int)

	// UseContext assigns the context value for the given id to v. If there is
	// no context value for the given id, v will be kept unchanged, so you
	// should initialize it to a sensible default. v must be a pointer to a
	// pointer to a value of the type you want to retrieve. You may
	// update it in-place and it will be saved for the next redraw cycle.
	//
	//	cctx := &scrollCtx{X: 0, Y: 0}
	//	ctx.UseContext("scrollcontainer", &cctx)
	//	cctx.X += 10 // persisted!
	UseContext(id string, v any)
}

type LayoutContext interface {
	MeasureContext
	interactivityContext
}

type RenderContext interface {
	MeasureContext
	interactivityContext

	// Until [PopScissor] is called, all rendering will be clipped to the
	// given rectangle. Multiple calls to Scissor will intersect the scissor
	// rectangle, and ClearScissor will clear them last-in-first-out.
	//
	// While scissoring, [GetMousePos] will return (-1, -1), as if the mouse
	// were outside the window, if the mouse is outside the scissor rectangle.
	Scissor(x, y, w, h int)
	PopScissor()

	RenderText(x, y int, text string, color color.Color)
	RenderTex(x, y, w, h int, path string, sx, sy int)
	RenderTexMixel(x, y, w, h int, path string, sx, sy, sw, sh int, nn bool)
	RenderTexPattern(x, y, w, h int, path string, sx, sy, sw, sh int)
	GetTexSize(path string) (x, y int)
	GetTexBorders(path string) (t, l, b, r int)

	// When called, the given texture will be loaded and cached after this
	// render pass is done. This is a no-op if the texture is already loaded.
	// It should be called prior to any calls to Texture-related functions.
	RequireTexture(path string)

	RenderOverlay(func())
}

type Size struct {
	W, H int
}

type Rect struct {
	X, Y, W, H int
}

func (rc Rect) Contains(x, y int) bool {
	return x >= rc.X && x < rc.X+rc.W && y >= rc.Y && y < rc.Y+rc.H
}

func (rc Rect) Intersects(other Rect) bool {
	return rc.X < other.X+other.W &&
		rc.X+rc.W > other.X &&
		rc.Y < other.Y+other.H &&
		rc.Y+rc.H > other.Y
}

type LayoutNode struct {
	Rect     Rect
	Children []LayoutNode
	Title    string
}

func ifelse[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
