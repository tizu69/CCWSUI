package components

import (
	"encoding/json"
	"image/color"
	"strconv"
	"strings"
)

type Literal struct {
	Text   string
	Color  color.RGBA
	Wrap   Wrap
	Shadow bool
}

func init() {
	RegisterWire("literal", LiteralFromWire)
}

func LiteralOf(text string) Literal {
	return Literal{Text: text, Color: color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}}
}

func (c Literal) WithHexColor(hex string) Literal {
	hex = strings.TrimPrefix(hex, "#")
	switch len(hex) {
	case 3:
		c.Color = color.RGBA{
			R: hexToByte(hex[0:1] + hex[0:1]),
			G: hexToByte(hex[1:2] + hex[1:2]),
			B: hexToByte(hex[2:3] + hex[2:3]),
			A: 0xff,
		}
	case 6:
		c.Color = color.RGBA{
			R: hexToByte(hex[0:2]),
			G: hexToByte(hex[2:4]),
			B: hexToByte(hex[4:6]),
			A: 0xff,
		}
	case 8:
		c.Color = color.RGBA{
			R: hexToByte(hex[0:2]),
			G: hexToByte(hex[2:4]),
			B: hexToByte(hex[4:6]),
			A: hexToByte(hex[6:8]),
		}
	}
	return c
}

func hexToByte(hex string) uint8 {
	b, err := strconv.ParseUint(hex, 16, 8)
	if err != nil {
		return 0
	}
	return uint8(b)
}

func (c Literal) WithColor(colr color.Color) Literal {
	r, g, b, a := colr.RGBA()
	c.Color = color.RGBA{
		R: uint8(r >> 8), G: uint8(g >> 8),
		B: uint8(b >> 8), A: uint8(a >> 8),
	}
	return c
}

func (c Literal) WithWrap(wrap Wrap) Literal {
	c.Wrap = wrap
	return c
}

func (c Literal) WithShadow() Literal {
	c.Shadow = true
	return c
}

func (c Literal) Measure(ctx MeasureContext, constraint Size) Size {
	return Size{W: ctx.GuessTextWidth(c.Text), H: ctx.GetLineHeight()}
}

func (c Literal) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{Rect: Rect{
		X: rect.X, Y: rect.Y,
		W: ctx.GuessTextWidth(c.Text), H: ctx.GetLineHeight(),
	}, Title: "Literal"}
}

func (c Literal) Render(ctx RenderContext, layout LayoutNode) {
	if c.Shadow {
		ctx.RenderText(layout.Rect.X+1, layout.Rect.Y+1, c.Text, color.RGBA{
			R: c.Color.R / 4, G: c.Color.G / 4,
			B: c.Color.B / 4, A: c.Color.A,
		})
	}
	ctx.RenderText(layout.Rect.X, layout.Rect.Y, c.Text, c.Color)
}

func (c Literal) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	return WireNode{Kind: "literal", Props: p}, err
}

func LiteralFromWire(n WireNode) (Native, error) {
	var p Literal
	err := json.Unmarshal(n.Props, &p)
	return p, err
}
