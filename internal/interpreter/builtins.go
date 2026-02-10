package interpreter

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"zork-basic/internal/ast"
)

// builtinFunc 内置函数类型
type builtinFunc func(i *Interpreter, node *ast.FunctionCall) Value

// builtinFuncs 内置函数分发表（O(1) map 查找替代 switch 线性匹配）
var builtinFuncs map[string]builtinFunc

func init() {
	builtinFuncs = map[string]builtinFunc{
		// 数学函数（单参数通用包装）
		"ABS": mathBuiltin1("ABS", math.Abs),
		"SIN": mathBuiltin1("SIN", math.Sin),
		"COS": mathBuiltin1("COS", math.Cos),
		"TAN": mathBuiltin1("TAN", math.Tan),
		"INT": mathBuiltin1("INT", math.Trunc),
		"EXP": mathBuiltin1("EXP", math.Exp),
		// 数学函数（需特殊处理）
		"SQR": (*Interpreter).builtinSQR,
		"LOG": (*Interpreter).builtinLOG,
		"RND": (*Interpreter).builtinRND,
		// 字符串函数
		"LEN":    (*Interpreter).builtinLEN,
		"LEFT$":  (*Interpreter).builtinLEFT,
		"RIGHT$": (*Interpreter).builtinRIGHT,
		"MID$":   (*Interpreter).builtinMID,
		"INSTR":  (*Interpreter).builtinINSTR,
		"UCASE$": (*Interpreter).builtinUCASE,
		"LCASE$": (*Interpreter).builtinLCASE,
		"SPACE$": (*Interpreter).builtinSPACE,
		"CHR$":   (*Interpreter).builtinCHR,
		"ASC":    (*Interpreter).builtinASC,
		// 常量支持
		"PI":    (*Interpreter).builtinPI,
		"EULER": (*Interpreter).builtinEULER,
	}
}

// mathBuiltin1 创建单参数数学函数的通用包装
// 减少重复的参数校验代码
func mathBuiltin1(name string, fn func(float64) float64) builtinFunc {
	return func(i *Interpreter, node *ast.FunctionCall) Value {
		if len(node.Args) != 1 {
			fmt.Fprintf(i.errOutput, "Error: %s requires 1 argument, got %d\n", name, len(node.Args))
			return NumberValue(0)
		}
		return NumberValue(fn(i.evaluateExpr(node.Args[0]).AsNumber()))
	}
}

func (i *Interpreter) builtinSQR(node *ast.FunctionCall) Value {
	if len(node.Args) != 1 {
		fmt.Fprintf(i.errOutput, "Error: SQR requires 1 argument, got %d\n", len(node.Args))
		return NumberValue(0)
	}
	value := i.evaluateExpr(node.Args[0]).AsNumber()
	if value < 0 {
		fmt.Fprintln(i.errOutput, "Error: SQR of negative number")
		return NumberValue(0)
	}
	return NumberValue(math.Sqrt(value))
}

func (i *Interpreter) builtinLOG(node *ast.FunctionCall) Value {
	if len(node.Args) != 1 {
		fmt.Fprintf(i.errOutput, "Error: LOG requires 1 argument, got %d\n", len(node.Args))
		return NumberValue(0)
	}
	value := i.evaluateExpr(node.Args[0]).AsNumber()
	if value <= 0 {
		fmt.Fprintln(i.errOutput, "Error: LOG of non-positive number")
		return NumberValue(0)
	}
	return NumberValue(math.Log(value))
}

func (i *Interpreter) builtinRND(node *ast.FunctionCall) Value {
	if len(node.Args) != 0 {
		fmt.Fprintf(i.errOutput, "Error: RND requires 0 arguments, got %d\n", len(node.Args))
		return NumberValue(0)
	}
	return NumberValue(rand.Float64())
}

func (i *Interpreter) builtinLEN(node *ast.FunctionCall) Value {
	if len(node.Args) != 1 {
		fmt.Fprintf(i.errOutput, "Error: LEN requires 1 argument, got %d\n", len(node.Args))
		return NumberValue(0)
	}
	return NumberValue(float64(len(i.evaluateExpr(node.Args[0]).String())))
}

func (i *Interpreter) builtinLEFT(node *ast.FunctionCall) Value {
	if len(node.Args) != 2 {
		fmt.Fprintf(i.errOutput, "Error: LEFT$ requires 2 arguments, got %d\n", len(node.Args))
		return StringValue("")
	}
	str := i.evaluateExpr(node.Args[0]).String()
	n := int(i.evaluateExpr(node.Args[1]).AsNumber())
	if n > len(str) {
		n = len(str)
	}
	if n < 0 {
		n = 0
	}
	return StringValue(str[:n])
}

func (i *Interpreter) builtinRIGHT(node *ast.FunctionCall) Value {
	if len(node.Args) != 2 {
		fmt.Fprintf(i.errOutput, "Error: RIGHT$ requires 2 arguments, got %d\n", len(node.Args))
		return StringValue("")
	}
	str := i.evaluateExpr(node.Args[0]).String()
	n := int(i.evaluateExpr(node.Args[1]).AsNumber())
	if n > len(str) {
		n = len(str)
	}
	if n < 0 {
		n = 0
	}
	return StringValue(str[len(str)-n:])
}

