package rteparser

import (
	"github.com/PRETgroup/easy-rte/rtedef"
	"github.com/PRETgroup/goFB/iec61499"
)

//pParse is the containing struct for the parsing code
type pParse struct {
	funcs []rtedef.EnforcedFunction

	items     []string
	itemIndex int

	currentLine int
	currentFile string
}

//getCurrentDebugInfo returns the debug info for the last popped item
func (t *pParse) getCurrentDebugInfo() iec61499.DebugInfo {
	return iec61499.DebugInfo{
		SourceLine: t.currentLine,
		SourceFile: t.currentFile,
	}
}

//isNameUnused will check all registered functions to see if a name can be used (as they need to be unique)
func (t *pParse) isNameUnused(name string) bool {
	for i := 0; i < len(t.funcs); i++ {
		if t.funcs[i].Name == name {
			return false
		}
	}
	return true
}

//pop gets the current element of the pParse internal items slice
// and increments the index
func (t *pParse) pop() string {
	if t.done() {
		return ""
	}
	s := t.items[t.itemIndex]
	t.itemIndex++

	if s == pNewline {
		t.currentLine++
		return t.pop()
	}
	return s
}

//peek gets the current element of the pParse internal items slice (or the next non-newline character)
// without incrementing the index
func (t *pParse) peek() string {
	if t.done() {
		return ""
	}
	for i := 0; i < len(t.items); i++ {
		if t.items[t.itemIndex+i] != pNewline {
			return t.items[t.itemIndex+i]
		}
	}
	return ""
}

//done checks to see if the pParse is completed (i.e. nothing left to parse)
func (t *pParse) done() bool {
	return t.itemIndex >= len(t.items)
}

//getIndexFromName will search the pParse slice of FBs for one that matches
// the provided name and return the index if found
func (t *pParse) getIndexFromName(name string) int {
	for i := 0; i < len(t.funcs); i++ {
		if t.funcs[i].Name == name {
			return i
		}
	}
	return -1
}
