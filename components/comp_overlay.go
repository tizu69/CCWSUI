package components

import (
	"encoding/json"
	"fmt"
)

type Overlay struct {
	Layers        []Native `json:"-"`
	HoverRequired bool
}

func init() {
	RegisterWire("overlay", OverlayFromWire)
}

func MkOverlay(layers ...Native) Overlay {
	return Overlay{Layers: layers}
}

func (c Overlay) WithChildren(children ...Native) Overlay {
	c.Layers = append(c.Layers, children...)
	return c
}

func (c Overlay) WithRequireHover() Overlay {
	c.HoverRequired = true
	return c
}

func (c Overlay) Measure(ctx MeasureContext, constraint Size) Size {
	size := Size{}
	for _, layer := range c.Layers {
		layerSize := layer.Measure(ctx, constraint)
		if layerSize.W > size.W {
			size.W = layerSize.W
		}
		if layerSize.H > size.H {
			size.H = layerSize.H
		}
	}
	return size
}

func (c Overlay) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	children := make([]LayoutNode, len(c.Layers))
	if c.HoverRequired && !rect.Contains(ctx.GetMousePos()) {
		children[0] = c.Layers[0].Layout(ctx, rect)
	} else {
		for i, layer := range c.Layers {
			children[i] = layer.Layout(ctx, rect)
		}
	}
	return LayoutNode{
		Rect: rect, Title: fmt.Sprintf("Overlay (%dx)", len(c.Layers)),
		Children: children,
	}
}

func (c Overlay) Render(ctx RenderContext, layout LayoutNode) {
	if len(c.Layers) == 0 {
		return
	}
	if c.HoverRequired && !layout.Rect.Contains(ctx.GetMousePos()) {
		c.Layers[0].Render(ctx, layout.Children[0])
		return
	}
	for i, layer := range c.Layers {
		layer.Render(ctx, layout.Children[i])
	}
}

func (c Overlay) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	children := make([]WireNode, len(c.Layers))
	for i, layer := range c.Layers {
		if children[i], err = layer.ToWire(); err != nil {
			return WireNode{}, err
		}
	}
	return WireNode{Kind: "overlay", Props: p, Children: children}, nil
}

func OverlayFromWire(n WireNode) (Native, error) {
	var p Overlay
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}

	p.Layers = make([]Native, len(n.Children))
	for i, child := range n.Children {
		var err error
		if p.Layers[i], err = FromWire(child); err != nil {
			return nil, err
		}
	}
	return p, nil
}
