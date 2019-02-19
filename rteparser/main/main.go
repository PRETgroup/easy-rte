package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/PRETgroup/easy-rte/rteparser"
)

var (
	inFileName  = flag.String("i", "", "Specifies the name of the source file (.erte) file to be compiled.")
	outFileName = flag.String("o", "out.xml", "Specifies the name of the output file (.erte.xml) files.")
)

var (
	xmlHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
)

var (
	inputExtension  = "erte"
	outputExtension = "xml"
)

func main() {
	flag.Parse()

	if *inFileName == "" {
		fmt.Println("You need to specify a file or directory name to compile! Check out -help for options")
		return
	}

	// fileInfo, err := os.Stat(*inFileName)
	// if err != nil {
	// 	fmt.Println("Error reading file statistics:", err.Error())
	// 	return
	// }

	// var fileNames []string

	// if fileInfo.IsDir() {
	// 	//Running in Dir mode
	// 	files, err := ioutil.ReadDir(*inFileName)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	for _, file := range files {
	// 		//only read the .tfb files
	// 		name := file.Name()
	// 		nameComponents := strings.Split(name, ".")
	// 		extension := nameComponents[len(nameComponents)-1]
	// 		if extension == inputExtension {
	// 			fileNames = append(fileNames, name)
	// 		}
	// 	}
	// } else {
	// 	//Running in Single mode
	// 	fileNames = append(fileNames, *inFileName)
	// }

	// var funcs []rtedef.EnforcedFunction

	// for _, name := range fileNames {
	// 	sourceFile, err := ioutil.ReadFile(fmt.Sprintf("%s%c%s", *inFileName, os.PathSeparator, name))
	// 	if err != nil {
	// 		fmt.Printf("Error reading file '%s' for conversion: %s\n", name, err.Error())
	// 		return
	// 	}

	// 	mfbs, parseErr := rteparser.ParseString(name, string(sourceFile))
	// 	if parseErr != nil {
	// 		fmt.Printf("Error during parsing file '%s': %s\n", name, parseErr.Error())
	// 		return
	// 	}

	// 	funcs = append(funcs, mfbs...)

	// }

	sourceFile, err := ioutil.ReadFile(*inFileName)
	if err != nil {
		fmt.Printf("Error reading file '%s' for conversion: %s\n", *inFileName, err.Error())
		return
	}
	mfbs, parseErr := rteparser.ParseString(*inFileName, string(sourceFile))
	if parseErr != nil {
		fmt.Printf("Error during parsing file '%s': %s\n", *inFileName, parseErr.Error())
		return
	}
	for _, fun := range mfbs {
		for i := 0; i < len(fun.Policies); i++ {
			fun.Policies[i].SortTransitionsViolationsToEnd()
		}

		// name := fun.Name
		// extn := outputExtension
		//TODO: work out what extension to use based on the fb.XMLname field
		bytes, err := xml.MarshalIndent(fun, "", "\t")
		if err != nil {
			fmt.Println("Error during marshal:", err.Error())
			return
		}
		//output := append(xmlHeader, fbTypeHeader...)
		output := append(xmlHeader, bytes...)

		fmt.Printf("Writing to %s\n", *outFileName)
		err = ioutil.WriteFile(*outFileName, output, 0644)
		if err != nil {
			fmt.Println("Error during file write:", err.Error())
			return
		}
	}

}
