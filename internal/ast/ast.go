package ast

import (
	"fmt"
	"strings"
)

// Node AST 节点接口
// 所有 AST 节点类型都必须实现此接口，提供字符串表示方法
type Node interface {
	String() string
}

// Program 表示一个完整的 BASIC 程序
// BASIC 程序由多行组成，每行都有行号
type Program struct {
	Lines      []*Line // 带行号的语句行，按行号排序
	Statements []Node  // 旧版兼容：用于向后兼容
}

// Line 表示 BASIC 程序中的一行
// 每行都有一个行号和零个或多个语句
type Line struct {
	LineNumber int    // 行号（10, 20, 30 等）
	Statements []Node // 该行包含的语句列表
}

// Assignment 表示变量赋值语句
// 语法: LET <变量名> = <表达式> 或 <变量名> = <表达式>
// 也支持数组元素赋值: <数组名>(<索引>) = <表达式>
type Assignment struct {
	Target Node   // 赋值目标（Identifier 或 ArrayAccess）
	Value  Node   // 要赋值的表达式
}

// PrintStmt 表示 PRINT 输出语句
// 语法: PRINT <表达式1>[,|;] <表达式2>[,|;] ... [;|,]
// 支持多个参数，用逗号或分号分隔
// 末尾的分隔符决定是否换行：分号或逗号表示不换行，无分隔符表示换行
type PrintStmt struct {
	Values     []Node   // 要输出的值列表
	Separators []string // 值之间的分隔符：";" 表示紧凑输出，"," 表示添加空格
	Trailer    string   // 末尾的分隔符：";", "," 或 ""
}

// IfStmt 表示 IF...THEN...ELSE...END IF 条件语句
// 语法: IF <条件> THEN <语句块> [ELSE <语句块>] END IF
type IfStmt struct {
	Condition  Node   // 条件表达式（比较或逻辑运算）
	ThenStmts  []Node // 条件为真时执行的语句块
	ElseStmts  []Node // 条件为假时执行的语句块（可选）
}

// ForStmt 表示 FOR...NEXT 循环语句
// 语法: FOR <变量> = <起始值> TO <结束值> [STEP <步长>]
type ForStmt struct {
	Var   string // 循环变量名
	Start Node   // 循环起始值
	End   Node   // 循环结束值
	Step  Node   // 循环步长（可选，nil 表示步长为 1）
}

// NextStmt 表示 NEXT 语句
// 用于终止 FOR 循环的一次迭代
// 语法: NEXT [<变量名>]
type NextStmt struct {
	Var string // 循环变量名（可选，空字符串表示未指定）
}

// GotoStmt 表示 GOTO 无条件跳转语句
// 语法: GOTO <行号>
type GotoStmt struct {
	LineNumber int // 要跳转到的目标行号
}

// GosubStmt 表示 GOSUB 子程序调用语句
// 语法: GOSUB <行号>
type GosubStmt struct {
	LineNumber int // 子程序开始的行号
}

// ReturnStmt 表示 RETURN 语句
// 用于从子程序返回
// 语法: RETURN
type ReturnStmt struct{}

// EndStmt 表示 END 程序结束语句
// 语法: END
type EndStmt struct{}

// RemStmt 表示 REM 注释语句
// 语法: REM <注释文本>
type RemStmt struct {
	Text string // 注释文本（包括 REM 关键字）
}

// DimStmt 表示 DIM 数组声明语句
// 语法: DIM <数组名>(<大小1>[, <大小2>, ...])
type DimStmt struct {
	Name string   // 数组名
	Sizes []Node  // 数组各维度的大小（表达式列表）
}

// InputStmt 表示 INPUT 输入语句
// 语法: INPUT ["提示字符串",] <变量名1>[, <变量名2>, ...]
type InputStmt struct {
	Prompt string   // 可选的提示字符串
	Vars   []string // 要接收输入的变量名列表（支持多个变量）
}

// BinaryOp 表示二元算术运算表达式
// 支持的运算符: +, -, *, /, ^
type BinaryOp struct {
	Left  Node   // 左操作数
	Op    string // 运算符 ("+", "-", "*", "/", "^")
	Right Node   // 右操作数
}

