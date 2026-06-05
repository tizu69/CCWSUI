package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"sync"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Room struct {
	ID       string
	hostConn *websocket.Conn
	Title    string
	Root     components.Native

	frozen     bool
	wantedSlug string

	clients   map[uuid.UUID]*roomClient
	clientsMu sync.Mutex
}

func NewRoom(hostConn *websocket.Conn, title string, root components.Native) *Room {
	return &Room{
		hostConn: hostConn,
		Title:    title,
		Root:     root,
		clients:  make(map[uuid.UUID]*roomClient),
	}
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
}

func (r *Room) HandleEvent(id uuid.UUID, ev webmsg.Event) {
	if r.hostConn != nil {
		if err := sendHostMsg(r.hostConn, HostMsgEvent, HostEventPayload{
			Client: id, Event: ev.ID, Data: ev.Event,
		}); err != nil {
			slog.Error("failed to send event message to host", "err", err)
		}
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

func (r *Room) sendUpdate(conn *websocket.Conn) error {
	root, err := r.Root.ToWire()
	if err != nil {
		return err
	}

	// TODO: host-submitted textures should be sent here
	if err := webmsg.SendMsg(conn, webmsg.TypeTexture, webmsg.Texture{
		ID:   "externallyloaded",
		Data: "iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAQAAADYv8WvAAAAAXNSR0IArs4c6QAAABJJREFUCJlj9Pz/kYHpIwM/AwAVIAM98dt1zAAAAABJRU5ErkJggg==",
	}); err != nil {
		return err
	}

	return webmsg.SendMsg(conn, webmsg.TypeUpdate, webmsg.Update{Root: root})
}

type roomClient struct {
	conn   *websocket.Conn
	UserID uuid.UUID
}

func (c *roomClient) Update(root components.WireNode) error {
	return webmsg.SendMsg(c.conn, webmsg.TypeUpdate, webmsg.Update{Root: root})
}

func (c *roomClient) Close() error {
	return c.conn.Close(websocket.StatusAbnormalClosure, "Host left")
}
