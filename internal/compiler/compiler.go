package compiler

import (
	"encoding/binary"
	"fmt"
	"strings"

	"zork-basic/internal/ast"
	"zork-basic/internal/bytecode"
	"zork-basic/internal/interpreter"
)

// forInfo tracks a FOR loop's compilation state
type forInfo struct {
	varName string // Loop variable name (uppercased)
	varIdx  int    // Index of loop variable in globals
	loopTop int    // Bytecode offset of the loop body start (after OpForInit)
}

// Compiler translates AST to bytecode
type Compiler struct {
	chunk       *bytecode.Chunk
	lineOffsets map[int]int    // map[BasicLineNumber]BytecodeOffset
	fixups      map[int][]int  // map[BasicLineNumber][]BytecodeOffsetToPatch
	currentLine int            // Current source line number being compiled
	globals     map[string]int // map[Name]Index (Global variables)
	arrays      map[string]int // map[Name]Index (Arrays)
	globalCount int
	arrayCount  int
	forStack    []forInfo // FOR loop stack for matching FOR/NEXT
}

// New creates a new Compiler
func New() *Compiler {
	return &Compiler{
		chunk:       bytecode.NewChunk(),
		lineOffsets: make(map[int]int),
		fixups:      make(map[int][]int),
		globals:     make(map[string]int),
		arrays:      make(map[string]int),
	}
}

// Compile compiles a program into a chunk
func (c *Compiler) Compile(prog *ast.Program) (*bytecode.Chunk, error) {
	for _, line := range prog.Lines {
		c.currentLine = line.LineNumber
		// Record the bytecode offset for this line
		c.lineOffsets[line.LineNumber] = len(c.chunk.Code)

		for _, stmt := range line.Statements {
			if err := c.compileStatement(stmt); err != nil {
				return nil, err
			}
		}
	}

	// Terminate program
	c.emit(bytecode.OpEnd)

	// Resolve fixups (GOTO/GOSUB targets)
	for lineNum, offsets := range c.fixups {
		targetOffset, ok := c.lineOffsets[lineNum]
		if !ok {
			return nil, fmt.Errorf("undefined line number %d", lineNum)
		}

		for _, offset := range offsets {
			// Write the 2-byte target offset (offset in bytecode, not line number)
			binary.BigEndian.PutUint16(c.chunk.Code[offset:], uint16(targetOffset))
		}
	}

	// Store counts in chunk
	c.chunk.GlobalCount = c.globalCount
	c.chunk.ArrayCount = c.arrayCount

	return c.chunk, nil
}

