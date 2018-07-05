package rtedef

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PRETgroup/goFB/goFB/stconverter"
)

//FBECCGuardToSTExpression converts a given FB's guard into a STExpression parsetree
func FBECCGuardToSTExpression(pName, guard string) ([]stconverter.STInstruction, *stconverter.STParseError) {
	return stconverter.ParseString(pName, guard)
}

//PSTTransition is a container struct for a PTransition and its ST translated guard
type PSTTransition struct {
	PTransition
	STGuard stconverter.STExpression
}

//A PEnforcerPolicy is what goes inside a PEnforcer, it is derived from a Policy
type PEnforcerPolicy struct {
	InternalVars []Variable
	States       []PState
	Transitions  []PSTTransition
}

//GetDTimers returns all DTIMERS in a PEnforcerPolicy
func (pol PEnforcerPolicy) GetDTimers() []Variable {
	dTimers := make([]Variable, 0)
	for _, v := range pol.InternalVars {
		if strings.ToLower(v.Type) == "dtimer_t" {
			dTimers = append(dTimers, v)
		}
	}
	return dTimers
}

//GetViolationTransitions returns a slice of all transitions in this PEnforcerPolicy
//that have their destinations set to "violation", ie. are violation transitions
func (pol PEnforcerPolicy) GetViolationTransitions() []PSTTransition {
	violTrans := make([]PSTTransition, 0)
	for _, tr := range pol.Transitions {
		if tr.Destination == "violation" {
			violTrans = append(violTrans, tr)
		}
	}
	return violTrans
}

//GetNonViolationTransitions returns a slice of all transitions in this PEnforcerPolicy
//that have their destinations not set to "violation", ie. are not violation transitions
func (pol PEnforcerPolicy) GetNonViolationTransitions() []PSTTransition {
	nviolTrans := make([]PSTTransition, 0)
	for _, tr := range pol.Transitions {
		if tr.Destination != "violation" {
			nviolTrans = append(nviolTrans, tr)
		}
	}
	return nviolTrans
}

//DoesExpressionInvolveTime returns true if a given expression uses time
func (pol PEnforcerPolicy) DoesExpressionInvolveTime(expr stconverter.STExpression) bool {
	op := expr.HasOperator()
	if op == nil {
		return VariablesContain(pol.GetDTimers(), expr.HasValue())
	}
	for _, arg := range expr.GetArguments() {
		if pol.DoesExpressionInvolveTime(arg) {
			return true
		}
	}
	return false
}

//A PEnforcer will store a given input and output policy and can derive the enforcers required to uphold them
type PEnforcer struct {
	interfaceList InterfaceList
	Name          string
	OutputPolicy  PEnforcerPolicy
	InputPolicy   PEnforcerPolicy
}

//MakePEnforcer will convert a given policy to an enforcer for that policy
func MakePEnforcer(il InterfaceList, p Policy) (*PEnforcer, error) {
	//make the enforcer
	enf := &PEnforcer{interfaceList: il, Name: p.Name}
	//first, convert policy transitions
	outpTr, err := p.GetPSTTransitions()
	if err != nil {
		return nil, err
	}
	splOutpTr := SplitPSTTransitions(outpTr)
	enf.OutputPolicy = PEnforcerPolicy{
		InternalVars: p.InternalVars,
		States:       p.States,
		Transitions:  splOutpTr,
	}
	enf.OutputPolicy.RemoveDuplicateTransitions()

	enf.InputPolicy = DeriveInputEnforcerPolicy(il, enf.OutputPolicy)
	enf.InputPolicy.RemoveNilTransitions()
	enf.InputPolicy.RemoveDuplicateTransitions()
	enf.InputPolicy.RemoveAlwaysTrueTransitions()

	return enf, nil
}

//RemoveNilTransitions will do a search through a policies transitions and remove any that have nil guards
func (pol *PEnforcerPolicy) RemoveNilTransitions() {
	for i := 0; i < len(pol.Transitions); i++ {
		for j := i + 1; j < len(pol.Transitions); j++ {
			if pol.Transitions[j].STGuard == nil {
				pol.Transitions = append(pol.Transitions[:j], pol.Transitions[j+1:]...)
				j--
			}
		}
	}
}

