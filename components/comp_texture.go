package components

import (
	"encoding/json"
	"fmt"
)

type Texture struct {
	Tex        string
	T, L, B, R int
	Child      Native `json:"-"`
}

func Textured(tex string, t, l, b, r int, child Native) Texture {
	return Texture{Tex: tex, T: t, L: l, B: b, R: r, Child: child}
}

func init() {
	RegisterWire("texture", TextureFromWire)
}

func (c Texture) Measure(ctx MeasureContext, constraint Size) Size {
	return c.Child.Measure(ctx, constraint)
}

func (c Texture) Layout(ctx LayoutContext, rect Rect) LayoutNode {
	return LayoutNode{
		Rect:     rect,
		Children: []LayoutNode{c.Child.Layout(ctx, rect)},
		Title:    fmt.Sprintf("Texture (%s)", c.Tex),
	}
}

func (c Texture) Render(ctx RenderContext, layout LayoutNode) {
	r := layout.Rect
	texW, texH := ctx.GetTexSize(c.Tex)
	if r.W <= 0 || r.H <= 0 || texW <= 0 || texH <= 0 {
		c.Child.Render(ctx, layout.Children[0])
		return
	}

	// destination border sizes
	leftW := min(c.L, r.W)
	rightW := min(c.R, max(0, r.W-leftW))
	topH := min(c.T, r.H)
	bottomH := min(c.B, max(0, r.H-topH))
	midW := max(0, r.W-leftW-rightW)
	midH := max(0, r.H-topH-bottomH)

	// source center size
	srcMidW := max(1, texW-c.L-c.R)
	srcMidH := max(1, texH-c.T-c.B)

	x0, x1, x2 := r.X, r.X+leftW, r.X+leftW+midW
	y0, y1, y2 := r.Y, r.Y+topH, r.Y+topH+midH

	// corners
	if leftW > 0 && topH > 0 {
		ctx.RenderTex(x0, y0, leftW, topH, c.Tex, 0, 0)
	}
	if rightW > 0 && topH > 0 {
		ctx.RenderTex(x2, y0, rightW, topH, c.Tex, texW-c.R, 0)
	}
	if leftW > 0 && bottomH > 0 {
		ctx.RenderTex(x0, y2, leftW, bottomH, c.Tex, 0, texH-c.B)
	}
	if rightW > 0 && bottomH > 0 {
		ctx.RenderTex(x2, y2, rightW, bottomH, c.Tex, texW-c.R, texH-c.B)
	}

	// edges
	for x := x1; x < x2; {
		w := min(srcMidW, x2-x)
		if topH > 0 {
			ctx.RenderTex(x, y0, w, topH, c.Tex, c.L, 0)
		}
		if bottomH > 0 {
			ctx.RenderTex(x, y2, w, bottomH, c.Tex, c.L, texH-c.B)
		}
		x += w
	}
	for y := y1; y < y2; {
		h := min(srcMidH, y2-y)
		if leftW > 0 {
			ctx.RenderTex(x0, y, leftW, h, c.Tex, 0, c.T)
		}
		if rightW > 0 {
			ctx.RenderTex(x2, y, rightW, h, c.Tex, texW-c.R, c.T)
		}
		y += h
	}

	// center
	for y := y1; y < y2; {
		h := min(srcMidH, y2-y)
		for x := x1; x < x2; {
			w := min(srcMidW, x2-x)
			ctx.RenderTex(x, y, w, h, c.Tex, c.L, c.T)
			x += w
		}
		y += h
	}

	c.Child.Render(ctx, layout.Children[0])
}

func (c Texture) ToWire() (WireNode, error) {
	p, err := json.Marshal(c)
	if err != nil {
		return WireNode{}, err
	}
	child, err := c.Child.ToWire()
	if err != nil {
		return WireNode{}, err
	}
	return WireNode{Kind: "texture", Props: p, Children: []WireNode{child}}, nil
}

func TextureFromWire(n WireNode) (Native, error) {
	var p Texture
	if err := json.Unmarshal(n.Props, &p); err != nil {
		return nil, err
	}
	var err error
	p.Child, err = FromWire(n.Children[0])
	return p, err
}
