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

		s := t.pop()           //this might be an openbracket
		if s == pOpenBracket { //for arrays
			initialValue += s //we need to keep the brackets in
			for {
				s := t.pop()
				if s == "" {
					return t.error(ErrUnexpectedEOF)
				}
				if s == pSemicolon {
					return t.errorUnexpectedWithExpected(s, pOpenBracket)
				}
				if s == pCloseBracket {
					initialValue += s //we need to keep the brackets in
					break
				}
				initialValue += s
			}
		} else { //wasn't an open bracket, must just be value
			initialValue = s
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