func (c *Compiler) compileStatement(stmt ast.Node) error {
	switch n := stmt.(type) {
	case *ast.Assignment:
		if err := c.compileExpression(n.Value); err != nil {
			return err
		}
		switch target := n.Target.(type) {
		case *ast.Identifier:
			name := strings.ToUpper(target.Name)
			idx := c.resolveGlobal(name)
			c.emit(bytecode.OpSetGlobal, byte(idx>>8), byte(idx))
		case *ast.ArrayAccess:
			// Compile indices
			for _, idxExpr := range target.Indices {
				if err := c.compileExpression(idxExpr); err != nil {
					return err
				}
			}
			name := strings.ToUpper(target.Name)
			idx := c.resolveArray(name)
			c.emit(bytecode.OpSetArray, byte(idx>>8), byte(idx), byte(len(target.Indices)))
		default:
			return fmt.Errorf("invalid assignment target: %T", target)
		}

	case *ast.PrintStmt:
		for i, val := range n.Values {
			// Print value
			if err := c.compileExpression(val); err != nil {
				return err
			}
			c.emit(bytecode.OpPrint)

			// Handle separator if this is not the last item, or if there is a trailing separator
			// AST structure: Values[i] followed by Separators[i] if it exists
			sep := ""
			if i < len(n.Separators) {
				sep = n.Separators[i]
			} else if i == len(n.Values)-1 {
				sep = n.Trailer
			}

			if sep == "," {
				// Comma: print a tab or space? Interpreter uses " ".
				spaceIdx := c.addConstant(interpreter.StringValue(" "))
				c.emit(bytecode.OpConstant, byte(spaceIdx>>8), byte(spaceIdx))
				c.emit(bytecode.OpPrint)
			}
			// Semicolon: no space, do nothing
		}

		// If no trailer, print newline
		if n.Trailer == "" {
			c.emit(bytecode.OpPrintNl)
		}

	case *ast.IfStmt:
		// Condition
		if err := c.compileExpression(n.Condition); err != nil {
			return err
		}

		// Jump if false -> to ELSE or END IF
		jumpIfFalseOffset := c.emitJump(bytecode.OpJumpIfFalse)

		// THEN block
		for _, s := range n.ThenStmts {
			if err := c.compileStatement(s); err != nil {
				return err
			}
		}

		// Jump -> to END IF (skip ELSE)
		jumpToExitOffset := c.emitJump(bytecode.OpJump)

		// Patch JumpIfFalse to here (start of ELSE)
		c.patchJump(jumpIfFalseOffset)

		// ELSE block
		for _, s := range n.ElseStmts {
			if err := c.compileStatement(s); err != nil {
				return err
			}
		}

		// Patch JumpToExit to here (end of IF)
		c.patchJump(jumpToExitOffset)

	case *ast.ForStmt:
		varName := strings.ToUpper(n.Var)
		idx := c.resolveGlobal(varName)

		// Compile: SET VAR = START
		if err := c.compileExpression(n.Start); err != nil {
			return err
		}
		c.emit(bytecode.OpSetGlobal, byte(idx>>8), byte(idx))

		// Compile end and step expressions, push them on stack for OpForInit
		if err := c.compileExpression(n.End); err != nil {
			return err
		}
		if err := c.compileExpression(n.Step); err != nil {
			return err
		}

		// Emit OpForInit: pops step and end from stack, stores in ForFrame
		c.emit(bytecode.OpForInit, byte(idx>>8), byte(idx))

		// Record loop body start (the ip AFTER OpForInit)
		loopTop := len(c.chunk.Code)

		// Push to compiler's for stack
		c.forStack = append(c.forStack, forInfo{
			varName: varName,
			varIdx:  idx,
			loopTop: loopTop,
		})

	case *ast.NextStmt:
		if len(c.forStack) == 0 {
			return fmt.Errorf("NEXT without FOR")
		}

		// Pop the matching FOR info
		frame := c.forStack[len(c.forStack)-1]
		c.forStack = c.forStack[:len(c.forStack)-1]

		// If NEXT specifies a variable, verify it matches
		if n.Var != "" {
			nextVar := strings.ToUpper(n.Var)
			if nextVar != frame.varName {
				return fmt.Errorf("NEXT %s does not match FOR %s", nextVar, frame.varName)
			}
		}

		// Emit OpNext with variable index and loop top offset
		idx := frame.varIdx
		loopTop := frame.loopTop
		c.emit(bytecode.OpNext, byte(idx>>8), byte(idx), byte(loopTop>>8), byte(loopTop))

	case *ast.GotoStmt:
		c.emit(bytecode.OpJump, 0, 0) // Placeholder
		offset := len(c.chunk.Code) - 2
		c.fixups[n.LineNumber] = append(c.fixups[n.LineNumber], offset)

	case *ast.GosubStmt:
		c.emit(bytecode.OpGosub, 0, 0) // Placeholder
		offset := len(c.chunk.Code) - 2
		c.fixups[n.LineNumber] = append(c.fixups[n.LineNumber], offset)

	case *ast.ReturnStmt:
		c.emit(bytecode.OpReturn)

	case *ast.EndStmt:
		c.emit(bytecode.OpEnd)

	case *ast.InputStmt:
		if n.Prompt != "" {
			idx := c.addConstant(interpreter.StringValue(n.Prompt))
			c.emit(bytecode.OpConstant, byte(idx>>8), byte(idx))
			c.emit(bytecode.OpPrint) // Print prompt
		}
		for _, varName := range n.Vars {
			name := strings.ToUpper(varName)
			idx := c.resolveGlobal(name)
			c.emit(bytecode.OpInput, byte(idx>>8), byte(idx))
		}

	case *ast.RemStmt:
		// Ignore comments

	case *ast.DimStmt:
		// Compile dimension expressions
		for _, sizeExpr := range n.Sizes {
			if err := c.compileExpression(sizeExpr); err != nil {
				return err
			}
		}

		name := strings.ToUpper(n.Name)
		idx := c.resolveArray(name)

		// Emit OpDim with name index and dimension count
		c.emit(bytecode.OpDim, byte(idx>>8), byte(idx), byte(len(n.Sizes)))

	default:
		return fmt.Errorf("unknown statement: %T", stmt)
	}
	return nil
}

