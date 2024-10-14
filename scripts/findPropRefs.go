package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func findPropRefs(schemaFile string) {
	// Command-line argument to implement
	// xsdFile := flag.String("xsd", "b2mml-batchml/B2MML-ConfirmBOD .xsd", "XSD file to parse")
	// jsonFile := flag.String("json", "b2mml-batchml/Schema/v2.0.0.base.schema.json", "JSON file to parse")
	// outputFile := flag.String("output", "testFile.json", "Output file name")
	// flag.Parse()

	// schemaFile := "/Users/mattprincev/Documents/Rhize/JSON Schema/b2mml-batchml/errorTest.json"
	baseFile := "/Users/mattprincev/Documents/Rhize/JSON Schema/b2mml-batchml/schemas/v2.0.0.batchInformation.schema.json"

	// Get the list of names from the XSD file
	propNames, err := getPropNames(schemaFile)
	if err != nil {
		fmt.Println("Error parsing schema file:", err)
		return
	}

	baseRefs, err := getBaseRefs(propNames, baseFile)
	if err != nil {
		fmt.Println("Error getting base references:", err)
		return
	}

	err = addDefs(baseRefs, schemaFile)
	if err != nil {
		fmt.Println("Error adding definitions:", err)
		return
	}
}

func getPropNames(jsonFile string) (map[string]interface{}, error) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var jsonData JSONStructure
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	var refProps = make(map[string]interface{}) // Initialize the map

	for key, value := range jsonData.Properties {
		refProps[key] = value
	}

	return refProps, nil
}

func getBaseRefs(propNames map[string]interface{}, baseSchema string) (map[string]interface{}, error) {
	data, err := os.ReadFile(baseSchema)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var jsonData JSONStructure
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	// Iterate over the list of keys and check if they exist in the map
	for key, _ := range propNames {
		if value, exists := jsonData.Defs[key]; exists {
			// Copy the matching key-value pair to the new map
			fmt.Printf("%s: %v\n", key, value)
			propNames[key] = value
		}
	}

	return propNames, nil
}

func addDefs(definitions map[string]interface{}, schemaFile string) error {
	data, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	var jsonData JSONStructure
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	// var refProps = make(map[string]interface{})
	for key, value := range definitions {
		jsonData.Defs[key] = value
	}

	updatedContent, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON data: %v", err)
	}

	// Step 5: Write the updated JSON content back to the file
	if err := os.WriteFile(schemaFile, updatedContent, 0644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
