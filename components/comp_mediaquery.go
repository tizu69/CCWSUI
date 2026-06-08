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

func (c MediaQuery) Measure(ctx MeasureContext, constraint Size) Size {
	if c.MinWidth > 0 && constraint.W < c.MinWidth {
		return Size{}
	}
	if c.MaxWidth > 0 && constraint.W > c.MaxWidth {
		return Size{}
	}
	if c.MinHeight > 0 && constraint.H < c.MinHeight {
		return Size{}
	}
	if c.MaxHeight > 0 && constraint.H > c.MaxHeight {
		return Size{}
	}
	return c.Child.Measure(ctx, constraint)
}

func (c MediaQuery) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	if c.MinWidth > 0 && rect.W < c.MinWidth {
		return LayoutNode{Rect: rect, Title: "MediaQuery (min width unmet)"}
	}
	if c.MaxWidth > 0 && rect.W > c.MaxWidth {
		return LayoutNode{Rect: rect, Title: "MediaQuery (max width unmet)"}
	}
	if c.MinHeight > 0 && rect.H < c.MinHeight {
		return LayoutNode{Rect: rect, Title: "MediaQuery (min height unmet)"}
	}
	if c.MaxHeight > 0 && rect.H > c.MaxHeight {
		return LayoutNode{Rect: rect, Title: "MediaQuery (max height unmet)"}
	}
	return LayoutNode{
		Rect: rect, Children: []LayoutNode{c.Child.Layout(ctx, rect)},
		Title: "MediaQuery",
	}
}

func (c MediaQuery) Render(ctx RenderContext, layout LayoutNode) {
	if c.MinWidth > 0 && layout.Rect.W < c.MinWidth {
		return
	}
	if c.MaxWidth > 0 && layout.Rect.W > c.MaxWidth {
		return
	}
	if c.MinHeight > 0 && layout.Rect.H < c.MinHeight {
		return
	}
	if c.MaxHeight > 0 && layout.Rect.H > c.MaxHeight {
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
