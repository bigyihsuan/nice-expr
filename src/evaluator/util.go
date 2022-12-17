package evaluator

//go:generate stringer -type=IdentifierType
type IdentifierType int

const (
	InvalidIdentifier IdentifierType = iota - 1
	ConstantIdentifier
	VariableIdentifier
	FunctionIdentifier
)

var (
	BuiltinFunctionNames = []string{"print", "println", "len"}
)
