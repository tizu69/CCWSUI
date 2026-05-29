package components

import (
	"encoding/json"
	"image/color"
	"strconv"
	"strings"
)

type Literal struct {
	Text      string
	Color     color.RGBA
	Wrap      bool
	Shadow    bool
	Alignment Alignment
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

func (c Literal) WithWrap(wrap bool) Literal {
	c.Wrap = wrap
	return c
}

func (c Literal) WithShadow(shadow bool) Literal {
	c.Shadow = shadow
	return c
}

func (c Literal) WithAlignment(alignment Alignment) Literal {
	c.Alignment = alignment
	return c
}

func (c Literal) Measure(ctx MeasureContext, constraint Size) Size {
	lines, longest := c.maybeWrap(ctx, constraint.W)
	return Size{W: longest, H: len(lines) * ctx.GetLineHeight()}
}

func (c Literal) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	lines, longest := c.maybeWrap(ctx, rect.W)
	return LayoutNode{Rect: Rect{
		X: rect.X, Y: rect.Y, W: longest, H: len(lines) * ctx.GetLineHeight(),
	}, Title: "Literal"}
}

func (c Literal) Render(ctx RenderContext, layout LayoutNode) {
	lines, _ := c.maybeWrap(ctx, layout.Rect.W)
	for i, line := range lines {
		offsetY := layout.Rect.Y + i*ctx.GetLineHeight()
		offsetX := layout.Rect.X
		switch c.Alignment {
		case AlignmentCenter:
			offsetX += (layout.Rect.W - ctx.GuessTextWidth(line)) / 2
		case AlignmentEnd:
			offsetX += layout.Rect.W - ctx.GuessTextWidth(line)
		}
		if c.Shadow {
			col := color.RGBA{
				R: c.Color.R / 4, G: c.Color.G / 4,
				B: c.Color.B / 4, A: c.Color.A,
			}
			if col.R+col.G+col.B < 10 { // (practically) black
				col.R = 128
				col.G = 128
				col.B = 128
			}
			ctx.RenderText(offsetX+1, offsetY+1, line, col)
		}
		ctx.RenderText(offsetX, offsetY, line, c.Color)
	}
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

// maybeWrap returns the wrapped lines of text and the longest line's width.
func (c Literal) maybeWrap(ctx MeasureContext, width int) (lines []string, longest int) {
	if !c.Wrap || width <= 0 {
		return []string{c.Text}, ctx.GuessTextWidth(c.Text)
	}

	words := strings.Fields(c.Text)
	if len(words) == 0 {
		return []string{""}, 0
	}

	var current string
	flush := func() {
		if current == "" {
			return
		}
		w := ctx.GuessTextWidth(current)
		if w > longest {
			longest = w
		}
		lines = append(lines, current)
		current = ""
	}

	for _, word := range words {
		if current == "" {
			current = word
			continue
		}

		candidate := current + " " + word
		if ctx.GuessTextWidth(candidate) <= width {
			current = candidate
		} else {
			flush()
			current = word
		}
	}

	flush()

	return lines, longest
}
