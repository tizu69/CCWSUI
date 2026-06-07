package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"sync"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/predefined"
	"g.tizu.dev/CCWSUI/web/webmsg"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Room struct {
	ID       string
	hostConn *websocket.Conn
	predef   predefined.Room

	frozen     bool
	wantedSlug string

	clients   map[uuid.UUID]*roomClient
	clientsMu sync.Mutex
}

func NewRemoteRoom(hostConn *websocket.Conn) *Room {
	return &Room{
		hostConn: hostConn,
		clients:  make(map[uuid.UUID]*roomClient),
	}
}

func NewPredefinedRoom(predef predefined.Room, slug string) *Room {
	r := &Room{
		ID:      slug,
		predef:  predef,
		clients: make(map[uuid.UUID]*roomClient),
		frozen:  true,
	}
	predef.SetUpdater(r)
	return r
}

func (r *Room) Frozen() bool { return r.frozen }

func (r *Room) Add(conn *websocket.Conn, userid uuid.UUID) (uuid.UUID, error) {
	if !r.frozen {
		return uuid.Nil, fmt.Errorf("room is not yet accepting clients")
	}

	if err := r.sendUpdate(conn); err != nil {
		conn.Close(websocket.StatusInternalError, err.Error())
		return uuid.Nil, err
	}

	id := uuid.New()
	r.clientsMu.Lock()
	r.clients[id] = &roomClient{conn: conn, UserID: userid}
	r.clientsMu.Unlock()

	if r.hostConn != nil {
		if err := sendHostMsg(r.hostConn, HostMsgHello, HostHelloPayload{
			Client: id,
			User:   userid,
		}); err != nil {
			return uuid.Nil, err
		}
	}
	if r.predef != nil {
		r.predef.Hello(id, userid)
	}

	return id, nil
}

func (r *Room) Remove(id uuid.UUID) {
	r.clientsMu.Lock()
	delete(r.clients, id)
	r.clientsMu.Unlock()

	if r.hostConn != nil {
		if err := sendHostMsg(r.hostConn, HostMsgLeave, HostLeavePayload{Client: id}); err != nil {
			slog.Error("failed to send leave message to host", "err", err)
		}
	}
	if r.predef != nil {
		r.predef.Leave(id)
	}
}

func (r *Room) HandleEvent(id uuid.UUID, ev webmsg.Event) {
	if r.hostConn != nil {
		if err := sendHostMsg(r.hostConn, HostMsgEvent, HostEventPayload{
			Client: id, Event: ev.ID, Data: ev.Event,
		}); err != nil {
			slog.Error("failed to send event message to host", "err", err)
		}
	}
	if r.predef != nil {
		r.predef.Event(id, ev.ID, ev.Event)
	}
}

func (r *Room) Get(id uuid.UUID) (*roomClient, bool) {
	r.clientsMu.Lock()
	c, ok := r.clients[id]
	r.clientsMu.Unlock()
	return c, ok
}

func (r *Room) SocketURL() string {
	return fmt.Sprintf("/r/%s/service", url.PathEscape(r.ID))
}

func (r *Room) UserURL() string {
	return fmt.Sprintf("/r/%s", url.PathEscape(r.ID))
}

var waitingForHost, _ = components.MkAlignCenter(
	components.MkTexture("tooltip",
		components.MkPadding(8, 8, 8, 8,
			components.MkLiteral("Waiting for Host"),
		),
	),
).ToWire()

func (r *Room) Update(client uuid.UUID, root components.Native) {
	wire, err := root.ToWire()
	if err != nil {
		slog.Error("failed to serialize root", "err", err)
		return
	}

	r.clientsMu.Lock()
	defer r.clientsMu.Unlock()
	if c, ok := r.clients[client]; ok {
		if err := c.Update(wire); err != nil {
			slog.Error("failed to send update message to client", "err", err)
		}
	}
}

func (r *Room) Redirect(client uuid.UUID, url string) {
	r.clientsMu.Lock()
	defer r.clientsMu.Unlock()
	if c, ok := r.clients[client]; ok {
		if err := webmsg.SendMsg(c.conn, webmsg.TypeRedirect, webmsg.Redirect{URL: url}); err != nil {
			slog.Error("failed to send redirect message to client", "err", err)
		}
	}
}

func (r *Room) sendUpdate(conn *websocket.Conn) error {
	// TODO: host-submitted textures should be sent here
	if err := webmsg.SendMsg(conn, webmsg.TypeTexture, webmsg.Texture{
		ID:   "externallyloaded",
		Data: "iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAQAAADYv8WvAAAAAXNSR0IArs4c6QAAABJJREFUCJlj9Pz/kYHpIwM/AwAVIAM98dt1zAAAAABJRU5ErkJggg==",
	}); err != nil {
		return err
	}

	return webmsg.SendMsg(conn, webmsg.TypeUpdate, webmsg.Update{Root: waitingForHost})
}

type roomClient struct {
	conn   *websocket.Conn
	UserID uuid.UUID
}

func (c *roomClient) Update(root components.WireNode) error {
	return webmsg.SendMsg(c.conn, webmsg.TypeUpdate, webmsg.Update{Root: root})
}

func (c *roomClient) UpdateMetadata(d HostMetadataPayload) error {
	return webmsg.SendMsg(c.conn, webmsg.TypeMetadata, webmsg.Metadata{
		Title: d.Title,
	})
}

func (c *roomClient) Close() error {
	return c.conn.Close(websocket.StatusAbnormalClosure, "Host left")
}
