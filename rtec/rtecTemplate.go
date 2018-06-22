package rtec

import "text/template"

const rtecTemplate = `{{define "_policyIn"}}{{$block := .}}
	//input policies
	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}
	{{if not $pfbEnf}}//{{$pol.Name}} is broken!
	{{else}}{{/* this is where the policy comes in */}}//INPUT POLICY {{$pol.Name}} BEGIN 
		switch(me->_policy_{{$pol.Name}}_state) {
			{{range $sti, $st := $pol.States}}case POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$st.Name}}:
				{{range $tri, $tr := $pfbEnf.InputPolicy.GetViolationTransitions}}{{if eq $tr.Source $st.Name}}{{/*
				*/}}
				if({{$cond := getCECCTransitionCondition $block $tr.Condition}}{{$cond.IfCond}}) {
					//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
					//select a transition to solve the problem
					{{$solution := $pfbEnf.SolveViolationTransition $tr true}}
					{{if $solution.Comment}}//{{$solution.Comment}}{{end}}
					{{if $solution.Expression}}{{$sol := getCECCTransitionCondition $block $solution.Expression}}{{$sol.IfCond}};{{end}}
				} {{end}}{{end}}
				
				break;

			{{end}}
		}
	{{end}}
	//INPUT POLICY {{$pol.Name}} END
	{{end}}
{{end}}

{{define "_policyOut"}}{{$block := .}}
	//output policies
	{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}
	{{if not $pfbEnf}}//{{$pol.Name}} is broken!
	{{else}}{{/* this is where the policy comes in */}}//OUTPUT POLICY {{$pol.Name}} BEGIN 
		switch(me->_policy_{{$pol.Name}}_state) {
			{{range $sti, $st := $pol.States}}case POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$st.Name}}:
				{{range $tri, $tr := $pfbEnf.OutputPolicy.GetViolationTransitions}}{{if eq $tr.Source $st.Name}}{{/*
				*/}}
				if({{$cond := getCECCTransitionCondition $block $tr.Condition}}{{$cond.IfCond}}) {
					//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
					//select a transition to solve the problem
					{{$solution := $pfbEnf.SolveViolationTransition $tr false}}
					{{if $solution.Comment}}//{{$solution.Comment}}{{end}}
					{{if $solution.Expression}}{{$sol := getCECCTransitionCondition $block $solution.Expression}}{{$sol.IfCond}};{{end}}
				} {{end}}{{end}}

				break;

			{{end}}
		}

		//advance timers
		{{range $varI, $var := $pfbEnf.OutputPolicy.GetDTimers}}
		me->{{$var.Name}}++;{{end}}

		//select transition to advance state
		switch(me->_policy_{{$pol.Name}}_state) {
			{{range $sti, $st := $pol.States}}case POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$st.Name}}:
				{{range $tri, $tr := $pfbEnf.OutputPolicy.GetNonViolationTransitions}}{{if eq $tr.Source $st.Name}}{{/*
				*/}}
				if({{$cond := getCECCTransitionCondition $block $tr.Condition}}{{$cond.IfCond}}) {
					//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
					me->_policy_{{$pol.Name}}_state = POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$tr.Destination}};
				} {{end}}{{end}}
				
				break;

			{{end}}
		}
	{{end}}
	//OUTPUT POLICY {{$pol.Name}} END
	{{end}}
{{end}}

{{define "functionH"}}{{$block := index .Functions .FunctionIndex}}{{$blocks := .Functions}}
//This file should be called F_{{$block.Name}}.h
//This is autogenerated code. Edit by hand at your peril!

#include <stdint.h>
#include <stdbool.h>

//For each policy, we need an enum type for the state machine
{{range $polI, $pol := $block.Policies}}
enum {{$block.Name}}_policy_{{$pol.Name}}_states { {{if len $pol.States}}{{range $index, $state := $pol.States}}
	POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$state}}{{if not $index}}, {{end}}{{end}}{{else}}POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_unknown{{end}} 
};
{{end}}

typedef struct {
	//inputs:
	{{range $index, $var := $block.Inputs}}{{$var.Type}} {{$var.Name}}{{if $var.ArraySize}}[{{$var.ArraySize}}]{{end}};
	{{end}}
	//outputs:
	{{range $index, $var := $block.Outputs}}{{$var.Type}} {{$var.Name}}{{if $var.ArraySize}}[{{$var.ArraySize}}]{{end}};
	{{end}}
	//policy state vars:
	{{range $polI, $pol := $block.Policies}}enum {{$block.Name}}_policy_{{$pol.Name}}_states _policy_{{$pol.Name}}_state;
	{{$pfbEnf := getPolicyEnfInfo $block $polI}}{{if not $pfbEnf}}//Policy is broken!{{else}}//internal vars
	{{range $vari, $var := $pfbEnf.OutputPolicy.InternalVars}}{{$var.Type}} {{$var.Name}}{{if $var.ArraySize}}[{{$var.ArraySize}}]{{end}};
	{{end}}{{end}}
	{{end}}
} io_{{$block.Name}}_t;
{{end}}

{{define "functionC"}}{{$block := index .Functions .FunctionIndex}}{{$blocks := .Functions}}
//This file should be called F_{{$block.Name}}.c
//This is autogenerated code. Edit by hand at your peril!
#include "F_{{$block.Name}}.h"

void {{$block.Name}}_run(io_{{$block.Name}}_t *me) {
	{{if $block.Policies}}{{template "_policyIn" $block}}{{end}}

	{{$block.Name}}(me);

	{{if $block.Policies}}{{template "_policyOut" $block}}{{end}}
}{{end}}`

var cTemplateFuncMap = template.FuncMap{
	"getCECCTransitionCondition": getCECCTransitionCondition,

	"getPolicyEnfInfo": getPolicyEnfInfo,
}

var cTemplates = template.Must(template.New("").Funcs(cTemplateFuncMap).Parse(rtecTemplate))