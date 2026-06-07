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
	totalPadding := (len(c.Children) - 1) * c.Padding
	if c.Direction == StackDirectionH {
		ideals := make([]int, 0, len(c.Children))
		maxH := 0
		nRest := 0

		for _, child := range c.Children {
			if isStackRest(child) {
				nRest++
				cs := child.Measure(ctx, Size{W: 0, H: constraint.H})
				maxH = max(maxH, cs.H)
				continue
			}
			cs := child.Measure(ctx, constraint)
			maxH = max(maxH, cs.H)
			ideals = append(ideals, cs.W)
		}

		availW := max(0, constraint.W-totalPadding)
		totalW := 0
		for _, w := range ideals {
			totalW += w
		}

		if totalW > availW {
			rigidW := 0
			shrinkN := 0
			for _, child := range c.Children {
				if isStackRest(child) {
					continue
				}
				if isShrinkable(child) {
					shrinkN++
				} else {
					rigidW += child.Measure(ctx, constraint).W
				}
			}

			remaining := availW - rigidW
			if remaining > 0 && shrinkN > 0 {
				each := remaining / shrinkN
				totalW = rigidW
				for _, child := range c.Children {
					if isStackRest(child) || !isShrinkable(child) {
						continue
					}
					sz := child.Measure(ctx, Size{W: each, H: constraint.H}).W
					totalW += min(sz, each)
				}
			} else {
				totalW = rigidW
			}
		}

		return Size{W: totalW + totalPadding, H: maxH}
	} else {
		ideals := make([]int, 0, len(c.Children))
		maxW := 0
		nRest := 0

		for _, child := range c.Children {
			if isStackRest(child) {
				nRest++
				cs := child.Measure(ctx, Size{W: constraint.W, H: 0})
				maxW = max(maxW, cs.W)
				continue
			}
			cs := child.Measure(ctx, constraint)
			maxW = max(maxW, cs.W)
			ideals = append(ideals, cs.H)
		}

		availH := max(0, constraint.H-totalPadding)
		totalH := 0
		for _, h := range ideals {
			totalH += h
		}

		if totalH > availH {
			rigidH := 0
			shrinkN := 0
			for _, child := range c.Children {
				if isStackRest(child) {
					continue
				}
				if isShrinkable(child) {
					shrinkN++
				} else {
					rigidH += child.Measure(ctx, constraint).H
				}
			}

			remaining := availH - rigidH
			if remaining > 0 && shrinkN > 0 {
				each := remaining / shrinkN
				totalH = rigidH
				for _, child := range c.Children {
					if isStackRest(child) || !isShrinkable(child) {
						continue
					}
					sz := child.Measure(ctx, Size{W: constraint.W, H: each}).H
					totalH += min(sz, each)
				}
			} else {
				totalH = rigidH
			}
		}

		return Size{W: maxW, H: totalH + totalPadding}
	}
}

func (c Stack) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	constraint := Size{W: max(0, rect.W), H: max(0, rect.H)}
	if c.Direction == StackDirectionH {
		return c.layoutAlong(ctx, rect, constraint,
			func(s Size) int { return s.W },
			func(s *Size, v int) { s.W = v },
			func(r Rect) int { return r.X },
			func(r *Rect, v int) { r.X = v },
			func(r *Rect, v int) { r.W = v },
		)
	}
	return c.layoutAlong(ctx, rect, constraint,
		func(s Size) int { return s.H },
		func(s *Size, v int) { s.H = v },
		func(r Rect) int { return r.Y },
		func(r *Rect, v int) { r.Y = v },
		func(r *Rect, v int) { r.H = v },
	)
}

func (c Stack) layoutAlong(ctx LayoutContext, rect Rect, constraint Size,
	mainDim func(Size) int,
	setMainDim func(*Size, int),
	getPos func(Rect) int,
	setPos func(*Rect, int),
	setDim func(*Rect, int),
) LayoutNode {
	children := make([]LayoutNode, 0, len(c.Children))
	totalPadding := (len(c.Children) - 1) * c.Padding

	fixedMain := 0
	totalRest := 0
	for _, child := range c.Children {
		if isStackRest(child) {
			totalRest++
			continue
		}
		fixedMain += mainDim(child.Measure(ctx, constraint))
	}

	mainConstraint := mainDim(constraint)
	availMain := max(0, mainConstraint-totalPadding)
	pos := getPos(rect)

	if fixedMain > availMain {
		sizes := make([]int, len(c.Children))
		rigidMain := 0
		shrinkN := 0
		for i, child := range c.Children {
			if isStackRest(child) {
				continue
			}
			sz := mainDim(child.Measure(ctx, constraint))
			sizes[i] = sz
			if isShrinkable(child) {
				shrinkN++
			} else {
				rigidMain += sz
			}
		}

		remaining := availMain - rigidMain
		if remaining > 0 && shrinkN > 0 {
			each := remaining / shrinkN
			for i, child := range c.Children {
				if isStackRest(child) || !isShrinkable(child) {
					continue
				}
				shrunkConstraint := constraint
				setMainDim(&shrunkConstraint, each)
				sizes[i] = min(mainDim(child.Measure(ctx, shrunkConstraint)), each)
			}
		}

		for i, child := range c.Children {
			m := 0
			if !isStackRest(child) {
				m = sizes[i]
			}
			childRect := rect
			setPos(&childRect, pos)
			setDim(&childRect, m)
			children = append(children, child.Layout(ctx, childRect))
			pos += m
			if i < len(c.Children)-1 {
				pos += c.Padding
			}
		}
	} else {
		remainingMain := max(0, mainConstraint-fixedMain-totalPadding)
		if totalRest == 0 && c.Align != 0 {
			extra := max(0, availMain-fixedMain)
			pos += int(float32(extra) * float32(c.Align))
		}
		for i, child := range c.Children {
			m := 0
			if isStackRest(child) {
				m = remainingMain / totalRest
			} else {
				m = mainDim(child.Measure(ctx, constraint))
			}
			childRect := rect
			setPos(&childRect, pos)
			setDim(&childRect, m)
			children = append(children, child.Layout(ctx, childRect))
			pos += m
			if i < len(c.Children)-1 {
				pos += c.Padding
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

func isShrinkable(child Native) bool {
	switch child.(type) {
	case Shrinkable, Expanded:
		return true
	}
	return false
}
