package components

type Direction string

const (
	DirectionH  Direction = "h"
	DirectionV  Direction = "v"
	DirectionHV Direction = "hv"
)

type StackDirection string

const (
	StackDirectionH StackDirection = "h"
	StackDirectionV StackDirection = "v"
)

type Alignment float32

const (
	AlignmentStart  Alignment = 0
	AlignmentCenter Alignment = .5
	AlignmentEnd    Alignment = 1
)
