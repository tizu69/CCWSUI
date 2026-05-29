package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
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
	if err := r.sendUpdate(conn); err != nil {
		conn.Close(websocket.StatusInternalError, err.Error())
		return
	}
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

func (r *Room) sendUpdate(conn *websocket.Conn) error {
	root, err := r.Root.ToWire()
	if err != nil {
		return err
	}

	var b []byte
	if b, err = json.Marshal(webmsg.Update{Root: root}); err != nil {
		return err
	}
	if b, err = json.Marshal(webmsg.Envelope{
		Type: webmsg.TypeUpdate, Data: b,
	}); err != nil {
		return err
	}
	return conn.Write(context.TODO(), websocket.MessageText, b)
}

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
		components.Overlaid(
			components.Expand(components.Textured("background", 0, 0, 0, 0, components.Filler())),
			components.Padded(8, 8, 8, 8, components.HStacked(
				components.Expand(
					components.AlignedCenter(
						components.Constrained(100, 100,
							components.Textured("stockkeeper", 20, 21, 17, 21,
								components.VStacked(
									components.Padded(3, 20, 3, 20,
										components.AlignedCenter(
											components.LiteralOf("Stock Keeper").
												WithHexColor("#000000")),
									),

									components.Expand(components.Blanked(1, 1)),
								),
							)))),
				components.Blanked(8, 0),
				components.Textured("stockkeeper", 20, 21, 17, 21,
					components.VStacked(
						components.Padded(3, 20, 4, 20,
							components.AlignedCenter(
								components.LiteralOf("Stock Keeper").
									WithHexColor("#000000")),
						),

						components.Padded(1, 22, 18, 22, components.VStacked(
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots()),
						),
					),
				),
			))))
	// components.AtRenderTime(func() components.Native {
	// 	return components.AlignedCenter(
	// 			components.VStacked(
	// 				components.AlignedCenter(
	// 					components.Padded(3, 20, 3, 20,
	// 						components.LiteralOf("Stock Keeper").
	// 							WithColor("#000000"),
	// 					)),

	// 				components.Padded(0, 20, 16, 20,
	// 					components.Clickable("awa", components.LiteralOf(strconv.Itoa(rand.Int()))),
	// 				),

	// 				components.Padded(0, 22, 18, 22, components.VStacked(
	// 					nineSlots(),
	// 					nineSlots(),
	// 					nineSlots(),
	// 					nineSlots()),
	// 				),

	// 				components.Padded(0, 22, 18, 22,
	// 					components.HStacked(
	// 						slot(),
	// 						components.Textured("lineguide-inset", 4, 1, 4, 1,
	// 							components.Padded(3, 3, 4, 3,
	// 								components.Inputable("search", "Foo Bar!")),
	// 						).WithPadBorder(false),
	// 					).WithPadding(4),
	// 				),
	// 			),
	// 		).WithPadBorder(false),
	// 	)
	// }))
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
	_ = count

	return components.VStacked(
		components.Textured("plain-inset", 1, 1, 1, 1, components.Padded(1, 1, 1, 1,
			components.ItemTextured(item.Name()[:len(item.Name())-4]))),
		// components.AbsoluteAt(-1, -1, components.LiteralOf(
		// 	strconv.Itoa(count)).WithShadow()),
	)
}
