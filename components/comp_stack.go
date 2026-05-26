package components

import (
	"encoding/json"
	"fmt"
)

type Stack struct {
	Direction StackDirection
	Children  []Native `json:"-"`
	Padding   int
}

func init() {
	RegisterWire("stack", StackFromWire)
}

func HStacked(children ...Native) Stack { return Stack{Direction: StackDirectionH, Children: children} }
func VStacked(children ...Native) Stack { return Stack{Direction: StackDirectionV, Children: children} }

func (c Stack) WithPadding(padding int) Stack {
	c.Padding = padding
	return c
}

func (c Stack) Measure(ctx MeasureContext, constraint Size) Size {
	if len(c.Children) == 0 {
		return Size{}
	}
	if c.Direction == StackDirectionH {
		totalW := 0
		maxH := 0
		for i, child := range c.Children {
			cs := child.Measure(ctx, constraint)
			totalW += cs.W
			if i > 0 {
				totalW += c.Padding
			}
			maxH = max(maxH, cs.H)
		}
		return Size{W: totalW, H: maxH}
	} else {
		totalH := 0
		maxW := 0
		for i, child := range c.Children {
			cs := child.Measure(ctx, constraint)
			totalH += cs.H
			if i > 0 {
				totalH += c.Padding
			}
			maxW = max(maxW, cs.W)
		}
		return Size{W: maxW, H: totalH}
	}
}

func (c Stack) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	children := make([]LayoutNode, 0, len(c.Children))
	if c.Direction == StackDirectionH {
		x := rect.X
		for i, child := range c.Children {
			cs := child.Measure(ctx, Size{W: max(0, rect.W), H: max(0, rect.H)})
			childRect := Rect{X: x, Y: rect.Y, W: cs.W, H: rect.H}
			children = append(children, child.Layout(ctx, childRect))
			x += cs.W
			if i < len(c.Children)-1 {
				x += c.Padding
			}
		}
	} else {
		y := rect.Y
		for i, child := range c.Children {
			cs := child.Measure(ctx, Size{W: max(0, rect.W), H: max(0, rect.H)})
			childRect := Rect{X: rect.X, Y: y, W: rect.W, H: cs.H}
			children = append(children, child.Layout(ctx, childRect))
			y += cs.H
			if i < len(c.Children)-1 {
				y += c.Padding
			}
		}
	}
	return LayoutNode{
		Rect: rect, Children: children,
		Title: fmt.Sprintf("Stack (%s x%d)", c.Direction, len(c.Children)),
	}
}

func (c Stack) Render(ctx RenderContext, layout LayoutNode) {
	for i, child := range c.Children {
		child.Render(ctx, layout.Children[i])
	}
}

func (c Stack) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	children := make([]WireNode, 0, len(c.Children))
	for _, child := range c.Children {
		w, err := child.ToWire()
		if err != nil {
			return WireNode{}, err
		}
		children = append(children, w)
	}
	return WireNode{Kind: "stack", Props: p, Children: children}, nil
}

func StackFromWire(n WireNode) (Native, error) {
	var p Stack
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	for _, child := range n.Children {
		c, err := FromWire(child)
		if err != nil {
			return nil, err
		}
		p.Children = append(p.Children, c)
	}
	return p, err
}
