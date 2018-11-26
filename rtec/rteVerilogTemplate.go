package rtec

import (
	"text/template"

	"github.com/PRETgroup/goFB/goFB/stconverter"
)

const rteVerilogTemplate = `{{define "_policyIn"}}{{$block := .}}
//input policies
module inputEditMux_{{$block.Name}}(
	//inputs (plant to controller){{range $index, $var := $block.InputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_in_inputmux,
	output wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_out_inputmux,
	{{end}}

	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}},
	{{end}}{{end}}{{end}}

	//state variables
	{{range $polI, $pol := $block.Policies}}{{if $polI}},
	{{end}}input wire {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state{{end}}
);

{{range $index, $var := $block.InputVars}}
{{getVerilogType $var.Type}} {{$var.Name}} {{if $var.InitialValue}}/* = {{$var.InitialValue}}*/{{end}};
{{end}}

always @* begin
	//capture synchronous inputs
	{{range $index, $var := $block.InputVars}}
		{{$var.Name}} = {{$var.Name}}_ptc_in_inputmux;
	{{end}}

	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}
	{{if not $pfbEnf}}//{{$pol.Name}} is broken!
	{{else}}{{/* this is where the policy comes in */}}//INPUT POLICY {{$pol.Name}} BEGIN 
		case({{$block.Name}}_policy_{{$pol.Name}}_state)
		{{range $sti, $st := $pol.States}}` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$st.Name}}: begin
				{{range $tri, $tr := $pfbEnf.InputPolicy.GetViolationTransitions}}{{if eq $tr.Source $st.Name}}{{/*
				*/}}
				if ({{$cond := getVerilogECCTransitionCondition $block (compileExpression $tr.STGuard)}}{{$cond.IfCond}}) begin
					//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
					//select a transition to solve the problem
					{{$solution := $pfbEnf.SolveViolationTransition $tr true}}
					{{if $solution.Comment}}//{{$solution.Comment}}{{end}}
					{{range $soleI, $sole := $solution.Expressions}}{{$sol := getVerilogECCTransitionCondition $block (compileExpression $sole)}}{{$sol.IfCond}};
					{{end}}
				end{{end}}{{end}}
			end
			{{end}}
		endcase
	{{end}}
	//INPUT POLICY {{$pol.Name}} END
	{{end}}

end

//emit outputs
{{range $index, $var := $block.InputVars}}
	assign {{$var.Name}}_ptc_out_inputmux = {{$var.Name}};
{{end}}

endmodule
{{end}}

{{define "_policyOut"}}{{$block := .}}
module outputEditMux_{{$block.Name}} (
	//inputs (plant to controller){{range $index, $var := $block.InputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_in_outputmux,
	{{end}}
	//outputs (controller to plant){{range $index, $var := $block.OutputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ctp_in_outputmux,
	output wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ctp_out_outputmux,
	{{end}}
	
	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_in,
	{{end}}{{end}}{{end}}

	//state variables
	{{range $polI, $pol := $block.Policies}}{{if $polI}},
	{{end}}input wire {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state{{end}}
);

{{range $index, $var := $block.InputVars}}
wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_out_inputmux;
{{getVerilogType $var.Type}} {{$var.Name}}{{if $var.InitialValue}}/* = {{$var.InitialValue}}*/{{end}};
{{end}}
{{range $index, $var := $block.OutputVars}}
{{getVerilogType $var.Type}} {{$var.Name}}{{if $var.InitialValue}}/* = {{$var.InitialValue}}*/{{end}};
{{end}}
{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}reg {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}};
{{end}}{{end}}{{end}}


