package rteparser

import (
	"errors"
	"fmt"
)

var (
	//ErrInternal means something went wrong and it's the transpiler's fault
	ErrInternal = errors.New("An internal error occured")

	//ErrUnexpectedEOF means the document ended unexpectedly
	ErrUnexpectedEOF = errors.New("Unexpected EOF")

	//ErrUnexpectedValue is used to indicate a parsed value was not what was expected (i.e. a word instead of a semicolon)
	ErrUnexpectedValue = errors.New("Unexpected value")

	//ErrUndefinedFunction is used to indicate a Function was referenced that can't be found (so probably a typo has occured)
	ErrUndefinedFunction = errors.New("Can't find Function with name")

	//ErrInvalidType is used when the type of a data variable is bad
	ErrInvalidType = errors.New("Invalid or missing data type")

	//ErrInvalidIOMeta is used when metadata for an I/O line is bad
	ErrInvalidIOMeta = errors.New("Invalid metadata for data/event line")

	//ErrNameAlreadyInUse is returned whenever something is named but the name is already in use elsewhere
	ErrNameAlreadyInUse = errors.New("This name is already defined elsewhere")
)

//ParseError is used to contain a helpful error message when parsing fails
type ParseError struct {
	LineNumber int
	Argument   string
	Reason     string
	Err        error
}

//Error makes ParseError fulfill error interface
func (p ParseError) Error() string {
	s := fmt.Sprintf("Error (Line %v): %s", p.LineNumber, p.Err.Error())
	if p.Argument != "" {
		s += " '" + p.Argument + "'"
	}
	if p.Reason != "" {
		s += ", (" + p.Reason + ")"
	}
	return s
}

// helper functions to help construct helpful error messages

func (t *pParse) errorWithArg(err error, arg string) *ParseError {
	return &ParseError{LineNumber: t.currentLine, Argument: arg, Reason: "", Err: err}
}

func (t *pParse) errorWithArgAndLineNumber(err error, arg string, line int) *ParseError {
	return &ParseError{LineNumber: line, Argument: arg, Reason: "", Err: err}
}

func (t *pParse) errorWithReason(err error, reason string) *ParseError {
	return &ParseError{LineNumber: t.currentLine, Argument: "", Reason: reason, Err: err}
}

func (t *pParse) error(err error) *ParseError {
	return &ParseError{LineNumber: t.currentLine, Argument: "", Reason: "", Err: err}
}

func (t *pParse) errorWithArgAndReason(err error, arg string, reason string) *ParseError {
	return &ParseError{LineNumber: t.currentLine, Argument: arg, Reason: reason, Err: err}
}

func (t *pParse) errorUnexpectedWithExpected(unexpected string, expected string) *ParseError {
	return &ParseError{LineNumber: t.currentLine, Argument: unexpected, Reason: "Expected: " + expected, Err: ErrUnexpectedValue}
}
