package rtec

import (
	"regexp"
	"strings"

	"github.com/PRETgroup/easy-rte/rtedef"
)

//VhdlECCTransition is used with getVhdlECCTransitionCondition to return results to the template
type VhdlECCTransition CECCTransition

//getVhdlECCTransitionCondition returns the C "if" condition to use in state machine next state logic and associated events
// returns "full condition", "associated events"
func getVhdlECCTransitionCondition(function rtedef.EnforcedFunction, trans string) VhdlECCTransition {
	var events []string

	re1 := regexp.MustCompile("([<>=!]+)")          //for capturing operators
	re2 := regexp.MustCompile("([a-zA-Z0-9_<>=]+)") //for capturing variable and event names and operators
	isNum := regexp.MustCompile("^[0-9.]+$")

	retVal := trans

	//rename AND and OR
	retVal = strings.Replace(retVal, "AND", "&&", -1)
	retVal = strings.Replace(retVal, "OR", "||", -1)

	//re1: add whitespace around operators
	retVal = re1.ReplaceAllStringFunc(retVal, func(in string) string {
		if in != "!" {
			return " " + in + " "
		}
		return " !"
	})

	//re2: add "me->" where appropriate
	retVal = re2.ReplaceAllStringFunc(retVal, func(in string) string {
		if strings.ToLower(in) == "and" || strings.ToLower(in) == "or" || strings.ContainsAny(in, "!><=") || strings.ToLower(in) == "true" || strings.ToLower(in) == "false" {
			//no need to make changes, these aren't variables or events
			return in
		}

		if isNum.MatchString(in) {
			//no need to make changes, it is a numerical value of some sort
			return in
		}

		//check to see if it is input data
		if function.InputVars != nil {
			for _, Var := range function.InputVars {
				if Var.Name == in {
					return "inputs->" + in
				}
			}
		}

		//check to see if it is output data
		if function.OutputVars != nil {
			for _, Var := range function.OutputVars {
				if Var.Name == in {
					return "outputs->" + in
				}
			}
		}

		//check to see if it is a policy internal var
		for i := 0; i < len(function.Policies); i++ {
			for _, Var := range function.Policies[i].InternalVars {
				if Var.Name == in {
					return "me->" + in
				}
				if Var.Name+"_i" == in {
					return "me->" + in
				}
			}
		}

		//else, return it (no idea what else to do!) - it might be a function call or strange text constant
		return in
	})

	//tidy the whitespace
	retVal = strings.Replace(retVal, "  ", " ", -1)

	return VhdlECCTransition{IfCond: retVal, AssEvents: events}
}

//getVhdlType returns the VHDL type to use with respect to an IEC61499 type
func getVhdlType(ctype string) string {
	vhdlType := ""

	switch strings.ToLower(ctype) {
	case "bool":
		vhdlType = "std_logic"
	case "char":
		vhdlType = "unsigned(7 downto 0)"
	case "uint8_t":
		vhdlType = "unsigned(7 downto 0)"
	case "uint16_t":
		vhdlType = "unsigned(15 downto 0)"
	case "uint32_t":
		vhdlType = "unsigned(31 downto 0)"
	case "uint64_t":
		vhdlType = "unsigned(63 downto 0)"
	case "int8_t":
		vhdlType = "signed(7 downto 0)"
	case "int16_t":
		vhdlType = "signed(15 downto 0)"
	case "int32_t":
		vhdlType = "signed(31 downto 0)"
	case "int64_t":
		vhdlType = "signed(63 downto 0)"
	case "float":
		panic("Float type not allowed in conversion")
	case "double":
		panic("Double type not allowed in conversion")
	case "dtimer_t":
		vhdlType = "unsigned(63 downto 0)"
	case "rtimer_t":
		panic("rtimer type not allowed in conversion")
	default:
		panic("Unknown type: " + ctype)
	}

	return vhdlType
}
