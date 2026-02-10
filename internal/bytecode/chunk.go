package bytecode

import (
	"encoding/binary"
	"fmt"
	"io"
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

// Write serializes the chunk to a writer
func (c *Chunk) Write(w io.Writer) error {
	// Magic header: ZBC (Zork Basic Compiled) + version 1
	if _, err := w.Write([]byte("ZBC\x01")); err != nil {
		return err
	}

	// Meta counts
	if err := binary.Write(w, binary.BigEndian, uint16(c.GlobalCount)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint16(c.ArrayCount)); err != nil {
		return err
	}

	// Constants
	if err := binary.Write(w, binary.BigEndian, uint16(len(c.Constants))); err != nil {
		return err
	}
	for _, val := range c.Constants {
		if val.IsNumber() {
			if _, err := w.Write([]byte{1}); err != nil {
				return err
			}
			if err := binary.Write(w, binary.BigEndian, val.AsNumber()); err != nil {
				return err
			}
		} else {
			if _, err := w.Write([]byte{2}); err != nil {
				return err
			}
			str := val.String()
			if err := binary.Write(w, binary.BigEndian, uint16(len(str))); err != nil {
				return err
			}
			if _, err := w.Write([]byte(str)); err != nil {
				return err
			}
		}
	}

	// Code
	if err := binary.Write(w, binary.BigEndian, uint32(len(c.Code))); err != nil {
		return err
	}
	if _, err := w.Write(c.Code); err != nil {
		return err
	}

	// Lines
	if err := binary.Write(w, binary.BigEndian, uint32(len(c.Lines))); err != nil {
		return err
	}
	for _, line := range c.Lines {
		if err := binary.Write(w, binary.BigEndian, uint32(line)); err != nil {
			return err
		}
	}

	return nil
}

// ReadChunk deserializes a chunk from a reader
func ReadChunk(r io.Reader) (*Chunk, error) {
	// Magic header
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}
	if string(header[:3]) != "ZBC" {
		return nil, fmt.Errorf("invalid bytecode header")
	}

	c := NewChunk()

	// Meta counts
	var globalCount, arrayCount uint16
	if err := binary.Read(r, binary.BigEndian, &globalCount); err != nil {
		return nil, err
	}
	c.GlobalCount = int(globalCount)
	if err := binary.Read(r, binary.BigEndian, &arrayCount); err != nil {
		return nil, err
	}
	c.ArrayCount = int(arrayCount)

	// Constants
	var constCount uint16
	if err := binary.Read(r, binary.BigEndian, &constCount); err != nil {
		return nil, err
	}
	c.Constants = make([]interpreter.Value, constCount)
	for i := 0; i < int(constCount); i++ {
		var typ byte
		if err := binary.Read(r, binary.BigEndian, &typ); err != nil {
			return nil, err
		}
		if typ == 1 {
			var num float64
			if err := binary.Read(r, binary.BigEndian, &num); err != nil {
				return nil, err
			}
			c.Constants[i] = interpreter.NumberValue(num)
		} else if typ == 2 {
			var length uint16
			if err := binary.Read(r, binary.BigEndian, &length); err != nil {
				return nil, err
			}
			strBuf := make([]byte, length)
			if _, err := io.ReadFull(r, strBuf); err != nil {
				return nil, err
			}
			c.Constants[i] = interpreter.StringValue(string(strBuf))
		}
	}

	// Code
	var codeLen uint32
	if err := binary.Read(r, binary.BigEndian, &codeLen); err != nil {
		return nil, err
	}
	c.Code = make([]byte, codeLen)
	if _, err := io.ReadFull(r, c.Code); err != nil {
		return nil, err
	}

	// Lines
	var linesLen uint32
	if err := binary.Read(r, binary.BigEndian, &linesLen); err != nil {
		return nil, err
	}
	c.Lines = make([]int, linesLen)
	for i := 0; i < int(linesLen); i++ {
		var line uint32
		if err := binary.Read(r, binary.BigEndian, &line); err != nil {
			return nil, err
		}
		c.Lines[i] = int(line)
	}

	return c, nil
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

	if offset > 0 && offset < len(c.Lines) && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Fprint(out, "   | ")
	} else if offset < len(c.Lines) {
		fmt.Fprintf(out, "%4d ", c.Lines[offset])
	} else {
		fmt.Fprint(out, "   ? ")
	}

	op := OpCode(c.Code[offset])
	def, err := Lookup(op)
	if err != nil {
		fmt.Fprintf(out, "Unknown opcode %d\n", op)
		return offset + 1
	}

	fmt.Fprintf(out, "%-16s", strings.TrimPrefix(def.Name, "Op"))

	offset++

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			val := binary.BigEndian.Uint16(c.Code[offset:])
			fmt.Fprintf(out, "%d ", val)

			// Special handling for instructions that reference pools
			if op == OpConstant && i == 0 {
				if int(val) < len(c.Constants) {
					constVal := c.Constants[val]
					if constVal.IsString() {
						fmt.Fprintf(out, "(\"%s\") ", constVal.String())
					} else {
						fmt.Fprintf(out, "(%s) ", constVal.String())
					}
				}
			} else if op == OpCallBuiltin && i == 0 {
				if int(val) < len(BuiltinNames) {
					fmt.Fprintf(out, "(%s) ", BuiltinNames[val])
				}
			}

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
