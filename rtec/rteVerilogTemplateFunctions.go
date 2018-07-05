package rtec

import (
	"fmt"
	"strings"

	"github.com/PRETgroup/easy-rte/rtedef"
)

//VerilogECCTransition is used with getVerilogECCTransitionCondition to return results to the template
type VerilogECCTransition CECCTransition

//getVerilogECCTransitionCondition returns the C "if" condition to use in state machine next state logic and associated events
// returns "full condition", "associated events"
func getVerilogECCTransitionCondition(function rtedef.EnforcedFunction, trans string) VerilogECCTransition {

	return VerilogECCTransition{IfCond: trans, AssEvents: nil}
}

//getVerilogType returns the VHDL type to use with respect to an IEC61499 type
func getVerilogType(ctype string) string {
	return "reg " + getVerilogWidthArrayForType(ctype)
}

func getVerilogWidthArrayForType(ctype string) string {
	verilogType := ""

	switch strings.ToLower(ctype) {
	case "bool":
		verilogType = ""
	case "char":
		verilogType = "[7:0]"
	case "uint8_t":
		verilogType = "unsigned [7:0]"
	case "uint16_t":
		verilogType = "unsigned [15:0]"
	case "uint32_t":
		verilogType = "unsigned [31:0]"
	case "uint64_t":
		verilogType = "unsigned [63:0]"
	case "int8_t":
		verilogType = "signed [7:0]"
	case "int16_t":
		verilogType = "signed [15:0]"
	case "int32_t":
		verilogType = "signed [31:0]"
	case "int64_t":
		verilogType = "signed [63:0]"
	case "float":
		panic("Float type not allowed in conversion")
	case "double":
		panic("Double type not allowed in conversion")
	case "dtimer_t":
		verilogType = "unsigned [63:0]"
	case "rtimer_t":
		panic("rtimer type not allowed in conversion")
	default:
		panic("Unknown type: " + ctype)
	}

	return verilogType
}

func getVerilogWidthArray(l int) string {
	cl2 := ceilLog2(uint64(l)) - 1
	if cl2 >= 1 {
		return fmt.Sprintf("[%v:0]", cl2)
	}
	return ""
}

var t = [6]uint64{
	0xFFFFFFFF00000000,
	0x00000000FFFF0000,
	0x000000000000FF00,
	0x00000000000000F0,
	0x000000000000000C,
	0x0000000000000002,
}

//ceilLog2 performs a log2 ceiling function quickly
func ceilLog2(x uint64) int {

	y := 0
	if (x & (x - 1)) != 0 {
		y = 1
	}
	j := 32
	var i int

	for i = 0; i < 6; i++ {
		k := 0
		if (x & t[i]) != 0 {
			k = j
		}
		y += k
		x >>= uint64(k)
		j >>= 1
	}

	return y
}
