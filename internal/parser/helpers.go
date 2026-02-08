package parser

import (
	"strings"

	"zork-basic/internal/ast"
)

// builtinFunctions 是内置函数名称的集合
var builtinFunctions = map[string]bool{
	// 数学函数
	"ABS":  true,
	"SIN":  true,
	"COS":  true,
	"TAN":  true,
	"INT":  true,
	"SQR":  true,
	"LOG":  true,
	"EXP":  true,
	"RND":  true,
	// 字符串函数
	"LEN":     true,
	"LEFT$":   true,
	"RIGHT$":  true,
	"MID$":    true,
	"INSTR":   true,
	"UCASE$":  true,
	"LCASE$":  true,
	"SPACE$":  true,
	"CHR$":    true,
	"ASC":     true,
}

// isBuiltinFunction 检查标识符是否是内置函数
func isBuiltinFunction(name string) bool {
	return builtinFunctions[strings.ToUpper(name)]
}

// toLineSlice converts a slice of interface{} to []*ast.Line
func toLineSlice(lines []interface{}) []*ast.Line {
	result := make([]*ast.Line, 0, len(lines))
	for _, line := range lines {
		if l, ok := line.(*ast.Line); ok {
			result = append(result, l)
		}
	}
	return result
}

// toLineSliceFromAny converts a slice of interface{} (from any) to []*ast.Line
func toLineSliceFromAny(lines any) []*ast.Line {
	if lines == nil {
		return []*ast.Line{}
	}
	slice, ok := lines.([]interface{})
	if !ok {
		return []*ast.Line{}
	}
	result := make([]*ast.Line, 0, len(slice))
	for _, line := range slice {
		if l, ok := line.(*ast.Line); ok {
			result = append(result, l)
		}
	}
	return result
}

// toNodeSlice converts a slice of interface{} to []ast.Node
func toNodeSlice(nodes []interface{}) []ast.Node {
	result := make([]ast.Node, 0, len(nodes))
	for _, node := range nodes {
		if n, ok := node.(ast.Node); ok {
			result = append(result, n)
		}
	}
	return result
}

// toNodeSliceFromAny converts a slice of interface{} (from any) to []ast.Node
func toNodeSliceFromAny(nodes any) []ast.Node {
	if nodes == nil {
		return []ast.Node{}
	}
	slice, ok := nodes.([]interface{})
	if !ok {
		return []ast.Node{}
	}
	result := make([]ast.Node, 0, len(slice))
	for _, node := range slice {
		if n, ok := node.(ast.Node); ok {
			result = append(result, n)
		}
	}
	return result
}

// buildBinaryOp is a helper function for building binary operations
// from PEG parser output (for additive, multiplicative, power)
func buildBinaryOp(left interface{}, rest interface{}) ast.Node {
	if rest == nil {
		return left.(ast.Node)
	}

	restSlice := rest.([]interface{})
	if len(restSlice) == 0 {
		return left.(ast.Node)
	}

	result := left.(ast.Node)
	for _, item := range restSlice {
		pair := item.([]interface{})
		op := string(pair[0].([]byte))
		right := pair[1].(ast.Node)
		result = &ast.BinaryOp{Left: result, Op: op, Right: right}
	}

	return result
}

// buildBinaryOpFromAny is a helper function for building binary operations
// from any (PEG parser output)
func buildBinaryOpFromAny(left interface{}, rest interface{}) ast.Node {
	if rest == nil {
		return left.(ast.Node)
	}

	restSlice, ok := rest.([]interface{})
	if !ok || len(restSlice) == 0 {
		return left.(ast.Node)
	}

	result := left.(ast.Node)
	for _, item := range restSlice {
		seq := item.([]interface{})
		// Sequence structure: [ ]* ('+' / '-') [ ]* Right:Multiplicative
		// When [ ]* matches zero times, the sequence has 2 elements: [op, right]
		// When [ ]* matches one or more times, the sequence has 4 elements: [spaces, op, spaces, right]
		var op string
		var right ast.Node
		if len(seq) == 2 {
			op = string(seq[0].([]byte))
			right = seq[1].(ast.Node)
		} else if len(seq) >= 4 {
			op = string(seq[1].([]byte))
			right = seq[3].(ast.Node)
		}
		result = &ast.BinaryOp{Left: result, Op: op, Right: right}
	}

	return result
}

// buildLogicalOp is a helper function for building logical operations
// from PEG parser output (for AND, OR)
func buildLogicalOp(left interface{}, rest interface{}, op string) ast.Node {
	if rest == nil {
		return left.(ast.Node)
	}

	restSlice := rest.([]interface{})
	if len(restSlice) == 0 {
		return left.(ast.Node)
	}

	result := left.(ast.Node)
	for _, item := range restSlice {
		right := item.(ast.Node)
		result = &ast.LogicalOp{Left: result, Op: op, Right: right}
	}

	return result
}

// buildLogicalOpFromAny is a helper function for building logical operations
// from any (PEG parser output)
func buildLogicalOpFromAny(left interface{}, rest interface{}, op string) ast.Node {
	if rest == nil {
		return left.(ast.Node)
	}

	restSlice, ok := rest.([]interface{})
	if !ok || len(restSlice) == 0 {
		return left.(ast.Node)
	}

	result := left.(ast.Node)
	for _, item := range restSlice {
		seq := item.([]interface{})
		// Sequence structure: [ ]* "OR"/"AND" [ ]* Right:*
		// When [ ]* matches zero times: 2 elements [keyword, right]
		// When [ ]* matches one or more times: 4 elements [spaces, keyword, spaces, right]
		var right ast.Node
		if len(seq) == 2 {
			right = seq[1].(ast.Node)
		} else if len(seq) >= 4 {
			right = seq[3].(ast.Node)
		}
		result = &ast.LogicalOp{Left: result, Op: op, Right: right}
	}

	return result
}
