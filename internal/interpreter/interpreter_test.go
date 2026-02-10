package interpreter_test

import (
	"io"
	"testing"
	"zork-basic/internal/ast"
	"zork-basic/internal/interpreter"
)

// BenchmarkSinLoop 测试 1,000,000 次 SIN 计算的性能
// 对应 PERFORMANCE.md 中的核心测试用例
func BenchmarkSinLoop(b *testing.B) {
	// 构造 AST:
	// 5 SUM = 0
	// 10 FOR I = 1 TO 1000
	// 20 SUM = SUM + SIN(I)
	// 30 NEXT I

	program := &ast.Program{
		Lines: []*ast.Line{
			{
				LineNumber: 5,
				Statements: []ast.Node{
					&ast.Assignment{
						Target: &ast.Identifier{Name: "SUM"},
						Value:  &ast.Number{Value: 0},
					},
				},
			},
			{
				LineNumber: 10,
				Statements: []ast.Node{
					&ast.ForStmt{
						Var:   "I",
						Start: &ast.Number{Value: 1.0},
						End:   &ast.Number{Value: 1000.0},
						Step:  &ast.Number{Value: 1.0},
					},
				},
			},
			{
				LineNumber: 20,
				Statements: []ast.Node{
					&ast.Assignment{
						Target: &ast.Identifier{Name: "SUM"},
						Value: &ast.BinaryOp{
							Left: &ast.Identifier{Name: "SUM"},
							Op:   "+",
							Right: &ast.FunctionCall{
								Name: "SIN",
								Args: []ast.Node{
									&ast.Identifier{Name: "I"},
								},
							},
						},
					},
				},
			},
			{
				LineNumber: 30,
				Statements: []ast.Node{
					&ast.NextStmt{Var: "I"},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		interp := interpreter.NewInterpreter(interpreter.WithOutput(io.Discard))
		interp.ExecuteProgram(program)
	}
}
