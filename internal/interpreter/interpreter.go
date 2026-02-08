package interpreter

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"zork-basic/internal/ast"
)

// Value 表示 BASIC 解释器中的任意值
// 支持数字和字符串两种类型
type Value struct {
	isNumber bool    // 是否为数字类型
	isString bool    // 是否为字符串类型
	number   float64 // 数字值
	string   string  // 字符串值
}

// NumberValue 创建一个数字类型的 Value
func NumberValue(v float64) Value {
	return Value{isNumber: true, number: v}
}

// StringValue 创建一个字符串类型的 Value
func StringValue(v string) Value {
	return Value{isString: true, string: v}
}

// String 返回值的字符串表示
func (v Value) String() string {
	if v.isNumber {
		return fmt.Sprintf("%g", v.number)
	}
	if v.isString {
		return v.string
	}
	return ""
}

// AsNumber 将值转换为数字类型返回
// 如果是数字则直接返回，如果是字符串则尝试解析，否则返回 0
func (v Value) AsNumber() float64 {
	if v.isNumber {
		return v.number
	}
	if v.isString {
		f, _ := strconv.ParseFloat(v.string, 64)
		return f
	}
	return 0
}

// IsTrue 判断值是否为真
// 数字：非零为真，零为假
// 字符串：非空为真，空字符串为假
func (v Value) IsTrue() bool {
	if v.isNumber {
		return v.number != 0
	}
	if v.isString {
		return v.string != ""
	}
	return false
}

// Interpreter BASIC 解释器
// 负责解析和执行 AST（抽象语法树）
type Interpreter struct {
	variables   map[string]Value         // 变量存储表
	arrays      map[string][]float64     // 数组存储表
	program     *ast.Program             // 当前加载的程序
	currentLine int                      // 当前执行到的行索引
	lineMap     map[int]int             // 行号 -> 程序行索引的映射表
	returnStack []int                   // GOSUB 返回地址栈
	forStack    []*ForFrame             // FOR 循环栈
}

// ForFrame 表示 FOR 循环的栈帧
// 用于存储循环状态，支持嵌套循环
type ForFrame struct {
	varName   string  // 循环变量名
	endValue  float64 // 循环结束值
	stepValue float64 // 循环步长
	lineIdx   int     // NEXT 后要返回到的行索引
}

// NewInterpreter 创建一个新的 BASIC 解释器实例
func NewInterpreter() *Interpreter {
	return &Interpreter{
		variables: make(map[string]Value),
		arrays:    make(map[string][]float64),
		lineMap:   make(map[int]int),
	}
}

// normalizeName 将名称转换为大写，用于统一变量名和函数名
// BASIC 传统上是不区分大小写的
func (i *Interpreter) normalizeName(name string) string {
	// 转换为大写
	return strings.ToUpper(name)
}

// LoadProgram 加载 BASIC 程序到解释器
// 建立行号到行索引的映射表，用于 GOTO/GOSUB 跳转
func (i *Interpreter) LoadProgram(program *ast.Program) {
	i.program = program
	i.lineMap = make(map[int]int)
	for idx, line := range program.Lines {
		i.lineMap[line.LineNumber] = idx
	}
}

// ExecuteProgram 执行 BASIC 程序
// 按行号顺序执行程序，支持 GOTO/GOSUB 改变执行流
func (i *Interpreter) ExecuteProgram(program *ast.Program) {
	i.LoadProgram(program)

	// 按顺序执行各行
	for i.currentLine = 0; i.currentLine < len(program.Lines); {
		line := program.Lines[i.currentLine]
		i.currentLine++ // 移动到下一行
		for _, stmt := range line.Statements {
			if i.executeStatement(stmt) {
				// GOTO/GOSUB/END/RETURN 改变了 currentLine，跳出内层循环
				// currentLine 已经被设置为正确的目标索引（下一行要执行的）
				break
			}
		}
		// 如果没有跳转，currentLine 已经指向下一行，继续循环
		// 如果有跳转，currentLine 已被设置为要执行的目标行，继续循环
	}
}

