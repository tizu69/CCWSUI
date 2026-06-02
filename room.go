package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"sync"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Room struct {
	conn  *websocket.Conn
	id    string
	Title string
	Root  components.Native

	clients   map[uuid.UUID]*roomClient
	clientsMu sync.Mutex
}

func NewRoom(conn *websocket.Conn, id, title string, root components.Native) (*Room, error) {
	r := &Room{
		conn: conn, id: id,
		Title: title, Root: root,
		clients: make(map[uuid.UUID]*roomClient),
	}
	if conn == nil {
		return r, nil
	}

	var err error
	var b []byte
	if b, err = json.Marshal(HostMessageReady{URL: r.UserURL()}); err != nil {
		return nil, err
	}
	if b, err = json.Marshal(HostMessageEnvelope{
		Type: HostMessageTypeReady, Data: b,
	}); err != nil {
		return nil, err
	}
	if err = r.conn.Write(context.TODO(), websocket.MessageText, b); err != nil {
		return nil, err
	}

	return r, err
}

func (r *Room) Add(conn *websocket.Conn) (uuid.UUID, error) {
	if err := r.sendUpdate(conn); err != nil {
		conn.Close(websocket.StatusInternalError, err.Error())
		return uuid.Nil, err
	}

	id := uuid.New()
	r.clientsMu.Lock()
	r.clients[id] = &roomClient{conn: conn}
	r.clientsMu.Unlock()

	if r.conn == nil {
		return id, nil
	}
	var err error
	var b []byte
	if b, err = json.Marshal(HostMessageHello{Client: id}); err != nil {
		return uuid.Nil, err
	}
	if b, err = json.Marshal(HostMessageEnvelope{
		Type: HostMessageTypeHello, Data: b,
	}); err != nil {
		return uuid.Nil, err
	}
	if err = r.conn.Write(context.TODO(), websocket.MessageText, b); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *Room) Remove(id uuid.UUID) {
	r.clientsMu.Lock()
	delete(r.clients, id)
	r.clientsMu.Unlock()
}

func (r *Room) Get(id uuid.UUID) (*roomClient, bool) {
	r.clientsMu.Lock()
	c, ok := r.clients[id]
	r.clientsMu.Unlock()
	return c, ok
}

func (r *Room) sendUpdate(conn *websocket.Conn) error {
	root, err := r.Root.ToWire()
	if err != nil {
		return err
	}

	var b []byte
	const data = "iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAQAAADYv8WvAAAAAXNSR0IArs4c6QAAABJJREFUCJlj9Pz/kYHpIwM/AwAVIAM98dt1zAAAAABJRU5ErkJggg=="
	if b, err = json.Marshal(webmsg.Texture{ID: "externallyloaded", Data: data}); err != nil {
		return err
	}
	if b, err = json.Marshal(webmsg.Envelope{
		Type: webmsg.TypeTexture, Data: b,
	}); err != nil {
		return err
	}
	if err = conn.Write(context.TODO(), websocket.MessageText, b); err != nil {
		return err
	}

	if b, err = json.Marshal(webmsg.Update{Root: root}); err != nil {
		return err
	}
	if b, err = json.Marshal(webmsg.Envelope{
		Type: webmsg.TypeUpdate, Data: b,
	}); err != nil {
		return err
	}
	if err = conn.Write(context.TODO(), websocket.MessageText, b); err != nil {
		return err
	}

	return nil
}

func (r *Room) SocketURL() string {
	return fmt.Sprintf("/r/%s/service", url.PathEscape(r.id))
}

func (r *Room) UserURL() string {
	return fmt.Sprintf("/r/%s", url.PathEscape(r.id))
}

func (r *Room) EventURL(ev string) string {
	return fmt.Sprintf("/r/%s/event/%s", url.PathEscape(r.id), url.PathEscape(ev))
}

type roomClient struct {
	conn *websocket.Conn
	root components.WireNode
}

