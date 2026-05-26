package components

type Wrap string

const (
	YesWrap Wrap = "wrap"
	NoWrap  Wrap = "nowrap"
)

func (w Wrap) CSS() string {
	switch w {
	case YesWrap:
		return "wrap"
	default:
		return "nowrap"
	}
}

type StackAlign string

const (
	StackAlignStart  StackAlign = "start"
	StackAlignCenter StackAlign = "center"
	StackAlignEnd    StackAlign = "end"
)

func (a StackAlign) CSS() string {
	switch a {
	case StackAlignEnd:
		return "end"
	case StackAlignCenter:
		return "center"
	default:
		return "start"
	}
}
