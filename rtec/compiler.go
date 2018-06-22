package rtec

import (
	"bytes"
	"encoding/xml"
	"errors"
	"text/template"

	"github.com/PRETgroup/easy-rte/rtedef"
)

//Converter is the struct we use to store all functions for conversion (and what we operate from)
type Converter struct {
	Funcs []rtedef.EnforcedFunction

	templates *template.Template
}

//New returns a new instance of a Converter based on the provided language
func New() (*Converter, error) {
	return &Converter{Funcs: make([]rtedef.EnforcedFunction, 0), templates: cTemplates}, nil

}

//AddFunction should be called for each Function in the project
func (c *Converter) AddFunction(functionbytes []byte) error {
	FB := rtedef.EnforcedFunction{}
	if err := xml.Unmarshal(functionbytes, &FB); err != nil {
		return errors.New("Couldn't unmarshal EnforcedFunction xml: " + err.Error())
	}

	c.Funcs = append(c.Funcs, FB)

	return nil
}

//OutputFile is used when returning the converted data from the iec61499
type OutputFile struct {
	Name      string
	Extension string
	Contents  []byte
}

//TemplateData is the structure used to hold data being passed into the templating engine
type TemplateData struct {
	FunctionIndex int
	Functions     []rtedef.EnforcedFunction
}

//ConvertAll converts iec61499 xml (stored as []FB) into vhdl []byte for each block (becomes []VHDLOutput struct)
//Returns nil error on success
func (c *Converter) ConvertAll() ([]OutputFile, error) {
	finishedConversions := make([]OutputFile, 0, len(c.Funcs))

	//convert all functions
	for i := 0; i < len(c.Funcs); i++ {

		output := &bytes.Buffer{}
		templateName := "functionRun"

		if err := c.templates.ExecuteTemplate(output, templateName, TemplateData{FunctionIndex: i, Functions: c.Funcs}); err != nil {
			return nil, errors.New("Couldn't format template (fb) of" + c.Funcs[i].Name + ": " + err.Error())
		}

		finishedConversions = append(finishedConversions, OutputFile{Name: "F_" + c.Funcs[i].Name, Extension: "c", Contents: output.Bytes()})
	}

	return finishedConversions, nil
}
