package interpreter

import (
	"github.com/cstuartroe/teko/src/checker"
)

func VerifyStdlibDeclarations() {
	for name, _ := range checker.BaseCheckerTypeFields {
		_, ok := BaseInterpreterFieldValues[name]
		if !ok {
			panic("Standard library object has type declared but no value: " + name)
		}
	}

	// I would use generics instead of repeating myself
	// https://tinyurl.com/ndwqwz4

	for name, _ := range BaseInterpreterFieldValues {
		_, ok := checker.BaseCheckerTypeFields[name]
		if !ok {
			panic("Standard library object's type was never declared: " + name)
		}
	}
}
