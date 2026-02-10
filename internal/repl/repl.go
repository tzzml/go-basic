// Package repl 提供 BASIC 交互式编程环境（REPL）
// 包括代码存储、命令处理、程序执行等功能
package repl

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"zork-basic/internal/ast"
	"zork-basic/internal/compiler"
	"zork-basic/internal/formatter"
	"zork-basic/internal/interpreter"
	"zork-basic/internal/parser"
	"zork-basic/internal/vm"
)

const (
	// Prompt 交互模式提示符
	Prompt = "READY"
)

// CodeStore 代码存储：行号 -> 代码内容
type CodeStore struct {
	lines     map[int]string
	cachedAST *ast.Program // AST 缓存
	isDirty   bool         // 缓存失效标记
}

// NewCodeStore 创建一个新的代码存储实例
func NewCodeStore() *CodeStore {
	return &CodeStore{
		lines:   make(map[int]string),
		isDirty: true,
	}
}

// Set 添加或更新代码行
func (cs *CodeStore) Set(lineNumber int, code string) {
	cs.lines[lineNumber] = code
	cs.isDirty = true
	cs.cachedAST = nil
}

// Delete 删除代码行
func (cs *CodeStore) Delete(lineNumber int) bool {
	if _, exists := cs.lines[lineNumber]; !exists {
		return false
	}
	delete(cs.lines, lineNumber)
	cs.isDirty = true
	cs.cachedAST = nil
	return true
}

// GetLineNumbers 获取所有行号（排序后）
func (cs *CodeStore) GetLineNumbers() []int {
	numbers := make([]int, 0, len(cs.lines))
	for num := range cs.lines {
		numbers = append(numbers, num)
	}
	sort.Ints(numbers)
	return numbers
}

// GetCode 获取代码内容（用于执行，包含行号）
func (cs *CodeStore) GetCode() string {
	if len(cs.lines) == 0 {
		return ""
	}

	numbers := cs.GetLineNumbers()
	var result strings.Builder
	for _, num := range numbers {
		fmt.Fprintf(&result, "%d %s\n", num, cs.lines[num])
	}
	return result.String()
}

// GetProgram 获取解析后的 AST（带缓存）
func (cs *CodeStore) GetProgram() (*ast.Program, error) {
	// 如果缓存有效，直接返回
	if !cs.isDirty && cs.cachedAST != nil {
		return cs.cachedAST, nil
	}

	// 否则重新解析
	code := cs.GetCode()
	if code == "" {
		return nil, nil
	}

	parsedAST, err := parser.Parse("memory", []byte(code))
	if err != nil {
		return nil, err
	}

	prog, ok := parsedAST.(*ast.Program)
	if !ok {
		return nil, fmt.Errorf("parsed result is not a program")
	}

	// 更新缓存
	cs.cachedAST = prog
	cs.isDirty = false
	return prog, nil
}

// Clear 清空所有代码
func (cs *CodeStore) Clear() {
	cs.lines = make(map[int]string)
	cs.isDirty = true
	cs.cachedAST = nil
}

// IsEmpty 是否为空
func (cs *CodeStore) IsEmpty() bool {
	return len(cs.lines) == 0
}

// Count 获取行数
func (cs *CodeStore) Count() int {
	return len(cs.lines)
}

