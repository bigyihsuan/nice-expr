package blockkind

//go:generate stringer -type=BlockKind
type BlockKind int

const (
	FreeFloating BlockKind = iota
	If
	Else
	For
)