//RemoveDuplicateTransitions will do a search through a policies transitions and remove any that are simple duplicates
//(i.e. every field the same and in the same order).
func (pol *PEnforcerPolicy) RemoveDuplicateTransitions() {
	for i := 0; i < len(pol.Transitions); i++ {
		for j := i + 1; j < len(pol.Transitions); j++ {
			if reflect.DeepEqual(pol.Transitions[i], pol.Transitions[j]) {
				pol.Transitions = append(pol.Transitions[:j], pol.Transitions[j+1:]...)
				j--
			}
		}
	}
}

//RemoveAlwaysTrueTransitions will do a search through a policies transitions and remove any that are just "true"
func (pol *PEnforcerPolicy) RemoveAlwaysTrueTransitions() {
	for i := 0; i < len(pol.Transitions); i++ {
		for j := i + 1; j < len(pol.Transitions); j++ {
			if val := pol.Transitions[j].STGuard.HasValue(); val == "true" || val == "1" {
				pol.Transitions = append(pol.Transitions[:j], pol.Transitions[j+1:]...)
				j--
			}
		}
	}
}

//STExpressionSolution stores a solution to a violation transition
type STExpressionSolution struct {
	Expressions []stconverter.STExpression
	Comment     string
}

//SolveViolationTransition will attempt to solve a given transition
//TODO: consider where people have been too explicit with their time variables, and have got non-violating time-based transitions
//1. Check to see if there is a non-violating transition with an equivalent guard to the violating transition
//2. Select first solution
func (enf *PEnforcer) SolveViolationTransition(tr PSTTransition, inputPolicy bool) STExpressionSolution {

	//check if a recovery was provided
	if len(tr.Recover) > 0 {
		solution := make([]stconverter.STExpression, 0)
		for _, recov := range tr.Recover {
			solution = append(solution, stconverter.STExpressionOperator{
				Operator: stconverter.FindOp(":="),
				Arguments: []stconverter.STExpression{
					stconverter.STExpressionValue{Value: recov.Value},
					stconverter.STExpressionValue{Value: recov.VarName},
				},
			})
		}
		// solutionExpressions := make([]string, len(solution))
		// for i, soln := range solution {
		// 	solutionExpressions[i] = stconverter.CCompileExpression(soln)
		// }
		return STExpressionSolution{Expressions: solution, Comment: fmt.Sprintf("Recovery instructions manually provided.")}
	}

	fmt.Printf("Automatically deriving a solution for violation transition \r\n\t%s -> %s on (%s)\r\n\t(If this is undesirable behaviour, use a 'recover' keyword in the erte file to manually specify solution)\r\n", tr.Source, tr.Destination, stconverter.CCompileExpression(tr.STGuard))

	posSolTrs := make([]PSTTransition, 0) //possible Solution Transitions
	var pol PEnforcerPolicy
	if inputPolicy {
		pol = enf.InputPolicy
	} else {
		pol = enf.OutputPolicy
	}
	for _, propTr := range pol.Transitions {
		if propTr.Destination == "violation" {
			continue
		}
		if propTr.Source != tr.Source {
			continue
		}
		if pol.DoesExpressionInvolveTime(propTr.STGuard) {
			continue
		}
		posSolTrs = append(posSolTrs, propTr)
	}

	// Make sure there's at least one solution
	if len(posSolTrs) == 0 {
		fmt.Printf("\tNOTE: No solution found!\r\n")
		return STExpressionSolution{Expressions: nil, Comment: "No possible solutions!"}
	}

	//1. Check to see if there is a non-violating transition with an equivalent guard to the violating transition
	for _, posSolTr := range posSolTrs {
		if reflect.DeepEqual(tr.STGuard, posSolTr.STGuard) {
			fmt.Printf("\tNOTE: (Certain) Solution found with no edits required! (Equivalent safe transition found)\r\n")
			return STExpressionSolution{Expressions: nil, Comment: fmt.Sprintf("Selected non-violation transition \"%s -> %s on %s\" which has an equivalent guard, so no action is required", posSolTr.Source, posSolTr.Destination, posSolTr.Condition)}
		}
	}

	//2. Select first solution
	posSolTr := posSolTrs[0]
	solutions := SolveSTExpression(enf.interfaceList, inputPolicy, tr, posSolTr.STGuard)
	if solutions == nil {
		fmt.Printf("\tNOTE: (Guess) Solution found with no edits required! (I think the input policy has solved this transition already)\r\n")
		return STExpressionSolution{Expressions: nil, Comment: fmt.Sprintf("Selected non-violation transition \"%s -> %s on %s\" and action was not required", posSolTr.Source, posSolTr.Destination, posSolTr.Condition)}
	}

	// solutionExpressions := make([]string, len(solutions))
	// for i, soln := range solutions {
	// 	solutionExpressions[i] = stconverter.CCompileExpression(soln)
	// }

	fmt.Printf("\tNOTE: (Guess) Solution found, and edits required! (I have selected a safe transition, and edited the I/O so that it can be taken)\r\n\tSelected transition: \"%s -> %s on %s\"\r\n", posSolTr.Source, posSolTr.Destination, posSolTr.Condition)
	fmt.Printf("\tNOTE: I will perform the following edits:\r\n")
	for _, solnE := range solutions {
		fmt.Printf("\t\t%s;\r\n", stconverter.STCompileExpression(solnE))
	}
	return STExpressionSolution{Expressions: solutions, Comment: fmt.Sprintf("Selected non-violation transition \"%s -> %s on %s\" and action is required", posSolTr.Source, posSolTr.Destination, posSolTr.Condition)}

}

