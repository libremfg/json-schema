package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

func parse() {
	// Command-line argument to implement
	// xsdFile := flag.String("xsd", "b2mml-batchml/B2MML-ConfirmBOD .xsd", "XSD file to parse")
	// jsonFile := flag.String("json", "b2mml-batchml/Schema/v2.0.0.base.schema.json", "JSON file to parse")
	// outputFile := flag.String("output", "testFile.json", "Output file name")
	// flag.Parse()

	xsdFile := "/Users/mattprincev/Documents/Rhize/JSON Schema/b2mml-batchml/BatchML-GeneralRecipeExtensions.xsd"
	jsonFile := "/Users/mattprincev/Documents/Rhize/JSON Schema/b2mml-batchml/oldSchema/v2.0.0.base.schema.json"
	outputFile := "v2.0.0.generalRecipeExtensions.schema.json"

	// Get the list of names from the XSD file
	names, err := getDefNames(xsdFile)
	if err != nil {
		fmt.Println("Error parsing XSD file:", err)
		return
	}

	// Check the names against keys in the JSON file
	matchingData, err := checkNamesInJSONDefs(names, jsonFile)
	if err != nil {
		fmt.Println("Error parsing JSON file:", err)
		return
	}

	// Construct the output JSON structure
	outputData := JSONStructure{
		ID:          "v2.0.0.--CHANGE--.schema.json",
		Schema:      map[string]interface{}{},
		Title:       "Rhize B2MML-JSON - --CHANGE-- Definitions",
		Description: "Copyright 2024, Libre Technologies Inc., Version v2.0.0\nAll Rights Reserved. http://www.rhize.com\nBased upon the ANSI/ISA-95.00.02-2018 Enterprise-Control System Integration Part 2: Object Model Attributes Standard and the ANSI/ISA-95.00.05-2018 Enterprise-Control System Integration Part 5: Business to Manufacturing Transactions.\nThe Business To Manufacturing Markup Language-JSON (B2MML-JSON) is used courtesy of MESA International.",
		Defs:        matchingData,
		Type:        "object",
		Properties:  map[string]interface{}{},
		AllOf:       []map[string]interface{}{{}},
		UnevalProps: false,
	}

	// Write the output JSON structure to the output file
	if err := writeJSONToFile(outputData, outputFile); err != nil {
		fmt.Println("Error writing to output file:", err)
	}
}

func getDefNames(fileToParse string) ([]string, error) {
	// Open the XSD file
	file, err := os.Open(fileToParse)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Define a regex to match the complexType and simpleType name attributes
	re := regexp.MustCompile(`<xsd:(?:complexType|simpleType)\s+name\s*=\s*"([^"]+)"`)

	// Scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var names []string

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			names = append(names, matches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return names, nil
}

func checkNamesInJSONDefs(names []string, jsonFile string) (map[string]interface{}, error) {
	// Read the JSON file
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	// Parse the JSON file
	var jsonData JSONStructure
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	// Collect matching keys and their nested values from $defs
	matchingData := make(map[string]interface{})
	for key, value := range jsonData.Defs {
		for _, name := range names {
			if key == name {
				fmt.Printf("Found matching key: %s\n", key)
				matchingData[key] = extractNestedData(value)
			}
		}
	}

	return matchingData, nil
}

// Recursively extract all nested data
func extractNestedData(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// Create a new map to hold nested data
		nestedData := make(map[string]interface{})
		for k, value := range v {
			nestedData[k] = extractNestedData(value)
		}
		return nestedData
	case []interface{}:
		// Create a new slice to hold nested data
		var sliceData []interface{}
		for _, value := range v {
			sliceData = append(sliceData, extractNestedData(value))
		}
		return sliceData
	default:
		return v
	}
}

// Write JSON structure to file
func writeJSONToFile(data JSONStructure, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding JSON to file: %v", err)
	}

	return nil
}
