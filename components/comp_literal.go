package components

import (
	"encoding/json"
)

type Literal struct {
	Text   string
	Color  string
	Wrap   Wrap
	Select TextSelect
	Shadow bool
}

func init() {
	RegisterWire("literal", LiteralFromWire)
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

func (c Literal) Measure(ctx MeasureContext, constraint Size) Size {
	return Size{W: ctx.GuessTextWidth(c.Text), H: ctx.GetLineHeight()}
}

func (c Literal) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{Rect: rect, Title: "Literal"}
}

func (c Literal) Render(ctx RenderContext, layout LayoutNode) {
	ctx.RenderText(layout.Rect.X, layout.Rect.Y, c.Text)
}

func (c Literal) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "literal", Props: p}, nil
}

func LiteralFromWire(n WireNode) (Native, error) {
	var p Literal
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	return p, nil
}
