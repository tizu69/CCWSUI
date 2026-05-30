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
	return c.Child.Measure(ctx, constraint)
}

func (c Align) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	childSize := c.Child.Measure(ctx, Size{W: rect.W, H: rect.H})
	x := rect.X + int(float32(rect.W-childSize.W)*float32(c.X))
	y := rect.Y + int(float32(rect.H-childSize.H)*float32(c.Y))
	childRect := Rect{X: x, Y: y, W: childSize.W, H: childSize.H}
	return LayoutNode{
		Rect: rect, Children: []LayoutNode{c.Child.Layout(ctx, childRect)},
		Title: fmt.Sprintf("Align (h=%v v=%v)", c.X, c.Y),
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
