package rteparser

import (
	"errors"
	"strings"
	"text/scanner"

	"github.com/PRETgroup/easy-rte/rtedef"
)

const (
	pNewline = "\n"

	pFunction     = "function"
	pInterface    = "interface"
	pArchitecture = "architecture"

	pOpenBrace    = "{"
	pCloseBrace   = "}"
	pOpenBracket  = "["
	pCloseBracket = "]"
	pComma        = ","
	pSemicolon    = ";"
	pColon        = ":"
	pInitEq       = ":="
	pAssigment    = ":="

	pFBpolicy = "policy"
	pOf       = "of"

	pIn  = "in"
	pOut = "out"

	pWith = "with"

	pTrans = "->"
	pOn    = "on"

	pRecover = "recover"

	pInternal   = "internal"
	pInternals  = "internals"
	pState      = "state"
	pStates     = "states"
	pAlgorithm  = "algorithm"
	pAlgorithms = "algorithms"
)

//ParseString takes an input string (i.e. filename) and input and returns all FBs in that string
func ParseString(name string, input string) ([]rtedef.EnforcedFunction, *ParseError) {
	//break up input string into all of its parts
	items := scanString(name, input)

	//now parse the items
	return parseItems(name, items)
}

func scanString(name string, input string) []string {
	var s scanner.Scanner

	s.Filename = name
	s.Init(strings.NewReader(input))

	//we don't want to ignore \n characters (we want to know what line we're on)
	s.Whitespace = 1<<'\t' | 0<<'\n' | 1<<'\r' | 1<<' '
	//we don't want scanner.ScanChars
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments

	//TODO: think about the scanner error function. Maybe we should provide one, that when an error occurs, halts scanning?

	var tok rune
	var items []string
	for tok != scanner.EOF {
		tok = s.Scan()
		items = append(items, s.TokenText())
	}

	//combine multi-character operators
	for i := 0; i < len(items)-1; i++ {
		if items[i] == "<" && items[i+1] == "-" {
			items[i] = "<-"
			items = append(items[:i+1], items[i+2:]...)
		}

		if items[i] == "-" && items[i+1] == ">" {
			items[i] = "->"
			items = append(items[:i+1], items[i+2:]...)
		}

		if items[i] == ":" && items[i+1] == "=" {
			items[i] = ":="
			items = append(items[:i+1], items[i+2:]...)
		}

		if items[i] == "=" && items[i+1] == "=" {
			items[i] = "=="
			items = append(items[:i+1], items[i+2:]...)
		}

		if items[i] == ">" && items[i+1] == "=" {
			items[i] = ">="
			items = append(items[:i+1], items[i+2:]...)
		}

		if items[i] == "<" && items[i+1] == "=" {
			items[i] = "<="
			items = append(items[:i+1], items[i+2:]...)
		}

		if items[i] == "!" {
			items[i] = "!" + items[i+1]
			items = append(items[:i+1], items[i+2:]...)
		}

		if items[i] == "&" && items[i+1] == "&" {
			items[i] = "&&"
			items = append(items[:i+1], items[i+2:]...)
		}

		if items[i] == "|" && items[i+1] == "|" {
			items[i] = "||"
			items = append(items[:i+1], items[i+2:]...)
		}
	}

	return items
}

//parseItems creates and runs a pparse struct
func parseItems(name string, items []string) ([]rtedef.EnforcedFunction, *ParseError) {
	t := pParse{items: items, currentLine: 1, currentFile: name}

	for !t.done() {
		s := t.pop()
		if t.done() {
			break
		}
		//have we defined a basicFB or compositeFB
		if s == pFunction {
			if err := t.parseFunction(s); err != nil {
				return nil, err
			}
			continue
		}

		//is this defining an interface for an fb
		if s == pInterface {
			if err := t.parseFunctionInterface(); err != nil {
				return nil, err
			}
			continue
		}

		//is this defining an architecture for an fb
		if s == pArchitecture || s == pFBpolicy {
			if err := t.parseFunctionArchitecture(s); err != nil {
				return nil, err
			}
			continue
		}
		return nil, t.errorWithArg(ErrUnexpectedValue, s)
	}

	return t.funcs, nil
}

//isValidType returns true if string s is one of the valid event/data types
func isValidType(s string) bool {
	s = strings.ToLower(s)
	if s == "bool" ||
		s == "char" ||
		s == "uint8_t" ||
		s == "uint16_t" ||
		s == "uint32_t" ||
		s == "uint64_t" ||
		s == "int8_t" ||
		s == "int16_t" ||
		s == "int32_t" ||
		s == "int64_t" ||
		s == "float" ||
		s == "double" ||
		s == "dtimer_t" ||
		s == "rtimer_t" {
		return true
	}
	return false
}

//parseFunction will create a new function and add it to the list of internal functions
func (t *pParse) parseFunction(typ string) *ParseError {
	var funcs []rtedef.EnforcedFunction
	for {
		name := t.pop()
		if !t.isNameUnused(name) {
			return t.errorWithArg(ErrNameAlreadyInUse, name)
		}

		if typ == pFunction {
			funcs = append(funcs, rtedef.NewEnforcedFunction(name))
		} else {
			return t.errorWithReason(ErrInternal, "I can't parse type "+typ)
		}

		if t.peek() == pComma {
			t.pop() //get rid of comma
			continue
		}
		break
	}

	s := t.pop()
	if s != pSemicolon {
		return t.errorUnexpectedWithExpected(s, pSemicolon)
	}

	t.funcs = append(t.funcs, funcs...)

	return nil
}

func (t *pParse) parseFunctionArchitecture(archType string) *ParseError {
	var s string
	var pName string

	//if this is a policy, the name is here
	if archType == pFBpolicy {
		s = t.pop()
		pName = s
	}

	//first word should be of
	s = t.pop()
	if s != pOf {
		return t.errorUnexpectedWithExpected(s, pOf)
	}

	//second word is fb name
	s = t.pop()
	fbIndex := t.getIndexFromName(s)
	if fbIndex == -1 {
		return t.errorWithArg(ErrUndefinedFunction, s)
	}

	if archType == pFBpolicy {
		t.funcs[fbIndex].AddPolicy(pName)
		return t.parsePolicyArchitecture(fbIndex)
	}
	return t.error(errors.New("can't parse unknown architecture type"))
}
