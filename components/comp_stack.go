package components

import (
	"encoding/json"
	"fmt"
)

type Stack struct {
	Direction StackDirection
	Children  []Native `json:"-"`
	Padding   int
	Align     Alignment
}

func init() {
	RegisterWire("stack", StackFromWire)
}

func MkStackH(children ...Native) Stack { return Stack{Direction: StackDirectionH, Children: children} }
func MkStackV(children ...Native) Stack { return Stack{Direction: StackDirectionV, Children: children} }

func (c Stack) WithChildren(children ...Native) Stack {
	c.Children = append(c.Children, children...)
	return c
}

func (c Stack) WithGap(padding int) Stack {
	c.Padding = padding
	return c
}

func (c Stack) WithAlign(align Alignment) Stack {
	c.Align = align
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
			if isStackRest(child) {
				continue
			}
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
			if isStackRest(child) {
				continue
			}
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
	constraint := Size{W: max(0, rect.W), H: max(0, rect.H)}
	if c.Direction == StackDirectionH {
		fixedW := 0
		totalRest := 0
		for _, child := range c.Children {
			if isStackRest(child) {
				totalRest++
				continue
			}
			cs := child.Measure(ctx, constraint)
			fixedW += cs.W
		}

		remainingW := max(0, rect.W-fixedW-max(0, len(c.Children)-1)*c.Padding)
		x := rect.X
		if totalRest == 0 && c.Align != 0 {
			used := fixedW + max(0, len(c.Children)-1)*c.Padding
			extra := max(0, rect.W-used)
			x += int(float32(extra) * float32(c.Align))
		}

		for i, child := range c.Children {
			w := 0
			if isStackRest(child) {
				w = remainingW / totalRest
			} else {
				w = child.Measure(ctx, constraint).W
			}
			childRect := Rect{X: x, Y: rect.Y, W: w, H: rect.H}
			children = append(children, child.Layout(ctx, childRect))
			x += w
			if i < len(c.Children)-1 {
				x += c.Padding
			}
		}
	} else {
		fixedH := 0
		totalRest := 0
		for _, child := range c.Children {
			if isStackRest(child) {
				totalRest++
				continue
			}
			cs := child.Measure(ctx, constraint)
			fixedH += cs.H
		}

		remainingH := max(0, rect.H-fixedH-max(0, len(c.Children)-1)*c.Padding)
		y := rect.Y
		if totalRest == 0 && c.Align != 0 {
			used := fixedH + max(0, len(c.Children)-1)*c.Padding
			extra := max(0, rect.H-used)
			y += int(float32(extra) * float32(c.Align))
		}

		for i, child := range c.Children {
			h := 0
			if isStackRest(child) {
				h = remainingH / totalRest
			} else {
				h = child.Measure(ctx, constraint).H
			}
			childRect := Rect{X: rect.X, Y: y, W: rect.W, H: h}
			children = append(children, child.Layout(ctx, childRect))
			y += h
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
	c.RenderClipped(ctx, layout, layout.Rect)
}

func (c Stack) RenderClipped(ctx RenderContext, layout LayoutNode, clip Rect) {
	for i, child := range c.Children {
		childLayout := layout.Children[i]
		if !childLayout.Rect.Intersects(clip) {
			continue
		}
		if clipped, ok := child.(ClippedRenderer); ok {
			clipped.RenderClipped(ctx, childLayout, clip)
			continue
		}
		child.Render(ctx, childLayout)
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

func isStackRest(child Native) bool {
	_, ok := child.(Expanded)
	return ok
}
