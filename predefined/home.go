package predefined

import (
	"encoding/json"
	"runtime/debug"
	"time"

	"g.tizu.dev/CCWSUI/components"
	"github.com/google/uuid"
	"github.com/mergestat/timediff"
)

type Home struct {
	Updater Updater
}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) Event(client uuid.UUID, id string, event json.RawMessage) {
	switch id {
	case "docs":
		h.Updater.Redirect(client, "/r/docs")
	}
}

func (h *Home) Hello(client uuid.UUID, user uuid.UUID) {
	h.Updater.Update(client, h.buildUI())
}

func (h *Home) buildUI() components.Native {
	return components.MkStackV(
		components.MkTexture("#202020",
			components.MkAlignCenter(
				components.MkPadding(8, 2, 8, 2,
					components.MkLiteral("CCWSUI").
						WithText(" Beta").WithHexColor("#aaa"),
				),
			),
		),

		components.MkExpand(
			components.MkAlignCenter(
				components.MkConstrain(200, 0,
					components.MkPadding(10, 10, 10, 10,
						components.MkStackV(
							components.MkAlignCenter(
								components.MkTexture("#fff;rounded=8", h.buildExample()).WithPad(),
							),
							components.MkAlignX(components.AlignmentCenter,
								components.MkLiteral("UIs that feel like home."),
							),

							components.MkLiteral("CCWSUI offers a simple, declarative way to build pixel-perfect web UIs from within ComputerCraft.").
								WithHexColor("#aaa").WithWrap().WithAlignment(components.AlignmentCenter),

							components.MkAlignX(components.AlignmentCenter,
								components.MkClickRegion("docs",
									components.MkTexture("plain-outset",
										components.MkPadding(4, 4, 4, 4,
											components.MkStackH(
												components.MkFiller(4, 0),
												components.MkAlignY(components.AlignmentCenter,
													components.MkLiteral("Docs").WithShadow(),
												),
												components.MkIcon("chevron").WithRotation(90).WithShadow(),
											),
										),
									),
								),
							),
						).WithGap(10),
					),
				),
			),
		),

		components.MkAlignX(components.AlignmentCenter,
			components.MkLiteral(buildinfo).WithHexColor("#444"),
		),
	)
}

var buildinfo = "dev"

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	var rev string
	var committime time.Time
	var modified bool
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value[:7]
		case "vcs.time":
			committime, _ = time.Parse(time.RFC3339, s.Value)
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}

	if modified {
		rev += "*"
	}
	buildinfo = "build " + rev + " (" + timediff.TimeDiff(committime) + ")"
}

func (Home) buildExample() components.Native {
	return components.MkStackV(
		components.MkTexture("title",
			components.MkPadding(2, 10, 2, 10,
				components.MkAlignCenter(
					components.MkLiteral("Thermostat"),
				),
			),
		).WithTintHex("#888b79"),

		components.MkPadding(-1, 2, -1, 2,
			components.MkTexture("#000000",
				components.MkPadding(1, 1, 1, 1,
					components.MkTexture("sided-inset",
						components.MkTexture("checkerboard",
							components.MkPadding(-1, 12, -1, 12,
								components.MkTexture("content-inset",
									components.MkPadding(6, 8, 6, 8,
										components.MkStackH(
											components.MkAlignY(components.AlignmentCenter,
												components.MkLiteral("Target"),
											),
											components.MkExpandH(components.MkBlank()),
											components.MkTexture("lineguide-inset",
												components.MkPadding(-1, 5, -1, 2,
													components.MkLiteral("22°C").
														WithShadow(),
												),
											).WithPad(),
										).WithGap(2),
									),
								),
							),
						),
					).WithPad(),
				),
			),
		),

		components.MkTexture("title",
			components.MkPadding(6, 7, 6, 7,
				components.MkStackH(
					components.MkTexture("plain-outset",
						components.MkIcon("trash"),
					).WithPad(),
					components.MkExpandH(components.MkBlank()),
					components.MkPadding(-4, 0, -4, 0,
						components.MkTexture("title-inset",
							components.MkBlank(),
						).WithPad(),
					),
					components.MkTexture("plain-outset",
						components.MkIcon("check"),
					).WithPad(),
				).WithGap(4).WithAlign(1),
			),
		),
	)
}

func (h *Home) Leave(client uuid.UUID) {
}

func (h *Home) SetUpdater(updater Updater) { h.Updater = updater }

var _ Room = (*Home)(nil)