// ComparisonOp 表示比较运算表达式
// 支持的运算符: =, <>, >, <, >=, <=
type ComparisonOp struct {
	Left  Node   // 左操作数
	Op    string // 比较运算符 ("=", "<>", ">", "<", ">=", "<=")
	Right Node   // 右操作数
}

// LogicalOp 表示逻辑运算表达式
// 支持的运算符: AND, OR
type LogicalOp struct {
	Left  Node   // 左操作数
	Op    string // 逻辑运算符 ("AND", "OR")
	Right Node   // 右操作数
}

// UnaryOp 表示一元运算表达式
// 支持的运算符: +, -（正负号）
type UnaryOp struct {
	Op    string // 运算符 ("+", "-")
	Right Node   // 操作数
}

// Identifier 表示变量标识符
// 变量名由字母、数字和下划线组成，必须以字母或下划线开头
type Identifier struct {
	Name string // 变量名
}

// FunctionCall 表示函数调用
// 语法: <函数名>(<参数1>, <参数2>, ...)
type FunctionCall struct {
	Name string // 函数名
	Args []Node // 参数列表（表达式）
}

// ArrayAccess 表示数组访问
// 语法: <数组名>(<索引1>[, <索引2>, ...])
type ArrayAccess struct {
	Name     string   // 数组名
	Indices  []Node   // 索引表达式列表
}

// Number 表示数字字面量
// 支持整数和浮点数
type Number struct {
	Value float64 // 数字值
}

// StringLiteral 表示字符串字面量
// 语法: "<文本>"
type StringLiteral struct {
	Value string // 字符串值（不包含引号）
}

// String 返回程序的字符串表示
// 格式: "Program:\n<行号>: <语句>\n..."
func (p *Program) String() string {
	if len(p.Lines) > 0 {
		result := "Program:\n"
		for _, line := range p.Lines {
			result += fmt.Sprintf("%4d: ", line.LineNumber)
			for i, stmt := range line.Statements {
				if i > 0 {
					result += ": "
				}
				result += stmt.String()
			}
			result += "\n"
		}
		return result
	}
	// 旧版兼容：使用 Statements 字段
	result := "Program:\n"
	for _, stmt := range p.Statements {
		result += "  " + stmt.String() + "\n"
	}
	return result
}

// String 返回行的字符串表示
// 格式: "<行号>: <语句1>: <语句2>: ..."
func (l *Line) String() string {
	result := fmt.Sprintf("%4d: ", l.LineNumber)
	for i, stmt := range l.Statements {
		if i > 0 {
			result += ": "
		}
		result += stmt.String()
	}
	return result
}

// String 返回赋值语句的字符串表示
// 格式: "LET <变量名> = <值>" 或 "<数组名>(<索引>) = <值>"
func (a *Assignment) String() string {
	return fmt.Sprintf("LET %s = %s", a.Target.String(), a.Value.String())
}

// String 返回 PRINT 语句的字符串表示
// 格式: "PRINT <值1>, <值2>, ..." 或 "PRINT"
func (p *PrintStmt) String() string {
	if len(p.Values) == 0 {
		return "PRINT"
	}
	result := "PRINT"
	for i, v := range p.Values {
		if i > 0 {
			result += ","
		}
		result += " " + v.String()
	}
	return result
}

// String 返回 IF 语句的字符串表示
// 格式多行，包含 THEN 和 ELSE 块
func (i *IfStmt) String() string {
	thenStr := ""
	for _, stmt := range i.ThenStmts {
		thenStr += "    " + stmt.String() + "\n"
	}
	result := fmt.Sprintf("IF %s THEN\n%s", i.Condition.String(), thenStr)
	if len(i.ElseStmts) > 0 {
		elseStr := ""
		for _, stmt := range i.ElseStmts {
			elseStr += "    " + stmt.String() + "\n"
		}
		result += "ELSE\n" + elseStr
	}
	result += "END IF"
	return result
}

// String 返回 FOR 语句的字符串表示
// 格式: "FOR <变量> = <起始值> TO <结束值> [STEP <步长>]"
// 步长为 1 时省略 STEP 部分
func (f *ForStmt) String() string {
	result := fmt.Sprintf("FOR %s = %s TO %s", f.Var, f.Start.String(), f.End.String())
	if step, ok := f.Step.(*Number); ok && step.Value == 1 {
		// 默认步长为 1，不显示
	} else {
		result += fmt.Sprintf(" STEP %s", f.Step.String())
	}
	return result
}

