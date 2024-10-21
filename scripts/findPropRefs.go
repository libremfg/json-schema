package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// After converting an XSD file to a JSON file, this script adds all references found
// within the specified base file to the designated schema file.
func findPropRefs(schemaFile string, baseFile string) (int, error) {

	// Gets the list of names from the XSD file
	propNames, err := getPropNames(schemaFile)
	if err != nil {
		return fmt.Println("Error parsing schema file:", err)
	}

	// Gets the name of the references needed to search for in the basefile
	baseRefs, err := getBaseRefs(propNames, baseFile)
	if err != nil {
		return fmt.Println("Error getting base references:", err)
	}

	// Adds each instance of a reference to the definition
	err = addDefs(baseRefs, schemaFile)
	if err != nil {
		return fmt.Println("Error adding definitions:", err)
	}

	return 0, nil
}

// Reads the property names that need to be searched for in the base file
func getPropNames(jsonFile string) (map[string]interface{}, error) {

	// Read the inputted schema file
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	// unmarshal the read data into the JSON struct
	var jsonData JSONStructure
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	// Initialize the map
	var refProps = make(map[string]interface{})

	for key, value := range jsonData.Properties {
		refProps[key] = value
	}

	return refProps, nil
}

// Gets the property names from the base schema and adds them to the json structure
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

// Adds the definitions found from the base file into the inputted schema file
func addDefs(definitions map[string]interface{}, schemaFile string) error {

	data, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	var jsonData JSONStructure
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	for key, value := range definitions {
		jsonData.Defs[key] = value
	}

	updatedContent, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON data: %v", err)
	}

	// Write the updated JSON content back to the file
	if err := os.WriteFile(schemaFile, updatedContent, 0644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
