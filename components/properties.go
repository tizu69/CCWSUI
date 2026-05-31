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

type Rotation int

const (
	RotationU Rotation = 0
	RotationR Rotation = 90
	RotationD Rotation = 180
	RotationL Rotation = 270
)

type Flip string

const (
	FlipNone Flip = ""
	FlipX    Flip = "x"
	FlipY    Flip = "y"
)
