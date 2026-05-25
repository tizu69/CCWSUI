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

func (c Literal) Render() Node {
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

func (c Padding) Render() Node {
	return Div(Data("ccwsui", "padding"), Styles{
		"padding-top:calc(var(--px)*%d)":    c.T,
		"padding-left:calc(var(--px)*%d)":   c.L,
		"padding-bottom:calc(var(--px)*%d)": c.B,
		"padding-right:calc(var(--px)*%d)":  c.R,
	}, c.Child.Render())
}

type Texture struct {
	Background string
	Border     string
	PadBorder  bool
	Repeat     BorderRepeat
	T, L, B, R int
	Child      Native
}

func Textured(bg, border string, t, l, b, r int, child Native) Texture {
	return Texture{
		Background: bg, Border: border, PadBorder: true,
		T: t, L: l, B: b, R: r, Child: child,
	}
}

func (c Texture) WithRepeat(r BorderRepeat) Texture {
	c.Repeat = r
	return c
}

func (c Texture) WithPadBorder(v bool) Texture {
	c.PadBorder = v
	return c
}

func (c Texture) Render() Node {
	return Div(Data("ccwsui", "texture"), Styles{
		"background-image:url('/static/%s.png')": c.Background,
		"background-size:calc(var(--px)*%d)":     2,
		"background-clip:padding-box":            nil,

		"border-top:calc(var(--px)*%d)solid red":    c.T,
		"border-right:calc(var(--px)*%d)solid red":  c.R,
		"border-bottom:calc(var(--px)*%d)solid red": c.B,
		"border-left:calc(var(--px)*%d)solid red":   c.L,
		"border-image-source:url('/static/%s.png')": c.Border,
		"border-image-slice:%d %d %d %d fill":       []any{c.T, c.R, c.B, c.L},
		"border-image-repeat:%s round":              c.Repeat.CSS(),

		"--first-child-margin-top:calc(var(--px)*%d)":    -c.T,
		"--first-child-margin-left:calc(var(--px)*%d)":   -c.L,
		"--first-child-margin-bottom:calc(var(--px)*%d)": -c.B,
		"--first-child-margin-right:calc(var(--px)*%d)":  -c.R,
	}, If(!c.PadBorder, Data("ccwsui-first-child-margin", "")),
		c.Child.Render())
}
