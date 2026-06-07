package predefined

import (
	"encoding/json"

	"g.tizu.dev/CCWSUI/components"
	"github.com/google/uuid"
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

							components.MkLiteral("CCWSUI offers a simple, declarative way to build pixel-perfect web UIs, controlled entirely from within your ComputerCraft* code.").
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
			components.MkPadding(4, 4, 4, 4,
				components.MkLiteral("* or anything, really. If it can connect to a WebSocket, it can serve a CCWSUI!").
					WithHexColor("#555").WithWrap().WithAlignment(components.AlignmentCenter),
			),
		),
	)
}

func (Home) buildExample() components.Native {
	return components.MkStackV(
		components.MkTexture("title",
			components.MkPadding(3, 10, 3, 10,
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
												components.MkPadding(0, 5, 0, 2,
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
