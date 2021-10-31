package shared

import (
	"io"
	"os"
)

const TekoErrorMessage string = "Error occurred while executing Teko process"

type TekoErrorClass string

const (
	LexerError          = "Lexer Error"
	SyntaxError         = "Syntax Error"
	NameError           = "Name Error"
	NotImplementedError = "Unimplemented (this is planned functionality for Teko)"
	TypeError           = "Type Error"
	ArgumentError       = "Argument Error"
	UnexpectedIssue     = "Unexpected issue (this is a mistake in the Teko implementation)"
	RuntimeError        = "Runtime Error"
)

var PrintDest io.Writer = os.Stdout

func SetPrintDest(new_print_dest io.Writer) {
	PrintDest = new_print_dest
}

var ErrorDest io.Writer = os.Stderr

func SetErrorDest(new_error_dest io.Writer) {
	ErrorDest = new_error_dest
}
