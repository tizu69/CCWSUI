package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Stack struct {
	Direction StackDirection
	Children  []Native
	Padding   int
}

func HStacked(children ...Native) Stack { return Stack{Direction: StackDirectionH, Children: children} }
func VStacked(children ...Native) Stack { return Stack{Direction: StackDirectionV, Children: children} }

func (c Stack) WithPadding(padding int) Stack {
	c.Padding = padding
	return c
}

func (c Stack) Render(ctx RenderContext) Node {
	children := make([]Node, len(c.Children))
	for i, child := range c.Children {
		children[i] = child.Render(ctx)
	}
	return Div(Data("ccwsui", "stack"), Styles{
		"display:flex":           nil,
		"flex-direction:%s":      c.Direction.CSS(),
		"width:100%":             nil,
		"height:100%":            nil,
		"gap:calc(var(--px)*%d)": c.Padding,
		"position:relative":      nil,
	}, Group(children), Data("ccwsui-snap", ""))
}

type Align struct {
	X, Y  Alignment
	Child Native
}

func Aligned(x, y Alignment, child Native) Align { return Align{X: x, Y: y, Child: child} }
func AlignedX(x Alignment, child Native) Align   { return Aligned(x, AlignmentCenter, child) }
func AlignedY(y Alignment, child Native) Align   { return Aligned(AlignmentCenter, y, child) }
func AlignedCenter(child Native) Align           { return Aligned(AlignmentCenter, AlignmentCenter, child) }

func (c Align) Render(ctx RenderContext) Node {
	return Div(Data("ccwsui", "align"), Styles{
		"display:flex":       nil,
		"justify-content:%s": c.X.CSS(),
		"align-items:%s":     c.Y.CSS(),
		"width:100%":         nil,
		"height:100%":        nil,
	}, c.Child.Render(ctx), Data("ccwsui-snap", ""))
}

type Expanded struct {
	Child Native
}

func Expandeded(child Native) Expanded {
	return Expanded{Child: child}
}

func (c Expanded) Render(ctx RenderContext) Node {
	return Div(Data("ccwsui", "expanded"), Styles{
		"flex:1": nil,
	}, c.Child.Render(ctx))
}

type RelativeParent struct {
	Child Native
}

type Absolute struct {
	Child Native
	X, Y  int
}

func AbsoluteAt(x, y int, child Native) Absolute {
	return Absolute{Child: child, X: x, Y: y}
}

func (c Absolute) Render(ctx RenderContext) Node {
	st := Styles{"position:absolute": nil}
	if c.X >= 0 {
		st["left:calc(var(--px)*%d)"] = c.X
	} else {
		st["right:calc(var(--px)*%d)"] = -c.X - 1
	}
	if c.Y >= 0 {
		st["top:calc(var(--px)*%d)"] = c.Y
	} else {
		st["bottom:calc(var(--px)*%d)"] = -c.Y - 1
	}

	return Div(Data("ccwsui", "absolutechild"), st,
		c.Child.Render(ctx), Data("ccwsui-snap", ""))
}
