package interpreter

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

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
		return strconv.FormatFloat(v.number, 'g', -1, 64)
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

// IsNumber 返回是否为数字类型
func (v Value) IsNumber() bool {
	return v.isNumber
}

// IsString 返回是否为字符串类型
func (v Value) IsString() bool {
	return v.isString
}

// Interpreter BASIC 解释器
// 负责解析和执行 AST（抽象语法树）
type Interpreter struct {
	variables    map[string]Value      // 变量存储表
	arrays       map[string]*ArrayInfo // 数组存储表
	program      *ast.Program          // 当前加载的程序
	currentLine  int                   // 当前执行到的行索引
	lineMap      map[int]int           // 行号 -> 程序行索引的映射表
	returnStack  []int                 // GOSUB 返回地址栈
	forStack     []*ForFrame           // FOR 循环栈
	indexBuf     []int                 // 数组索引复用缓冲区（优化）
	nameCache    map[string]string     // 名称规范化缓存（优化）
	forFramePool *sync.Pool            // 循环帧对象池（优化）
	output       io.Writer             // 正常输出（PRINT 语句等）
	errOutput    io.Writer             // 错误输出
	input        io.Reader             // 输入源（INPUT 语句）
}

// Option 是解释器的配置选项函数
type Option func(*Interpreter)

// WithOutput 设置正常输出目标
func WithOutput(w io.Writer) Option {
	return func(i *Interpreter) {
		i.output = w
	}
}

// WithErrOutput 设置错误输出目标
func WithErrOutput(w io.Writer) Option {
	return func(i *Interpreter) {
		i.errOutput = w
	}
}

// WithInput 设置输入源
func WithInput(r io.Reader) Option {
	return func(i *Interpreter) {
		i.input = r
	}
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
	dims      []int     // 各维度的大小
	Data      []float64 // 扁平化存储的数组数据
	totalSize int       // 总元素数量
}