always @* begin
	//capture synchronous inputs
	{{range $index, $var := $block.InputVars}}
		{{$var.Name}} = {{$var.Name}}_ptc_in_outputmux;
	{{end}}{{range $index, $var := $block.OutputVars}}
		{{$var.Name}} = {{$var.Name}}_ctp_in_outputmux;
	{{end}}
	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $varI, $var := $pfbEnf.OutputPolicy.InternalVars}}
	{{$var.Name}} = {{$var.Name}}_in {{if $var.IsDTimer}} + 1{{end}};
	{{end}}{{end}}{{end}}

	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}
	{{if not $pfbEnf}}//{{$pol.Name}} is broken!
	{{else}}{{/* this is where the policy comes in */}}//OUTPUT POLICY {{$pol.Name}} BEGIN 
		
		case({{$block.Name}}_policy_{{$pol.Name}}_state)
		{{range $sti, $st := $pol.States}}` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$st.Name}}: begin
				{{range $tri, $tr := $pfbEnf.OutputPolicy.GetViolationTransitions}}{{if eq $tr.Source $st.Name}}{{/*
				*/}}
				if ({{$cond := getVerilogECCTransitionCondition $block (compileExpression $tr.STGuard)}}{{$cond.IfCond}}) begin
					//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
					//select a transition to solve the problem
					{{$solution := $pfbEnf.SolveViolationTransition $tr false}}
					{{if $solution.Comment}}//{{$solution.Comment}}{{end}}
					{{range $soleI, $sole := $solution.Expressions}}{{$sol := getVerilogECCTransitionCondition $block (compileExpression $sole)}}{{$sol.IfCond}};
					{{end}}
				end {{end}}{{end}}
			end
			{{end}}
		endcase		
	{{end}}
	//OUTPUT POLICY {{$pol.Name}} END
	{{end}}

end

//emit outputs
{{range $index, $var := $block.OutputVars}}
	assign {{$var.Name}}_ctp_out_outputmux = {{$var.Name}};
{{end}}

endmodule
{{end}}

{{define "_nextStateFunction"}}{{$block := .}}
module nextStateFunction_{{$block.Name}} (
	//inputs (plant to controller){{range $index, $var := $block.InputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}},
	{{end}}
	//outputs (controller to plant){{range $index, $var := $block.OutputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}},
	{{end}}
	
	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_in,
	{{end}}{{end}}{{end}}

	//state variables
	{{range $polI, $pol := $block.Policies}}{{if $polI}},
	{{end}}input wire {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state_in,
	output wire {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state_next{{end}}
);

{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}reg {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}};
{{end}}{{end}}{{end}}

//For each policy, we need a reg for the state machine
{{range $polI, $pol := $block.Policies}}reg {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state;
reg transTaken_{{$block.Name}}_policy_{{$pol.Name}}; //EBMC liveness check register flag (will be optimised away in normal compiles)
{{end}}

always @* begin
	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//policy is broken!{{else}}
	{{range $varI, $var := $pfbEnf.OutputPolicy.InternalVars}}
	{{$var.Name}} = {{$var.Name}}_in {{if $var.IsDTimer}} + 1{{end}};
	{{end}}

	//mark no transition taken
	transTaken_{{$block.Name}}_policy_{{$pol.Name}} = 0;
	{{$block.Name}}_policy_{{$pol.Name}}_state = ` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_violation;

	//select transition to advance state
	case({{$block.Name}}_policy_{{$pol.Name}}_state_in)
	{{range $sti, $st := $pol.States}}` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$st.Name}}: begin
		{{range $tri, $tr := $pfbEnf.OutputPolicy.GetTransitionsForSource $st.Name}}{{/*
		*/}}
		{{if $tri}}else {{end}}if ({{$cond := getVerilogECCTransitionCondition $block (compileExpression $tr.STGuard)}}{{$cond.IfCond}}) begin
			//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
			{{$block.Name}}_policy_{{$pol.Name}}_state = ` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$tr.Destination}};
			//set expressions
			{{range $exi, $ex := $tr.Expressions}}
			{{$ex.VarName}} = {{$ex.Value}};{{end}}
			transTaken_{{$block.Name}}_policy_{{$pol.Name}} = 1;
		end {{end}}
	end
	{{end}}
	endcase
	{{end}}{{end}}

	//For each policy, ensure correctness (systemverilog only) and liveness
	{{range $polI, $pol := $block.Policies}}//assert property ({{$block.Name}}_policy_{{$pol.Name}}_state != ` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_violation);
	//assert property (transTaken_{{$block.Name}}_policy_{{$pol.Name}} == 1);
	{{end}}