func (c *Compiler) compileExpression(expr ast.Node) error {
	switch n := expr.(type) {
	case *ast.Number:
		idx := c.addConstant(interpreter.NumberValue(n.Value))
		c.emit(bytecode.OpConstant, byte(idx>>8), byte(idx))

	case *ast.StringLiteral:
		idx := c.addConstant(interpreter.StringValue(n.Value))
		c.emit(bytecode.OpConstant, byte(idx>>8), byte(idx))

	case *ast.Identifier:
		name := strings.ToUpper(n.Name)
		idx := c.resolveGlobal(name)
		c.emit(bytecode.OpGetGlobal, byte(idx>>8), byte(idx))

	case *ast.FunctionCall:
		for _, arg := range n.Args {
			if err := c.compileExpression(arg); err != nil {
				return err
			}
		}
		name := strings.ToUpper(n.Name)
		builtinID := bytecode.GetBuiltinID(name)
		if builtinID < 0 {
			return fmt.Errorf("unknown builtin function: %s", name)
		}
		c.emit(bytecode.OpCallBuiltin, byte(builtinID>>8), byte(builtinID), byte(len(n.Args)))

	case *ast.BinaryOp:
		if err := c.compileExpression(n.Left); err != nil {
			return err
		}
		if err := c.compileExpression(n.Right); err != nil {
			return err
		}
		switch n.Op {
		case "+":
			c.emit(bytecode.OpAdd)
		case "-":
			c.emit(bytecode.OpSub)
		case "*":
			c.emit(bytecode.OpMul)
		case "/":
			c.emit(bytecode.OpDiv)
		case "^":
			c.emit(bytecode.OpPow)
		case "MOD":
			c.emit(bytecode.OpMod)
		default:
			return fmt.Errorf("unknown binary op: %s", n.Op)
		}

	case *ast.ComparisonOp:
		if err := c.compileExpression(n.Left); err != nil {
			return err
		}
		if err := c.compileExpression(n.Right); err != nil {
			return err
		}
		switch n.Op {
		case "=":
			c.emit(bytecode.OpEq)
		case "<>":
			c.emit(bytecode.OpNeq)
		case ">":
			c.emit(bytecode.OpGt)
		case "<":
			c.emit(bytecode.OpLt)
		case ">=":
			c.emit(bytecode.OpGte)
		case "<=":
			c.emit(bytecode.OpLte)
		default:
			return fmt.Errorf("unknown comparison op: %s", n.Op)
		}

	case *ast.LogicalOp:
		if err := c.compileExpression(n.Left); err != nil {
			return err
		}
		if err := c.compileExpression(n.Right); err != nil {
			return err
		}
		switch n.Op {
		case "AND":
			c.emit(bytecode.OpAnd)
		case "OR":
			c.emit(bytecode.OpOr)
		default:
			return fmt.Errorf("unknown logical op: %s", n.Op)
		}

	case *ast.UnaryOp:
		if err := c.compileExpression(n.Right); err != nil {
			return err
		}
		switch n.Op {
		case "-":
			c.emit(bytecode.OpNeg)
		case "NOT":
			c.emit(bytecode.OpNot)
		case "+": // no-op
		default:
			return fmt.Errorf("unknown unary op: %s", n.Op)
		}

	default:
		return fmt.Errorf("unknown expression: %T", expr)
	}
	return nil
}

func (c *Compiler) resolveGlobal(name string) int {
	if idx, ok := c.globals[name]; ok {
		return idx
	}
	idx := c.globalCount
	c.globals[name] = idx
	c.globalCount++
	return idx
}

func (c *Compiler) resolveArray(name string) int {
	if idx, ok := c.arrays[name]; ok {
		return idx
	}
	idx := c.arrayCount
	c.arrays[name] = idx
	c.arrayCount++
	return idx
}

func (c *Compiler) emit(op bytecode.OpCode, operands ...byte) {
	c.chunk.WriteByte(byte(op), c.currentLine)
	for _, b := range operands {
		c.chunk.WriteByte(b, c.currentLine)
	}
}

func (c *Compiler) emitJump(op bytecode.OpCode) int {
	c.emit(op, 0xff, 0xff)
	return len(c.chunk.Code) - 2
}

func (c *Compiler) patchJump(offset int) {
	target := len(c.chunk.Code)
	if target > 0xFFFF {
		panic("program too large")
	}
	c.chunk.Code[offset] = byte(target >> 8)
	c.chunk.Code[offset+1] = byte(target)
}

// addConstant adds a constant to the pool with deduplication.
// If an identical constant already exists, returns its index instead of adding a duplicate.
func (c *Compiler) addConstant(val interpreter.Value) int {
	// Check for existing identical constant
	for i, existing := range c.chunk.Constants {
		if existing == val {
			return i
		}
	}
	return c.chunk.AddConstant(val)
}
