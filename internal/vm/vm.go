package vm

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"zork-basic/internal/bytecode"
	"zork-basic/internal/interpreter"
)

const (
	StackSize = 2048
)

// Pre-cached boolean values to avoid allocation in hot path
var (
	valZero = interpreter.NumberValue(0)
	valOne  = interpreter.NumberValue(1)
)

// ForFrame stores the state of a FOR loop
type ForFrame struct {
	varIdx    int     // Index of loop variable in globals
	endValue  float64 // Loop end value
	stepValue float64 // Loop step value
	loopTop   int     // Bytecode offset to jump back to (after OpForInit)
}

// VM is the virtual machine
type VM struct {
	chunk   *bytecode.Chunk
	ip      int // Instruction pointer
	stack   []interpreter.Value
	sp      int // Stack pointer
	globals []interpreter.Value
	arrays  []*interpreter.ArrayInfo

	// I/O
	output    io.Writer
	errOutput io.Writer
	input     io.Reader

	// Call stack for GOSUB/RETURN
	returnStack []int

	// FOR loop stack
	forStack []ForFrame

	// Reusable buffer for array indices
	indexBuf []int
}

// Option represents a configuration option for the VM
type Option func(*VM)

func WithOutput(w io.Writer) Option {
	return func(vm *VM) { vm.output = w }
}

func WithErrOutput(w io.Writer) Option {
	return func(vm *VM) { vm.errOutput = w }
}

func WithInput(r io.Reader) Option {
	return func(vm *VM) { vm.input = r }
}

// New creates a new VM
func New(c *bytecode.Chunk, opts ...Option) *VM {
	// Initialize globals and arrays based on chunk counts
	globals := make([]interpreter.Value, c.GlobalCount)
	arrays := make([]*interpreter.ArrayInfo, c.ArrayCount)

	vm := &VM{
		chunk:       c,
		ip:          0,
		stack:       make([]interpreter.Value, StackSize),
		sp:          0,
		globals:     globals,
		arrays:      arrays,
		output:      os.Stdout,          // Default output
		errOutput:   os.Stderr,          // Default error output
		input:       os.Stdin,           // Default input
		returnStack: make([]int, 0, 16), // Pre-allocate capacity
		forStack:    make([]ForFrame, 0, 8),
	}
	for _, opt := range opts {
		opt(vm)
	}
	return vm
}

// pushUnchecked pushes a value onto the stack without bounds checking.
// Use only when we know the stack has space (e.g., after popping 2 values and pushing 1).
func (vm *VM) pushUnchecked(val interpreter.Value) {
	vm.stack[vm.sp] = val
	vm.sp++
}

// pushBool pushes a pre-cached boolean value (0 or 1) without bounds checking.
func (vm *VM) pushBool(cond bool) {
	if cond {
		vm.stack[vm.sp] = valOne
	} else {
		vm.stack[vm.sp] = valZero
	}
	vm.sp++
}

// getIndexBuf returns a reusable index buffer of the given size.
func (vm *VM) getIndexBuf(size int) []int {
	if cap(vm.indexBuf) >= size {
		return vm.indexBuf[:size]
	}
	vm.indexBuf = make([]int, size)
	return vm.indexBuf
}