//DeriveInputEnforcerPolicy will derive an Input Policy from a given Output Policy
func DeriveInputEnforcerPolicy(il InterfaceList, outPol PEnforcerPolicy) PEnforcerPolicy {
	inpEnf := PEnforcerPolicy{
		States: outPol.States,
	}

	//inpEnf.InternalVars = nil
	//just realised that theres no internalVars that can't be managed by externalVars?
	inpEnf.InternalVars = make([]Variable, len(outPol.InternalVars))
	// for i := 0; i < len(outPol.InternalVars; i++) {}
	copy(inpEnf.InternalVars, outPol.InternalVars)

	//convert transitions and internal var names in transitions
	for i := 0; i < len(outPol.Transitions); i++ {
		inpEnf.Transitions = append(inpEnf.Transitions, ConvertPSTTransitionForInputPolicy(il, true, outPol.Transitions[i]))
	}

	// //convert internal var names on enforcer policy
	// for i := 0; i < len(inpEnf.InternalVars); i++ {
	// 	inpEnf.InternalVars[i].Name = inpEnf.InternalVars[i].Name + "_i"
	// }

	return inpEnf
}

//ConvertPSTTransitionForInputPolicy will convert a single PSTTransition from an Output Policy to its Input Policy Deriviation
func ConvertPSTTransitionForInputPolicy(il InterfaceList, inputPolicy bool, outpTrans PSTTransition) PSTTransition {
	var nonAcceptableNames []string
	if inputPolicy {
		nonAcceptableNames = make([]string, len(il.OutputVars))
		for i, v := range il.OutputVars {
			nonAcceptableNames[i] = v.Name
		}
	} else {
		nonAcceptableNames = make([]string, len(il.InputVars))
		for i, v := range il.InputVars {
			nonAcceptableNames[i] = v.Name
		}
	}
	//fmt.Printf("calling with %s\r\n", outpTrans.Condition)

	retSTGuard := ConvertSTExpressionForPolicy(il, nonAcceptableNames, true, outpTrans.STGuard)

	retTrans := outpTrans
	retTrans.STGuard = retSTGuard
	retTrans.Condition = stconverter.STCompileExpression(retSTGuard)
	//fmt.Printf("returning %s\r\n", retTrans.Condition)
	return retTrans
}

