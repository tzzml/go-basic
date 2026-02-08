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
	variables   map[string]Value          // 变量存储表
	arrays      map[string]*ArrayInfo     // 数组存储表
	program     *ast.Program              // 当前加载的程序
	currentLine int                       // 当前执行到的行索引
	lineMap     map[int]int              // 行号 -> 程序行索引的映射表
	returnStack []int                    // GOSUB 返回地址栈
	forStack    []*ForFrame              // FOR 循环栈
	nameCache   map[string]string        // 名称规范化缓存（优化）
}

// ForFrame 表示 FOR 循环的栈帧
// 用于存储循环状态，支持嵌套循环
// 优化：缓存循环变量值，减少 map 查找
type ForFrame struct {
	varName   string  // 循环变量名
	endValue  float64 // 循环结束值
	stepValue float64 // 循环步长
	lineIdx   int     // NEXT 后要返回到的行索引
	value     float64 // 循环变量当前值（缓存）
}

// ArrayInfo 表示数组信息
// 用于存储多维数组的维度信息和数据
type ArrayInfo struct {
	dims    []int      // 各维度的大小
	data    []float64  // 扁平化存储的数组数据
	totalSize int      // 总元素数量
}

// NewArrayInfo 创建一个新的数组
func NewArrayInfo(dims []int) *ArrayInfo {
	totalSize := 1
	for _, d := range dims {
		totalSize *= d
	}
	return &ArrayInfo{
		dims:      dims,
		data:      make([]float64, totalSize),
		totalSize: totalSize,
	}
}

// CalculateIndex 计算多维索引的一维位置
// 将 (i1, i2, ..., in) 转换为扁平化索引
func (a *ArrayInfo) CalculateIndex(indices []int) int {
	if len(indices) != len(a.dims) {
		return -1
	}
	index := 0
	multiplier := 1
	// 从最后一维开始计算（行优先顺序）
	for i := len(indices) - 1; i >= 0; i-- {
		if indices[i] < 0 || indices[i] >= a.dims[i] {
			return -1 // 索引越界
		}
		index += indices[i] * multiplier
		multiplier *= a.dims[i]
	}
	return index
}

// NewInterpreter 创建一个新的 BASIC 解释器实例
func NewInterpreter() *Interpreter {
	return &Interpreter{
		variables: make(map[string]Value),
		arrays:    make(map[string]*ArrayInfo),
		lineMap:   make(map[int]int),
		nameCache: make(map[string]string), // 优化：初始化名称缓存
	}
}