// executeStatement 执行单个语句
// 返回 true 表示执行了 GOTO 需要改变执行流
func (i *Interpreter) executeStatement(stmt ast.Node) bool {
	switch n := stmt.(type) {
	case *ast.Assignment:
		// 赋值语句：支持变量赋值和数组元素赋值
		value := i.evaluateExpr(n.Value)
		switch target := n.Target.(type) {
		case *ast.Identifier:
			// 普通变量赋值 - 使用大写的变量名
			normalizedName := i.normalizeName(target.Name)
			i.variables[normalizedName] = value
		case *ast.ArrayAccess:
			// 数组元素赋值 - 使用大写的数组名
			normalizedName := i.normalizeName(target.Name)
			index := int(i.evaluateExpr(target.Index).AsNumber())
			arr, ok := i.arrays[normalizedName]
			if !ok {
				fmt.Printf("Error: Array '%s' not declared\n", target.Name)
				return false
			}
			if index < 0 || index >= len(arr) {
				fmt.Printf("Error: Array index %d out of bounds (0-%d)\n", index, len(arr)-1)
				return false
			}
			arr[index] = value.AsNumber()
		default:
			fmt.Printf("Error: Invalid assignment target type: %T\n", target)
		}
		return false

	case *ast.PrintStmt:
		// PRINT 语句：输出多个值
		// 分号分隔符：紧凑输出，值之间不添加空格
		// 逗号分隔符：值之间添加空格
		for j, val := range n.Values {
			if j > 0 {
				// 检查分隔符类型
				if j <= len(n.Separators) && n.Separators[j-1] == "," {
					// 逗号分隔符：添加空格
					fmt.Print(" ")
				}
				// 分号分隔符：不添加空格
			}
			v := i.evaluateExpr(val)
			fmt.Print(v.String())
		}
		// 只有在没有末尾分号或逗号时才换行
		if n.Trailer == "" {
			fmt.Println()
		}
		return false

	case *ast.IfStmt:
		// IF...THEN...ELSE...END IF 条件语句
		cond := i.evaluateExpr(n.Condition)
		if cond.IsTrue() {
			// 条件为真，执行 THEN 块
			for _, s := range n.ThenStmts {
				i.executeStatement(s)
			}
		} else if len(n.ElseStmts) > 0 {
			// 条件为假，执行 ELSE 块（如果有）
			for _, s := range n.ElseStmts {
				i.executeStatement(s)
			}
		}
		return false

	case *ast.ForStmt:
		// FOR...NEXT 循环语句
		startVal := i.evaluateExpr(n.Start).AsNumber()
		endVal := i.evaluateExpr(n.End).AsNumber()
		stepVal := i.evaluateExpr(n.Step).AsNumber()

		// 初始化循环变量（使用大写的变量名）
		normalizedName := i.normalizeName(n.Var)
		i.variables[normalizedName] = NumberValue(startVal)

		// 将循环帧压入栈中（使用大写的变量名）
		i.forStack = append(i.forStack, &ForFrame{
			varName:   normalizedName,
			endValue:  endVal,
			stepValue: stepVal,
			lineIdx:   i.currentLine,
		})
		return false

	case *ast.NextStmt:
		// NEXT 语句：检查循环条件并决定是否继续循环
		if len(i.forStack) == 0 {
			fmt.Println("Error: NEXT without FOR")
			return false
		}

		frame := i.forStack[len(i.forStack)-1]
		currentVal := i.variables[frame.varName].AsNumber()
		newVal := currentVal + frame.stepValue
		i.variables[frame.varName] = NumberValue(newVal)

		// 检查是否应该继续循环
		shouldContinue := false
		if frame.stepValue > 0 && newVal <= frame.endValue {
			// 正步长：当前值 <= 结束值时继续
			shouldContinue = true
		} else if frame.stepValue < 0 && newVal >= frame.endValue {
			// 负步长：当前值 >= 结束值时继续
			shouldContinue = true
		}

		if shouldContinue {
			// 跳转回 FOR 语句的下一行
			i.currentLine = frame.lineIdx
		} else {
			// 循环结束，弹出循环帧
			i.forStack = i.forStack[:len(i.forStack)-1]
		}
		return false

	case *ast.GotoStmt:
		// GOTO 无条件跳转语句
		if idx, ok := i.lineMap[n.LineNumber]; ok {
			i.currentLine = idx // 直接设置为目标行索引
		} else {
			fmt.Printf("Error: Line %d not found\n", n.LineNumber)
		}
		return true

	case *ast.GosubStmt:
		// GOSUB 子程序调用语句
		if idx, ok := i.lineMap[n.LineNumber]; ok {
			// 将返回地址压入栈（currentLine 已经被递增，指向调用行的下一行）
			i.returnStack = append(i.returnStack, i.currentLine)
			i.currentLine = idx
		} else {
			fmt.Printf("Error: Line %d not found\n", n.LineNumber)
		}
		return true

	case *ast.ReturnStmt:
		// RETURN 从子程序返回
		if len(i.returnStack) == 0 {
			fmt.Println("Error: RETURN without GOSUB")
			return false
		}
		retLine := i.returnStack[len(i.returnStack)-1]
		i.returnStack = i.returnStack[:len(i.returnStack)-1]
		i.currentLine = retLine
		return true

	case *ast.EndStmt:
		// END 程序结束语句
		i.currentLine = len(i.program.Lines)
		return true

	case *ast.RemStmt:
		// REM 注释语句：不做任何事
		return false

	case *ast.DimStmt:
		// DIM 数组声明语句：创建数组（使用大写的数组名）
		size := int(i.evaluateExpr(n.Size).AsNumber())
		if size < 0 {
			fmt.Printf("Error: Array size must be non-negative, got %d\n", size)
			return false
		}
		// BASIC 数组索引通常从 0 或 1 开始，这里实现为 0-based
		// DIM A(10) 创建 A(0) 到 A(9)，共 10 个元素
		normalizedName := i.normalizeName(n.Name)
		i.arrays[normalizedName] = make([]float64, size)
		return false

	case *ast.InputStmt:
		// INPUT 输入语句：从用户读取输入并存储到变量
		// 支持提示字符串和多个变量
		prompt := "? "
		if n.Prompt != "" {
			prompt = n.Prompt
		}

		for idx, varName := range n.Vars {
			// 多个变量时，后续变量显示序号
			if len(n.Vars) > 1 {
				fmt.Printf("%s [%d]: ", prompt, idx+1)
			} else {
				fmt.Print(prompt)
			}

			var input string
			fmt.Scanln(&input)
			num, err := strconv.ParseFloat(input, 64)
			// 使用大写的变量名
			normalizedName := i.normalizeName(varName)
			if err != nil {
				// 解析失败，作为字符串存储
				i.variables[normalizedName] = StringValue(input)
			} else {
				// 解析成功，作为数字存储
				i.variables[normalizedName] = NumberValue(num)
			}
		}
		return false

	default:
		fmt.Printf("Warning: unhandled statement type: %T\n", stmt)
		return false
	}
}