//VariablesContain returns true if a list of variables contains a given name
func VariablesContain(vars []Variable, name string) bool {
	for i := 0; i < len(vars); i++ {
		if vars[i].Name == name {
			return true
		}
	}
	return false
}

func stringSliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

//ConvertSTExpressionForPolicy will remove from a single STExpression
//all instances of vars located in []varNames if acceptableNames == false
//a == input
//b == output
//"a" becomes "a"
//"b" becomes "true" (technically becomes "true or not true")
//"a and b" becomes "a"
//"func(a, b)" becomes "func(a, true)"
//"!b" becomes "true" (technically becomes "not(true or not true)")
//TODO: a transition based only on time becomes nil?
func ConvertSTExpressionForPolicy(il InterfaceList, varNames []string, removeVarNames bool, expr stconverter.STExpression) stconverter.STExpression {
	//options
	//1. It is just a value
	//	  --if input or value, return
	//    --if output, return true
	//2. It is an operator
	//    Foreach arg
	//	      If arg

	//consider not(b)
	//op == not, args == []{b}
	//should return "true"

	op := expr.HasOperator()
	if op == nil { //if it's just a value, return if that value
		if removeVarNames && stringSliceContains(varNames, expr.HasValue()) ||
			!removeVarNames && !stringSliceContains(varNames, expr.HasValue()) {

			return nil
		}
		// if VariablesContain(intl, expr.HasValue()) {
		// 	return stconverter.STExpressionValue{Value: expr.HasValue() + "_i"}
		// }
		return expr
	}

	args := expr.GetArguments()
	acceptableArgIs := make([]bool, 0)
	numAcceptable := 0
	acceptableArgs := make([]stconverter.STExpression, 0)
	//for each argument, we want to check if it is "acceptable", which here means
	//"is not a value that is an output var"
	//and
	//"if it is an operator, convert it via this function, and see if it is acceptable"
	for i := 0; i < len(args); i++ {
		arg := args[i]
		argOp := arg.HasOperator()

		if argOp == nil {
			//it is a value
			argV := stconverter.STExpressionValue{Value: arg.HasValue()}
			//see if it is acceptable
			if removeVarNames && stringSliceContains(varNames, argV.HasValue()) ||
				!removeVarNames && !stringSliceContains(varNames, argV.HasValue()) {

				acceptableArgIs = append(acceptableArgIs, false)
				acceptableArgs = append(acceptableArgs, nil)
			} else {
				// if VariablesContain(intl, argV.HasValue()) {
				// 	argV.Value = argV.Value + "_i"
				// }
				acceptableArgIs = append(acceptableArgIs, true)
				acceptableArgs = append(acceptableArgs, argV)
				numAcceptable++
			}
			continue
		} else {
			//it is an operator, run the operator through this function and see if it is acceptable
			convArg := ConvertSTExpressionForPolicy(il, varNames, removeVarNames, args[i])
			if convArg != nil {
				acceptableArgIs = append(acceptableArgIs, true)
				acceptableArgs = append(acceptableArgs, convArg)
				numAcceptable++
			} else {
				acceptableArgIs = append(acceptableArgIs, false)
				acceptableArgs = append(acceptableArgs, nil)
			}
		}
	}

	//now we need to come up with a new STExpression to represent this expression and its arguments

	if numAcceptable < len(args) {
		//if less than the total args are acceptable, and only one argument is acceptable, then it is easy,
		//we can just return that one argument as an independent value, as long as the operator was a combinator like "and" or "or"

		if !stconverter.OpTokenIsCombinator(op.GetToken()) {
			//operator was not a comparison, e.g. "x >= 5" is not an "and" or an "or"
			//so no point returning "5"
			return nil
		}

		//e.g. "(a and b)" becomes "a"
		if numAcceptable == 1 {
			for i := 0; i < len(acceptableArgIs); i++ {
				if acceptableArgIs[i] == true {
					return acceptableArgs[i]
				}
			}
		}
	}
	if numAcceptable == 0 {
		//if nothing at all is acceptable then it is easy, we just return nil
		return nil
	}

	//if we are still here, then it means that there is no easy answer, so we'll just make a new
	//STExpressionOperator, which has the same operator as we're currently examining
	//then, all unacceptable (i.e. nil) arguments should be replaced with simple value "true"
	actualArgs := make([]stconverter.STExpression, len(acceptableArgs))
	validArgs := 0
	lastValidArg := 0
	for i := 0; i < len(actualArgs); i++ {
		if acceptableArgs[i] != nil {
			actualArgs[i] = acceptableArgs[i]
			validArgs++
			lastValidArg = i
		} else {
			actualArgs[i] = stconverter.STExpressionValue{Value: "true"}
		}
	}

	if validArgs == 1 && numAcceptable > 1 {
		return actualArgs[lastValidArg]
	}

	ret := stconverter.STExpressionOperator{
		Operator:  op,
		Arguments: actualArgs,
	}

	return ret
}