func (r *roomClient) Update(root components.WireNode) error {
	var err error
	var b []byte
	if b, err = json.Marshal(webmsg.Update{
		Root: root,
	}); err != nil {
		return err
	}
	if b, err = json.Marshal(webmsg.Envelope{
		Type: webmsg.TypeUpdate, Data: b,
	}); err != nil {
		return err
	}
	if err = r.conn.Write(context.TODO(), websocket.MessageText, b); err != nil {
		return err
	}
	r.root = root
	return nil
}

func getCoreRooms() map[string]*Room {
	return map[string]*Room{
		"home": getCoreRoomHome(),
	}
}

func getCoreRoomHome() *Room {
	r, _ := NewRoom(nil, "home", "CCWSUI!",
		components.Overlaid(
			components.Expand(components.Textured("background", components.Filler())),
			components.HStacked(
				stockkeeper("stockkeeper"),
				stockkeeper("stockkeeper2"),
				components.Expand(
					components.AlignedCenter(
						components.Constrained(200, 0,
							components.VStacked(
								components.Textured("title",
									components.Padded(3, 20, 3, 20,
										components.AlignedCenter(
											components.LiteralOf("Stock Keeper").
												WithHexColor("#000000")),
									),
								).TintedHex("#cba6f7"),

								components.Padded(0, 2, 0, 2,
									components.Textured("#000",
										components.Padded(0, 1, 0, 1,
											components.Textured("sided-inset",
												components.Textured("checkerboard",
													components.Padded(6, 18, 6, 18,
														components.Textured("content-inset",
															components.Padded(8, 8, 8, 8,
																components.LiteralOf("Stock Keeper").WithHexColor("#cba6f7").WithShadow(true).
																	Add(" is a simple inventory management system that allows you to easily").
																	Add(" keep track of all your items!").WithHexColor("#cba6f7").
																	WithWrap(true).WithAlignment(components.AlignmentCenter),
															)))),
											).SetPad(true)))),

								components.Textured("title",
									components.Padded(2, 2, 2, 2,
										components.AlignedX(components.AlignmentEnd,
											components.HStacked(
												components.Overlaid(
													components.Textured("plain-outset",
														components.Icony("cross").Flip(components.FlipX).TintedHex("#f38ba8"),
													).SetPad(true),
													components.FollowsMouse(components.AlignmentStart, components.AlignmentEnd,
														components.Padded(0, 3, 0, 3,
															components.Textured("tooltip",
																components.Padded(3, 4, 4, 4,
																	components.LiteralOf("Abort now"),
																),
															).TintedHex("#f38ba8"),
														),
													).FlipIfOverflowing(true),
												).RequireHover(true),
											),
										),
									),
								).SetPad(true),
							),
						))),
			).WithPadding(8)))
	return r
}

func stockkeeper(id string) components.Native {
	return components.AlignedCenter(
		components.Constrained(0, 197, components.Textured("stockkeeper",
			components.VStacked(
				components.Padded(3, 20, 4, 20,
					components.AlignedCenter(
						components.LiteralOf("Stock Keeper").
							WithHexColor("#000000")),
				),

				components.ExpandV(
					components.Padded(1, 22, 18, 22, components.Scrolling(
						id, components.DirectionV, components.VStacked(
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots(),
							nineSlots()),
					).SetStep(18))),
			))),
	)
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

	return components.Overlaid(
		components.Clickable("slot",
			components.Overlaid(
				components.Textured("plain-inset", components.Padded(1, 1, 1, 1,
					components.ItemTextured(item.Name()[:len(item.Name())-4]))),
				components.Aligned(components.AlignmentEnd, components.AlignmentEnd,
					components.LiteralOf(strconv.Itoa(count)).WithShadow(true)),
			),
		),
		components.FollowsMouse(components.AlignmentStart, components.AlignmentEnd,
			components.Padded(0, 3, 0, 3,
				components.Textured("tooltip",
					components.Padded(3, 4, 4, 4,
						components.LiteralOf(item.Name()[:len(item.Name())-4]+"\n"+"Count: "+strconv.Itoa(count)).
							WithWrap(true).WithShadow(true)),
				),
			),
		).FlipIfOverflowing(true),
	).RequireHover(true)
}