// String 返回 NEXT 语句的字符串表示
// 格式: "NEXT <变量名>" 或 "NEXT"
func (n *NextStmt) String() string {
	if n.Var != "" {
		return fmt.Sprintf("NEXT %s", n.Var)
	}
	return "NEXT"
}

// String 返回 GOTO 语句的字符串表示
// 格式: "GOTO <行号>"
func (g *GotoStmt) String() string {
	return fmt.Sprintf("GOTO %d", g.LineNumber)
}

// String 返回 GOSUB 语句的字符串表示
// 格式: "GOSUB <行号>"
func (g *GosubStmt) String() string {
	return fmt.Sprintf("GOSUB %d", g.LineNumber)
}

// String 返回 RETURN 语句的字符串表示
// 格式: "RETURN"
func (r *ReturnStmt) String() string {
	return "RETURN"
}

// String 返回 END 语句的字符串表示
// 格式: "END"
func (e *EndStmt) String() string {
	return "END"
}

// String 返回 REM 注释语句的字符串表示
// 格式: "REM <注释文本>"
func (r *RemStmt) String() string {
	return fmt.Sprintf("REM %s", r.Text)
}

// String 返回 DIM 数组声明语句的字符串表示
// 格式: "DIM <数组名>(<大小1>[, <大小2>, ...])"
func (d *DimStmt) String() string {
	sizes := make([]string, len(d.Sizes))
	for i, size := range d.Sizes {
		sizes[i] = size.String()
	}
	return fmt.Sprintf("DIM %s(%s)", d.Name, strings.Join(sizes, ", "))
}

// String 返回 INPUT 输入语句的字符串表示
// 格式: "INPUT [<提示字符串>,] <变量名1>[, <变量名2>, ...]"
func (i *InputStmt) String() string {
	result := "INPUT"
	if i.Prompt != "" {
		result += fmt.Sprintf(" \"%s\",", i.Prompt)
	}
	for idx, v := range i.Vars {
		if idx > 0 {
			result += ", "
		}
		result += v
	}
	return result
}

// String 返回二元运算的字符串表示
// 格式: "(<左操作数> <运算符> <右操作数>)"
func (b *BinaryOp) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Op, b.Right.String())
}

// String 返回比较运算的字符串表示
// 格式: "(<左操作数> <运算符> <右操作数>)"
func (c *ComparisonOp) String() string {
	return fmt.Sprintf("(%s %s %s)", c.Left.String(), c.Op, c.Right.String())
}

// String 返回逻辑运算的字符串表示
// 格式: "(<左操作数> <运算符> <右操作数>)"
func (l *LogicalOp) String() string {
	return fmt.Sprintf("(%s %s %s)", l.Left.String(), l.Op, l.Right.String())
}

// String 返回一元运算的字符串表示
// 格式: "(<运算符><操作数>)"
func (u *UnaryOp) String() string {
	return fmt.Sprintf("(%s%s)", u.Op, u.Right.String())
}

// String 返回标识符的字符串表示
// 格式: "<变量名>"
func (i *Identifier) String() string {
	return i.Name
}

// String 返回函数调用的字符串表示
// 格式: "<函数名>(<参数1>, <参数2>, ...)"
func (f *FunctionCall) String() string {
	result := f.Name + "("
	for i, arg := range f.Args {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	return result
}

// String 返回数组访问的字符串表示
// 格式: "<数组名>(<索引1>[, <索引2>, ...])"
func (a *ArrayAccess) String() string {
	indices := make([]string, len(a.Indices))
	for i, idx := range a.Indices {
		indices[i] = idx.String()
	}
	return a.Name + "(" + strings.Join(indices, ", ") + ")"
}

// String 返回数字的字符串表示
// 使用最短的有效格式（整数不带小数点，浮点数根据需要显示）
func (n *Number) String() string {
	return fmt.Sprintf("%g", n.Value)
}

// String 返回字符串字面量的字符串表示
// 格式: "\"<文本>\""
func (s *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}