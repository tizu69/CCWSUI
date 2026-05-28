package components

import (
	"encoding/json"
	"fmt"
)

type Constrain struct {
	W, H  int
	Child Native `json:"-"`
}

func init() {
	RegisterWire("constrain", ConstrainFromWire)
}

func Constrained(w, h int, child Native) Constrain {
	return Constrain{W: w, H: h, Child: child}
}

func (c Constrain) Measure(ctx MeasureContext, constraint Size) Size {
	child := c.Child.Measure(ctx, Size{
		W: min(c.W, constraint.W),
		H: min(c.H, constraint.H),
	})
	return Size{
		W: min(child.W, c.W, constraint.W),
		H: min(child.H, c.H, constraint.H),
	}
}

func (c Constrain) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	childRect := Rect{
		X: rect.X, Y: rect.Y, W: min(c.W, rect.W), H: min(c.H, rect.H),
	}
	return LayoutNode{
		Rect:     rect,
		Children: []LayoutNode{c.Child.Layout(ctx, childRect)},
		Title:    fmt.Sprintf("Constrain (w=%d, h=%d)", c.W, c.H),
	}
}

func (c Constrain) Render(ctx RenderContext, layout LayoutNode) {
	c.Child.Render(ctx, layout.Children[0])
}

func (c Constrain) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "constrain", Props: p, Children: []WireNode{child}}, nil
}

func ConstrainFromWire(n WireNode) (Native, error) {
	var p Constrain
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
