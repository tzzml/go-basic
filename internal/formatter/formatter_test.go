package formatter_test

import (
	"testing"
	"zork-basic/internal/ast"
	"zork-basic/internal/formatter"
)

func TestFormatInputStmt(t *testing.T) {
	tests := []struct {
		name     string
		stmt     *ast.InputStmt
		expected string
	}{
		{
			name: "Prompt and One Var",
			stmt: &ast.InputStmt{
				Prompt: "Enter value:",
				Vars:   []string{"A"},
			},
			expected: "INPUT \"Enter value:\", A",
		},
		{
			name: "No Prompt and One Var",
			stmt: &ast.InputStmt{
				Vars: []string{"A"},
			},
			expected: "INPUT A",
		},
		{
			name: "Prompt and Two Vars",
			stmt: &ast.InputStmt{
				Prompt: "Enter coordinates:",
				Vars:   []string{"X", "Y"},
			},
			expected: "INPUT \"Enter coordinates:\", X Y",
		},
		{
			name: "No Prompt and Two Vars",
			stmt: &ast.InputStmt{
				Vars: []string{"A", "B"},
			},
			expected: "INPUT A B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatter.FormatInputStmt(tt.stmt)
			if got != tt.expected {
				t.Errorf("FormatInputStmt() = %q, want %q", got, tt.expected)
			}
		})
	}
}
