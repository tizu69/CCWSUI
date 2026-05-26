package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Literal struct {
	Text  string
	Color string
	Wrap  Wrap
}

func LiteralOf(text string) Literal {
	return Literal{Text: text}
}

func (c Literal) WithColor(color string) Literal {
	c.Color = color
	return c
}

func (c Literal) WithWrap(wrap Wrap) Literal {
	c.Wrap = wrap
	return c
}

func (c Literal) Render(ctx RenderContext) Node {
	return Span(Data("ccwsui", "literal"), Styles{
		"color:%s":     c.Color,
		"text-wrap:%s": c.Wrap.CSS(),
	}, Text(c.Text))
}

type Padding struct {
	T, L, B, R int
	Child      Native
}

func Padded(t, l, b, r int, child Native) Padding {
	return Padding{T: t, L: l, B: b, R: r, Child: child}
}

func (c Padding) Render(ctx RenderContext) Node {
	return Div(Data("ccwsui", "padding"), Styles{
		"padding-top:calc(var(--px)*%d)":    c.T,
		"padding-left:calc(var(--px)*%d)":   c.L,
		"padding-bottom:calc(var(--px)*%d)": c.B,
		"padding-right:calc(var(--px)*%d)":  c.R,
	}, c.Child.Render(ctx))
}

type Blank struct {
	W, H int
}

func Blanked(w, h int) Blank {
	return Blank{W: w, H: h}
}

func (c Blank) Render(ctx RenderContext) Node {
	return Div(Data("ccwsui", "blank"), Styles{
		"width:calc(var(--px)*%d)":  c.W,
		"height:calc(var(--px)*%d)": c.H,
	})
}

type Texture struct {
	Tex        string
	PadBorder  bool
	T, L, B, R int
	Child      Native
}

func Textured(tex string, t, l, b, r int, child Native) Texture {
	return Texture{
		Tex: tex, PadBorder: true,
		T: t, L: l, B: b, R: r, Child: child,
	}
}

func (c Texture) WithPadBorder(v bool) Texture {
	c.PadBorder = v
	return c
}

func (c Texture) Render(ctx RenderContext) Node {
	return Div(Data("ccwsui", "texture"), Styles{
		"border-top:calc(var(--px)*%d)solid red":        c.T,
		"border-right:calc(var(--px)*%d)solid red":      c.R,
		"border-bottom:calc(var(--px)*%d)solid red":     c.B,
		"border-left:calc(var(--px)*%d)solid red":       c.L,
		"border-image-source:url('/static/tex/%s.png')": c.Tex,
		"border-image-slice:%d %d %d %d fill":           []any{c.T, c.R, c.B, c.L},
		"border-image-repeat:round":                     nil,

		"--first-child-margin-top:calc(var(--px)*%d)":    -c.T,
		"--first-child-margin-left:calc(var(--px)*%d)":   -c.L,
		"--first-child-margin-bottom:calc(var(--px)*%d)": -c.B,
		"--first-child-margin-right:calc(var(--px)*%d)":  -c.R,
	}, If(!c.PadBorder, Data("ccwsui-first-child-margin", "")),
		c.Child.Render(ctx))
}

type ClickRegion struct {
	Event string
	Child Native
}

func Clickable(event string, child Native) ClickRegion {
	return ClickRegion{Event: event, Child: child}
}

func (c ClickRegion) Render(ctx RenderContext) Node {
	return Div(Data("ccwsui", "clickregion"), c.Child.Render(ctx),
		Attr("role", "button"),
		Attr("hx-post", ctx.EventURL(c.Event)),
		Attr("hx-trigger", "click"))
}

type RenderTime struct {
	Child func() Native
}

func AtRenderTime(child func() Native) RenderTime  { return RenderTime{Child: child} }
func (c RenderTime) Render(ctx RenderContext) Node { return c.Child().Render(ctx) }
