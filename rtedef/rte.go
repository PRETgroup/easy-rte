package rtedef

import (
	"errors"
	"strconv"
	"strings"
)

//An EnforcedFunction is what we're here for, it's what we're going to wrap with our policies!
type EnforcedFunction struct {
	Name string `xml:"Name,attr"`

	InterfaceList

	Policies []Policy `xml:"Policy"`
}

//InterfaceList stores the IO
type InterfaceList struct {
	InputVars  []Variable `xml:"Interface>Input"`
	OutputVars []Variable `xml:"Interface>Output"`
}

//HasIONamed will check a given InterfaceList to see if it has an output of that name
func (il InterfaceList) HasIONamed(input bool, s string) bool {
	if input {
		for i := 0; i < len(il.InputVars); i++ {
			if il.InputVars[i].Name == s {
				return true
			}
		}
		return false
	}
	for i := 0; i < len(il.OutputVars); i++ {
		if il.OutputVars[i].Name == s {
			return true
		}
	}
	return false
}

//A Variable is used to store I/O or internal var data
type Variable struct {
	Name         string `xml:"Name,attr"`
	Type         string `xml:"Type,attr"`
	ArraySize    string `xml:"ArraySize,attr,omitempty"`
	InitialValue string `xml:"InitialValue,attr,omitempty"`
	Comment      string `xml:"Comment,attr"`
}

//GetInitialArray returns a formatted initial array if there is one to do so
func (v Variable) GetInitialArray() []string {
	//if cannot parse an array size then give up
	_, err := strconv.Atoi(v.ArraySize)
	if err != nil {
		return nil
	}

	//remove everything except commas and values
	raw := v.InitialValue
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")

	raws := strings.Split(raw, ",")
	for i := 0; i < len(raws); i++ {
		raws[i] = strings.Trim(raws[i], " ")
	}
	return raws
}

//Policy stores a policy, i.e. the vars that must be kept
type Policy struct {
	Name         string        `xml:"Name,attr"`
	InternalVars []Variable    `xml:"InternalVars>VarDeclaration,omitempty"`
	States       []PState      `xml:"Machine>PState"`
	Transitions  []PTransition `xml:"Machine>PTransition,omitempty"`
}

//PState is a state in the policy specification of an enforcerFB
type PState string

//Name returns the name of a PState
func (p PState) Name() string {
	return string(p)
}

//PTransition is a transition between PState in a Policy (mealy machine transitions)
type PTransition struct {
	Source      PState
	Destination PState
	Condition   string
	Expressions []PExpression //output expressions associated with this transition
	Recover     []PExpression
}

//PExpression is used to assign a var a value based on a PTransitions
type PExpression struct {
	VarName string
	Value   string
}

//NewEnforcedFunction creates a new EnforcedFunction struct
func NewEnforcedFunction(name string) EnforcedFunction {
	return EnforcedFunction{Name: name}
}

//AddIO adds the provided IO to a given EnforcedFunction, while checking to make sure that each name is unique in the interface
func (f *EnforcedFunction) AddIO(isInput bool, intNames []string, typ string, size string, initialValue string) error {
	seenNames := make(map[string]bool)
	for _, inp := range f.InputVars {
		seenNames[inp.Name] = true
	}
	for _, outp := range f.OutputVars {
		seenNames[outp.Name] = true
	}

	vars := make([]Variable, len(intNames))
	for i, name := range intNames {
		if seenNames[name] == true {
			return errors.New("The name " + name + " is already in use")
		}
		seenNames[name] = true
		vars[i] = Variable{
			Name:         name,
			Type:         typ,
			ArraySize:    size,
			InitialValue: initialValue,
		}
	}
	if isInput {
		f.InputVars = append(f.InputVars, vars...)
		return nil
	}
	f.OutputVars = append(f.OutputVars, vars...)
	return nil
}

//AddPolicy adds a Policy to an EnforcedFunction
func (f *EnforcedFunction) AddPolicy(name string) {
	f.Policies = append(f.Policies, Policy{Name: name})
}

//AddDataInternals adds data internals to a efb, and adds the InternalVars section if it is nil
func (efb *Policy) AddDataInternals(intNames []string, typ string, size string, initialValue string) *Policy {
	for _, iname := range intNames {
		efb.InternalVars = append(efb.InternalVars, Variable{Name: iname, Type: typ, ArraySize: size, InitialValue: initialValue})
	}
	return efb
}

//AddState adds a state to a bfb
func (efb *Policy) AddState(name string) error {
	efb.States = append(efb.States, PState(name))
	return nil //TODO: add check (make sure name is unique)
}

//AddTransition adds a state transition to a bfb
func (efb *Policy) AddTransition(source string, dest string, cond string, expressions []PExpression, recoveries []PExpression) error {
	efb.Transitions = append(efb.Transitions, PTransition{
		Source:      PState(source),
		Destination: PState(dest),
		Condition:   cond,
		Expressions: expressions,
		Recover:     recoveries,
	})
	return nil //TODO: make sure [source] and [dest] can be found, make sure [cond] is valid, make sure [expressions] is valid
}
