package bytecode

import (
	"encoding/binary"
	"fmt"
	"strings"
	"zork-basic/internal/interpreter"
)

// Chunk represents a sequence of bytecode instructions and data
type Chunk struct {
	Code        []byte
	Constants   []interpreter.Value
	Lines       []int // Map bytecode offset to source line number
	GlobalCount int   // Number of global variables used
	ArrayCount  int   // Number of arrays used
}

// NewChunk creates a new Chunk
func NewChunk() *Chunk {
	return &Chunk{
		Code:        make([]byte, 0),
		Constants:   make([]interpreter.Value, 0),
		Lines:       make([]int, 0),
		GlobalCount: 0,
		ArrayCount:  0,
	}
}

// WriteByte writes a single byte to the chunk
func (c *Chunk) WriteByte(b byte, line int) {
	c.Code = append(c.Code, b)
	c.Lines = append(c.Lines, line)
}

// AddConstant adds a constant to the pool and returns its index
func (c *Chunk) AddConstant(value interpreter.Value) int {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

// Disassemble returns a string representation of the chunk for debugging
func (c *Chunk) Disassemble(name string) string {
	var out strings.Builder
	fmt.Fprintf(&out, "== %s ==\n", name)

	offset := 0
	for offset < len(c.Code) {
		offset = c.disassembleInstruction(&out, offset)
	}

	return out.String()
}

func (c *Chunk) disassembleInstruction(out *strings.Builder, offset int) int {
	fmt.Fprintf(out, "%04d ", offset)

	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Fprint(out, "   | ")
	} else {
		fmt.Fprintf(out, "%4d ", c.Lines[offset])
	}

	op := OpCode(c.Code[offset])
	def, err := Lookup(op)
	if err != nil {
		fmt.Fprintf(out, "Unknown opcode %d\n", op)
		return offset + 1
	}

	fmt.Fprintf(out, "%-16s", def.Name)

	offset++

	for _, width := range def.OperandWidths {
		switch width {
		case 2:
			val := binary.BigEndian.Uint16(c.Code[offset:])
			fmt.Fprintf(out, "%d ", val)
			offset += 2
		case 1:
			val := c.Code[offset]
			fmt.Fprintf(out, "%d ", val)
			offset++
		}
	}

	fmt.Fprint(out, "\n")
	return offset
}
