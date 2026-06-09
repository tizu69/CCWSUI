package components

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"unicode/utf8"
)

type literalPiece struct {
	Text       string
	Color      color.RGBA
	Shadow     bool
	ClickEvent string
}

type Literal struct {
	Pieces    []literalPiece
	Wrap      bool
	Alignment Alignment
}

type literalLine struct {
	Pieces []literalPiece
	Width  int
}

func init() {
	RegisterWire("literal", LiteralFromWire)
}

func MkLiteral(text string) Literal {
	return Literal{Pieces: []literalPiece{{
		Text: text, Color: color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}}}
}

func (c Literal) WithText(text string) Literal {
	c.Pieces = append(c.Pieces, literalPiece{
		Text: text, Color: color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	})
	return c
}

func (c Literal) WithHexColor(hex string) Literal {
	c.Pieces[len(c.Pieces)-1].Color = colorFromHex(hex)
	return c
}

func colorFromHex(hex string) color.RGBA {
	hex = strings.TrimPrefix(hex, "#")
	switch len(hex) {
	case 3:
		return color.RGBA{
			R: hexToByte(hex[0:1] + hex[0:1]),
			G: hexToByte(hex[1:2] + hex[1:2]),
			B: hexToByte(hex[2:3] + hex[2:3]),
			A: 0xff,
		}
	case 6:
		return color.RGBA{
			R: hexToByte(hex[0:2]),
			G: hexToByte(hex[2:4]),
			B: hexToByte(hex[4:6]),
			A: 0xff,
		}
	case 8:
		return color.RGBA{
			R: hexToByte(hex[0:2]),
			G: hexToByte(hex[2:4]),
			B: hexToByte(hex[4:6]),
			A: hexToByte(hex[6:8]),
		}
	}
	return color.RGBA{}
}

func hexToByte(hex string) uint8 {
	b, err := strconv.ParseUint(hex, 16, 8)
	if err != nil {
		return 0
	}
	return uint8(b)
}

func (c Literal) WithColor(r, g, b, a uint8) Literal {
	c.Pieces[len(c.Pieces)-1].Color = color.RGBA{R: r, G: g, B: b, A: a}
	return c
}

func (c Literal) WithWrap() Literal {
	c.Wrap = true
	return c
}

func (c Literal) WithShadow() Literal {
	c.Pieces[len(c.Pieces)-1].Shadow = true
	return c
}

func (c Literal) WithClickEvent(event string) Literal {
	c.Pieces[len(c.Pieces)-1].ClickEvent = event
	return c
}

func (c Literal) WithAlignment(alignment Alignment) Literal {
	c.Alignment = alignment
	return c
}

func (c Literal) String() string {
	var s strings.Builder
	for _, piece := range c.Pieces {
		s.WriteString(piece.Text)
	}
	return s.String()
}

func (c Literal) Measure(ctx MeasureContext, constraint Size) Size {
	lines, longest := c.maybeWrap(ctx, constraint.W)
	return Size{W: longest, H: len(lines) * ctx.GetLineHeight()}
}

func (c Literal) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	lines, longest := c.maybeWrap(ctx, rect.W)

	for y, l := range lines {
		offsetX := rect.X
		switch c.Alignment {
		case AlignmentCenter:
			offsetX += (rect.W - l.Width) / 2
		case AlignmentEnd:
			offsetX += rect.W - l.Width
		}
		for _, p := range l.Pieces {
			ox, w := offsetX, ctx.GuessTextWidth(p.Text)
			offsetX += w
			if p.ClickEvent == "" {
				continue
			}
			rect := Rect{
				X: ox, Y: rect.Y + y*ctx.GetLineHeight(),
				W: w, H: ctx.GetLineHeight(),
			}
			if ctx.GetMouseDown() && rect.Contains(ctx.GetMousePos()) {
				ctx.SendEvent(p.ClickEvent, clickRegionEvent{
					Shift: ctx.GetShiftDown(),
					Ctrl:  ctx.GetCtrlDown(),
					Alt:   ctx.GetAltDown(),
				})
			}
		}
	}

	return LayoutNode{Rect: Rect{
		X: rect.X, Y: rect.Y, W: longest, H: len(lines) * ctx.GetLineHeight(),
	}, Title: fmt.Sprintf("Literal (%q)", c.String())}
}