// evaluateExpr 计算表达式的值
// 支持数字、字符串、变量、二元运算、比较运算、逻辑运算、一元运算
func (i *Interpreter) evaluateExpr(node ast.Node) Value {
	switch n := node.(type) {
	case *ast.Number:
		// 数字字面量
		return NumberValue(n.Value)

	case *ast.StringLiteral:
		// 字符串字面量
		return StringValue(n.Value)

	case *ast.Identifier:
		// 变量：从变量表中查找（使用大写的变量名）
		normalizedName := i.normalizeName(n.Name)
		if val, ok := i.variables[normalizedName]; ok {
			return val
		}
		fmt.Printf("Error: Undefined variable '%s'\n", n.Name)
		return NumberValue(0)

	case *ast.FunctionCall:
		// 函数调用
		return i.evaluateFunctionCall(n)

	case *ast.ArrayAccess:
		// 数组访问：获取数组元素的值（使用大写的数组名）
		normalizedName := i.normalizeName(n.Name)
		index := int(i.evaluateExpr(n.Index).AsNumber())
		arr, ok := i.arrays[normalizedName]
		if !ok {
			fmt.Printf("Error: Array '%s' not declared\n", n.Name)
			return NumberValue(0)
		}
		if index < 0 || index >= len(arr) {
			fmt.Printf("Error: Array index %d out of bounds (0-%d)\n", index, len(arr)-1)
			return NumberValue(0)
		}
		return NumberValue(arr[index])

	case *ast.BinaryOp:
		// 二元算术运算：+, -, *, /, ^, MOD
		left := i.evaluateExpr(n.Left).AsNumber()
		right := i.evaluateExpr(n.Right).AsNumber()
		switch n.Op {
		case "+":
			return NumberValue(left + right)
		case "-":
			return NumberValue(left - right)
		case "*":
			return NumberValue(left * right)
		case "/":
			if right == 0 {
				fmt.Println("Error: Division by zero")
				return NumberValue(0)
			}
			return NumberValue(left / right)
		case "^":
			return NumberValue(pow(left, right))
		case "MOD":
			return NumberValue(math.Mod(left, right))
		}
		return NumberValue(0)

	case *ast.ComparisonOp:
		// 比较运算：=, <>, >, <, >=, <=
		left := i.evaluateExpr(n.Left).AsNumber()
		right := i.evaluateExpr(n.Right).AsNumber()
		var result bool
		switch n.Op {
		case "=":
			result = left == right
		case "<>":
			result = left != right
		case ">":
			result = left > right
		case "<":
			result = left < right
		case ">=":
			result = left >= right
		case "<=":
			result = left <= right
		default:
			result = false
		}
		// BASIC 中布尔值用数字表示：真=1，假=0
		if result {
			return NumberValue(1)
		}
		return NumberValue(0)

	case *ast.LogicalOp:
		// 逻辑运算：AND, OR
		left := i.evaluateExpr(n.Left).IsTrue()
		right := i.evaluateExpr(n.Right).IsTrue()
		var result bool
		switch n.Op {
		case "AND":
			result = left && right
		case "OR":
			result = left || right
		default:
			result = false
		}
		if result {
			return NumberValue(1)
		}
		return NumberValue(0)

	case *ast.UnaryOp:
		// 一元运算：+, -, NOT
		if n.Op == "NOT" {
			// NOT 是逻辑运算，返回布尔值
			right := i.evaluateExpr(n.Right).IsTrue()
			if !right {
				return NumberValue(1)
			}
			return NumberValue(0)
		}
		// + 和 - 是算术运算
		right := i.evaluateExpr(n.Right).AsNumber()
		if n.Op == "+" {
			return NumberValue(right)
		}
		return NumberValue(-right)

	default:
		fmt.Printf("Warning: unhandled expression type: %T\n", node)
		return NumberValue(0)
	}
}

