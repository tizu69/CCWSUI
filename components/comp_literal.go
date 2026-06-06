package components

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

type literalPiece struct {
	Text   string
	Color  color.RGBA
	Shadow bool
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

func LiteralOf(text string) Literal {
	return Literal{Pieces: []literalPiece{{
		Text: text, Color: color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}}}
}

func (c Literal) Add(text string) Literal {
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

func (c Literal) WithColor(colr color.Color) Literal {
	r, g, b, a := colr.RGBA()
	c.Pieces[len(c.Pieces)-1].Color = color.RGBA{
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
	c.Pieces[len(c.Pieces)-1].Shadow = shadow
	return c
}

func (c Literal) WithAlignment(alignment Alignment) Literal {
	c.Alignment = alignment
	return c
}

func (c Literal) String() string {
	var s string
	for _, piece := range c.Pieces {
		s += piece.Text
	}
	return s
}

func (c Literal) Measure(ctx MeasureContext, constraint Size) Size {
	lines, longest := c.maybeWrap(ctx, constraint.W)
	return Size{W: longest, H: len(lines) * ctx.GetLineHeight()}
}

func (c Literal) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	lines, longest := c.maybeWrap(ctx, rect.W)
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

// maybeWrap returns the wrapped lines of text and the longest line's width.
func (c Literal) maybeWrap(ctx MeasureContext, width int) (lines []literalLine, longest int) {
	for _, paragraph := range c.paragraphs() {
		l, w := c.maybeWrapParagraph(paragraph, ctx, width)
		lines = append(lines, l...)
		if w > longest {
			longest = w
		}
	}
	return lines, longest
}

func (c Literal) paragraphs() [][]literalPiece {
	paragraphs := [][]literalPiece{{}}
	for _, piece := range c.Pieces {
		parts := strings.Split(piece.Text, "\n")
		for i, part := range parts {
			if i > 0 {
				paragraphs = append(paragraphs, []literalPiece{})
			}
			if part == "" {
				continue
			}
			piece.Text = part
			last := len(paragraphs) - 1
			paragraphs[last] = appendLiteralPiece(paragraphs[last], piece)
		}
	}
	return paragraphs
}

func (c Literal) maybeWrapParagraph(pieces []literalPiece, ctx MeasureContext, width int) (lines []literalLine, longest int) {
	if !c.Wrap || width <= 0 {
		line := literalLine{Pieces: pieces, Width: piecesWidth(ctx, pieces)}
		return []literalLine{line}, line.Width
	}

	words := literalWords(pieces)
	if len(words) == 0 {
		return []literalLine{{}}, 0
	}

	var current literalLine
	flush := func() {
		if len(current.Pieces) == 0 {
			return
		}
		if current.Width > longest {
			longest = current.Width
		}
		lines = append(lines, current)
		current = literalLine{}
	}

	spaceWidth := ctx.GuessTextWidth(" ")
	for _, word := range words {
		wordWidth := ctx.GuessTextWidth(word.Text)
		if len(current.Pieces) == 0 {
			current.Pieces = appendLiteralPiece(current.Pieces, word)
			current.Width = wordWidth
			continue
		}

		space := word
		space.Text = " "
		candidateWidth := current.Width + spaceWidth + wordWidth
		if candidateWidth <= width {
			current.Pieces = appendLiteralPiece(current.Pieces, space)
			current.Pieces = appendLiteralPiece(current.Pieces, word)
			current.Width = candidateWidth
		} else {
			flush()
			current.Pieces = appendLiteralPiece(current.Pieces, word)
			current.Width = wordWidth
		}
	}

	flush()

	return lines, longest
}

func literalWords(pieces []literalPiece) []literalPiece {
	var words []literalPiece
	for _, piece := range pieces {
		for word := range strings.FieldsSeq(piece.Text) {
			piece.Text = word
			words = append(words, piece)
		}
	}
	return words
}

func piecesWidth(ctx MeasureContext, pieces []literalPiece) int {
	var width int
	for _, piece := range pieces {
		width += ctx.GuessTextWidth(piece.Text)
	}
	return width
}

func appendLiteralPiece(pieces []literalPiece, piece literalPiece) []literalPiece {
	if piece.Text == "" {
		return pieces
	}
	last := len(pieces) - 1
	if last >= 0 && pieces[last].Color == piece.Color && pieces[last].Shadow == piece.Shadow {
		pieces[last].Text += piece.Text
		return pieces
	}
	return append(pieces, piece)
}
