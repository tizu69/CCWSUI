package predefined

import (
	"embed"
	"encoding/json"
	"path/filepath"
	"regexp"
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
		contents = []byte("Page not found: " + page)
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
	for compound := range strings.SplitSeq(strings.TrimSpace(string(contents)), "\n\n") {
		switch {
		case strings.HasPrefix(compound, "> "):
			l := components.MkLiteral("").WithWrap()
			color := "#AAAAAA"
			switch alert := strings.SplitN(compound, "\n> ", 2); alert[0] {
			case "> [!NOTE]":
				color = "#99B2F2"
				compound = "> " + alert[1]
				l = l.WithText("Note! ").WithHexColor(color)
			case "> [!CAUTION]":
				color = "#CC4C4C"
				compound = "> " + alert[1]
				l = l.WithText("Caution! ").WithHexColor(color)
			}
			txt := strings.ReplaceAll(compound[2:], "\n> ", " ")
			doc = doc.WithChildren(components.MkStackH(
				components.MkTexture(color, components.MkFiller(2, 0)),
				components.MkTexture(color+"20",
					components.MkPadding(4, 6, 4, 6,
						h.linkify(l, txt, "#aaa"),
					),
				),
			))
		case compound == "---":
			doc = doc.WithChildren(components.MkPadding(4, 16, 4, 16,
				components.MkTexture("#444", components.MkFiller(0, 1)),
			))
		case compound == "<ccwsui-download />":
			doc = doc.WithChildren(components.MkLiteral("TODO").WithHexColor("#ff0000"))
		default:
			txt := strings.ReplaceAll(compound, "\n", " ")
			doc = doc.WithChildren(h.linkify(components.MkLiteral(""), txt, "#f0f0f0").WithWrap())
		}
	}

	return components.MkStackH(
		components.MkTexture("#202020",
			components.MkPadding(12, 8, 12, 16,
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
				components.MkScroll("page:"+page, components.DirectionV,
					components.MkPadding(24, 8, 24, 8,
						components.MkConstrain(300, 0,
							components.MkOverlay(components.MkFiller(300, 0), doc),
						),
					),
				).WithStep(32),
			),
		),
	)
}

var linkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

func (h *Docs) linkify(l components.Literal, text, color string) components.Literal {
	const linkcolor = "#36c"
	lastEnd := 0
	for _, m := range linkRe.FindAllStringSubmatchIndex(text, -1) {
		if m[0] > lastEnd {
			l = l.WithText(text[lastEnd:m[0]]).WithHexColor(color)
		}
		l = l.WithText(text[m[2]:m[3]]).WithHexColor(linkcolor).WithClickEvent(text[m[4]:m[5]])
		lastEnd = m[1]
	}
	if lastEnd < len(text) {
		l = l.WithText(text[lastEnd:]).WithHexColor(color)
	}
	return l
}

func (h *Docs) Leave(client uuid.UUID) {
}

func (h *Docs) SetUpdater(updater Updater) { h.Updater = updater }

var _ Room = (*Docs)(nil)