end

//For each policy, emit state
{{range $polI, $pol := $block.Policies}}assign {{$block.Name}}_policy_{{$pol.Name}}_state_next = {{$block.Name}}_policy_{{$pol.Name}}_state;
{{end}}

endmodule
{{end}}

{{define "_combinatorialVerilog"}}{{$block := .}}

module combinatorialVerilog_{{$block.Name}} (
	//inputs (plant to controller){{range $index, $var := $block.InputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_in,
	output wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_out,
	{{end}}
	//outputs (controller to plant){{range $index, $var := $block.OutputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ctp_in,
	output wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ctp_out,
	{{end}}
	
	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}{{if not $var.Constant}}input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}},
	{{end}}{{end}}{{end}}{{end}}

	//state variables
	{{range $polI, $pol := $block.Policies}}{{if $polI}},
	{{end}}input wire {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state_in,
	output wire {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state_next{{end}}
);

{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}{{if $var.Constant}}{{getVerilogType $var.Type}} {{$var.Name}} = {{$var.InitialValue}};
{{end}}{{end}}{{end}}{{end}}

inputEditMux_{{$block.Name}} inputEditMux (
	//inputs (plant to controller){{range $index, $var := $block.InputVars}}
	.{{$var.Name}}_ptc_in_inputmux({{$var.Name}}_ptc_in),
	.{{$var.Name}}_ptc_out_inputmux({{$var.Name}}_ptc_out),
	{{end}}

	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}.{{$var.Name}}({{$var.Name}}),
	{{end}}{{end}}{{end}}

	{{range $polI, $pol := $block.Policies}}{{if $polI}},
	{{end}}.{{$block.Name}}_policy_{{$pol.Name}}_state({{$block.Name}}_policy_{{$pol.Name}}_state_in){{end}}
);

outputEditMux_{{$block.Name}} outputEditMux (
	//inputs (plant to controller){{range $index, $var := $block.InputVars}}
	.{{$var.Name}}_ptc_in_outputmux({{$var.Name}}_ptc_out),
	{{end}}

	//outputs (controller to plant){{range $index, $var := $block.OutputVars}}
	.{{$var.Name}}_ctp_in_outputmux({{$var.Name}}_ctp_in),
	.{{$var.Name}}_ctp_out_outputmux({{$var.Name}}_ctp_out),
	{{end}}

	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}.{{$var.Name}}_in({{$var.Name}}),
	{{end}}{{end}}{{end}}

	{{range $polI, $pol := $block.Policies}}{{if $polI}},
	{{end}}.{{$block.Name}}_policy_{{$pol.Name}}_state({{$block.Name}}_policy_{{$pol.Name}}_state_in){{end}}
);

nextStateFunction_{{$block.Name}} nextStateFunction (
	//inputs (plant to controller){{range $index, $var := $block.InputVars}}
	.{{$var.Name}}({{$var.Name}}_ptc_out),
	{{end}}
	//outputs (controller to plant){{range $index, $var := $block.OutputVars}}
	.{{$var.Name}}({{$var.Name}}_ctp_out),
	{{end}}
	
	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}.{{$var.Name}}_in({{$var.Name}}),
	{{end}}{{end}}{{end}}

	//state variables
	{{range $polI, $pol := $block.Policies}}{{if $polI}},
	{{end}}.{{$block.Name}}_policy_{{$pol.Name}}_state_in({{$block.Name}}_policy_{{$pol.Name}}_state_in),
	.{{$block.Name}}_policy_{{$pol.Name}}_state_next({{$block.Name}}_policy_{{$pol.Name}}_state_next){{end}}
);

