package main

import (
	"flag"
	"fmt"
)

// Struct to represent the JSON structure with generic map
type JSONStructure struct {
	ID          string                 `json:"$id"`
	Schema      map[string]interface{} `json:"schema"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Defs        map[string]interface{} `json:"$defs"`
	Type        string                 `json:"type"`
	Properties  map[string]interface{} `json:"properties"`
	AllOf       interface{}            `json:"allOf,omitempty"`
	UnevalProps bool                   `json:"unevaluatedProperties"`
}

func main() {
	// Define command-line flags

	// Flag variables for baseToDt.go
	// example command: go run ./ --dtNewFile=myDataType.json --baseFile=myBaseFile.json
	dtNewFile := flag.String("dtNewFile", "", "Path to the dataType file")
	dTBaseFile := flag.String("baseFile", "", "Path to the dataType base file")

	// Flag variables for findPropRefs.go
	// example command: go run ./ --propRefNewFile=myDataType.json --propRefBaseFile=myBaseFile.json
	propRefNewFile := flag.String("propRefNewFile", "", "Path to the file that the property references will be added to")
	propRefBaseFile := flag.String("propRefBaseFile", "", "Path to the base file that contains the property references")

	// Flag variables for refFinder.go
	// example command: go run ./ --refFinderNewFile=myDataType.json --refFinderBaseFile=myBaseFile.json
	refFinderNewFile := flag.String("refFinderNewFile", "", "Path to the file that the conatins the reference names that will be searched for")
	refFinderBaseFile := flag.String("refFinderBaseFile", "", "Path to the directory containing the schema files to be referenced.")

	// Flag variables for xsdToJson.go
	// example command: go run ./ --dtNewFile=myDataType.json --baseFile=myBaseFile.json
	xsdFile := flag.String("xsdFile", "", "Path to the xsd schema file")
	baseJsonFile := flag.String("baseJsonFile", "", "Path to the base json schema file")
	outputFile := flag.String("outputFile", "", "Name of the output file that will be created")

	// Parse the command-line flags
	flag.Parse()

	// Check if DTchangeReferences should be called
	if *dtNewFile != "" && *dTBaseFile != "" {
		fmt.Println("Running baseToDt.go")
		if _, err := DTchangeReferences(*dtNewFile, *dTBaseFile); err != nil {
			return
		}
		// fmt.Println(*dtNewFile, *dTBaseFile)
	} else {
		fmt.Println("Flags for dataType file are empty")
	}

	if *propRefNewFile != "" && *propRefBaseFile != "" {
		fmt.Println("Running findPropRefs.go")
		if _, err := findPropRefs(*propRefNewFile, *propRefBaseFile); err != nil {
			return
		}
		// fmt.Println(*propRefNewFile, *propRefBaseFile)
	} else {
		fmt.Println("Flags for property reference finder are empty")
	}

	if *refFinderNewFile != "" && *refFinderBaseFile != "" {
		fmt.Println("Running refFinder.go")
		if _, err := changeReferences(*refFinderNewFile, *refFinderBaseFile); err != nil {
			return
		}
		// fmt.Println(*refFinderNewFile, *refFinderBaseFile)
	} else {
		fmt.Println("Flags for reference finding are empty")
	}

	if *xsdFile != "" && *baseJsonFile != "" && *outputFile != "" {
		fmt.Println("Running refFinder.go")
		if _, err := parse(*xsdFile, *baseJsonFile, *outputFile); err != nil {
			return
		}
		// fmt.Println(*xsdFile, *baseJsonFile, *outputFile)
	} else {
		fmt.Println("Flags for xsd to JSON conversion are empty")
	}
}
