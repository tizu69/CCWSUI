package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"

	"g.tizu.dev/CCWSUI/components"
	"github.com/coder/websocket"
)

type Room struct {
	id    string
	Title string
	Root  components.Native

	listeners []*websocket.Conn
}

func NewRoom(id, title string, root components.Native) *Room {
	return &Room{
		id:    id,
		Title: title,
		Root:  root,

		listeners: make([]*websocket.Conn, 0),
	}
}

func (r *Room) Add(conn *websocket.Conn) {
	r.listeners = append(r.listeners, conn)
}

func (r *Room) Remove(conn *websocket.Conn) {
	for i, c := range r.listeners {
		if c == conn {
			r.listeners = append(r.listeners[:i], r.listeners[i+1:]...)
			break
		}
	}
}

var _ components.RenderContext = (*Room)(nil)

func (r *Room) SocketURL() string {
	return fmt.Sprintf("/r/%s/service", url.PathEscape(r.id))
}

func (r *Room) EventURL(ev string) string {
	return fmt.Sprintf("/r/%s/event/%s", url.PathEscape(r.id), url.PathEscape(ev))
}

func getCoreRooms() map[string]*Room {
	return map[string]*Room{
		"home": getCoreRoomHome(),
	}
}

func getCoreRoomHome() *Room {
	return NewRoom("home", "CCWSUI!",
		components.AtRenderTime(func() components.Native {
			return components.AlignedCenter(
				components.Textured("stockkeeper", 20, 21, 17, 21,
					components.VStacked(
						components.AlignedCenter(
							components.Padded(3, 20, 3, 20,
								components.LiteralOf("Stock Keeper").
									WithColor("#000000"),
							)),

						components.Padded(0, 20, 16, 20,
							components.Clickable("awa", components.LiteralOf(strconv.Itoa(rand.Int()))),
						),

						components.Padded(0, 22, 18, 22, components.VStacked(
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots()),
						),
					),
				).WithPadBorder(false),
			)
		}))
}

func nineSlots() components.Native {
	return components.HStacked(slot(), slot(), slot(),
		slot(), slot(), slot(),
		slot(), slot(), slot())
}

func slot() components.Native {
	items, _ := staticFS.ReadDir("static/item")
	item := items[rand.Intn(len(items))]
	count := rand.Intn(64) + 1

	return components.VStacked(
		components.Textured("plain-inset", 1, 1, 1, 1,
			components.ItemTextured(item.Name()[:len(item.Name())-4])),
		components.AbsoluteAt(-2, -2, components.LiteralOf(
			strconv.Itoa(count)).WithShadow()),
	)
}
