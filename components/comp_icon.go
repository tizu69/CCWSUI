package components

import (
	"encoding/json"
	"fmt"
	"image/color"
)

type Icon struct {
	Icon   string
	Shadow bool
}

func init() {
	RegisterWire("icon", IconFromWire)
}

func Icony(icon string) Icon { return Icon{Icon: icon} }

func (c Icon) Tinted(tint color.RGBA) Icon {
	c.Icon += fmt.Sprintf(";tint=#%02x%02x%02x%02x", tint.R, tint.G, tint.B, tint.A)
	return c
}

func (c Icon) TintedHex(hex string) Icon {
	return c.Tinted(colorFromHex(hex))
}

func (c Icon) Rotate(rot Rotation) Icon {
	c.Icon += fmt.Sprintf(";rotate=%d", rot)
	return c
}

func (c Icon) Flip(flip Flip) Icon {
	c.Icon += fmt.Sprintf(";flip=%s", flip)
	return c
}

func (c Icon) SetShadow(v bool) Icon {
	c.Shadow = true
	return c
}

func (c Icon) Measure(ctx MeasureContext, constraint Size) Size {
	w, h := ctx.GetTexSize("@icon/" + c.Icon)
	return Size{W: w, H: h}
}

func (c Icon) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{Rect: rect, Title: fmt.Sprintf("Icon (%s)", c.Icon)}
}

func (c Icon) Render(ctx RenderContext, layout LayoutNode) {
	path := "@icon/" + c.Icon
	ctx.RequireTexture(path)
	w, h := ctx.GetTexSize(path)
	if c.Shadow {
		ctx.RequireTexture(path + ";shadow")
		ctx.RenderTex(layout.Rect.X+1, layout.Rect.Y+1, w, h, path+";shadow", 0, 0)
	}
	ctx.RenderTex(layout.Rect.X, layout.Rect.Y, w, h, path, 0, 0)
}

func (c Icon) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	return WireNode{Kind: "icon", Props: p}, err
}

func IconFromWire(n WireNode) (Native, error) {
	var p Icon
	err := json.Unmarshal(n.Props, &p)
	return p, err
}