// Run 启动交互模式
func Run(version string, mode string) {
	printWelcome(version, mode)

	store := NewCodeStore()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("\n%s> ", Prompt)

		if !scanner.Scan() {
			// EOF (Ctrl+D)
			fmt.Println("\nGoodbye!")
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if input == "" {
			continue
		}

		// 处理命令
		if handleCommand(input, store, scanner, mode) {
			continue
		}

		// 尝试解析为 BASIC 代码行
		if lineNumber, code, ok := ParseBasicLine(input); ok {
			store.Set(lineNumber, code)
			fmt.Printf("Line %d updated\n", lineNumber)
		} else {
			fmt.Println("Error: Invalid input. Type 'HELP' for available commands.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}

// handleCommand 处理交互命令
func handleCommand(input string, store *CodeStore, scanner *bufio.Scanner, mode string) bool {
	trimmed := strings.TrimSpace(input)
	upper := strings.ToUpper(trimmed)

	// 提取命令的第一部分（第一个单词）
	firstWord := upper
	if spaceIdx := strings.Index(upper, " "); spaceIdx > 0 {
		firstWord = upper[:spaceIdx]
	}

	switch firstWord {
	case "LIST", "L":
		return cmdList(store)
	case "RUN", "R":
		return cmdRun(store, mode)
	case "CLEAR":
		store.Clear()
		fmt.Println("Program cleared")
		return true
	case "HELP", "H", "?":
		printInteractiveHelp()
		return true
	case "EXIT", "QUIT", "Q":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "NEW":
		store.Clear()
		fmt.Println("Ready for new program")
		return true
	case "DELETE", "D":
		// 从原始输入中提取参数
		parts := strings.Fields(trimmed)
		if len(parts) >= 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				if store.Delete(num) {
					fmt.Printf("Line %d deleted\n", num)
				} else {
					fmt.Printf("Line %d not found\n", num)
				}
				return true
			}
		}
		fmt.Println("Usage: DELETE <line_number>")
		return true
	case "EDIT", "E":
		parts := strings.Fields(trimmed)
		if len(parts) >= 2 {
			return cmdEdit(store, parts[1], scanner)
		}
		fmt.Println("Usage: EDIT <line_number>")
		return true
	case "SAVE":
		parts := strings.Fields(trimmed)
		if len(parts) >= 2 {
			// 保留文件名的原始大小写
			if idx := strings.Index(trimmed, " "); idx > 0 {
				filename := strings.TrimSpace(trimmed[idx+1:])
				return cmdSave(store, filename)
			}
		}
		fmt.Println("Usage: SAVE <filename>")
		return true
	case "LOAD":
		parts := strings.Fields(trimmed)
		if len(parts) >= 2 {
			// 保留文件名的原始大小写
			if idx := strings.Index(trimmed, " "); idx > 0 {
				filename := strings.TrimSpace(trimmed[idx+1:])
				return cmdLoad(store, filename)
			}
		}
		fmt.Println("Usage: LOAD <filename>")
		return true
	case "FORMAT", "F":
		return cmdFormat(store)
	}

	return false
}

// cmdList LIST 命令
func cmdList(store *CodeStore) bool {
	if store.IsEmpty() {
		fmt.Println("(No program in memory)")
		return true
	}

	numbers := store.GetLineNumbers()
	for _, num := range numbers {
		fmt.Printf("%d %s\n", num, store.lines[num])
	}
	return true
}

// cmdRun RUN 命令
func cmdRun(store *CodeStore, mode string) bool {
	if store.IsEmpty() {
		fmt.Println("Error: No program to run")
		return true
	}

	// 获取 AST（利用缓存）
	prog, err := store.GetProgram()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return true
	}

	if prog == nil {
		return true
	}

	// 执行程序
	if mode == "vm" {
		comp := compiler.New()
		chunk, err := comp.Compile(prog)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			return true
		}
		vm := vm.New(chunk)
		if err := vm.Run(); err != nil {
			fmt.Printf("Runtime error: %v\n", err)
		}
	} else {
		interp := interpreter.NewInterpreter()
		interp.ExecuteProgram(prog)
	}
	fmt.Println("\nProgram complete.")
	return true
}

// cmdEdit EDIT 命令
func cmdEdit(store *CodeStore, lineNumStr string, scanner *bufio.Scanner) bool {
	num, err := strconv.Atoi(lineNumStr)
	if err != nil {
		fmt.Println("Error: Invalid line number")
		return true
	}

	if code, exists := store.lines[num]; exists {
		fmt.Printf("Current line %d: %s\n", num, code)
		fmt.Print("Enter new line (or press Enter to cancel): ")

		if !scanner.Scan() {
			return true
		}

		newLine := strings.TrimSpace(scanner.Text())
		if newLine == "" {
			fmt.Println("Edit cancelled")
			return true
		}

		if lineNumber, code, ok := ParseBasicLine(newLine); ok {
			if lineNumber != num {
				fmt.Printf("Warning: Line number changed from %d to %d\n", num, lineNumber)
			}
			store.Set(lineNumber, code)
			fmt.Printf("Line %d updated\n", lineNumber)
		} else {
			fmt.Println("Error: Invalid BASIC syntax")
		}
	} else {
		fmt.Printf("Line %d not found\n", num)
	}

	return true
}

// cmdSave SAVE 命令
func cmdSave(store *CodeStore, filename string) bool {
	if store.IsEmpty() {
		fmt.Println("Error: No program to save")
		return true
	}

	code := store.GetCode()
	err := os.WriteFile(filename, []byte(code), 0644)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return true
	}

	fmt.Printf("Program saved to %s (%d lines)\n", filename, store.Count())
	return true
}

// cmdLoad LOAD 命令
func cmdLoad(store *CodeStore, filename string) bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return true
	}

	store.Clear()
	lineCount := 0

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if lineNumber, code, ok := ParseBasicLine(line); ok {
			store.Set(lineNumber, code)
			lineCount++
		}
	}

	fmt.Printf("Loaded %d lines from %s\n", lineCount, filename)
	return true
}

