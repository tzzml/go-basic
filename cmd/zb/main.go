// zork-basic - BASIC 语言解释器
// 使用 Go 语言编写，基于 PEG 解析器生成器
package main

import (
	"flag"
	"fmt"
	"os"

	"zork-basic/internal/ast"
	"zork-basic/internal/compiler"
	"zork-basic/internal/parser"
	"zork-basic/internal/repl"
)

const (
	// Version 版本信息
	Version = "1.1.0"
)

func main() {
	// 定义命令行参数
	interactive := flag.Bool("i", false, "Interactive mode")
	interactiveLong := flag.Bool("interactive", false, "Interactive mode")
	version := flag.Bool("v", false, "Show version")
	versionLong := flag.Bool("version", false, "Show version")
	help := flag.Bool("h", false, "Show help")
	helpLong := flag.Bool("help", false, "Show help")
	modePtr := flag.String("mode", "vm", "Execution mode: ast or vm")
	outputFile := flag.String("o", "", "Output bytecode file (for offline execution)")

	flag.Parse()

	// 显示帮助
	if *help || *helpLong {
		printHelp()
		return
	}

	// 显示版本
	if *version || *versionLong {
		fmt.Printf("zork-basic v%s\n", Version)
		fmt.Println("A high-performance BASIC interpreter written in Go")
		return
	}

	// 确定运行模式
	isInteractive := *interactive || *interactiveLong
	mode := *modePtr

	// 检查参数
	args := flag.Args()
	if len(args) == 0 && !isInteractive {
		// 无参数且不是交互模式，默认进入交互模式
		isInteractive = true
	}

	if isInteractive {
		if *outputFile != "" {
			fmt.Println("Warning: -o flag is ignored in interactive mode")
		}
		repl.Run(Version, mode)
	} else {
		// 执行模式
		if len(args) < 1 {
			fmt.Println("Usage: zork-basic [options] <program.bas>")
			fmt.Println("Run 'zork-basic -h' for more information")
			os.Exit(1)
		}

		if *outputFile != "" {
			compileFileToBytecode(args[0], *outputFile)
		} else {
			runFileMode(args[0], mode)
		}
	}
}

// compileFileToBytecode 编译文件为字节码并保存
func compileFileToBytecode(inputFile, outputFile string) {
	fmt.Printf("Compiling %s to %s...\n", inputFile, outputFile)
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	parsedAST, err := parser.Parse(inputFile, data)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		os.Exit(1)
	}

	prog, ok := parsedAST.(*ast.Program)
	if !ok {
		fmt.Println("Error: not a valid BASIC program")
		os.Exit(1)
	}

	comp := compiler.New()
	chunk, err := comp.Compile(prog)
	if err != nil {
		fmt.Printf("Compilation error: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := chunk.Write(f); err != nil {
		fmt.Printf("Error writing bytecode: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Compilation successful.")
}

// runFileMode 执行文件模式
func runFileMode(filename string, mode string) {
	fmt.Printf("Running Basic program: %s (mode: %s)\n", filename, mode)

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	repl.ExecuteProgram(string(data), filename, mode)
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("zork-basic - A high-performance BASIC interpreter written in Go")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  zork-basic [options] <program.bas>")
	fmt.Println("  zork-basic [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -i, --interactive    Run in interactive mode")
	fmt.Println("  -v, --version        Show version information")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Println("  -mode <ast|vm>       Execution engine (default: vm)")
	fmt.Println("  -o <file>            Compile to a bytecode file instead of running")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  zork-basic program.bas      Execute a BASIC program")
	fmt.Println("  zork-basic -o prog.zbc p.bas Compile p.bas to prog.zbc")
	fmt.Println("  zork-basic -i               Start interactive mode")
}
