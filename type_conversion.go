package main

import (
	"fmt"
	"go/ast"
)

func getInputType(param ast.Expr) string {
	var inputType string
	switch t := param.(type) {
	case *ast.ArrayType:
		if v, ok := t.Elt.(*ast.Ident); ok {
			inputType = `[]` + v.Name
		} else {
			panic(fmt.Sprintf(`parameter ast type %T not supported`, t.Elt))
		}
	case *ast.Ident:
		inputType = t.Name
	default:
		panic(fmt.Sprintf(`parameter ast type %T not supported`, t))
	}

	return inputType
}

// Return a string which converts the input variable (of type []byte) to the given inputType
// Supported types are Go native fuzzing types (https://go.dev/security/fuzz/)
func getInputCode(inputType string) (string, int, []string) {
	var inputCode string
	var inputLen int
	var imprts []string
	switch inputType {
	case `[]byte`:
		inputCode = `input`
	case `string`:
		inputCode = `string(input)`
	case `int`:
		inputLen = 4
		inputCode = `int(binary.BigEndian.Uint32(input))`
		imprts = []string{"encoding/binary"}
	case `int8`:
		inputLen = 1
		inputCode = `int8(input[0])`
	case `int16`:
		inputLen = 2
		inputCode = `int16(binary.BigEndian.Uint16(input))`
		imprts = []string{"encoding/binary"}
	case `int32`, `rune`:
		inputLen = 4
		inputCode = `int32(binary.BigEndian.Uint32(input))`
		imprts = []string{"encoding/binary"}
	case `int64`:
		inputLen = 8
		inputCode = `int64(binary.BigEndian.Uint64(input))`
		imprts = []string{"encoding/binary"}
	case `uint`:
		inputLen = 4
		inputCode = `uint(binary.BigEndian.Uint32(input))`
		imprts = []string{"encoding/binary"}
	case `uint8`, `byte`:
		inputLen = 1
		inputCode = `input[0]`
	case `uint16`:
		inputLen = 2
		inputCode = `uint16(binary.BigEndian.Uint16(input))`
		imprts = []string{"encoding/binary"}
	case `uint32`:
		inputLen = 4
		inputCode = `uint32(binary.BigEndian.Uint32(input))`
		imprts = []string{"encoding/binary"}
	case `uint64`:
		inputLen = 8
		inputCode = `uint64(binary.BigEndian.Uint64(input))`
		imprts = []string{"encoding/binary"}
	case `float32`:
		inputLen = 4
		inputCode = `math.Float32frombits(binary.BigEndian.Uint32(input))`
		imprts = []string{"encoding/binary", "math"}
	case `float64`:
		inputLen = 8
		inputCode = `math.Float64frombits(binary.BigEndian.Uint64(input))`
		imprts = []string{"encoding/binary", "math"}
	case `bool`:
		inputLen = 1
		inputCode = `input[0] == 1`
	default:
		panic(fmt.Sprintf(`parameter type %s not supported`, inputType))
	}

	return inputCode, inputLen, imprts
}

func getCheckCode(inputLen int) string {
	if inputLen == 0 {
		return ``
	}
	return fmt.Sprintf(`if len(input) != %d { return -1 }`, inputLen)
}
