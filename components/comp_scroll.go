package components

import (
	"encoding/json"
	"fmt"
)

type Scroll struct {
	Direction  Direction
	Step, X, Y int
	Child      Native `json:"-"`
}

func init() {
	RegisterWire("scroll", ScrollFromWire)
}

func Scrolling(direction Direction, child Native) *Scroll {
	return &Scroll{Direction: direction, Child: child, Step: 8}
}

func (c Scroll) SetStep(step int) *Scroll {
	c.Step = step
	return &c
}

func (c Scroll) SetCurrent(x, y int) *Scroll {
	c.X = x
	c.Y = y
	return &c
}

func (c Scroll) Measure(ctx MeasureContext, constraint Size) Size {
	unconstrained := Size{W: constraint.W, H: constraint.H}
	if c.scrollsH() {
		unconstrained.W = 1<<31 - 1
	}
	if c.scrollsV() {
		unconstrained.H = 1<<31 - 1
	}
	child := c.Child.Measure(ctx, unconstrained)
	return Size{
		W: min(child.W, constraint.W),
		H: min(child.H, constraint.H),
	}
}

func (c *Scroll) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	childSize := c.Child.Measure(ctx, Size{W: rect.W, H: rect.H})

	if rect.Contains(ctx.GetMousePos()) {
		dx, dy := ctx.GetMouseScroll()
		dx, dy = sign(dx, c.Step), sign(dy, c.Step)
		if dx != 0 || dy != 0 {
			c.X = max(min(c.X+dx, childSize.W-rect.W), 0)
			c.Y = max(min(c.Y+dy, childSize.H-rect.H), 0)
		}
	}

	childRect := Rect{
		X: rect.X - c.X, Y: rect.Y - c.Y, W: childSize.W, H: childSize.H,
	}
	return LayoutNode{
		Rect:     rect,
		Children: []LayoutNode{c.Child.Layout(ctx, childRect)},
		Title:    fmt.Sprintf("Scroll (%s)", c.Direction),
	}
}

func (c *Scroll) Render(ctx RenderContext, layout LayoutNode) {
	ctx.Scissor(layout.Rect.X, layout.Rect.Y, layout.Rect.W, layout.Rect.H)
	childLayout := layout.Children[0]
	if clipped, ok := c.Child.(ClippedRenderer); ok {
		clipped.RenderClipped(ctx, childLayout, layout.Rect)
	} else {
		c.Child.Render(ctx, childLayout)
	}
	ctx.PopScissor()
}

func sign(x, n int) int {
	if x > 0 {
		return n
	} else if x < 0 {
		return -n
	}
	return 0
}

func (c Scroll) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "scroll", Props: p, Children: []WireNode{child}}, nil
}

func ScrollFromWire(n WireNode) (Native, error) {
	var p Scroll
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return &p, err
}

func (c Scroll) scrollsH() bool {
	return c.Direction == DirectionH || c.Direction == DirectionHV
}

func (c Scroll) scrollsV() bool {
	return c.Direction == DirectionV || c.Direction == DirectionHV
}
