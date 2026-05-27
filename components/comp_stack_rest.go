package components

import "encoding/json"

type StackRest struct {
	Child Native `json:"-"`
}

func init() {
	RegisterWire("stackrest", StackRestFromWire)
}

func StackUseRest(child Native) StackRest {
	return StackRest{Child: child}
}

func (c StackRest) Measure(ctx MeasureContext, constraint Size) Size {
	return Size{}
}

func (c StackRest) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{
		Rect: rect, Children: []LayoutNode{c.Child.Layout(ctx, rect)},
		Title: "StackRest",
	}
}

func (c StackRest) Render(ctx RenderContext, layout LayoutNode) {
	c.Child.Render(ctx, layout.Children[0])
}

func (c StackRest) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "stackrest", Props: p, Children: []WireNode{child}}, nil
}

func StackRestFromWire(n WireNode) (Native, error) {
	var p StackRest
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}

var _ StackChild = (*StackRest)(nil)

func (c StackRest) StackRest() int {
	return 1
}
