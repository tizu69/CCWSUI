package components

type Wrap bool

func (w Wrap) CSS() string {
	switch w {
	case true:
		return "wrap"
	default:
		return "nowrap"
	}
}

type StackDirection string

const (
	StackDirectionH StackDirection = "h"
	StackDirectionV StackDirection = "v"
)

func (a StackDirection) CSS() string {
	switch a {
	case StackDirectionH:
		return "row"
	default:
		return "column"
	}
}

type Alignment string

const (
	AlignmentStart  Alignment = "start"
	AlignmentCenter Alignment = "center"
	AlignmentEnd    Alignment = "end"
)

func (a Alignment) CSS() string {
	switch a {
	case AlignmentEnd:
		return "end"
	case AlignmentCenter:
		return "center"
	default:
		return "start"
	}
}

type TextSelect bool

func (w TextSelect) CSS() string {
	switch w {
	case true:
		return "text"
	default:
		return "none"
	}
}

type GreedyGrow bool

func (w GreedyGrow) CSS() string {
	switch w {
	case true:
		return "100%"
	default:
		return "auto"
	}
}