//GetPSTTransitions will convert all internal PTransitions into PSTTransitions (i.e. PTransitions with a ST symbolic tree condition)
func (p *Policy) GetPSTTransitions() ([]PSTTransition, error) {
	stTrans := make([]PSTTransition, len(p.Transitions))
	for i := 0; i < len(p.Transitions); i++ {
		stguard, err := FBECCGuardToSTExpression(p.Name, p.Transitions[i].Condition)
		if err != nil {
			return nil, err
		}
		if len(stguard) != 1 {
			return nil, fmt.Errorf("Incompatible policy guard (wrong number of expressions)")
		}
		expr, ok := stguard[0].(stconverter.STExpression)
		if !ok {
			return nil, fmt.Errorf("Incompatible policy guard (not an expression)")
		}
		stTrans[i] = PSTTransition{
			PTransition: p.Transitions[i],
			STGuard:     expr,
		}
	}
	return stTrans, nil
}

//SplitPSTTransitions will take a slice of PSTTRansitions and then split transitions which have OR clauses
//into multiple transitions
//it relies on the SplitExpressionsOnOr function below
func SplitPSTTransitions(cTrans []PSTTransition) []PSTTransition {
	brTrans := make([]PSTTransition, 0)

	for i := 0; i < len(cTrans); i++ {
		cTran := cTrans[i]
		splitTrans := SplitExpressionsOnOr(cTran.STGuard)
		for j := 0; j < len(splitTrans); j++ {
			newTrans := PSTTransition{
				PTransition: cTran.PTransition,
			}
			//recompile the condition
			newTrans.PTransition.Condition = stconverter.STCompileExpression(splitTrans[len(splitTrans)-j-1])
			newTrans.STGuard = splitTrans[len(splitTrans)-j-1]

			brTrans = append(brTrans, newTrans)
		}
	}

	//reformat all the guards based off the transactions
	return brTrans
}

