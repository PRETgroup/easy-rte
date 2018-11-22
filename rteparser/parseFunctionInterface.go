package rteparser

//parseFunctionInterface will add an interface to an existing internal function
func (t *pParse) parseFunctionInterface() *ParseError {
	var s string
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

	//third should be open brace
	s = t.pop()
	if s != pOpenBrace {
		return t.errorUnexpectedWithExpected(s, pOpenBrace)
	}

	//now we run until closed brace
	for {
		//inside the interface, we have
		//[enforce] in|out event|bool|byte|word|dword|lword|sint|usint|int|uint|dint|udint|lint|ulint|real|lreal|time|any (name[, name]) [on (name[, name]) - non-event types only]; //(comments)
		//repeated over and over again
		s = t.pop()
		if s == "" {
			return t.error(ErrUnexpectedEOF)
		}
		if s == pCloseBrace {
			return nil //we're done here
		}

		if s == pIn || s == pOut {
			if err := t.addFunctionIo(s == pIn, fbIndex); err != nil {
				return err
			}
		}
	}
}

//addFunctionIo adds a line of in/out event/data [:= default] to the FB interface
// some error checking is done
//isInput: set to TRUE if you want to add this to the Inputs rather than the Outputs
//fbIndex: the index of the Function we are working on inside t
func (t *pParse) addFunctionIo(isInput bool, fbIndex int) *ParseError {
	fb := &t.funcs[fbIndex]

	//next s is type
	typ := t.pop()
	if !isValidType(typ) {
		return t.errorWithArgAndReason(ErrInvalidType, typ, "Expected valid type")
	}

	var intNames []string

	//there might be an array size next
	size := ""
	if t.peek() == pOpenBracket {
		t.pop() // get rid of open bracket
		size = t.pop()
		if s := t.peek(); s != pCloseBracket {
			return t.errorUnexpectedWithExpected(s, pCloseBracket)
		}
		t.pop() //get rid of close bracket
	}

	//this could be an array of names, so we'll loop while we are finding commas
	for {
		name := t.pop()

		intNames = append(intNames, name)
		if t.peek() == pComma {
			t.pop() //get rid of the pComma
			continue
		}
		break
	}

	//there might be a default value next
	initialValue := ""
	if t.peek() == pInitEq {
		t.pop() //get rid of pInitial

		bracketOpen := 0
		for {
			s := t.peek()
			if s == "" {
				return t.error(ErrUnexpectedEOF)
			}
			//deal with brackets, if we have an open bracket we must have a close bracket, etc
			if s == pOpenBracket && bracketOpen == 0 {
				bracketOpen = 1
			} else if s == pOpenBracket && bracketOpen != 0 {
				return t.errorUnexpectedWithExpected(s, "[Value]")
			}
			if s == pCloseBracket && bracketOpen == 1 {
				bracketOpen = 2
			} else if s == pCloseBracket && bracketOpen != 1 {
				return t.errorUnexpectedWithExpected(s, pSemicolon)
			}
			if s == pSemicolon && bracketOpen == 1 { //can't return if brackets are open
				return t.errorUnexpectedWithExpected(s, pCloseBracket)
			}
			if s == pSemicolon {
				break
			}
			initialValue += s
			t.pop() //pop whatever we were just peeking at
		}
	}

	//clear out last semicolon
	if s := t.pop(); s != pSemicolon {
		return t.errorUnexpectedWithExpected(s, pSemicolon)
	}

	//we now have everything we need to add the io to the interface

	if err := fb.AddIO(isInput, intNames, typ, size, initialValue); err != nil {
		return t.errorWithArg(ErrNameAlreadyInUse, err.Error())
	}

	return nil
}
