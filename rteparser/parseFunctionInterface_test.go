package rteparser

import (
	"testing"

	"github.com/PRETgroup/easy-rte/rtedef"
)

var interfaceTests = []ParseTest{

	{
		Name: "events typo 1",
		Input: `function testBlock;
					interface of testBlock {
						in bool inEvent;
						out outEvent;
					}`,
		Err: ErrInvalidType,
	},
	{
		Name: "events typo 2",
		Input: `function testBlock;
					interface of testBlock {
						in bool inEvent;
						out bool outEvent
					}`,
		Err: ErrUnexpectedValue,
	},
	{
		Name: "data typo 1",
		Input: `function testBlock;
					interface of testBlock {
						in bool inEvent;
						in asdasd inData;
						out bool outEvent;
					}`,
		Err: ErrInvalidType,
	},
	{
		Name: "data input 1",
		Input: `function testBlock;
					interface of testBlock {
						in bool inEvent;
						out bool outEvent;
					}`,
		Output: []rtedef.EnforcedFunction{
			rtedef.EnforcedFunction{
				Name: "testBlock",
				InterfaceList: rtedef.InterfaceList{
					InputVars: []rtedef.Variable{
						rtedef.Variable{Name: "inEvent", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
					},
					OutputVars: []rtedef.Variable{
						rtedef.Variable{Name: "outEvent", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
					},
				},
				Policies: []rtedef.Policy(nil)},
		},
		Err: nil,
	},
	{
		Name: "data input 2",
		Input: `function testBlock;
					interface of testBlock {
						in bool inEvent;
						in bool[3] inData;
						out bool outEvent;
					}`,
		Output: []rtedef.EnforcedFunction{
			rtedef.EnforcedFunction{
				Name: "testBlock",
				InterfaceList: rtedef.InterfaceList{
					InputVars: []rtedef.Variable{
						rtedef.Variable{Name: "inEvent", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
						rtedef.Variable{Name: "inData", Type: "bool", ArraySize: "3", InitialValue: "", Comment: ""},
					},
					OutputVars: []rtedef.Variable{
						rtedef.Variable{Name: "outEvent", Type: "bool", ArraySize: "", InitialValue: "", Comment: ""},
					},
				},
				Policies: []rtedef.Policy(nil),
			},
		},
		Err: nil,
	},
	{
		Name: "data input array typo 1",
		Input: `basicFB testBlock;
					interface of testBlock {
						in event inEvent;
						in bool[3 inData;
						out event outEvent;
					}`,
		Err: ErrUnexpectedValue,
	},
	{
		Name: "data input 3",
		Input: `function testBlock;
					interface of testBlock {
						in int8_t inEvent;
						in bool[3] inData := [0,1,0];
						out char outEvent;
					}`,
		Output: []rtedef.EnforcedFunction{
			rtedef.EnforcedFunction{
				Name: "testBlock",
				InterfaceList: rtedef.InterfaceList{
					InputVars: []rtedef.Variable{
						rtedef.Variable{Name: "inEvent", Type: "int8_t", ArraySize: "", InitialValue: "", Comment: ""},
						rtedef.Variable{Name: "inData", Type: "bool", ArraySize: "3", InitialValue: "[0,1,0]", Comment: ""},
					},
					OutputVars: []rtedef.Variable{
						rtedef.Variable{Name: "outEvent", Type: "char", ArraySize: "", InitialValue: "", Comment: ""},
					},
				},
				Policies: []rtedef.Policy(nil)}},
		Err: nil,
	},
	{
		Name: "data default typo 1",
		Input: `basicFB testBlock;
					interface of testBlock {
						in bool inEvent;
						in bool[3] inData := 0,1,0;
						out bool outEvent;
					}`,
		Err: ErrUnexpectedValue,
	},
	{
		Name: "Unexpected EOF",
		Input: `basicFB testBlock;
					interface of testBlock {
						in bool inEvent;
						in bool inEvent2;
						out int32_t `,
		Err: ErrUnexpectedValue,
	},
}

func TestParseStringInterface(t *testing.T) {
	runParseTests(t, interfaceTests)
}
