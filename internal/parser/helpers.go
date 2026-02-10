package parser

import (
	"strings"

	"zork-basic/internal/ast"
)

// builtinFunctions 是内置函数名称的集合
var builtinFunctions = map[string]bool{
	// 数学函数
	"ABS": true,
	"SIN": true,
	"COS": true,
	"TAN": true,
	"INT": true,
	"SQR": true,
	"LOG": true,
	"EXP": true,
	"RND": true,
	// 字符串函数
	"LEN":    true,
	"LEFT$":  true,
	"RIGHT$": true,
	"MID$":   true,
	"INSTR":  true,
	"UCASE$": true,
	"LCASE$": true,
	"SPACE$": true,
	"CHR$":   true,
	"ASC":    true,
}

// isBuiltinFunction 检查标识符是否是内置函数
func isBuiltinFunction(name string) bool {
	return builtinFunctions[strings.ToUpper(name)]
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
		// Sequence structure: [ ]* ('+' / '-' / KW_MOD) [ ]* Right:Power
		// When [ ]* matches zero times, the sequence has 2 elements: [op, right]
		// When [ ]* matches one or more times, the sequence has 4 elements: [spaces, op, spaces, right]
		var op string
		var right ast.Node
		if len(seq) == 2 {
			op = extractOpString(seq[0])
			right = seq[1].(ast.Node)
		} else if len(seq) >= 4 {
			op = extractOpString(seq[1])
			right = seq[3].(ast.Node)
		}
		result = &ast.BinaryOp{Left: result, Op: strings.ToUpper(op), Right: right}
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
		// Sequence structure: [ ]* KW_OR/KW_AND [ ]* Right:*
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

// extractOpString 从 pigeon 返回的各种类型中提取操作符字符串
// 支持 []byte（字面量匹配）、string、[]interface{}（规则引用如 KW_MOD）等类型
func extractOpString(v interface{}) string {
	switch t := v.(type) {
	case []byte:
		return string(t)
	case string:
		return t
	case []interface{}:
		// KW_MOD 等关键字规则返回 []interface{}，递归提取第一个文本元素
		var result string
		for _, item := range t {
			result += extractOpString(item)
		}
		return result
	default:
		return ""
	}
}