func (i *Interpreter) builtinMID(node *ast.FunctionCall) Value {
	if len(node.Args) < 2 || len(node.Args) > 3 {
		fmt.Fprintf(i.errOutput, "Error: MID$ requires 2 or 3 arguments, got %d\n", len(node.Args))
		return StringValue("")
	}
	str := i.evaluateExpr(node.Args[0]).String()
	start := int(i.evaluateExpr(node.Args[1]).AsNumber())
	if start < 1 {
		start = 1
	}
	n := len(str) - start + 1
	if len(node.Args) == 3 {
		n = int(i.evaluateExpr(node.Args[2]).AsNumber())
	}
	startIdx := start - 1
	endIdx := startIdx + n
	if endIdx > len(str) {
		endIdx = len(str)
	}
	if startIdx >= len(str) || startIdx < 0 {
		return StringValue("")
	}
	return StringValue(str[startIdx:endIdx])
}

func (i *Interpreter) builtinINSTR(node *ast.FunctionCall) Value {
	if len(node.Args) < 2 || len(node.Args) > 3 {
		fmt.Fprintf(i.errOutput, "Error: INSTR requires 2 or 3 arguments, got %d\n", len(node.Args))
		return NumberValue(0)
	}
	var start int = 1
	var str, substr string
	if len(node.Args) == 2 {
		str = i.evaluateExpr(node.Args[0]).String()
		substr = i.evaluateExpr(node.Args[1]).String()
	} else {
		start = int(i.evaluateExpr(node.Args[0]).AsNumber())
		str = i.evaluateExpr(node.Args[1]).String()
		substr = i.evaluateExpr(node.Args[2]).String()
	}
	if start < 1 {
		start = 1
	}
	pos := strings.Index(str[start-1:], substr)
	if pos == -1 {
		return NumberValue(0)
	}
	return NumberValue(float64(start + pos))
}

func (i *Interpreter) builtinUCASE(node *ast.FunctionCall) Value {
	if len(node.Args) != 1 {
		fmt.Fprintf(i.errOutput, "Error: UCASE$ requires 1 argument, got %d\n", len(node.Args))
		return StringValue("")
	}
	return StringValue(strings.ToUpper(i.evaluateExpr(node.Args[0]).String()))
}

func (i *Interpreter) builtinLCASE(node *ast.FunctionCall) Value {
	if len(node.Args) != 1 {
		fmt.Fprintf(i.errOutput, "Error: LCASE$ requires 1 argument, got %d\n", len(node.Args))
		return StringValue("")
	}
	return StringValue(strings.ToLower(i.evaluateExpr(node.Args[0]).String()))
}

func (i *Interpreter) builtinSPACE(node *ast.FunctionCall) Value {
	if len(node.Args) != 1 {
		fmt.Fprintf(i.errOutput, "Error: SPACE$ requires 1 argument, got %d\n", len(node.Args))
		return StringValue("")
	}
	n := int(i.evaluateExpr(node.Args[0]).AsNumber())
	if n < 0 {
		n = 0
	}
	return StringValue(strings.Repeat(" ", n))
}

func (i *Interpreter) builtinCHR(node *ast.FunctionCall) Value {
	if len(node.Args) != 1 {
		fmt.Fprintf(i.errOutput, "Error: CHR$ requires 1 argument, got %d\n", len(node.Args))
		return StringValue("")
	}
	code := int(i.evaluateExpr(node.Args[0]).AsNumber())
	if code < 0 || code > 255 {
		fmt.Fprintf(i.errOutput, "Error: CHR$ argument must be between 0 and 255, got %d\n", code)
		return StringValue("")
	}
	return StringValue(string(rune(code)))
}

func (i *Interpreter) builtinASC(node *ast.FunctionCall) Value {
	if len(node.Args) != 1 {
		fmt.Fprintf(i.errOutput, "Error: ASC requires 1 argument, got %d\n", len(node.Args))
		return NumberValue(0)
	}
	str := i.evaluateExpr(node.Args[0]).String()
	if len(str) == 0 {
		fmt.Fprintln(i.errOutput, "Error: ASC argument is an empty string")
		return NumberValue(0)
	}
	return NumberValue(float64(str[0]))
}

// builtinPI 返回 π 值
func (i *Interpreter) builtinPI(node *ast.FunctionCall) Value {
	if len(node.Args) != 0 {
		fmt.Fprintf(i.errOutput, "Error: PI requires 0 arguments, got %d\n", len(node.Args))
		return NumberValue(0)
	}
	return NumberValue(math.Pi)
}

// builtinEULER 返回自然常数 e 值
func (i *Interpreter) builtinEULER(node *ast.FunctionCall) Value {
	if len(node.Args) != 0 {
		fmt.Fprintf(i.errOutput, "Error: EULER requires 0 arguments, got %d\n", len(node.Args))
		return NumberValue(0)
	}
	return NumberValue(math.E)
}
