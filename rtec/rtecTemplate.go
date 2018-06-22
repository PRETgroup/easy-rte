package rtec

import "text/template"

const rtecTemplate = `{{define "_policyIn"}}{{$block := .}}
//input policies
{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}
{{if not $pfbEnf}}//{{$pol.Name}} is broken!
{{else}}{{/* this is where the policy comes in */}}//INPUT POLICY {{$pol.Name}} BEGIN 
	{{range $tri, $tr := $pfbEnf.InputPolicy.GetViolationTransitions}}{{/*
	*/}}{{if $tri}}else {{end}}if((me->_policy_{{$pol.Name}}_state == POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$tr.Source}}) && 
		({{$cond := getCECCTransitionCondition $block $tr.Condition}}{{$cond.IfCond}})) {
		//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
		//select a transition to solve the problem
		{{$solution := $pfbEnf.SolveViolationTransition $tr true}}
		{{if $solution.Comment}}//{{$solution.Comment}}{{end}}
		{{if $solution.Expression}}{{$sol := getCECCTransitionCondition $block $solution.Expression}}{{$sol.IfCond}};{{end}}
	} {{end}}
{{end}}
//INPUT POLICY {{$pol.Name}} END
{{end}}
{{end}}

{{define "_policyOut"}}{{$block := .}}
//output policies
{{range $polI, $pol := $block.Policies}}{{$pfbEnf := getPolicyEnfInfo $block $polI}}
{{if not $pfbEnf}}//{{$pol.Name}} is broken!
{{else}}{{/* this is where the policy comes in */}}//OUTPUT POLICY {{$pol.Name}} BEGIN 
	{{range $tri, $tr := $pfbEnf.OutputPolicy.GetViolationTransitions}}{{/*
	*/}}{{if $tri}}else {{end}}if((me->_policy_{{$pol.Name}}_state == POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$tr.Source}}) && 
		({{$cond := getCECCTransitionCondition $block $tr.Condition}}{{$cond.IfCond}})) {
		//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
		//select a transition to solve the problem
		{{$solution := $pfbEnf.SolveViolationTransition $tr false}}
		{{if $solution.Comment}}//{{$solution.Comment}}{{end}}
		{{if $solution.Expression}}{{$sol := getCECCTransitionCondition $block $solution.Expression}}{{$sol.IfCond}};{{end}}
	} {{end}}

	//advance timers
	{{range $varI, $var := $pfbEnf.OutputPolicy.GetDTimers}}
	me->{{$var.Name}}++;{{end}}

	//select transition to advance state
	{{range $tri, $tr := $pfbEnf.OutputPolicy.GetNonViolationTransitions}}{{/*
	*/}}{{if $tri}}else {{end}}if((me->_policy_{{$pol.Name}}_state == POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$tr.Source}}) && 
		({{$cond := getCECCTransitionCondition $block $tr.Condition}}{{$cond.IfCond}})) {
		//transition {{$tr.Source}} -> {{$tr.Destination}} on {{$tr.Condition}}
		me->_policy_{{$pol.Name}}_state = POLICY_STATE_{{$block.Name}}_{{$pol.Name}}_{{$tr.Destination}};
	} {{end}}
{{end}}
//OUTPUT POLICY {{$pol.Name}} END
{{end}}
{{end}}

{{define "functionRun"}}{{$block := index .Functions .FunctionIndex}}{{$blocks := .Functions}}
void {{$block.Name}}_run({{$block.Name}}_t *me) {
	{{if $block.Policies}}{{template "_policyIn" $block}}{{end}}

	CALL_FUNCTION

	{{if $block.Policies}}{{template "_policyOut" $block}}{{end}}
}{{end}}`

var cTemplateFuncMap = template.FuncMap{
	"getCECCTransitionCondition": getCECCTransitionCondition,

	"getPolicyEnfInfo": getPolicyEnfInfo,
}

var cTemplates = template.Must(template.New("").Funcs(cTemplateFuncMap).Parse(rtecTemplate))
