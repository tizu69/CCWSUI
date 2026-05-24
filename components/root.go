package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Root(title string, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []Node{
			Link(Rel("stylesheet"), Href("/static/ccwsui.css")),
		},
		Body: []Node{
			Class("bg-gradient-to-b from-white to-indigo-100 bg-no-repeat"),
			Div(Class("min-h-screen flex flex-col justify-between"),
				Div(Class("node"),
					Group(children),
					P(Text("Hello Worldd!")),
				),
			),
		},
	})
}
