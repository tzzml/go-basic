package bytecode

import "strings"

// BuiltinID represents the index of a builtin function
type BuiltinID int

// List of builtin function names in fixed order.
// This order MUST match the implementation array in the VM.
var BuiltinNames = []string{
	"ABS",
	"SIN",
	"COS",
	"TAN",
	"INT",
	"EXP",
	"SQR",
	"LOG",
	"RND",
	"LEN",
	"LEFT$",
	"RIGHT$",
	"MID$",
	"INSTR",
	"UCASE$",
	"LCASE$",
	"SPACE$",
	"CHR$",
	"ASC",
}

var builtinMap map[string]int

func init() {
	builtinMap = make(map[string]int)
	for i, name := range BuiltinNames {
		builtinMap[name] = i
	}
}

// GetBuiltinID returns the ID for a given builtin function name.
// Returns -1 if not found.
func GetBuiltinID(name string) int {
	if idx, ok := builtinMap[strings.ToUpper(name)]; ok {
		return idx
	}
	return -1
}