// NewArrayInfo 创建一个新的数组
func NewArrayInfo(dims []int) *ArrayInfo {
	totalSize := 1
	for _, d := range dims {
		totalSize *= d
	}
	return &ArrayInfo{
		dims:      dims,
		Data:      make([]float64, totalSize),
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
// 默认输出到 os.Stdout，错误到 os.Stderr，输入从 os.Stdin
// 可通过 Option 函数自定义
func NewInterpreter(opts ...Option) *Interpreter {
	i := &Interpreter{
		variables: make(map[string]Value),
		arrays:    make(map[string]*ArrayInfo),
		lineMap:   make(map[int]int),
		nameCache: make(map[string]string),
		forFramePool: &sync.Pool{
			New: func() interface{} {
				return &ForFrame{}
			},
		},
		output:    os.Stdout,
		errOutput: os.Stderr,
		input:     os.Stdin,
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
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
	i.nameCache = make(map[string]string) // 重置名称缓存，避免无限增长
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
				fmt.Fprintf(i.errOutput, "Error: Array '%s' not declared\n", target.Name)
				return false
			}
			// 计算多维索引
			indices := i.getIndexBuf(len(target.Indices))
			for idx, idxExpr := range target.Indices {
				indices[idx] = int(i.evaluateExpr(idxExpr).AsNumber())
			}
			flatIndex := arr.CalculateIndex(indices)
			if flatIndex < 0 {
				fmt.Fprintf(i.errOutput, "Error: Array index out of bounds\n")
				return false
			}
			arr.Data[flatIndex] = value.AsNumber()
		default:
			fmt.Fprintf(i.errOutput, "Error: Invalid assignment target type: %T\n", target)
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
					fmt.Fprint(i.output, " ")
				}
				// 分号分隔符：不添加空格
			}
			v := i.evaluateExpr(val)
			fmt.Fprint(i.output, v.String())
		}
		// 只有在没有末尾分号或逗号时才换行
		if n.Trailer == "" {
			fmt.Fprintln(i.output)
		}
		return false

	case *ast.IfStmt:
		// IF...THEN...ELSE...END IF 条件语句
		cond := i.evaluateExpr(n.Condition)
		if cond.IsTrue() {
			// 条件为真，执行 THEN 块
			for _, s := range n.ThenStmts {
				if i.executeStatement(s) {
					return true // 传播 GOTO/GOSUB/END/RETURN 的控制流变更
				}
			}
		} else if len(n.ElseStmts) > 0 {
			// 条件为假，执行 ELSE 块（如果有）
			for _, s := range n.ElseStmts {
				if i.executeStatement(s) {
					return true // 传播 GOTO/GOSUB/END/RETURN 的控制流变更
				}
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

		// Get frame from pool
		frame := i.forFramePool.Get().(*ForFrame)
		frame.varName = normalizedName
		frame.endValue = endVal
		frame.stepValue = stepVal
		frame.lineIdx = i.currentLine
		frame.value = startVal

		// 将循环帧压入栈中，缓存循环变量值
		i.forStack = append(i.forStack, frame)
		return false

	case *ast.NextStmt:
		// NEXT 语句：检查循环条件并决定是否继续循环
		if len(i.forStack) == 0 {
			fmt.Fprintln(i.errOutput, "Error: NEXT without FOR")
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
			// Release frame back to pool
			i.forFramePool.Put(frame)
		}
		return false

	case *ast.GotoStmt:
		// GOTO 无条件跳转语句
		if idx, ok := i.lineMap[n.LineNumber]; ok {
			i.currentLine = idx // 直接设置为目标行索引
		} else {
			fmt.Fprintf(i.errOutput, "Error: Line %d not found\n", n.LineNumber)
		}
		return true

	case *ast.GosubStmt:
		// GOSUB 子程序调用语句
		if idx, ok := i.lineMap[n.LineNumber]; ok {
			// 将返回地址压入栈（currentLine 已经被递增，指向调用行的下一行）
			i.returnStack = append(i.returnStack, i.currentLine)
			i.currentLine = idx
		} else {
			fmt.Fprintf(i.errOutput, "Error: Line %d not found\n", n.LineNumber)
		}
		return true

	case *ast.ReturnStmt:
		// RETURN 从子程序返回
		if len(i.returnStack) == 0 {
			fmt.Fprintln(i.errOutput, "Error: RETURN without GOSUB")
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
				fmt.Fprintf(i.errOutput, "Error: Array dimension %d must be non-negative, got %d\n", idx+1, size)
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
				fmt.Fprintf(i.output, "%s [%d]: ", prompt, idx+1)
			} else {
				fmt.Fprint(i.output, prompt)
			}

			var input string
			fmt.Fscanln(i.input, &input)
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
		fmt.Fprintf(i.errOutput, "Warning: unhandled statement type: %T\n", stmt)
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
		fmt.Fprintf(i.errOutput, "Error: Undefined variable '%s'\n", n.Name)
		return NumberValue(0)

	case *ast.FunctionCall:
		// 函数调用
		return i.evaluateFunctionCall(n)

	case *ast.ArrayAccess:
		// 数组访问：获取数组元素的值（使用大写的数组名）
		normalizedName := i.normalizeName(n.Name)
		arr, ok := i.arrays[normalizedName]
		if !ok {
			fmt.Fprintf(i.errOutput, "Error: Array '%s' not declared\n", n.Name)
			return NumberValue(0)
		}
		// 计算多维索引
		indices := i.getIndexBuf(len(n.Indices))
		for idx, idxExpr := range n.Indices {
			indices[idx] = int(i.evaluateExpr(idxExpr).AsNumber())
		}
		flatIndex := arr.CalculateIndex(indices)
		if flatIndex < 0 {
			fmt.Fprintf(i.errOutput, "Error: Array index out of bounds\n")
			return NumberValue(0)
		}
		return NumberValue(arr.Data[flatIndex])

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
				fmt.Fprintln(i.errOutput, "Error: Division by zero")
				return NumberValue(0)
			}
			return NumberValue(left / right)
		case "^":
			return NumberValue(math.Pow(left, right))
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
		// 逻辑运算：AND, OR（支持短路求值）
		left := i.evaluateExpr(n.Left).IsTrue()
		switch n.Op {
		case "AND":
			if !left {
				return NumberValue(0) // 短路：左侧为假，跳过右侧
			}
			if i.evaluateExpr(n.Right).IsTrue() {
				return NumberValue(1)
			}
			return NumberValue(0)
		case "OR":
			if left {
				return NumberValue(1) // 短路：左侧为真，跳过右侧
			}
			if i.evaluateExpr(n.Right).IsTrue() {
				return NumberValue(1)
			}
			return NumberValue(0)
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
		fmt.Fprintf(i.errOutput, "Warning: unhandled expression type: %T\n", node)
		return NumberValue(0)
	}
}

// evaluateFunctionCall 计算函数调用的值
// 使用函数分发表实现 O(1) 查找（定义在 builtins.go）
func (i *Interpreter) evaluateFunctionCall(node *ast.FunctionCall) Value {
	normalizedName := i.normalizeName(node.Name)
	if fn, ok := builtinFuncs[normalizedName]; ok {
		return fn(i, node)
	}
	fmt.Fprintf(i.errOutput, "Error: Unknown function '%s'\n", node.Name)
	return NumberValue(0)
}

// getIndexBuf 获取可复用的索引缓冲区，避免每次数组访问分配新切片
func (i *Interpreter) getIndexBuf(size int) []int {
	if cap(i.indexBuf) >= size {
		return i.indexBuf[:size]
	}
	i.indexBuf = make([]int, size)
	return i.indexBuf
}
