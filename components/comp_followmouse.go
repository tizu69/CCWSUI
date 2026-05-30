package components

import (
	"encoding/json"
)

type FollowMouse struct {
	X, Y           Alignment
	Child          Native `json:"-"`
	FlipIfOverflow bool
}

func init() {
	RegisterWire("followmouse", FollowMouseFromWire)
}

func FollowsMouse(x, y Alignment, child Native) FollowMouse {
	return FollowMouse{X: x, Y: y, Child: child}
}

func (c FollowMouse) FlipIfOverflowing(v bool) FollowMouse {
	c.FlipIfOverflow = v
	return c
}

func (c FollowMouse) Measure(ctx MeasureContext, constraint Size) Size {
	// we do not take up any parent space.
	return Size{}
}

func (c FollowMouse) Layout(ctx LayoutContext, _ Rect) LayoutNode {
	w, h := ctx.GetDimensions()
	x, y := ctx.GetMousePos()
	wanted := c.Child.Measure(ctx, Size{W: w - x, H: h - y})

	rect := Rect{X: x, Y: y, W: wanted.W, H: wanted.H}
	rect.X -= int(float32(rect.W) * float32(c.X))
	rect.Y -= int(float32(rect.H) * float32(c.Y))

	if wanted.W > w-rect.X {
		rect.X = ifelse(c.FlipIfOverflow, x, w) - wanted.W
	}
	if wanted.H > h-rect.Y {
		rect.Y = ifelse(c.FlipIfOverflow, y, h) - wanted.H
	}

	child := c.Child.Layout(ctx, rect)
	return LayoutNode{
		Rect: rect, Title: "FollowMouse", Children: []LayoutNode{child},
	}
}

func (c FollowMouse) Render(ctx RenderContext, layout LayoutNode) {
	if x, y := ctx.GetMousePos(); x == -1 && y == -1 {
		return
	}
	ctx.RenderOverlay(func() { c.Child.Render(ctx, layout.Children[0]) })
}

func (c FollowMouse) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "followmouse", Props: p, Children: []WireNode{child}}, nil
}

func FollowMouseFromWire(n WireNode) (Native, error) {
	var p FollowMouse
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