// Run executes the bytecode
func (vm *VM) Run() error {
	code := vm.chunk.Code
	constants := vm.chunk.Constants
	globals := vm.globals

	for vm.ip < len(code) {
		op := bytecode.OpCode(code[vm.ip])
		vm.ip++

		switch op {
		case bytecode.OpConstant:
			constIdx := vm.readUint16()
			if err := vm.push(constants[constIdx]); err != nil {
				return err
			}

		case bytecode.OpPop:
			vm.pop()

		case bytecode.OpAdd:
			right := vm.pop()
			left := vm.pop()
			if left.IsString() || right.IsString() {
				vm.pushUnchecked(interpreter.StringValue(left.String() + right.String()))
			} else {
				vm.pushUnchecked(interpreter.NumberValue(left.AsNumber() + right.AsNumber()))
			}

		case bytecode.OpSub:
			right := vm.pop()
			left := vm.pop()
			vm.pushUnchecked(interpreter.NumberValue(left.AsNumber() - right.AsNumber()))

		case bytecode.OpMul:
			right := vm.pop()
			left := vm.pop()
			vm.pushUnchecked(interpreter.NumberValue(left.AsNumber() * right.AsNumber()))

		case bytecode.OpDiv:
			right := vm.pop()
			left := vm.pop()
			if right.AsNumber() == 0 {
				return fmt.Errorf("division by zero")
			}
			vm.pushUnchecked(interpreter.NumberValue(left.AsNumber() / right.AsNumber()))

		case bytecode.OpPow:
			right := vm.pop()
			left := vm.pop()
			vm.pushUnchecked(interpreter.NumberValue(math.Pow(left.AsNumber(), right.AsNumber())))

		case bytecode.OpMod:
			right := vm.pop()
			left := vm.pop()
			vm.pushUnchecked(interpreter.NumberValue(math.Mod(left.AsNumber(), right.AsNumber())))

		case bytecode.OpEq:
			right := vm.pop()
			left := vm.pop()
			if left.IsString() || right.IsString() {
				vm.pushBool(left.String() == right.String())
			} else {
				vm.pushBool(left.AsNumber() == right.AsNumber())
			}

		case bytecode.OpNeq:
			right := vm.pop()
			left := vm.pop()
			if left.IsString() || right.IsString() {
				vm.pushBool(left.String() != right.String())
			} else {
				vm.pushBool(left.AsNumber() != right.AsNumber())
			}

		case bytecode.OpGt:
			right := vm.pop()
			left := vm.pop()
			if left.IsString() || right.IsString() {
				vm.pushBool(left.String() > right.String())
			} else {
				vm.pushBool(left.AsNumber() > right.AsNumber())
			}

		case bytecode.OpGte:
			right := vm.pop()
			left := vm.pop()
			if left.IsString() || right.IsString() {
				vm.pushBool(left.String() >= right.String())
			} else {
				vm.pushBool(left.AsNumber() >= right.AsNumber())
			}

		case bytecode.OpLt:
			right := vm.pop()
			left := vm.pop()
			if left.IsString() || right.IsString() {
				vm.pushBool(left.String() < right.String())
			} else {
				vm.pushBool(left.AsNumber() < right.AsNumber())
			}

		case bytecode.OpLte:
			right := vm.pop()
			left := vm.pop()
			if left.IsString() || right.IsString() {
				vm.pushBool(left.String() <= right.String())
			} else {
				vm.pushBool(left.AsNumber() <= right.AsNumber())
			}

		case bytecode.OpAnd:
			right := vm.pop()
			left := vm.pop()
			vm.pushBool(left.IsTrue() && right.IsTrue())

		case bytecode.OpOr:
			right := vm.pop()
			left := vm.pop()
			vm.pushBool(left.IsTrue() || right.IsTrue())

		case bytecode.OpNeg:
			val := vm.pop()
			if !val.IsNumber() {
				return fmt.Errorf("operand must be a number")
			}
			vm.pushUnchecked(interpreter.NumberValue(-val.AsNumber()))

		case bytecode.OpNot:
			val := vm.pop()
			vm.pushBool(!val.IsTrue())

		case bytecode.OpPrint:
			val := vm.pop()
			fmt.Fprint(vm.output, val.String())

		case bytecode.OpPrintNl:
			fmt.Fprintln(vm.output)

		case bytecode.OpInput:
			nameIdx := vm.readUint16()
			var input string
			fmt.Fscanln(vm.input, &input)

			// Try parse number
			num, err := strconv.ParseFloat(input, 64)
			var val interpreter.Value
			if err != nil {
				val = interpreter.StringValue(input)
			} else {
				val = interpreter.NumberValue(num)
			}

			if int(nameIdx) >= len(globals) {
				return fmt.Errorf("global index out of bounds: %d", nameIdx)
			}
			globals[int(nameIdx)] = val

		case bytecode.OpJump:
			offset := vm.readUint16()
			vm.ip = int(offset)

		case bytecode.OpJumpIfFalse:
			offset := vm.readUint16()
			cond := vm.pop()
			if !cond.IsTrue() {
				vm.ip = int(offset)
			}

		case bytecode.OpGetGlobal:
			nameIdx := vm.readUint16()
			if err := vm.push(globals[int(nameIdx)]); err != nil {
				return err
			}

		case bytecode.OpSetGlobal:
			nameIdx := vm.readUint16()
			val := vm.pop()
			globals[int(nameIdx)] = val

		case bytecode.OpGetArray:
			nameIdx := vm.readUint16()
			dimCount := int(vm.readUint8())

			indices := vm.getIndexBuf(dimCount)
			for i := dimCount - 1; i >= 0; i-- {
				indices[i] = int(vm.pop().AsNumber())
			}

			arr := vm.arrays[int(nameIdx)]
			if arr == nil {
				return fmt.Errorf("array not declared (index %d)", nameIdx)
			}

			flatIdx := arr.CalculateIndex(indices)
			if flatIdx < 0 {
				return fmt.Errorf("array index out of bounds")
			}
			if err := vm.push(interpreter.NumberValue(arr.Data[flatIdx])); err != nil {
				return err
			}

		case bytecode.OpSetArray:
			nameIdx := vm.readUint16()
			dimCount := int(vm.readUint8())
			val := vm.pop()

			indices := vm.getIndexBuf(dimCount)
			for i := dimCount - 1; i >= 0; i-- {
				indices[i] = int(vm.pop().AsNumber())
			}

			arr := vm.arrays[int(nameIdx)]
			if arr == nil {
				return fmt.Errorf("array not declared (index %d)", nameIdx)
			}

			flatIdx := arr.CalculateIndex(indices)
			if flatIdx < 0 {
				return fmt.Errorf("array index out of bounds")
			}
			arr.Data[flatIdx] = val.AsNumber()

		case bytecode.OpGosub:
			target := vm.readUint16()
			// Push return address (current ip)
			vm.returnStack = append(vm.returnStack, vm.ip)
			vm.ip = int(target)

		case bytecode.OpReturn:
			if len(vm.returnStack) == 0 {
				return fmt.Errorf("return without gosub")
			}
			addr := vm.returnStack[len(vm.returnStack)-1]
			vm.returnStack = vm.returnStack[:len(vm.returnStack)-1]
			vm.ip = addr

		case bytecode.OpEnd:
			return nil

		case bytecode.OpForInit:
			varIdx := vm.readUint16()
			stepVal := vm.pop().AsNumber()
			endVal := vm.pop().AsNumber()

			// Push FOR frame with current ip as loop top
			vm.forStack = append(vm.forStack, ForFrame{
				varIdx:    int(varIdx),
				endValue:  endVal,
				stepValue: stepVal,
				loopTop:   vm.ip,
			})

		case bytecode.OpNext:
			varIdx := int(vm.readUint16())
			loopTop := int(vm.readUint16())

			if len(vm.forStack) == 0 {
				return fmt.Errorf("NEXT without FOR")
			}

			frame := &vm.forStack[len(vm.forStack)-1]

			// Safety check: variable index must match
			if frame.varIdx != varIdx {
				return fmt.Errorf("NEXT variable mismatch")
			}

			// Increment loop variable
			newVal := globals[varIdx].AsNumber() + frame.stepValue
			globals[varIdx] = interpreter.NumberValue(newVal)

			// Check loop condition
			shouldContinue := false
			if frame.stepValue > 0 && newVal <= frame.endValue {
				shouldContinue = true
			} else if frame.stepValue < 0 && newVal >= frame.endValue {
				shouldContinue = true
			}

			if shouldContinue {
				vm.ip = loopTop
			} else {
				// Pop frame
				vm.forStack = vm.forStack[:len(vm.forStack)-1]
			}

		case bytecode.OpDim:
			nameIdx := vm.readUint16()
			dimCount := int(vm.readUint8())

			dims := make([]int, dimCount)
			for i := dimCount - 1; i >= 0; i-- {
				val := vm.pop()
				dim := int(val.AsNumber())
				if dim < 0 {
					return fmt.Errorf("negative array dimension: %d", dim)
				}
				dims[i] = dim
			}

			if int(nameIdx) >= len(vm.arrays) {
				return fmt.Errorf("array index out of bounds: %d", nameIdx)
			}

			// Create new array info
			vm.arrays[int(nameIdx)] = interpreter.NewArrayInfo(dims)

		case bytecode.OpCallBuiltin:
			builtinIdx := vm.readUint16()
			argCount := int(vm.readUint8())

			if int(builtinIdx) >= len(builtinImpls) {
				return fmt.Errorf("unknown builtin function index: %d", builtinIdx)
			}
			fn := builtinImpls[int(builtinIdx)]

			// Optimize argument passing: Use slice of stack directly!
			startIdx := vm.sp - argCount
			if startIdx < 0 {
				return fmt.Errorf("stack underflow for builtin call")
			}

			// Get arguments from stack without allocation (just slice header)
			args := vm.stack[startIdx:vm.sp]

			// Call function
			res, err := fn(vm, args)
			if err != nil {
				return err
			}

			// Pop arguments and push result (net -argCount+1)
			vm.sp = startIdx
			if err := vm.push(res); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown opcode %d", op)
		}
	}
	return nil
}

func (vm *VM) readUint8() uint8 {
	val := vm.chunk.Code[vm.ip]
	vm.ip++
	return val
}

func (vm *VM) push(val interpreter.Value) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = val
	vm.sp++
	return nil
}

func (vm *VM) pop() interpreter.Value {
	if vm.sp == 0 {
		panic("stack underflow")
	}
	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *VM) readUint16() uint16 {
	val := binary.BigEndian.Uint16(vm.chunk.Code[vm.ip:])
	vm.ip += 2
	return val
}
