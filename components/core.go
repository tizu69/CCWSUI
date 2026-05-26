package components

import (
	"strings"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Literal struct {
	Text   string
	Color  string
	Wrap   Wrap
	Select TextSelect
	Shadow bool
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

func (c Literal) WithSelect(sel TextSelect) Literal {
	c.Select = sel
	return c
}

func (c Literal) WithShadow() Literal {
	c.Shadow = true
	return c
}

func (c Literal) Render(ctx RenderContext) Node {
	return Span(Data("ccwsui", "literal"), Styles{
		"color:%s":       c.Color,
		"text-wrap:%s":   c.Wrap.CSS(),
		"user-select:%s": c.Select.CSS(),
		"text-shadow:%s": ifelse(c.Shadow, "color-mix(in srgb, currentColor 25%, #000000) var(--px) var(--px)", "none"),
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
		"border-top:calc(var(--px)*%d)solid #202020":    c.T,
		"border-right:calc(var(--px)*%d)solid #202020":  c.R,
		"border-bottom:calc(var(--px)*%d)solid #202020": c.B,
		"border-left:calc(var(--px)*%d)solid #202020":   c.L,
		"border-image-source:url('/static/tex/%s.png')": c.Tex,
		"border-image-slice:%d %d %d %d fill":           []any{c.T, c.R, c.B, c.L},
		"border-image-repeat:round":                     nil,
	}, Div(If(!c.PadBorder, Styles{
		"margin-top:calc(var(--px)*%d)":    -c.T,
		"margin-left:calc(var(--px)*%d)":   -c.L,
		"margin-bottom:calc(var(--px)*%d)": -c.B,
		"margin-right:calc(var(--px)*%d)":  -c.R,
	}), c.Child.Render(ctx)))
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

type ItemTexture struct {
	Item string
}

func ItemTextured(item string) ItemTexture { return ItemTexture{Item: item} }
func (c ItemTexture) Render(ctx RenderContext) Node {
	id := strings.ReplaceAll(c.Item, ":", "__")
	return Div(Data("ccwsui", "itemtexture"), Styles{
		"width:calc(var(--px)*16)":                    nil,
		"height:calc(var(--px)*16)":                   nil,
		"background-image:url('/static/item/%s.png')": id,
		"background-size:100%":                        nil,
		"background-repeat:no-repeat":                 nil,
	})
}
