package components

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ItemTexture struct {
	Item string
}

func init() {
	RegisterWire("itemtexture", ItemTextureFromWire)
}

func ItemTextured(item string) ItemTexture { return ItemTexture{Item: item} }

func (c ItemTexture) Measure(ctx MeasureContext, constraint Size) Size {
	return Size{W: 16, H: 16}
}

func (c ItemTexture) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{Rect: rect, Title: fmt.Sprintf("ItemTexture (%s)", c.Item)}
}

func (c ItemTexture) Render(ctx RenderContext, layout LayoutNode) {
	path := "@item/" + strings.ReplaceAll(c.Item, ":", "__")
	ctx.RequireTexture(path)
	w, h := ctx.GetTexSize(path)
	ctx.RenderTexMixel(layout.Rect.X, layout.Rect.Y, 16, 16, path, 0, 0, w, h, false)
}

func (c ItemTexture) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	return WireNode{Kind: "itemtexture", Props: p}, err
}

func ItemTextureFromWire(n WireNode) (Native, error) {
	var p ItemTexture
	err := json.Unmarshal(n.Props, &p)
	return p, err
}