//For each policy, ensure correctness (systemverilog only) and liveness
{{range $polI, $pol := $block.Policies}}assert property ({{$block.Name}}_policy_{{$pol.Name}}_state_in < ` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_violation |-> {{$block.Name}}_policy_{{$pol.Name}}_state_next != ` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_violation);
{{end}}

endmodule

{{end}}

{{define "functionVerilog"}}{{$block := index .Functions .FunctionIndex}}{{$blocks := .Functions}}
//This file should be called F_{{$block.Name}}.sv
//This is autogenerated code. Edit by hand at your peril!

//To check this file using EBMC, run the following command:
//$ ebmc ebmc_F_{{$block.Name}}.sv

//For each policy, we need define types for the state machines
{{range $polI, $pol := $block.Policies}}
{{if len $pol.States}}{{range $index, $state := $pol.States}}
` + "`" + `define POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$state}} {{$index}}{{end}}{{else}}POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_unknown 0{{end}}
` + "`" + `define POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_violation {{if len $pol.States}}{{len $pol.States}}{{else}}1{{end}}
{{end}}

{{if $block.Policies}}{{template "_policyIn" $block}}{{end}}
{{if $block.Policies}}{{template "_policyOut" $block}}{{end}}
{{template "_nextStateFunction" $block}}
{{template "_combinatorialVerilog" $block}}

module F_{{$block.Name}} (
	//inputs (plant to controller){{range $index, $var := $block.InputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_in,
	output wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_out,
	{{end}}
	//outputs (controller to plant){{range $index, $var := $block.OutputVars}}
	input wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ctp_in,
	output wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ctp_out,
	{{end}}

	input wire CLOCK
);

