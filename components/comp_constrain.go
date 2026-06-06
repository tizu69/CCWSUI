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

func MkConstrain(w, h int, child Native) Constrain {
	return Constrain{W: w, H: h, Child: child}
}

func (c Constrain) Measure(ctx MeasureContext, constraint Size) Size {
	cw, ch := ifelse(c.W != 0, c.W, constraint.W), ifelse(c.H != 0, c.H, constraint.H)
	child := c.Child.Measure(ctx, Size{
		W: min(cw, constraint.W),
		H: min(ch, constraint.H),
	})
	return Size{
		W: min(child.W, cw, constraint.W),
		H: min(child.H, ch, constraint.H),
	}
}

func (c Constrain) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	cw, ch := ifelse(c.W != 0, c.W, rect.W), ifelse(c.H != 0, c.H, rect.H)
	childRect := Rect{
		X: rect.X, Y: rect.Y, W: min(cw, rect.W), H: min(ch, rect.H),
	}
	return LayoutNode{
		Rect:     rect,
		Children: []LayoutNode{c.Child.Layout(ctx, childRect)},
		Title:    fmt.Sprintf("Constrain (w=%d h=%d)", cw, ch),
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
