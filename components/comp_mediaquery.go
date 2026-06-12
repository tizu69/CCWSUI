package components

import "encoding/json"

type MediaQuery struct {
	MinWidth  int    `json:"minWidth"`
	MaxWidth  int    `json:"maxWidth"`
	MinHeight int    `json:"minHeight"`
	MaxHeight int    `json:"maxHeight"`
	Child     Native `json:"-"`
}

func init() {
	RegisterWire("mediaquery", MediaQueryFromWire)
}

func MkMediaQuery(child Native) MediaQuery { return MediaQuery{Child: child} }

func (c MediaQuery) WithMinWidth(minWidth int) MediaQuery {
	c.MinWidth = minWidth
	return c
}

func (c MediaQuery) WithMaxWidth(maxWidth int) MediaQuery {
	c.MaxWidth = maxWidth
	return c
}

func (c MediaQuery) WithMinHeight(minHeight int) MediaQuery {
	c.MinHeight = minHeight
	return c
}

func (c MediaQuery) WithMaxHeight(maxHeight int) MediaQuery {
	c.MaxHeight = maxHeight
	return c
}

func (c MediaQuery) Fits(w, h int) (reason string, ok bool) {
	if c.MinWidth > 0 && w < c.MinWidth {
		return "minimum width", false
	}
	if c.MaxWidth > 0 && w > c.MaxWidth {
		return "maximum width", false
	}
	if c.MinHeight > 0 && h < c.MinHeight {
		return "minimum height", false
	}
	if c.MaxHeight > 0 && h > c.MaxHeight {
		return "maximum height", false
	}
	return "", true
}

func (c MediaQuery) Measure(ctx MeasureContext, constraint Size) Size {
	if _, ok := c.Fits(ctx.GetDimensions()); !ok {
		return Size{}
	}
	return c.Child.Measure(ctx, constraint)
}

func (c MediaQuery) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	if reason, ok := c.Fits(ctx.GetDimensions()); !ok {
		return LayoutNode{Rect: rect, Title: "MediaQuery (unmet " + reason + ")"}
	}
	return LayoutNode{
		Rect: rect, Children: []LayoutNode{c.Child.Layout(ctx, rect)},
		Title: "MediaQuery",
	}
}

func (c MediaQuery) Render(ctx RenderContext, layout LayoutNode) {
	if _, ok := c.Fits(ctx.GetDimensions()); !ok {
		return
	}
	c.Child.Render(ctx, layout.Children[0])
}

func (c MediaQuery) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "mediaquery", Props: p, Children: []WireNode{child}}, nil
}

func MediaQueryFromWire(n WireNode) (Native, error) {
	var p MediaQuery
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
