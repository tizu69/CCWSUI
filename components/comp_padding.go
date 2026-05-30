package components

import (
	"encoding/json"
	"fmt"
)

type Padding struct {
	T, L, B, R int
	Child      Native `json:"-"`
}

func init() {
	RegisterWire("padding", PaddingFromWire)
}

func Padded(t, l, b, r int, child Native) Padding {
	return Padding{T: t, L: l, B: b, R: r, Child: child}
}

func (c Padding) Measure(ctx MeasureContext, constraint Size) Size {
	childConstraint := Size{
		W: max(0, constraint.W-c.L-c.R),
		H: max(0, constraint.H-c.T-c.B),
	}
	child := c.Child.Measure(ctx, childConstraint)
	return Size{W: child.W + c.L + c.R, H: child.H + c.T + c.B}
}

func (c Padding) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	childRect := Rect{
		X: rect.X + c.L,
		Y: rect.Y + c.T,
		W: max(0, rect.W-c.L-c.R),
		H: max(0, rect.H-c.T-c.B),
	}
	return LayoutNode{
		Rect:     rect,
		Children: []LayoutNode{c.Child.Layout(ctx, childRect)},
		Title: fmt.Sprintf("Padding (t=%d l=%d b=%d r=%d)",
			c.T, c.L, c.B, c.R),
	}
}

func (c Padding) Render(ctx RenderContext, layout LayoutNode) {
	c.Child.Render(ctx, layout.Children[0])
}

func (c Padding) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "padding", Props: p, Children: []WireNode{child}}, nil
}

func PaddingFromWire(n WireNode) (Native, error) {
	var p Padding
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
