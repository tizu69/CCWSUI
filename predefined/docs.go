package predefined

import (
	"embed"
	"encoding/json"
	"path/filepath"
	"strings"

	"g.tizu.dev/CCWSUI/components"
	"github.com/google/uuid"
)

//go:embed docs/*.md
var docsFS embed.FS
var docsPages []string

func init() {
	entries, _ := docsFS.ReadDir("docs")
	// the file ordering seems to be consistent, so we rely on that here. if it
	// isn't, we'll need to sort the entries by name, but thats for later me ig
	for _, entry := range entries {
		name := entry.Name()
		ref := strings.TrimSuffix(name, filepath.Ext(name))
		docsPages = append(docsPages, ref)
	}
}

type Docs struct {
	Updater     Updater
	clientpages map[uuid.UUID]string
}

func NewDocs() *Docs {
	return &Docs{
		clientpages: make(map[uuid.UUID]string),
	}
}

func (h *Docs) Event(client uuid.UUID, id string, event json.RawMessage) {
	switch {
	case strings.HasPrefix(id, "page:"):
		h.clientpages[client] = id[5:]
		h.Updater.Update(client, h.buildUI(h.clientpages[client]))
	}
}

func (h *Docs) Hello(client uuid.UUID, user uuid.UUID) {
	h.clientpages[client] = "000-Getting-Started"
	h.Updater.Update(client, h.buildUI(h.clientpages[client]))
}

func (h *Docs) buildUI(page string) components.Native {
	contents, _ := docsFS.ReadFile("docs/" + page + ".md")
	if contents == nil {
		contents = []byte("Page not found")
	}

	pages := make([]components.Native, 0)
	for _, p := range docsPages {
		l := components.MkLiteral(strings.ReplaceAll(p[4:], "-", " ")).WithWrap()
		if p != page {
			l = l.WithHexColor("#aaa")
		}
		pages = append(pages, components.MkClickRegion("page:"+p, l))
	}

	doc := components.MkStackV().WithGap(8)
	for _, compound := range strings.Split(string(contents[:len(contents)-1]), "\n\n") {
		switch {
		case strings.HasPrefix(compound, "> "):
			txt := strings.ReplaceAll(compound[2:], "\n> ", " ")
			doc = doc.WithChildren(
				components.MkStackH(
					components.MkTexture("#aaa", components.MkFiller(2, 0)),
					components.MkPadding(4, 8, 4, 8, components.MkLiteral(txt).WithWrap()),
				),
			)
		default:
			txt := strings.ReplaceAll(compound, "\n", " ")
			doc = doc.WithChildren(components.MkLiteral(txt).WithWrap())
		}
	}

	return components.MkStackH(
		components.MkTexture("#202020",
			components.MkPadding(12, 8, 12, 8,
				components.MkStackV(
					append([]components.Native{
						components.MkLiteral("CCWSUI").
							WithText(" Docs").WithHexColor("#aaa"),
						components.MkFiller(0, 8),
					}, pages...)...,
				).WithGap(2),
			),
		),

		components.MkExpandH(
			components.MkAlignX(components.AlignmentCenter,
				components.MkPadding(24, 8, 24, 8,
					components.MkConstrain(300, 0, doc),
				),
			),
		),
	)
}

func (h *Docs) Leave(client uuid.UUID) {
}

func (h *Docs) SetUpdater(updater Updater) { h.Updater = updater }

var _ Room = (*Docs)(nil)
