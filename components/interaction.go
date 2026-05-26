package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ClickRegion struct {
	Event string
	Child Native
}

func Clickable(event string, child Native) ClickRegion {
	return ClickRegion{Event: event, Child: child}
}

func (c ClickRegion) Render(ctx RenderContext) Node {
	return Div(Data("ccwsui", "clickregion"), c.Child.Render(ctx),
		ID("ccwsui-clickregion-"+c.Event),
		Attr("role", "button"),
		Attr("hx-post", ctx.EventURL(c.Event)),
		Attr("hx-trigger", "click"))
}

type InputRegion struct {
	Event string
	Value string
}

func Inputable(event, value string) InputRegion {
	return InputRegion{Event: event, Value: value}
}

func (c InputRegion) Render(ctx RenderContext) Node {
	return Input(Data("ccwsui", "inputregion"),
		ID("ccwsui-inputregion-"+c.Event),
		Value(c.Value), Attr("role", "textbox"),
		Attr("hx-post", ctx.EventURL(c.Event)),
		Attr("hx-trigger", "input delay:2000ms"), Styles{
			"text-shadow:%s": ifelse(true, "color-mix(in srgb, currentColor 25%, #000000) var(--px) var(--px)", "none"),
		})
}