// cmdFormat FORMAT 命令 - 重新格式化程序行号
func cmdFormat(store *CodeStore) bool {
	if store.IsEmpty() {
		fmt.Println("Error: No program to format")
		return true
	}

	// 获取 AST（利用缓存）
	prog, err := store.GetProgram()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return true
	}

	if prog == nil {
		return true
	}

	// 收集所有旧行号
	oldLineNumbers := store.GetLineNumbers()

	// 创建旧行号到新行号的映射（10, 20, 30, ...）
	lineNumberMap := make(map[int]int)
	newLineNumber := 10
	for _, oldNum := range oldLineNumbers {
		lineNumberMap[oldNum] = newLineNumber
		newLineNumber += 10
	}

	// 创建新代码存储
	newStore := NewCodeStore()

	// 遍历所有语句，更新行号和行号引用，并计算缩进
	currentIndent := 0
	for _, line := range prog.Lines {
		newNum := lineNumberMap[line.LineNumber]

		// 计算当前行开始时的缩进（处理 NEXT 等需要提前缩进的情况）
		before, after := formatter.GetIndentDelta(line)

		displayIndent := currentIndent
		if before < 0 {
			displayIndent += before
		}
		if displayIndent < 0 {
			displayIndent = 0
		}

		formattedCode := formatter.FormatLine(line, lineNumberMap, displayIndent)
		newStore.Set(newNum, formattedCode)

		// 更新下一行的缩进
		currentIndent += after
		if currentIndent < 0 {
			currentIndent = 0
		}
	}

	// 替换旧存储
	store.lines = newStore.lines

	fmt.Printf("Program formatted: %d lines renumbered\n", store.Count())
	return true
}

// ParseBasicLine 解析 BASIC 代码行，提取行号和代码
func ParseBasicLine(input string) (int, string, bool) {
	// 提取行号
	i := 0
	for i < len(input) && input[i] >= '0' && input[i] <= '9' {
		i++
	}

	if i == 0 {
		return 0, "", false
	}

	lineNumber, err := strconv.Atoi(input[:i])
	if err != nil {
		return 0, "", false
	}

	// 提取代码部分（跳过行号后的空格）
	code := strings.TrimSpace(input[i:])
	if code == "" {
		return 0, "", false
	}

	// 注意：PEG 解析器已原生支持大小写不敏感关键字，无需预处理

	return lineNumber, code, true
}

// ExecuteProgram 执行 BASIC 程序
func ExecuteProgram(code string, source string, mode string) {
	// 使用 pigeon 生成的解析器
	parsedAST, err := parser.Parse(source, []byte(code))
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	prog, ok := parsedAST.(*ast.Program)
	if !ok {
		fmt.Println("Parse error: not a program")
		return
	}

	// 执行程序
	if mode == "vm" {
		comp := compiler.New()
		chunk, err := comp.Compile(prog)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			return
		}
		vm := vm.New(chunk)
		if err := vm.Run(); err != nil {
			fmt.Printf("Runtime error: %v\n", err)
		}
	} else {
		interp := interpreter.NewInterpreter()
		interp.ExecuteProgram(prog)
	}
	fmt.Println("\nProgram complete.")
}

// printWelcome 打印欢迎信息
func printWelcome(version string, mode string) {
	fmt.Println("=====================================")
	fmt.Println("   zork-basic BASIC Interpreter")
	fmt.Printf("   Version %s (Mode: %s)\n", version, mode)
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("Interactive mode. Type 'HELP' for commands.")
	fmt.Println("Enter BASIC statements directly or use commands.")
}

// printInteractiveHelp 打印交互模式帮助
func printInteractiveHelp() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  BASIC Code     - Enter BASIC statements (e.g., '10 PRINT X')")
	fmt.Println("  LIST, L        - List all program lines")
	fmt.Println("  RUN, R         - Execute the program")
	fmt.Println("  EDIT <n>       - Edit line number n")
	fmt.Println("  DELETE <n>     - Delete line number n")
	fmt.Println("  FORMAT, F      - Format program (renumber lines, uppercase keywords)")
	fmt.Println("  CLEAR          - Clear all program lines")
	fmt.Println("  NEW            - Start a new program")
	fmt.Println("  SAVE <file>    - Save program to file")
	fmt.Println("  LOAD <file>    - Load program from file")
	fmt.Println("  HELP, ?, H     - Show this help message")
	fmt.Println("  EXIT, QUIT, Q  - Exit the interpreter")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  10 PRINT \"Hello World\"")
	fmt.Println("  20 X = 10")
	fmt.Println("  RUN")
	fmt.Println("  LIST")
	fmt.Println("  FORMAT")
}