//SplitExpressionsOnOr will take a given STExpression and return a slice of STExpressions which are
//split over the "or" operators, e.g.
//[a] should become [a]
//[or a b] should become [a] [b]
//[or a [b and c]] should become [a] [b and c]
//[[a or b] and [c or d]] should become [a and c] [a and d] [b and c] [b and d]
func SplitExpressionsOnOr(expr stconverter.STExpression) []stconverter.STExpression {
	//IF IS OR
	//	BREAK APART
	//IF IS VALUE
	//	RETURN CURRENT
	//IF IS OTHER OPERATOR
	//	MARK LOCATION AND RECURSE

	// broken := breakIfOr(expr)
	// if len(broken) == 1 {
	// 	return broken
	// }

	op := expr.HasOperator()
	if op == nil { //if it's just a value, return
		return []stconverter.STExpression{expr}
	}
	if op.GetToken() == "or" { //if it's an "or", return the arguments
		rets := make([]stconverter.STExpression, 0)
		args := expr.GetArguments()
		for i := 0; i < len(args); i++ { //for each argument of the "or", return it, unless it is itself an "or" (in which case, expand further)
			arg := args[i]
			argOp := arg.HasOperator()
			if argOp == nil || argOp.GetToken() != "or" {
				rets = append(rets, arg)
				continue
			}
			args = append(args, arg.GetArguments()...)
		}
		return rets
	}

	//otherwise, things are more interesting

	//make the thing we're returning
	rets := make([]stconverter.STExpressionOperator, 0)

	//build a new expression
	var nExpr stconverter.STExpressionOperator

	//operator is op, arguments are args
	nExpr.Operator = op
	args := expr.GetArguments()
	nExpr.Arguments = make([]stconverter.STExpression, len(args))

	rets = append(rets, nExpr)
	//for each argument in the expression operator
	for i, arg := range args {
		//get arguments to operator by calling SplitExpressionsOnOr again
		argT := SplitExpressionsOnOr(arg)
		//if argT has more than one value, it indicates that this argument was "split", and we should return two nExpr, one with each argument
		//we will increase the size of rets by a multiplyFactor, which is the size of argT
		//i.e. if we receive two arguments, and we already had two elements in rets, it indicates we need to return 4 values
		//for instance, if our original command was "(a or b) and (c or d)" we'd need to return 4 elements (a and c) (a and d) (b and c) (b and d)
		multiplyFactor := len(argT)
		//for each factor in multiplyFactor, duplicate rets[n]
		//e.g. multiplyFactor 2 on [1 2 3] becomes [1 1 2 2 3 3]
		//e.g. multiplyFactor 3 on [1 2 3] becomes [1 1 1 2 2 2 3 3 3]
		for y := 0; y < len(rets); y++ {
			for z := 1; z < multiplyFactor; z++ {

				var newElem stconverter.STExpressionOperator
				copyElem := rets[y]
				newElem.Operator = copyElem.Operator
				newElem.Arguments = make([]stconverter.STExpression, len(copyElem.Arguments))
				copy(newElem.Arguments, copyElem.Arguments)

				rets = append(rets, stconverter.STExpressionOperator{})
				copy(rets[y+1:], rets[y:])
				rets[y] = newElem
				y++
			}
		}

		//for each argument, copy it into the return elements at the appropriate locations
		//(if we have multiple arguments, they will be chosen in a round-robin fashion)
		for j := 0; j < len(argT); j++ {
			at := argT[j]
			for k := j; k < len(rets); k += len(argT) {
				rets[k].Arguments[i] = at
			}
		}

		//expected, _ := json.MarshalIndent(rets, "\t", "\t")
		//fmt.Printf("Current:\n\t%s\n\n", expected)
	}

	//conversion for returning
	actualRets := make([]stconverter.STExpression, len(rets))
	for i := 0; i < len(rets); i++ {
		actualRets[i] = rets[i]
	}
	return actualRets

}

//DeepGetValues recursively gets all values from a given stconverter.STExpression
func DeepGetValues(expr stconverter.STExpression) []string {
	if expr == nil {
		return nil
	}
	if val := expr.HasValue(); val != "" {
		return []string{val}
	}
	vals := make([]string, 0)
	for _, arg := range expr.GetArguments() {
		vals = append(vals, DeepGetValues(arg)...)
	}
	return vals
}

