package components

import (
	"encoding/json"
	"fmt"
)

type ClickRegion struct {
	Event string
	Child Native `json:"-"`
}

func init() {
	RegisterWire("clickregion", ClickRegionFromWire)
}

func Clickable(event string, child Native) ClickRegion {
	return ClickRegion{Event: event, Child: child}
}

func (c ClickRegion) Measure(ctx MeasureContext, constraint Size) Size {
	return c.Child.Measure(ctx, constraint)
}

type clickRegionEvent struct {
	Shift bool `json:"shift"`
	Ctrl  bool `json:"ctrl"`
	Alt   bool `json:"alt"`
}

func (c ClickRegion) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	if ctx.GetMouseDown() && rect.Contains(ctx.GetMousePos()) {
		ctx.SendEvent(c.Event, clickRegionEvent{
			Shift: ctx.GetShiftDown(),
			Ctrl:  ctx.GetCtrlDown(),
			Alt:   ctx.GetAltDown(),
		})
	}
	return LayoutNode{
		Rect: rect, Title: fmt.Sprintf("ClickRegion (%s)", c.Event),
		Children: []LayoutNode{c.Child.Layout(ctx, rect)},
	}
}

func (c ClickRegion) Render(ctx RenderContext, layout LayoutNode) {
	c.Child.Render(ctx, layout.Children[0])
}

func (c ClickRegion) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "clickregion", Props: p, Children: []WireNode{child}}, nil
}

func ClickRegionFromWire(n WireNode) (Native, error) {
	var p ClickRegion
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
