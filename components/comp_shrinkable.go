package components

import "encoding/json"

type Shrinkable struct {
	Child Native `json:"-"`
}

func init() {
	RegisterWire("shrinkable", ShrinkableFromWire)
}

func MkShrinkable(child Native) Shrinkable { return Shrinkable{Child: child} }

func (c Shrinkable) Measure(ctx MeasureContext, constraint Size) Size {
	return c.Child.Measure(ctx, constraint)
}

func (c Shrinkable) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{
		Rect: rect, Children: []LayoutNode{c.Child.Layout(ctx, rect)},
		Title: "Shrinkable",
	}
}

func (c Shrinkable) Render(ctx RenderContext, layout LayoutNode) {
	c.Child.Render(ctx, layout.Children[0])
}

func (c Shrinkable) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "shrinkable", Props: p, Children: []WireNode{child}}, nil
}

func ShrinkableFromWire(n WireNode) (Native, error) {
	var p Shrinkable
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