//SolveSTExpression will solve simple STExpressions
//It will project the solutionTransition onto the problemTransition
//Then, it will use the resulting transition with STMakeSolutionAssignments to
//convert the comparison into an assignment
func SolveSTExpression(il InterfaceList, inputPolicy bool, problemTransition PSTTransition, solutionTransition stconverter.STExpression) []stconverter.STExpression {

	//first we need to project the solutionTransition over the problemTransition

	//lets get all mentioned values in the problemTransition
	problemTransitionExpr := problemTransition.STGuard
	problemVals := DeepGetValues(problemTransitionExpr)

	//now lets classify them
	problemInputs := make([]string, 0)
	problemOutputs := make([]string, 0)
	problemInternals := make([]string, 0)

	for _, problemVal := range problemVals {
		if il.HasIONamed(true, problemVal) {
			problemInputs = append(problemInputs, problemVal)
			continue
		}
		if il.HasIONamed(false, problemVal) {
			problemOutputs = append(problemOutputs, problemVal)
			continue
		}
		problemInternals = append(problemInternals, problemVal)
	}

	//now let's do the projection
	var proposedSolution stconverter.STExpression
	if inputPolicy {
		if len(problemInputs) == 0 {
			//this is a time-based problem on the inputs, so we'll use all inputs to fix it
			nonAcceptableNames := make([]string, len(il.OutputVars))
			for i, v := range il.OutputVars {
				nonAcceptableNames[i] = v.Name
			}
			proposedSolution = ConvertSTExpressionForPolicy(il, nonAcceptableNames, true, solutionTransition)
		} else {
			acceptableNames := append(problemInputs, problemInternals...)
			proposedSolution = ConvertSTExpressionForPolicy(il, acceptableNames, false, solutionTransition)
			if proposedSolution == nil {
				//fmt.Printf("Well, that didn't work (1)\r\nacceptableNames:%v\r\nproblemTransition:%v\r\n", acceptableNames, stconverter.CCompileExpression(problemTransition))
				return nil
			}
		}
	} else {
		if len(problemOutputs) == 0 && len(problemInputs) == 0 { //if problemInputs != 0, then it is likely this was fixed already
			//this is a time-based problem on the outputs, so we'll use all outputs to fix it
			nonAcceptableNames := make([]string, len(il.InputVars))
			for i, v := range il.InputVars {
				nonAcceptableNames[i] = v.Name
			}
			proposedSolution = ConvertSTExpressionForPolicy(il, nonAcceptableNames, true, solutionTransition)
		} else {
			acceptableNames := append(problemOutputs, problemInternals...)
			proposedSolution = ConvertSTExpressionForPolicy(il, acceptableNames, false, solutionTransition)
			if proposedSolution == nil {
				//fmt.Printf("Well, that didn't work (2)\r\nacceptableNames:%v\r\nproblemTransition:%v\r\nsolutionTransition:%v\r\n", acceptableNames, stconverter.CCompileExpression(problemTransitionExpr), stconverter.CCompileExpression(solutionTransition))
				return nil
			}
		}
	}

	if proposedSolution == nil {
		//fmt.Printf("Well, that didn't work (3)\r\nsolutionTransition:%v\r\nproblemTransition:%v\r\n", stconverter.CCompileExpression(solutionTransition), stconverter.CCompileExpression(problemTransition))
		return nil
	}

	// problem := solutionTransition
	// if !inputPolicy {
	// 	//project problem on outputs to solve (as we can only effect vars in input or output depending on our problem scope)
	// 	acceptableNames := make([]string, len(il.InputVars))
	// 	for i, v := range il.InputVars {
	// 		acceptableNames[i] = v.Name
	// 	}
	// 	problem = ConvertSTExpressionForPolicy(il, acceptableNames, true, solutionTransition)
	// }

	//TODO: remove TIMERS from problem space if present

	return STMakeSolutionAssignments(proposedSolution)

	//return nil
}

