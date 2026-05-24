package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Literal struct {
	Text string
}

func LiteralOf(text string) Literal {
	return Literal{Text: text}
}

func (c Literal) Render() Node {
	return Span(Text(c.Text))
}

type Padding struct {
	T, L, B, R int
	Child      Native
}

func Padded(t, l, b, r int, child Native) Padding {
	return Padding{T: t, L: l, B: b, R: r, Child: child}
}

func (c Padding) Render() Node {
	return Div(Styles{
		"padding-top:calc(var(--px)*%d)":    c.T,
		"padding-left:calc(var(--px)*%d)":   c.L,
		"padding-bottom:calc(var(--px)*%d)": c.B,
		"padding-right:calc(var(--px)*%d)":  c.R,
	}, c.Child.Render())
}

type Margin struct {
	T, L, B, R int
	Child      Native
}

func Margined(t, l, b, r int, child Native) Margin {
	return Margin{T: t, L: l, B: b, R: r, Child: child}
}

func (c Margin) Render() Node {
	return Div(Styles{
		"margin-top:calc(var(--px)*%d)":    c.T,
		"margin-left:calc(var(--px)*%d)":   c.L,
		"margin-bottom:calc(var(--px)*%d)": c.B,
		"margin-right:calc(var(--px)*%d)":  c.R,
	}, c.Child.Render())
}
