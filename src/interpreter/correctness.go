package interpreter

import (
	"github.com/cstuartroe/teko/src/checker"
)

func VerifyStdlibDeclarations() {
	for name := range checker.StdlibSymbolTable {
		_, ok := StdLibFieldValues[name]
		if !ok {
			panic("Standard library object has type declared but no value: " + name)
		}
	}

	// I would use generics instead of repeating myself
	// https://tinyurl.com/ndwqwz4

	for name := range StdLibFieldValues {
		_, ok := checker.StdlibSymbolTable[name]
		if !ok {
			panic("Standard library object's type was never declared: " + name)
		}
	}
}
