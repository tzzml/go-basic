package bytecode

import "fmt"

// OpCode represents a bytecode instruction
type OpCode byte

const (
	// OpConstant pushes a constant from the constant pool onto the stack
	// Operand: 2 bytes (index in constant pool)
	OpConstant OpCode = iota

	// OpPop pops the top element from the stack
	OpPop

	// Arithmetic operations
	OpAdd // +
	OpSub // -
	OpMul // *
	OpDiv // /
	OpPow // ^
	OpMod // MOD

	// Unary operations
	OpNeg // - (unary)
	OpNot // NOT

	// Logical operations
	OpAnd // AND
	OpOr  // OR

	// Comparison operations
	OpEq  // =
	OpNeq // <>
	OpGt  // >
	OpGte // >=
	OpLt  // <
	OpLte // <=

	// Control flow
	OpJump        // Unconditional jump. Operand: 2 bytes (offset)
	OpJumpIfFalse // Jump if stack top is false. Operand: 2 bytes (offset)
	OpGosub       // Call subroutine. Operand: 2 bytes (line index)
	OpReturn      // Return from subroutine
	OpEnd         // Terminate program
	OpForInit     // Initialize FOR loop. Operand: 2 bytes (loop variable index in globals). Pops step, end from stack.
	OpNext        // FOR loop next iteration. Operands: 2 bytes (variable index), 2 bytes (loop top offset)

	// Variable access
	OpGetGlobal // Get global variable. Operand: 2 bytes (index in global names)
	OpSetGlobal // Set global variable. Operand: 2 bytes (index in global names)

	// Array access
	OpGetArray // Get array element. Operands: 2 bytes (array name index), 1 byte (dimensions count)
	OpSetArray // Set array element. Operands: 2 bytes (array name index), 1 byte (dimensions count)

	// I/O
	OpPrint   // Print stack top (no newline)
	OpPrintNl // Print newline
	OpInput   // Read input. Operand: 2 bytes (variable name index)

	// Builtin calls
	OpCallBuiltin // Call builtin function. Operands: 2 bytes (name index), 1 byte (arg count)

	// Array declaration
	OpDim // Declare array. Operands: 2 bytes (array name index), 1 byte (dimensions count)
)

// OpDefinition defines the properties of an opcode
type OpDefinition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*OpDefinition{
	OpConstant:    {"OpConstant", []int{2}},
	OpPop:         {"OpPop", []int{}},
	OpAdd:         {"OpAdd", []int{}},
	OpSub:         {"OpSub", []int{}},
	OpMul:         {"OpMul", []int{}},
	OpDiv:         {"OpDiv", []int{}},
	OpPow:         {"OpPow", []int{}},
	OpMod:         {"OpMod", []int{}},
	OpNeg:         {"OpNeg", []int{}},
	OpNot:         {"OpNot", []int{}},
	OpAnd:         {"OpAnd", []int{}},
	OpOr:          {"OpOr", []int{}},
	OpEq:          {"OpEq", []int{}},
	OpNeq:         {"OpNeq", []int{}},
	OpGt:          {"OpGt", []int{}},
	OpGte:         {"OpGte", []int{}},
	OpLt:          {"OpLt", []int{}},
	OpLte:         {"OpLte", []int{}},
	OpJump:        {"OpJump", []int{2}},
	OpJumpIfFalse: {"OpJumpIfFalse", []int{2}},
	OpGosub:       {"OpGosub", []int{2}}, // Index in line map (not byte offset)
	OpReturn:      {"OpReturn", []int{}},
	OpEnd:         {"OpEnd", []int{}},
	OpForInit:     {"OpForInit", []int{2}},
	OpNext:        {"OpNext", []int{2, 2}},
	OpGetGlobal:   {"OpGetGlobal", []int{2}},
	OpSetGlobal:   {"OpSetGlobal", []int{2}},
	OpGetArray:    {"OpGetArray", []int{2, 1}},
	OpSetArray:    {"OpSetArray", []int{2, 1}},
	OpPrint:       {"OpPrint", []int{}},
	OpPrintNl:     {"OpPrintNl", []int{}},
	OpInput:       {"OpInput", []int{2}},
	OpCallBuiltin: {"OpCallBuiltin", []int{2, 1}},
	OpDim:         {"OpDim", []int{2, 1}},
}

// Lookup returns the definition for an opcode
func Lookup(op OpCode) (*OpDefinition, error) {
	def, ok := definitions[op]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}
