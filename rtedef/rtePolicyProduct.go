package rtedef

import (
	"errors"

	"github.com/PRETgroup/goFB/goFB/stconverter"
)

//PolicyProduct takes the product as defined by the EMSOFT '17 paper of two policies
//in an EnforcedFunction
//
//Given two discrete TAs A1 = (L1 , l01 , lv1 , Σ1 , V1 , ∆1 , F1 ) and
//A2 = (L2 , l02 , lv2 , Σ2 , V2 , ∆2 , F2 ) with disjoint sets of integer clocks,
//their product is the DTA A1 × A2
//A = (L, l0 , lv , Σ, V , ∆, F ) where
// L = L1 × L2 ,
// l0 = (l01 , l02 ),
// lv = (lv1 , lv2 ),
// V = V1 ∪ V2 ,
// F = F1 × F2 ,
// and ∆ ⊆ L × G(V ) × R × Σ × L is the transition relation,
// with ((l1 , l2 ), g1 ∧ g2 , R1 ∪ R2 , a, (l1' , l2' )) ∈ ∆
// if (l1 , g1 , R1 , a, l1' ) ∈ ∆1 and (l2 , g2 , R2 , a, l2' ) ∈ ∆2 .
// In the product DTA, all locations in (L1 × lv2 ) ∪ (lv1 × L2 ) are trap locations,
// and all the outgoing transitions for these locations can be replaced with
// self loops.
// We consider merging all the trap locations into a single location lv where any outgoing
// transition from any location in L \ (L1 × lv2  ) ∪ (lv1 × L2 ) to a location in (L1 × lv2 ) ∪ (lv1 × L2 ) goes to
// lv instead.
//
// The product of DTAs is useful when we want to enforce multiple properties. Given two determin-
// istic and complete DTAs A1 and A2 the DTA A obtained by computing their product recognizes
// the language L(A1) ∩ L(A2), and is also deterministic and complete.
//
func (ef EnforcedFunction) PolicyProduct(polA Policy, polB Policy) (Policy, error) {
	polProduct := Policy{
		Name: "Product_of_" + polA.Name + "_and_" + polB.Name,
	}

	//combine the locations
	// L = L1 × L2
	// l0 = (l01 , l02 ) [the initial state] is handled implicitly, as the new state 0 is states i=0, j=0
	// lv = (lv1 , lv2 ) [the violation state] is handled implicitly, as there is no "violation state" actually represented in the states slices
	// F = F1 × F2 [the accepting states] are handled implicitly, as we don't really explicitly enumerate these in this system
	for i := 0; i < len(polA.States); i++ {
		for j := 0; j < len(polB.States); j++ {
			newStateName := polA.States[i].Name() + "_comma_" + polB.States[j].Name()
			polProduct.States = append(polProduct.States, PState(newStateName))
		}
	}

	//combine the internal variables
	// V = V1 ∪ V2
	for i := 0; i < len(polA.InternalVars); i++ {
		for j := 0; j < len(polB.InternalVars); j++ {
			if polA.InternalVars[i].Name == polB.InternalVars[j].Name {
				return Policy{}, errors.New("Taking the product of policies is not possible if they share internal variable names")
			}
		}
	}
	polProduct.InternalVars = append(polA.InternalVars, polB.InternalVars...)

	//combine the transitions
	//with ((l1 , l2 ), g1 ∧ g2 , R1 ∪ R2 , a, (l1' , l2' )) ∈ ∆
	// if (l1 , g1 , R1 , a, l1' ) ∈ ∆1 and (l2 , g2 , R2 , a, l2' ) ∈ ∆2 .
	//now, a little hiccup
	//we don't store g or R separately, instead they are bundled into the transition condition altogether
	//but, we can emulate (g1 ∧ g2 , R1 ∪ R2) by just taking the AND of the transitions
	//i.e. to go from l1 to l1' takes condition g1, R1, stored as T1
	//to get from l2 to l2' takes condition g2, R2, stored as T2
	//so to get from (l1, l2) to (l1',l2') takes condition T1 AND T2
	//it's a little different for violation transitions
	//if half a transition combination is violation, we ignore the non-violation half
	//then we will add it as it was before but with the new source
	polAPSTTrans, err := polA.GetPSTTransitions()
	if err != nil {
		return Policy{}, err
	}
	polBPSTTrans, err := polB.GetPSTTransitions()
	if err != nil {
		return Policy{}, err
	}
	//out:
	for i := 0; i < len(polAPSTTrans); i++ {
		for j := 0; j < len(polBPSTTrans); j++ {

			newTransGuard := stconverter.STExpressionOperator{
				Operator:  stconverter.FindOp("and"),
				Arguments: []stconverter.STExpression{polAPSTTrans[i].STGuard, polBPSTTrans[j].STGuard},
			}
			source := polAPSTTrans[i].Source + "_comma_" + polBPSTTrans[j].Source

			if polAPSTTrans[i].Destination == "violation" || polBPSTTrans[j].Destination == "violation" {
				if polAPSTTrans[i].Destination == "violation" {
					polProduct.Transitions = append(polProduct.Transitions, PTransition{
						Source:      source,
						Destination: "violation",
						Condition:   stconverter.STCompileExpression(polAPSTTrans[i].STGuard),
						Expressions: polAPSTTrans[i].Expressions,
						Recover:     polAPSTTrans[i].Recover,
					})
					//continue out
				}
				if polBPSTTrans[j].Destination == "violation" {
					polProduct.Transitions = append(polProduct.Transitions, PTransition{
						Source:      source,
						Destination: "violation",
						Condition:   stconverter.STCompileExpression(polBPSTTrans[j].STGuard),
						Expressions: polBPSTTrans[j].Expressions,
						Recover:     polBPSTTrans[j].Recover,
					})
					continue
				}
			} else {
				destination := polAPSTTrans[i].Destination + "_comma_" + polBPSTTrans[j].Destination
				polProduct.Transitions = append(polProduct.Transitions, PTransition{
					Source:      source,
					Destination: destination,
					Condition:   stconverter.STCompileExpression(newTransGuard),
					Expressions: append(polAPSTTrans[i].Expressions, polBPSTTrans[j].Expressions...),
					Recover:     append(polAPSTTrans[i].Recover, polBPSTTrans[j].Recover...),
				})
			}

		}
	}

	return polProduct, nil
}