// evaluateFunctionCall 计算函数调用的值
// 支持内置数学函数：ABS, SIN, COS, TAN, INT, RND, SQR, LOG, EXP
func (i *Interpreter) evaluateFunctionCall(node *ast.FunctionCall) Value {
	// 使用大写的函数名，使函数名不区分大小写
	normalizedName := i.normalizeName(node.Name)
	switch normalizedName {
	case "ABS":
		// 绝对值
		if len(node.Args) != 1 {
			fmt.Printf("Error: ABS requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0]).AsNumber()
		return NumberValue(math.Abs(value))

	case "SIN":
		// 正弦（弧度）
		if len(node.Args) != 1 {
			fmt.Printf("Error: SIN requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0]).AsNumber()
		return NumberValue(math.Sin(value))

	case "COS":
		// 余弦（弧度）
		if len(node.Args) != 1 {
			fmt.Printf("Error: COS requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0]).AsNumber()
		return NumberValue(math.Cos(value))

	case "TAN":
		// 正切（弧度）
		if len(node.Args) != 1 {
			fmt.Printf("Error: TAN requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0]).AsNumber()
		return NumberValue(math.Tan(value))

	case "INT":
		// 取整
		if len(node.Args) != 1 {
			fmt.Printf("Error: INT requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0]).AsNumber()
		return NumberValue(math.Trunc(value))

	case "SQR":
		// 平方根
		if len(node.Args) != 1 {
			fmt.Printf("Error: SQR requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0]).AsNumber()
		if value < 0 {
			fmt.Println("Error: SQR of negative number")
			return NumberValue(0)
		}
		return NumberValue(math.Sqrt(value))

	case "LOG":
		// 自然对数
		if len(node.Args) != 1 {
			fmt.Printf("Error: LOG requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0]).AsNumber()
		if value <= 0 {
			fmt.Println("Error: LOG of non-positive number")
			return NumberValue(0)
		}
		return NumberValue(math.Log(value))

	case "EXP":
		// e 的 x 次方
		if len(node.Args) != 1 {
			fmt.Printf("Error: EXP requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0]).AsNumber()
		return NumberValue(math.Exp(value))

	case "RND":
		// 随机数（0 到 1 之间）
		if len(node.Args) != 0 {
			fmt.Printf("Error: RND requires 0 arguments, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		return NumberValue(rand.Float64())

	default:
		fmt.Printf("Error: Unknown function '%s'\n", node.Name)
		return NumberValue(0)
	}
}

// pow 计算 x 的 y 次方
// BASIC 中的幂运算符
func pow(x, y float64) float64 {
	result := 1.0
	for i := 0; i < int(y); i++ {
		result *= x
	}
	return result
}
