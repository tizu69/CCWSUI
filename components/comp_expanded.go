package components

import "encoding/json"

type Expanded struct {
	Child     Native `json:"-"`
	Direction Direction
}

func init() {
	RegisterWire("expanded", ExpandedFromWire)
}

func MkExpand(child Native) Expanded  { return Expanded{Child: child, Direction: DirectionHV} }
func MkExpandH(child Native) Expanded { return Expanded{Child: child, Direction: DirectionH} }
func MkExpandV(child Native) Expanded { return Expanded{Child: child, Direction: DirectionV} }

func (c Expanded) Measure(ctx MeasureContext, constraint Size) Size {
	switch c.Direction {
	case DirectionH:
		constraint.H = c.Child.Measure(ctx, constraint).H
	case DirectionV:
		constraint.W = c.Child.Measure(ctx, constraint).W
	}
	return constraint
}

func (c Expanded) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{
		Rect: rect, Children: []LayoutNode{c.Child.Layout(ctx, rect)},
		Title: "Expanded",
	}
}

func (c Expanded) Render(ctx RenderContext, layout LayoutNode) {
	c.Child.Render(ctx, layout.Children[0])
}

func (c Expanded) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "expanded", Props: p, Children: []WireNode{child}}, nil
}

func ExpandedFromWire(n WireNode) (Native, error) {
	var p Expanded
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
