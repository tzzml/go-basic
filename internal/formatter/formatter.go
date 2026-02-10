// Package formatter 提供 BASIC 代码格式化功能
// 包括行号重编号、关键字大写化和跳转目标更新
package formatter

import (
	"fmt"
	"strings"

	"zork-basic/internal/ast"
)

// FormatLine 格式化单行代码，更新行号引用
func FormatLine(line *ast.Line, lineNumberMap map[int]int) string {
	var result strings.Builder

	for _, stmt := range line.Statements {
		result.WriteString(FormatStatement(stmt, lineNumberMap))
		result.WriteString(" ")
	}

	return strings.TrimSpace(result.String())
}

// FormatStatement 格式化单个语句，更新行号引用
func FormatStatement(stmt ast.Node, lineNumberMap map[int]int) string {
	switch s := stmt.(type) {
	case *ast.GotoStmt:
		// 更新 GOTO 目标行号
		if newNum, ok := lineNumberMap[s.LineNumber]; ok {
			return fmt.Sprintf("GOTO %d", newNum)
		}
		return fmt.Sprintf("GOTO %d", s.LineNumber)

	case *ast.GosubStmt:
		// 更新 GOSUB 目标行号
		if newNum, ok := lineNumberMap[s.LineNumber]; ok {
			return fmt.Sprintf("GOSUB %d", newNum)
		}
		return fmt.Sprintf("GOSUB %d", s.LineNumber)

	case *ast.IfStmt:
		// 格式化 IF 语句
		thenPart := FormatStatements(s.ThenStmts, lineNumberMap)
		elsePart := ""
		if len(s.ElseStmts) > 0 {
			elsePart = " ELSE " + FormatStatements(s.ElseStmts, lineNumberMap)
		}

		// 如果是简单的单行 IF（THEN 和 ELSE 块都只有一个简单语句），使用紧凑格式
		isSimpleThen := len(s.ThenStmts) == 1
		isSimpleElse := len(s.ElseStmts) <= 1

		if isSimpleThen && isSimpleElse {
			// 检查是否是嵌套的 IF
			for _, stmt := range s.ThenStmts {
				if _, isIf := stmt.(*ast.IfStmt); isIf {
					isSimpleThen = false
					break
				}
			}
			for _, stmt := range s.ElseStmts {
				if _, isIf := stmt.(*ast.IfStmt); isIf {
					isSimpleElse = false
					break
				}
			}

			if isSimpleThen && isSimpleElse {
				return fmt.Sprintf("IF %s THEN %s%s", s.Condition.String(), thenPart, elsePart)
			}
		}

		// 否则使用多行格式
		return fmt.Sprintf("IF %s THEN %sEND IF", s.Condition.String(), thenPart+elsePart)

	case *ast.RemStmt:
		// RemStmt.Text 已经包含 "REM" 前缀，直接使用
		return s.Text

	case *ast.PrintStmt:
		return FormatPrintStmt(s)

	case *ast.Assignment:
		return fmt.Sprintf("%s = %s", s.Target.String(), s.Value.String())

	case *ast.InputStmt:
		return FormatInputStmt(s)

	case *ast.DimStmt:
		sizes := make([]string, len(s.Sizes))
		for i, size := range s.Sizes {
			sizes[i] = size.String()
		}
		return fmt.Sprintf("DIM %s(%s)", s.Name, strings.Join(sizes, ", "))

	case *ast.ForStmt:
		result := fmt.Sprintf("FOR %s = %s TO %s", s.Var, s.Start.String(), s.End.String())
		if s.Step != nil {
			if step, ok := s.Step.(*ast.Number); !ok || step.Value != 1 {
				result += fmt.Sprintf(" STEP %s", s.Step.String())
			}
		}
		return result

	case *ast.NextStmt:
		if s.Var != "" {
			return fmt.Sprintf("NEXT %s", s.Var)
		}
		return "NEXT"

	case *ast.ReturnStmt:
		return "RETURN"

	case *ast.EndStmt:
		return "END"

	default:
		return stmt.String()
	}
}

// FormatStatements 格式化语句列表
func FormatStatements(stmts []ast.Node, lineNumberMap map[int]int) string {
	var result strings.Builder
	for i, stmt := range stmts {
		if i > 0 {
			result.WriteString(": ")
		}
		result.WriteString(FormatStatement(stmt, lineNumberMap))
	}
	return result.String()
}

// FormatPrintStmt 格式化 PRINT 语句，保留原始分隔符
func FormatPrintStmt(stmt *ast.PrintStmt) string {
	if len(stmt.Values) == 0 {
		return "PRINT"
	}
	var result strings.Builder
	result.WriteString("PRINT ")
	for i, v := range stmt.Values {
		if i > 0 {
			// 使用原始分隔符（逗号或分号）
			if i-1 < len(stmt.Separators) {
				result.WriteString(stmt.Separators[i-1])
				result.WriteString(" ")
			} else {
				result.WriteString("; ")
			}
		}
		result.WriteString(v.String())
	}
	// 保留末尾分隔符
	if stmt.Trailer != "" {
		result.WriteString(stmt.Trailer)
	}
	return result.String()
}

// FormatInputStmt 格式化 INPUT 语句
func FormatInputStmt(stmt *ast.InputStmt) string {
	var result strings.Builder
	result.WriteString("INPUT")
	if stmt.Prompt != "" {
		result.WriteString(" \"")
		result.WriteString(stmt.Prompt)
		result.WriteString("\", ")
	}

	for i, v := range stmt.Vars {
		// 变量之间以及 INPUT/Prompt 与第一个变量之间需要分隔
		if i > 0 || stmt.Prompt == "" {
			result.WriteString(" ")
		}
		result.WriteString(v)
	}
	return result.String()
}
