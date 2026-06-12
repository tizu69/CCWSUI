package predefined

import (
	"embed"
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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

type docsUser struct {
	Page string
	Time time.Time
}

type Docs struct {
	Updater       Updater
	clientpages   map[uuid.UUID]string
	clientsidebar map[uuid.UUID]bool
	clientuser    map[uuid.UUID]uuid.UUID
	users         map[uuid.UUID]docsUser
}

func NewDocs() *Docs {
	return &Docs{
		clientpages:   make(map[uuid.UUID]string),
		clientsidebar: make(map[uuid.UUID]bool),
		clientuser:    make(map[uuid.UUID]uuid.UUID),
		users:         make(map[uuid.UUID]docsUser),
	}
}

func (h *Docs) Event(client uuid.UUID, id string, event json.RawMessage) {
	switch {
	case strings.HasPrefix(id, "page:"):
		h.clientpages[client] = id[5:]
		h.users[h.clientuser[client]] = docsUser{Page: id[5:], Time: time.Now()}
		delete(h.clientsidebar, client)
		h.Updater.Update(client, h.buildUI(h.clientpages[client], false))
	case id == "menu":
		h.clientsidebar[client] = true
		h.Updater.Update(client, h.buildUI(h.clientpages[client], true))
	}
}

func (h *Docs) Hello(client uuid.UUID, user uuid.UUID) {
	h.clientuser[client] = user
	h.clientpages[client] = "000-Getting-Started"
	if u, ok := h.users[user]; ok && time.Since(u.Time) < time.Hour {
		h.clientpages[client] = u.Page
	}
	h.users[h.clientuser[client]] = docsUser{Page: h.clientpages[client], Time: time.Now()}
	h.Updater.Update(client, h.buildUI(h.clientpages[client], false))
}

func (h *Docs) buildUI(page string, sidebaropen bool) components.Native {
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
		case strings.HasPrefix(compound, "```"):
			lang, code, _ := strings.Cut(compound[3:], "\n")
			lang, code = strings.TrimSpace(lang), strings.TrimSpace(code[:len(code)-3])
			_ = lang // TODO: syntax highlighting for other languages? maybe?
			doc = doc.WithChildren(components.MkTexture("#202020",
				components.MkScroll("codeblock:"+code, components.DirectionH,
					components.MkPadding(8, 8, 8, 8, highlightLua(code)),
				),
			))
		default:
			txt := strings.ReplaceAll(compound, "\n", " ")
			doc = doc.WithChildren(h.linkify(components.MkLiteral(""), txt, "#f0f0f0").WithWrap())
		}
	}

	sidebar := components.MkTexture("#202020",
		components.MkPadding(12, 8, 12, 16,
			components.MkStackV(
				append([]components.Native{
					components.MkLiteral("CCWSUI").
						WithText(" Docs").WithHexColor("#aaa"),
					components.MkFiller(0, 8),
				}, pages...)...,
			).WithGap(2),
		),
	)
	if sidebaropen {
		return sidebar
	}

	return components.MkStackH(
		components.MkMediaQuery(sidebar).WithMinWidth(400),

		components.MkExpandH(
			components.MkOverlay(
				components.MkAlignX(components.AlignmentCenter,
					components.MkScroll("page:"+page, components.DirectionV,
						components.MkPadding(24, 8, 24, 8,
							components.MkConstrain(300, 0,
								components.MkOverlay(components.MkFiller(300, 0), doc),
							),
						),
					),
				),
				components.MkAlign(components.AlignmentStart, components.AlignmentStart,
					components.MkMediaQuery(
						components.MkClickRegion("menu",
							components.MkPadding(2, 2, 2, 2,
								components.MkTexture("#202020;rounded=1",
									components.MkIcon("rows"),
								),
							),
						),
					).WithMaxWidth(399),
				),
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
	delete(h.clientuser, client)
	delete(h.clientpages, client)
	delete(h.clientsidebar, client)
}

func (h *Docs) SetUpdater(updater Updater) { h.Updater = updater }

var _ Room = (*Docs)(nil)
