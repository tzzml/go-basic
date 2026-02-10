package vm

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"zork-basic/internal/interpreter"
)

type BuiltinFunc func(vm *VM, args []interpreter.Value) (interpreter.Value, error)

// Use slice instead of map for O(1) access
var builtinImpls = []BuiltinFunc{
	builtinAbs,   // ABS
	builtinSin,   // SIN
	builtinCos,   // COS
	builtinTan,   // TAN
	builtinInt,   // INT
	builtinExp,   // EXP
	builtinSqr,   // SQR
	builtinLog,   // LOG
	builtinRnd,   // RND
	builtinLen,   // LEN
	builtinLeft,  // LEFT$
	builtinRight, // RIGHT$
	builtinMid,   // MID$
	builtinInstr, // INSTR
	builtinUcase, // UCASE$
	builtinLcase, // LCASE$
	builtinSpace, // SPACE$
	builtinChr,   // CHR$
	builtinAsc,   // ASC
	builtinPi,    // PI
	builtinEuler, // E (EULER)
}

func builtinAbs(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("ABS requires 1 argument")
	}
	return interpreter.NumberValue(math.Abs(args[0].AsNumber())), nil
}

func builtinSin(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("SIN requires 1 argument")
	}
	return interpreter.NumberValue(math.Sin(args[0].AsNumber())), nil
}

func builtinCos(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("COS requires 1 argument")
	}
	return interpreter.NumberValue(math.Cos(args[0].AsNumber())), nil
}

func builtinTan(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("TAN requires 1 argument")
	}
	return interpreter.NumberValue(math.Tan(args[0].AsNumber())), nil
}

func builtinInt(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("INT requires 1 argument")
	}
	return interpreter.NumberValue(math.Trunc(args[0].AsNumber())), nil
}

func builtinExp(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("EXP requires 1 argument")
	}
	return interpreter.NumberValue(math.Exp(args[0].AsNumber())), nil
}

func builtinSqr(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("SQR requires 1 argument")
	}
	val := args[0].AsNumber()
	if val < 0 {
		return interpreter.NumberValue(0), fmt.Errorf("SQR of negative number")
	}
	return interpreter.NumberValue(math.Sqrt(val)), nil
}

func builtinLog(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("LOG requires 1 argument")
	}
	val := args[0].AsNumber()
	if val <= 0 {
		return interpreter.NumberValue(0), fmt.Errorf("LOG of non-positive number")
	}
	return interpreter.NumberValue(math.Log(val)), nil
}

func builtinRnd(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 0 {
		return interpreter.NumberValue(0), fmt.Errorf("RND requires 0 arguments")
	}
	return interpreter.NumberValue(rand.Float64()), nil
}

func builtinLen(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("LEN requires 1 argument")
	}
	return interpreter.NumberValue(float64(len(args[0].String()))), nil
}

func builtinLeft(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return interpreter.NumberValue(0), fmt.Errorf("LEFT$ requires 2 arguments")
	}
	str := args[0].String()
	n := int(args[1].AsNumber())
	if n > len(str) {
		n = len(str)
	}
	if n < 0 {
		n = 0
	}
	return interpreter.StringValue(str[:n]), nil
}

func builtinRight(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return interpreter.NumberValue(0), fmt.Errorf("RIGHT$ requires 2 arguments")
	}
	str := args[0].String()
	n := int(args[1].AsNumber())
	if n > len(str) {
		n = len(str)
	}
	if n < 0 {
		n = 0
	}
	return interpreter.StringValue(str[len(str)-n:]), nil
}

func builtinMid(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return interpreter.NumberValue(0), fmt.Errorf("MID$ requires 2 or 3 arguments")
	}
	str := args[0].String()
	start := int(args[1].AsNumber())
	if start < 1 {
		start = 1
	}
	n := len(str) - start + 1
	if len(args) == 3 {
		n = int(args[2].AsNumber())
	}
	startIdx := start - 1
	endIdx := startIdx + n
	if endIdx > len(str) {
		endIdx = len(str)
	}
	if startIdx >= len(str) || startIdx < 0 {
		return interpreter.StringValue(""), nil
	}
	return interpreter.StringValue(str[startIdx:endIdx]), nil
}

func builtinInstr(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return interpreter.NumberValue(0), fmt.Errorf("INSTR requires 2 or 3 arguments")
	}
	var start int = 1
	var str, substr string
	if len(args) == 2 {
		str = args[0].String()
		substr = args[1].String()
	} else {
		start = int(args[0].AsNumber())
		str = args[1].String()
		substr = args[2].String()
	}
	if start < 1 {
		start = 1
	}
	if start > len(str) {
		return interpreter.NumberValue(0), nil
	}
	pos := strings.Index(str[start-1:], substr)
	if pos == -1 {
		return interpreter.NumberValue(0), nil
	}
	return interpreter.NumberValue(float64(start + pos)), nil
}

func builtinUcase(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("UCASE$ requires 1 argument")
	}
	return interpreter.StringValue(strings.ToUpper(args[0].String())), nil
}

func builtinLcase(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("LCASE$ requires 1 argument")
	}
	return interpreter.StringValue(strings.ToLower(args[0].String())), nil
}

func builtinSpace(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("SPACE$ requires 1 argument")
	}
	n := int(args[0].AsNumber())
	if n < 0 {
		n = 0
	}
	return interpreter.StringValue(strings.Repeat(" ", n)), nil
}

func builtinChr(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("CHR$ requires 1 argument")
	}
	code := int(args[0].AsNumber())
	if code < 0 || code > 255 {
		return interpreter.NumberValue(0), fmt.Errorf("CHR$ argument must be between 0 and 255")
	}
	return interpreter.StringValue(string(rune(code))), nil
}

func builtinAsc(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return interpreter.NumberValue(0), fmt.Errorf("ASC requires 1 argument")
	}
	str := args[0].String()
	if len(str) == 0 {
		return interpreter.NumberValue(0), fmt.Errorf("ASC argument is an empty string")
	}
	return interpreter.NumberValue(float64(str[0])), nil
}

func builtinPi(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 0 {
		return interpreter.NumberValue(0), fmt.Errorf("PI requires 0 arguments")
	}
	return interpreter.NumberValue(math.Pi), nil
}

func builtinEuler(vm *VM, args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 0 {
		return interpreter.NumberValue(0), fmt.Errorf("EULER requires 0 arguments")
	}
	return interpreter.NumberValue(math.E), nil
}