//STMakeSolutionAssignments will convert a comparison stExpression into an assignment that meets the comparison
//The top level should be one of the following
//if VARIABLE ONLY, 			return VARIABLE = 1
//if NOT(VARIABLE) ONLY, 		return VARIABLE = 0
//if VARIABLE == EXPRESSION, 	return VARIABLE = EXPRESSION
//if VARIABLE > EXPRESSION, 	return VARIABLE = EXPRESSION + 1
//if VARIABLE >= EXPRESSION, 	return VARIABLE = EXPRESSION
//if VARIABLE < EXPRESSION, 	return VARIABLE = EXPRESSION - 1
//if VARIABLE <= EXPRESSION, 	return VARIABLE = EXPRESSION
//if VARIABLE != EXPRESSION,	return VARIABLE = EXPRESSION + 1
//otherwise, return nil (can't solve)
func STMakeSolutionAssignments(soln stconverter.STExpression) []stconverter.STExpression {
	op := soln.HasOperator()
	//if VARIABLE ONLY, 			return VARIABLE = 1
	if op == nil {
		return []stconverter.STExpression{
			stconverter.STExpressionOperator{
				Operator: stconverter.FindOp(":="),
				Arguments: []stconverter.STExpression{
					stconverter.STExpressionValue{Value: "1"},
					stconverter.STExpressionValue{Value: soln.HasValue()},
				},
			},
		}
	}

	if op.GetToken() == "and" {
		solns := make([]stconverter.STExpression, 0)
		for _, arg := range soln.GetArguments() {
			tempSoln := STMakeSolutionAssignments(arg)
			if arg != nil {
				solns = append(solns, tempSoln...)
			}
		}
		return solns
	}

	args := soln.GetArguments()

	//if NOT(VARIABLE) ONLY, 		return VARIABLE = 1
	if op.GetToken() == "not" && len(args) == 1 {
		return []stconverter.STExpression{stconverter.STExpressionOperator{
			Operator: stconverter.FindOp(":="),
			Arguments: []stconverter.STExpression{
				stconverter.STExpressionValue{Value: "0"},
				stconverter.STExpressionValue{Value: args[0].HasValue()},
			}}}
	}

	//if VARIABLE == EXPRESSION, 	return VARIABLE = EXPRESSION
	if op.GetToken() == "=" && len(args) == 2 {
		return []stconverter.STExpression{stconverter.STExpressionOperator{
			Operator:  stconverter.FindOp(":="),
			Arguments: args,
		}}
	}

	//if VARIABLE > EXPRESSION, 	return VARIABLE = EXPRESSION + 1
	if op.GetToken() == ">" && len(args) == 2 {
		return []stconverter.STExpression{stconverter.STExpressionOperator{
			Operator: stconverter.FindOp(":="),
			Arguments: []stconverter.STExpression{
				stconverter.STExpressionOperator{
					Operator: stconverter.FindOp("+"),
					Arguments: []stconverter.STExpression{
						stconverter.STExpressionValue{Value: "1"},
						args[1],
					},
				},
				stconverter.STExpressionValue{Value: args[0].HasValue()},
			},
		}}
	}

	//if VARIABLE >= EXPRESSION, 	return VARIABLE = EXPRESSION
	if op.GetToken() == ">=" && len(args) == 2 {
		return []stconverter.STExpression{stconverter.STExpressionOperator{
			Operator:  stconverter.FindOp(":="),
			Arguments: args,
		}}
	}

	//if VARIABLE < EXPRESSION, 	return VARIABLE = EXPRESSION - 1
	if op.GetToken() == ">" && len(args) == 2 {
		return []stconverter.STExpression{stconverter.STExpressionOperator{
			Operator: stconverter.FindOp(":="),
			Arguments: []stconverter.STExpression{
				stconverter.STExpressionOperator{
					Operator: stconverter.FindOp("-"),
					Arguments: []stconverter.STExpression{
						stconverter.STExpressionValue{Value: "1"},
						args[1],
					},
				},
				stconverter.STExpressionValue{Value: args[0].HasValue()},
			},
		}}
	}

	//if VARIABLE <= EXPRESSION, 	return VARIABLE = EXPRESSION
	if op.GetToken() == "<=" && len(args) == 2 {
		return []stconverter.STExpression{stconverter.STExpressionOperator{
			Operator:  stconverter.FindOp(":="),
			Arguments: args,
		}}
	}

	//if VARIABLE != EXPRESSION,	return VARIABLE = EXPRESSION + 1
	if op.GetToken() == "<>" && len(args) == 2 {
		return []stconverter.STExpression{stconverter.STExpressionOperator{
			Operator: stconverter.FindOp(":="),
			Arguments: []stconverter.STExpression{
				stconverter.STExpressionOperator{
					Operator: stconverter.FindOp("+"),
					Arguments: []stconverter.STExpression{
						stconverter.STExpressionValue{Value: "1"},
						args[1],
					},
				},
				stconverter.STExpressionValue{Value: args[0].HasValue()},
			},
		}}
	}

	//If still here, we don't know what to do
	fmt.Println("WARNING: I couldn't solve guard \"", stconverter.STCompileExpression(soln), "\"")
	return nil
}