func (c Literal) Render(ctx RenderContext, layout LayoutNode) {
	lines, _ := c.maybeWrap(ctx, layout.Rect.W)
	for i, line := range lines {
		offsetY := layout.Rect.Y + i*ctx.GetLineHeight()
		offsetX := layout.Rect.X
		switch c.Alignment {
		case AlignmentCenter:
			offsetX += (layout.Rect.W - line.Width) / 2
		case AlignmentEnd:
			offsetX += layout.Rect.W - line.Width
		}

		for _, piece := range line.Pieces {
			if piece.Shadow {
				col := color.RGBA{
					R: piece.Color.R / 4, G: piece.Color.G / 4,
					B: piece.Color.B / 4, A: piece.Color.A,
				}
				if col.R+col.G+col.B < 10 { // (practically) black
					col.R = 128
					col.G = 128
					col.B = 128
				}
				ctx.RenderText(offsetX+1, offsetY+1, piece.Text, col)
			}
			ctx.RenderText(offsetX, offsetY, piece.Text, piece.Color)
			offsetX += ctx.GuessTextWidth(piece.Text)
		}
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

func (c Literal) maybeWrap(ctx MeasureContext, maxWidth int) ([]literalLine, int) {
	if !c.Wrap || maxWidth <= 0 {
		total := 0
		for _, p := range c.Pieces {
			total += ctx.GuessTextWidth(p.Text)
		}
		return []literalLine{{Pieces: c.Pieces, Width: total}}, total
	}

	var lines []literalLine
	var cur literalLine
	curW := 0

	emit := func() {
		if len(cur.Pieces) > 0 {
			cur.Width = curW
			lines = append(lines, cur)
			cur = literalLine{}
			curW = 0
		}
	}

	for _, p := range c.Pieces {
		text := p.Text
		segStart := 0
		segW := 0
		lastSpace := -1
		lastSpaceSegW := 0
		lastSpaceEnd := 0

		for i := 0; i < len(text); {
			r, size := utf8.DecodeRuneInString(text[i:])
			if r == '\n' {
				if i > segStart {
					cur.Pieces = append(cur.Pieces, literalPiece{
						Text: text[segStart:i], Color: p.Color,
						Shadow: p.Shadow, ClickEvent: p.ClickEvent,
					})
					curW += segW
				}
				emit()
				i += size
				segStart = i
				segW = 0
				lastSpace = -1
				continue
			}

			if r == ' ' {
				lastSpace = i
				lastSpaceSegW = segW
				lastSpaceEnd = i + size
			}

			rw := ctx.GuessTextWidth(string(r))
			if curW+segW+rw > maxWidth && segStart < i {
				if lastSpace != -1 {
					// word break
					cur.Pieces = append(cur.Pieces, literalPiece{
						Text: text[segStart:lastSpace], Color: p.Color,
						Shadow: p.Shadow, ClickEvent: p.ClickEvent,
					})
					curW += lastSpaceSegW
					emit()
					i = lastSpaceEnd
					segStart = i
					segW = 0
					lastSpace = -1
				} else {
					// hard break
					cur.Pieces = append(cur.Pieces, literalPiece{
						Text: text[segStart:i], Color: p.Color,
						Shadow: p.Shadow, ClickEvent: p.ClickEvent,
					})
					curW += segW
					emit()
					segStart = i
					segW = 0
					lastSpace = -1
				}
				continue
			}

			segW += rw
			i += size
		}

		if segStart < len(text) {
			cur.Pieces = append(cur.Pieces, literalPiece{
				Text: text[segStart:], Color: p.Color,
				Shadow: p.Shadow, ClickEvent: p.ClickEvent,
			})
			curW += segW
		}
	}
	emit()

	longest := 0
	for _, l := range lines {
		if l.Width > longest {
			longest = l.Width
		}
	}
	return lines, longest
}
