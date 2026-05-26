package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type VStack struct {
	Align    StackAlign
	Children []Native
}

func VStacked(children ...Native) VStack {
	return VStack{Children: children}
}

func (c VStack) WithAlign(align StackAlign) VStack {
	c.Align = align
	return c
}

func (c VStack) Render(ctx RenderContext) Node {
	children := make([]Node, len(c.Children))
	for i, child := range c.Children {
		children[i] = child.Render(ctx)
	}
	return Div(Data("ccwsui", "vstack"), Styles{
		"display:flex":          nil,
		"flex-direction:column": nil,
		"align-items:%s":        c.Align.CSS(),
	}, Group(children), Data("ccwsui-snap", ""))
}

type HStack struct {
	Align          StackAlign
	Children       []Native
	MinW, MaxW     int
	OverflowScroll bool
}

func HStacked(children ...Native) HStack {
	return HStack{Children: children}
}

func (c HStack) WithAlign(align StackAlign) HStack {
	c.Align = align
	return c
}

func (c HStack) Render(ctx RenderContext) Node {
	children := make([]Node, len(c.Children))
	for i, child := range c.Children {
		children[i] = child.Render(ctx)
	}
	return Div(Data("ccwsui", "hstack"), Styles{
		"display:flex":       nil,
		"flex-direction:row": nil,
		"justify-content:%s": c.Align.CSS(),
	}, Group(children), Data("ccwsui-snap", ""))
}