//For each policy, we need a reg for the state machine
{{range $polI, $pol := $block.Policies}}reg {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state = 0;
wire {{getVerilogWidthArray (add (len $pol.States) 1)}} {{$block.Name}}_policy_{{$pol.Name}}_state_next;
{{end}}

{{range $index, $var := $block.InputVars}}
wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ptc_out_inputmux;
{{getVerilogType $var.Type}} {{$var.Name}} {{if $var.InitialValue}} = {{$var.InitialValue}}{{end}};
{{end}}{{range $index, $var := $block.OutputVars}}
wire {{getVerilogWidthArrayForType $var.Type}} {{$var.Name}}_ctp_out_outputmux;
{{getVerilogType $var.Type}} {{$var.Name}} {{if $var.InitialValue}} = {{$var.InitialValue}}{{end}};
{{end}}{{range $polI, $pol := $block.Policies}}
{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}{{if not $var.Constant}}{{getVerilogType $var.Type}}{{else}}localparam{{end}} {{$var.Name}}{{if $var.InitialValue}} = {{$var.InitialValue}}{{end}};
{{end}}{{end}}{{end}}

	inputEditMux_{{$block.Name}} inputEditMux (
		//inputs (plant to controller){{range $index, $var := $block.InputVars}}
		.{{$var.Name}}_ptc_in_inputmux({{$var.Name}}_ptc_in),
		.{{$var.Name}}_ptc_out_inputmux({{$var.Name}}_ptc_out_inputmux),
		{{end}}

		{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
		{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}.{{$var.Name}}({{$var.Name}}),{{end}}{{end}}{{end}}

		{{range $polI, $pol := $block.Policies}}{{if $polI}},
		{{end}}.{{$block.Name}}_policy_{{$pol.Name}}_state({{$block.Name}}_policy_{{$pol.Name}}_state){{end}}
	);

	outputEditMux_{{$block.Name}} outputEditMux (
		//inputs (plant to controller){{range $index, $var := $block.InputVars}}
		.{{$var.Name}}_ptc_in_outputmux({{$var.Name}}_ptc_out_inputmux),
		{{end}}

		//outputs (controller to plant){{range $index, $var := $block.OutputVars}}
		.{{$var.Name}}_ctp_in_outputmux({{$var.Name}}_ctp_in),
		.{{$var.Name}}_ctp_out_outputmux({{$var.Name}}_ctp_out_outputmux),
		{{end}}

		{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
		{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}.{{$var.Name}}_in({{$var.Name}}),{{end}}{{end}}{{end}}

		{{range $polI, $pol := $block.Policies}}{{if $polI}},
		{{end}}.{{$block.Name}}_policy_{{$pol.Name}}_state({{$block.Name}}_policy_{{$pol.Name}}_state){{end}}
	);

	nextStateFunction_{{$block.Name}} nextStateFunction (
		//inputs (plant to controller){{range $index, $var := $block.InputVars}}
		.{{$var.Name}}({{$var.Name}}_ptc_out_inputmux),
		{{end}}
		//outputs (controller to plant){{range $index, $var := $block.OutputVars}}
		.{{$var.Name}}({{$var.Name}}_ctp_out_outputmux),
		{{end}}
		
		{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
		{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}.{{$var.Name}}_in({{$var.Name}}),
		{{end}}{{end}}{{end}}

		//state variables
		{{range $polI, $pol := $block.Policies}}{{if $polI}},
		{{end}}.{{$block.Name}}_policy_{{$pol.Name}}_state_in({{$block.Name}}_policy_{{$pol.Name}}_state),
		.{{$block.Name}}_policy_{{$pol.Name}}_state_next({{$block.Name}}_policy_{{$pol.Name}}_state_next){{end}}
	);

	always@(posedge CLOCK) begin
		//capture synchronous inputs
	{{range $index, $var := $block.InputVars}}
		{{$var.Name}} = {{$var.Name}}_ptc_out_inputmux;
	{{end}}
	{{range $index, $var := $block.OutputVars}}
		{{$var.Name}} = {{$var.Name}}_ctp_out_outputmux;
	{{end}}

		//embc inputs to introduce nondeterminism
		{{range $polI, $pol := $block.Policies}}//{{$block.Name}}_policy_{{$pol.Name}}_state = {{$block.Name}}_policy_{{$pol.Name}}_state_embc_in % {{len $pol.States}};
		{{end}}

		{{range $polI, $pol := $block.Policies}}//internal vars
		{{range $vari, $var := $pol.InternalVars}}//{{if not $var.Constant}}{{$var.Name}} = {{$var.Name}}_embc_in;
		{{end}}{{end}}{{end}}
		
		//advance state and timers
		{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}
		{{range $varI, $var := $pfbEnf.OutputPolicy.GetDTimers}}
		{{$var.Name}} = {{$var.Name}} + 1;{{end}}

		{{$block.Name}}_policy_{{$pol.Name}}_state = {{$block.Name}}_policy_{{$pol.Name}}_state_next;
		{{end}}{{end}}
	end

	//For each policy, ensure correctness (systemverilog only) and liveness
	{{range $polI, $pol := $block.Policies}}assert property ({{$block.Name}}_policy_{{$pol.Name}}_state != ` + "`" + `POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_violation);
	//assert property (transTaken_{{$block.Name}}_policy_{{$pol.Name}} == 1);
	{{end}}
	
	//emit outputs
{{range $index, $var := $block.InputVars}}
	assign {{$var.Name}}_ptc_out = {{$var.Name}};
{{end}}
{{range $index, $var := $block.OutputVars}}
	assign {{$var.Name}}_ctp_out = {{$var.Name}};
{{end}}

endmodule{{end}}`

var verilogTemplateFuncMap = template.FuncMap{
	"getVerilogECCTransitionCondition": getVerilogECCTransitionCondition,
	"getVerilogType":                   getVerilogType,
	"getPolicyEnfInfo":                 getPolicyEnfInfo,
	"getVerilogWidthArray":             getVerilogWidthArray,
	"getVerilogWidthArrayForType":      getVerilogWidthArrayForType,
	"add1IfClock":                      add1IfClock,

	"compileExpression": stconverter.VerilogCompileExpression,

	"add": add,
}

var verilogTemplates = template.Must(template.New("").Funcs(verilogTemplateFuncMap).Parse(rteVerilogTemplate))
