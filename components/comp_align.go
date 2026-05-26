package components

import (
	"encoding/json"
	"fmt"
)

type Align struct {
	X, Y  Alignment
	Child Native `json:"-"`
}

func init() {
	RegisterWire("align", AlignFromWire)
}

func Aligned(x, y Alignment, child Native) Align { return Align{X: x, Y: y, Child: child} }
func AlignedX(x Alignment, child Native) Align   { return Aligned(x, AlignmentCenter, child) }
func AlignedY(y Alignment, child Native) Align   { return Aligned(AlignmentCenter, y, child) }
func AlignedCenter(child Native) Align           { return Aligned(AlignmentCenter, AlignmentCenter, child) }

func (c Align) Measure(ctx MeasureContext, constraint Size) Size {
	return constraint
	// do we want this?: return c.Child.Measure(ctx, constraint)
}

func (c Align) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	childSize := c.Child.Measure(ctx, Size{W: rect.W, H: rect.H})
	x := rect.X
	y := rect.Y
	switch c.X {
	case AlignmentCenter:
		x = rect.X + (rect.W-childSize.W)/2
	case AlignmentEnd:
		x = rect.X + (rect.W - childSize.W)
	}
	switch c.Y {
	case AlignmentCenter:
		y = rect.Y + (rect.H-childSize.H)/2
	case AlignmentEnd:
		y = rect.Y + (rect.H - childSize.H)
	}
	childRect := Rect{X: x, Y: y, W: childSize.W, H: childSize.H}
	return LayoutNode{
		Rect: rect, Children: []LayoutNode{c.Child.Layout(ctx, childRect)},
		Title: fmt.Sprintf("Align (h=%s v=%s)", c.X, c.Y),
	}
}

func (c Align) Render(ctx RenderContext, layout LayoutNode) {
	c.Child.Render(ctx, layout.Children[0])
}

func (c Align) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "align", Props: p, Children: []WireNode{child}}, nil
}

func AlignFromWire(n WireNode) (Native, error) {
	var p Align
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