// normalizeName 将名称转换为大写，用于统一变量名和函数名
// BASIC 传统上是不区分大小写的
// 优化：使用缓存避免重复的 strings.ToUpper 调用
func (i *Interpreter) normalizeName(name string) string {
	// 尝试从缓存获取
	if cached, ok := i.nameCache[name]; ok {
		return cached
	}
	// 转换为大写并存入缓存
	normalized := strings.ToUpper(name)
	i.nameCache[name] = normalized
	return normalized
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
			arr, ok := i.arrays[normalizedName]
			if !ok {
				fmt.Printf("Error: Array '%s' not declared\n", target.Name)
				return false
			}
			// 计算多维索引
			indices := make([]int, len(target.Indices))
			for idx, idxExpr := range target.Indices {
				indices[idx] = int(i.evaluateExpr(idxExpr).AsNumber())
			}
			flatIndex := arr.CalculateIndex(indices)
			if flatIndex < 0 {
				fmt.Printf("Error: Array index out of bounds\n")
				return false
			}
			arr.data[flatIndex] = value.AsNumber()
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

		// 将循环帧压入栈中，缓存循环变量值
		i.forStack = append(i.forStack, &ForFrame{
			varName:   normalizedName,
			endValue:  endVal,
			stepValue: stepVal,
			lineIdx:   i.currentLine,
			value:     startVal, // 缓存初始值
		})
		return false

	case *ast.NextStmt:
		// NEXT 语句：检查循环条件并决定是否继续循环
		if len(i.forStack) == 0 {
			fmt.Println("Error: NEXT without FOR")
			return false
		}

		frame := i.forStack[len(i.forStack)-1]

		// 优化：直接使用缓存的值，避免 map 查找
		newVal := frame.value + frame.stepValue

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
			// 更新缓存的值
			frame.value = newVal
			// 同步到 map，以便循环体内可以访问新值
			i.variables[frame.varName] = NumberValue(newVal)
			// 跳转回 FOR 语句的下一行
			i.currentLine = frame.lineIdx
		} else {
			// 循环结束，弹出循环帧
			// 将最终值写回 map（保持一致性）
			i.variables[frame.varName] = NumberValue(newVal)
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
		// DIM 数组声明语句：创建多维数组（使用大写的数组名）
		dims := make([]int, len(n.Sizes))
		for idx, sizeExpr := range n.Sizes {
			size := int(i.evaluateExpr(sizeExpr).AsNumber())
			if size < 0 {
				fmt.Printf("Error: Array dimension %d must be non-negative, got %d\n", idx+1, size)
				return false
			}
			dims[idx] = size
		}
		// BASIC 数组索引从 0 开始
		// DIM A(10) 创建 A(0) 到 A(9)，共 10 个元素
		// DIM B(3, 4) 创建 3x4 的二维数组，共 12 个元素
		normalizedName := i.normalizeName(n.Name)
		i.arrays[normalizedName] = NewArrayInfo(dims)
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
		arr, ok := i.arrays[normalizedName]
		if !ok {
			fmt.Printf("Error: Array '%s' not declared\n", n.Name)
			return NumberValue(0)
		}
		// 计算多维索引
		indices := make([]int, len(n.Indices))
		for idx, idxExpr := range n.Indices {
			indices[idx] = int(i.evaluateExpr(idxExpr).AsNumber())
		}
		flatIndex := arr.CalculateIndex(indices)
		if flatIndex < 0 {
			fmt.Printf("Error: Array index out of bounds\n")
			return NumberValue(0)
		}
		return NumberValue(arr.data[flatIndex])

	case *ast.BinaryOp:
		// 二元算术运算：+, -, *, /, ^, MOD
		leftVal := i.evaluateExpr(n.Left)
		rightVal := i.evaluateExpr(n.Right)

		// 处理字符串连接运算符 (+)
		if n.Op == "+" {
			// 如果任一操作数是字符串，则进行字符串连接
			if leftVal.isString || rightVal.isString {
				return StringValue(leftVal.String() + rightVal.String())
			}
			// 否则进行数字加法
			return NumberValue(leftVal.AsNumber() + rightVal.AsNumber())
		}

		// 其他运算只支持数字
		left := leftVal.AsNumber()
		right := rightVal.AsNumber()
		switch n.Op {
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
		leftVal := i.evaluateExpr(n.Left)
		rightVal := i.evaluateExpr(n.Right)
		var result bool

		// 如果任一操作数是字符串，则进行字符串比较
		if leftVal.isString || rightVal.isString {
			leftStr := leftVal.String()
			rightStr := rightVal.String()
			switch n.Op {
			case "=":
				result = leftStr == rightStr
			case "<>":
				result = leftStr != rightStr
			case ">":
				result = leftStr > rightStr
			case "<":
				result = leftStr < rightStr
			case ">=":
				result = leftStr >= rightStr
			case "<=":
				result = leftStr <= rightStr
			default:
				result = false
			}
		} else {
			// 数字比较
			left := leftVal.AsNumber()
			right := rightVal.AsNumber()
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
// 支持内置字符串函数：LEN, LEFT$, RIGHT$, MID$, INSTR, UCASE$, LCASE$
func (i *Interpreter) evaluateFunctionCall(node *ast.FunctionCall) Value {
	// 使用大写的函数名，使函数名不区分大小写
	normalizedName := i.normalizeName(node.Name)
	switch normalizedName {
	// 字符串函数
	case "LEN":
		// 返回字符串长度
		if len(node.Args) != 1 {
			fmt.Printf("Error: LEN requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		value := i.evaluateExpr(node.Args[0])
		return NumberValue(float64(len(value.String())))

	case "LEFT$":
		// 返回字符串左边 n 个字符
		if len(node.Args) != 2 {
			fmt.Printf("Error: LEFT$ requires 2 arguments, got %d\n", len(node.Args))
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

	case "RIGHT$":
		// 返回字符串右边 n 个字符
		if len(node.Args) != 2 {
			fmt.Printf("Error: RIGHT$ requires 2 arguments, got %d\n", len(node.Args))
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

	case "MID$":
		// 返回字符串从位置 start 开始的 n 个字符
		// MID$(str, start[, n])
		if len(node.Args) < 2 || len(node.Args) > 3 {
			fmt.Printf("Error: MID$ requires 2 or 3 arguments, got %d\n", len(node.Args))
			return StringValue("")
		}
		str := i.evaluateExpr(node.Args[0]).String()
		start := int(i.evaluateExpr(node.Args[1]).AsNumber())
		// BASIC 中字符串位置从 1 开始
		if start < 1 {
			start = 1
		}
		// 计算长度参数
		n := len(str) - start + 1 // 默认到字符串末尾
		if len(node.Args) == 3 {
			n = int(i.evaluateExpr(node.Args[2]).AsNumber())
		}
		// 转换为 0-based 索引
		startIdx := start - 1
		endIdx := startIdx + n
		if endIdx > len(str) {
			endIdx = len(str)
		}
		if startIdx >= len(str) || startIdx < 0 {
			return StringValue("")
		}
		return StringValue(str[startIdx:endIdx])

	case "INSTR":
		// 返回子串在字符串中的位置
		// INSTR([start,] str, substr)
		if len(node.Args) < 2 || len(node.Args) > 3 {
			fmt.Printf("Error: INSTR requires 2 or 3 arguments, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		var start int = 1
		var str, substr string
		if len(node.Args) == 2 {
			// INSTR(str, substr)
			str = i.evaluateExpr(node.Args[0]).String()
			substr = i.evaluateExpr(node.Args[1]).String()
		} else {
			// INSTR(start, str, substr)
			start = int(i.evaluateExpr(node.Args[0]).AsNumber())
			str = i.evaluateExpr(node.Args[1]).String()
			substr = i.evaluateExpr(node.Args[2]).String()
		}
		if start < 1 {
			start = 1
		}
		// 转换为 0-based 索引
		pos := strings.Index(str[start-1:], substr)
		if pos == -1 {
			return NumberValue(0) // 未找到返回 0
		}
		return NumberValue(float64(start + pos)) // 返回 1-based 位置

	case "UCASE$":
		// 将字符串转换为大写
		if len(node.Args) != 1 {
			fmt.Printf("Error: UCASE$ requires 1 argument, got %d\n", len(node.Args))
			return StringValue("")
		}
		str := i.evaluateExpr(node.Args[0]).String()
		return StringValue(strings.ToUpper(str))

	case "LCASE$":
		// 将字符串转换为小写
		if len(node.Args) != 1 {
			fmt.Printf("Error: LCASE$ requires 1 argument, got %d\n", len(node.Args))
			return StringValue("")
		}
		str := i.evaluateExpr(node.Args[0]).String()
		return StringValue(strings.ToLower(str))

	case "SPACE$":
		// 返回 n 个空格的字符串
		if len(node.Args) != 1 {
			fmt.Printf("Error: SPACE$ requires 1 argument, got %d\n", len(node.Args))
			return StringValue("")
		}
		n := int(i.evaluateExpr(node.Args[0]).AsNumber())
		if n < 0 {
			n = 0
		}
		return StringValue(strings.Repeat(" ", n))

	case "CHR$":
		// 将 ASCII 码转换为字符
		if len(node.Args) != 1 {
			fmt.Printf("Error: CHR$ requires 1 argument, got %d\n", len(node.Args))
			return StringValue("")
		}
		code := int(i.evaluateExpr(node.Args[0]).AsNumber())
		if code < 0 || code > 255 {
			fmt.Printf("Error: CHR$ argument must be between 0 and 255, got %d\n", code)
			return StringValue("")
		}
		return StringValue(string(rune(code)))

	case "ASC":
		// 返回字符的 ASCII 码
		if len(node.Args) != 1 {
			fmt.Printf("Error: ASC requires 1 argument, got %d\n", len(node.Args))
			return NumberValue(0)
		}
		str := i.evaluateExpr(node.Args[0]).String()
		if len(str) == 0 {
			fmt.Println("Error: ASC argument is an empty string")
			return NumberValue(0)
		}
		return NumberValue(float64(str[0]))

	// 数学函数
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
