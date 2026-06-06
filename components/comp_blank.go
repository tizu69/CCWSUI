package components

import "encoding/json"

type Blank struct {
	W, H int
}

func init() {
	RegisterWire("blank", BlankFromWire)
}

func MkBlank() Blank          { return Blank{} }
func MkFiller(w, h int) Blank { return Blank{W: w, H: h} }

func (c Blank) Measure(ctx MeasureContext, constraint Size) Size {
	return Size{W: c.W, H: c.H}
}

func (c Blank) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{Rect: rect, Title: "Blank"}
}

func (c Blank) Render(ctx RenderContext, layout LayoutNode) {}

func (c Blank) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	return WireNode{Kind: "blank", Props: p}, err
}

func BlankFromWire(n WireNode) (Native, error) {
	var p Blank
	err := json.Unmarshal(n.Props, &p)
	return p, err
}
