package rtec

import (
	"fmt"
	"strings"

	"github.com/PRETgroup/easy-rte/rtedef"
)

//VhdlECCTransition is used with getVhdlECCTransitionCondition to return results to the template
type VhdlECCTransition CECCTransition

//getVhdlECCTransitionCondition returns the C "if" condition to use in state machine next state logic and associated events
// returns "full condition", "associated events"
func getVhdlECCTransitionCondition(function rtedef.EnforcedFunction, trans string) VhdlECCTransition {

	return VhdlECCTransition{IfCond: trans, AssEvents: nil}
}

//getVhdlType returns the VHDL type to use with respect to an IEC61499 type
func getVhdlType(ctype string) string {
	vhdlType := ""

	switch strings.ToLower(ctype) {
	case "bool":
		vhdlType = "integer range 0 to 1"
	case "char":
		vhdlType = "integer range 0 to (2**8 - 1)"
	case "uint8_t":
		vhdlType = "integer range 0 to (2**8 - 1)"
	case "uint16_t":
		vhdlType = "integer range 0 to (2**16 - 1)"
	case "uint32_t":
		fmt.Printf("WARNING: uint32_t type constrainted to max value 2^31-1 in VHDL")
		vhdlType = "integer range 0 to (2**31 - 1)"
	case "uint64_t":
		panic("uint16_t type not allowed in conversion")
	case "int8_t":
		vhdlType = "integer range -(2**7) to (2**7-1)"
	case "int16_t":
		vhdlType = "integer range -(2**15) to (2**15-1)"
	case "int32_t":
		vhdlType = "integer"
	case "int64_t":
		panic("int64_t type not allowed in conversion")
	case "float":
		panic("Float type not allowed in conversion")
	case "double":
		panic("Double type not allowed in conversion")
	case "dtimer_t":
		vhdlType = "integer range 0 to (2**31 - 1)"
	case "rtimer_t":
		panic("rtimer type not allowed in conversion")
	default:
		panic("Unknown type: " + ctype)
	}

	return vhdlType
}
