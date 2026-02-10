package main

import (
	"fmt"
	"math"
	"time"

	"zork-basic/internal/ast"
	"zork-basic/internal/bytecode"
	"zork-basic/internal/compiler"
	"zork-basic/internal/interpreter"
	"zork-basic/internal/parser"
	"zork-basic/internal/vm"
)

const (
	Iterations = 10000000
	NativeIter = 10000000
)

func main() {
	fmt.Println("=== BASIC Interpreter Performance Benchmark ===")
	fmt.Printf("Task: Calculate sum of SIN(i) for i = 1 to %d\n\n", Iterations)

	// 1. Native Go Benchmark
	runNative()

	// Prepare BASIC Code (GOTO loop version)
	gotoCode := fmt.Sprintf(`10 LET SUM = 0
20 LET I = 1
30 IF I > %d THEN GOTO 70
40 LET SUM = SUM + SIN(I)
50 LET I = I + 1
60 GOTO 30
70 PRINT "SUM =", SUM
80 END
`, Iterations)

	// Prepare BASIC Code (FOR/NEXT loop version)
	forCode := fmt.Sprintf(`10 LET SUM = 0
20 FOR I = 1 TO %d
30 LET SUM = SUM + SIN(I)
40 NEXT I
50 PRINT "SUM =", SUM
60 END
`, Iterations)

	// Parse programs
	gotoProg := mustParse(gotoCode)
	forProg := mustParse(forCode)

	// 2. AST Interpreter Benchmark (GOTO)
	runAST("AST Interpreter (GOTO)", gotoProg)

	// 3. AST Interpreter Benchmark (FOR/NEXT)
	runAST("AST Interpreter (FOR/NEXT)", forProg)

	// 4. VM Benchmark (GOTO)
	runVM("Bytecode VM (GOTO)", gotoProg)

	// 5. VM Benchmark (FOR/NEXT)
	runVM("Bytecode VM (FOR/NEXT)", forProg)
}

func mustParse(code string) *ast.Program {
	prog, err := parser.Parse("benchmark", []byte(code))
	if err != nil {
		panic(err)
	}
	return prog.(*ast.Program)
}

func runNative() {
	fmt.Println("--- Native Go ---")
	start := time.Now()
	sum := 0.0
	for i := 1; i <= NativeIter; i++ {
		sum += math.Sin(float64(i))
	}
	elapsed := time.Since(start)
	fmt.Printf("Result: %g\n", sum)
	fmt.Printf("Time:   %v\n", elapsed)
	if elapsed > 0 {
		ops := float64(NativeIter) / elapsed.Seconds()
		fmt.Printf("Speed:  %.2f ops/sec\n", ops)
	}
	fmt.Println()
}

func runAST(name string, prog *ast.Program) {
	fmt.Printf("--- %s ---\n", name)
	interp := interpreter.NewInterpreter()

	start := time.Now()
	interp.ExecuteProgram(prog)
	elapsed := time.Since(start)

	fmt.Printf("Time:   %v\n", elapsed)
	if elapsed > 0 {
		ops := float64(Iterations) / elapsed.Seconds()
		fmt.Printf("Speed:  %.2f ops/sec\n", ops)
	}
	fmt.Println()
}

func runVM(name string, prog *ast.Program) {
	fmt.Printf("--- %s ---\n", name)
	// Compile
	comp := compiler.New()
	chunk, err := comp.Compile(prog)
	if err != nil {
		panic(err)
	}

	// Dump bytecode size
	fmt.Printf("Bytecode Size: %d bytes\n", len(chunk.Code))
	dumpChunk(chunk)

	// execution
	vm := vm.New(chunk)
	start := time.Now()
	if err := vm.Run(); err != nil {
		panic(err)
	}
	elapsed := time.Since(start)

	fmt.Printf("Time:   %v\n", elapsed)
	if elapsed > 0 {
		ops := float64(Iterations) / elapsed.Seconds()
		fmt.Printf("Speed:  %.2f ops/sec\n", ops)
	}
	fmt.Println()
}

func dumpChunk(c *bytecode.Chunk) {
	// Optional: print disassembly if small enough
	// fmt.Println(c.Disassemble("Benchmark Chunk"))
}
