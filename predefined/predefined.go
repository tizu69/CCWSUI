// Package predefined contains predefined rooms.
package predefined

import (
	"encoding/json"

	"g.tizu.dev/CCWSUI/components"
	"github.com/google/uuid"
)

type Room interface {
	Event(client uuid.UUID, id string, event json.RawMessage)
	Hello(client, user uuid.UUID)
	Leave(client uuid.UUID)
	SetUpdater(updater Updater)
}

type Updater interface {
	Update(client uuid.UUID, root components.Native)
	Redirect(client uuid.UUID, url string)
}
