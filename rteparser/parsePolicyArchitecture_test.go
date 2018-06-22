package rteparser

import (
	"testing"

	"github.com/PRETgroup/easy-rte/rtedef"
)

var efbArchitectureTests = []ParseTest{
	{
		Name: "missing brace after s1",
		Input: `function testBlock;
				interface of testBlock{
				}
				policy of testBlock {
					states {
						s1 

					}
				}`,
		Err: ErrUnexpectedValue,
	},
	{
		Name: "AEIPolicy",
		Input: `function AEIPolicy;
				interface of AEIPolicy {
					in bool AS, VS; //in here means that they're going from PLANT to CONTROLLER
					out bool AP, VP;//out here means that they're going from CONTROLLER to PLANT
				
					in uint64_t AEI_ns := 900000000;
				}
				policy AEI of AEIPolicy {
					internals {
						dtimer tAEI; //DTIMER increases in DISCRETE TIME continuously
					}
				
					//P3: AS or AP must be true within AEI after a ventricular event VS or VP.
				
					states {
						s1 {
							//-> <destination> [on guard] [: output expression][, output expression...] ;
							-> s2 on (VS or VP): tAEI := 0;
						}
				
						s2 {
							-> s1 on (AS or AP);
							-> violation on (tAEI > AEI_ns);
						}
					} 
				}`,
		Output: []rtedef.EnforcedFunction{
			rtedef.EnforcedFunction{
				Name: "AEIPolicy",
				Inputs: []rtedef.Variable{
					rtedef.Variable{Name: "AS", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
					rtedef.Variable{Name: "VS", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
					rtedef.Variable{Name: "AEI_ns", Type: "uint64_t", ArraySize: "", InitialValue: "900000000", Comment: ""},
				},
				Outputs: []rtedef.Variable{
					rtedef.Variable{Name: "AP", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
					rtedef.Variable{Name: "VP", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
				},
				Policies: []rtedef.Policy{
					rtedef.Policy{
						Name: "AEI",
						InternalVars: []rtedef.Variable{
							rtedef.Variable{Name: "tAEI", Type: "DTIMER", ArraySize: "", InitialValue: "", Comment: ""},
						},
						States: []rtedef.PState{"s1", "s2"},
						Transitions: []rtedef.PTransition{
							rtedef.PTransition{Source: "s1", Destination: "s2", Condition: "( VS or VP )", Expressions: []rtedef.PExpression{rtedef.PExpression{VarName: "tAEI", Value: "0"}}},
							rtedef.PTransition{Source: "s2", Destination: "s1", Condition: "( AS or AP )", Expressions: []rtedef.PExpression(nil)},
							rtedef.PTransition{Source: "s2", Destination: "violation", Condition: "( tAEI > AEI_ns )", Expressions: []rtedef.PExpression(nil)},
						},
					},
				},
			},
		},
	},
	{
		Name: "AB5Policy",
		Input: `function AB5Policy;
			interface of AB5Policy {
				in bool A;  //in here means that they're going from PLANT to CONTROLLER
				out bool B; //out here means that they're going from CONTROLLER to PLANT
			}
			
			policy AB5 of AB5Policy {
				internals {
					dtimer v;
				}
			
				states {
					s0 {														//first state is initial, and represents "We're waiting for an A"
						-> s0 on (!A and !B): v := 0;							//if we receive neither A nor B, do nothing
						-> s1 on (A and !B): v := 0;							//if we receive an A only, head to state s1
						-> violation on ((!A and B) or (A and B));				//if we receive a B, or an A and a B (i.e. if we receive a B) then VIOLATION
					}
			
					s1 {														//s1 is "we're waiting for a B, and it needs to get here within 5 ticks"
						-> s1 on (!A and !B and v < 5);							//if we receive nothing, and we aren't over-time, then we do nothing
						-> s0 on (!A and B);									//if we receive a B only, head to state s0
						-> violation on ((v >= 5) or (A and B) or (A and !B));	//if we go overtime, or we receive another A, then VIOLATION
					}
				}
			}`,
		Output: []rtedef.EnforcedFunction{
			rtedef.EnforcedFunction{
				Name: "AB5Policy",
				Inputs: []rtedef.Variable{
					rtedef.Variable{Name: "A", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
				},
				Outputs: []rtedef.Variable{
					rtedef.Variable{Name: "B", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
				},
				Policies: []rtedef.Policy{
					rtedef.Policy{
						Name: "AB5",
						InternalVars: []rtedef.Variable{
							rtedef.Variable{Name: "v", Type: "DTIMER", ArraySize: "", InitialValue: "", Comment: ""},
						},
						States: []rtedef.PState{"s0", "s1"},
						Transitions: []rtedef.PTransition{
							rtedef.PTransition{Source: "s0", Destination: "s0", Condition: "( !A and !B )", Expressions: []rtedef.PExpression{rtedef.PExpression{VarName: "v", Value: "0"}}},
							rtedef.PTransition{Source: "s0", Destination: "s1", Condition: "( A and !B )", Expressions: []rtedef.PExpression{rtedef.PExpression{VarName: "v", Value: "0"}}},
							rtedef.PTransition{Source: "s0", Destination: "violation", Condition: "( ( !A and B ) or ( A and B ) )", Expressions: []rtedef.PExpression(nil)},
							rtedef.PTransition{Source: "s1", Destination: "s1", Condition: "( !A and !B and v < 5 )", Expressions: []rtedef.PExpression(nil)},
							rtedef.PTransition{Source: "s1", Destination: "s0", Condition: "( !A and B )", Expressions: []rtedef.PExpression(nil)},
							rtedef.PTransition{Source: "s1", Destination: "violation", Condition: "( ( v >= 5 ) or ( A and B ) or ( A and !B ) )", Expressions: []rtedef.PExpression(nil)},
						},
					},
				},
			},
		},
	},
}

func TestParsePFBArchitecture(t *testing.T) {
	runParseTests(t, efbArchitectureTests)
}
