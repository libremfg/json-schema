package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// This script is designed to recursively search for properties referenced within other schemas.
// It traverses the schemas deeply to ensure that all references are accurately identified.
func changeReferences(jsonFile string, schemaDirectory string) (int, error) {

	// Extract $ref values
	refValues, err := extractRefValues(jsonFile)
	if err != nil {
		return fmt.Println("Error extracting $ref values:", err)
	}

	fmt.Println("Extracted $ref values:", refValues)
	fmt.Println("Length of extracted $ref values:", len(refValues))

	// Find definition files based on the extracted $ref values
	definitionRefs, err := findRefFiles(refValues, schemaDirectory)
	if err != nil {
		return fmt.Println("Error finding definition references:", err)
	}

	// Print the results
	fmt.Println("Length of definition references:", len(definitionRefs))
	// fmt.Println("List of definition files:", definitionFiles)
	for name, path := range definitionRefs {
		fmt.Printf("Definition '%s'\n", name)
		fmt.Printf("found in file: %s\n", path[0])
	}

	err = replacesRefs(jsonFile, definitionRefs)
	if err != nil {
		return fmt.Println("Error replacing references:", err)
	}

	return 0, nil
}

// Finds the name of the reference values needed to be searched for
func extractRefValues(jsonFile string) ([]string, error) {

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var jsonData JSONStructure
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	var refValues []string
	for _, value := range jsonData.Defs {
		refValues = append(refValues, extractRefsRead(value)...)
	}

	return refValues, nil
}

// Recursively does the extracting of references
func extractRefsRead(data interface{}) []string {

	var refs []string
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if key == "$ref" {
				if refValue, ok := value.(string); ok {
					name := extractName(refValue)
					if refValue != "" {
						refs = append(refs, name)
					}
				}
			} else {
				refs = append(refs, extractRefsRead(value)...)
			}
		}
	case []interface{}:
		for _, item := range v {
			refs = append(refs, extractRefsRead(item)...)
		}
	}
	return refs
}

// Extracts a specific part of the string that comes after the /
func extractName(ref string) string {

	lastSlashIndex := strings.LastIndex(ref, "/")
	if lastSlashIndex == -1 {
		// Return the original string if '/' is not found because we only want the string that comes after it
		return ref
	}
	return ref[lastSlashIndex+1:]
}

// Finds definition files based on $ref names in a directory
func findRefFiles(definitions []string, directory string) (map[string][]string, error) {

	results := make(map[string][]string)
	// Walk the directory to go through each schema file
	err := filepath.WalkDir(directory, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// checks to see if the specified file is a json file
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			fileData, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var jsonData JSONStructure
			if err := json.Unmarshal(fileData, &jsonData); err != nil {
				return err
			}

			// Directly use jsonData.Defs which is of type map[string]interface{}
			// If the reference exists within a definition in the directory, then it will add the path as a reference
			for _, name := range definitions {
				if _, exists := jsonData.Defs[name]; exists {
					if paths, found := results[name]; found {
						results[name] = append(paths, path)
					} else {
						results[name] = []string{path}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Unmarshals the data from the map and calls the helper functions
func replacesRefs(jsonFile string, definitionRefs map[string][]string) error {

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	var jsonData JSONStructure
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	// Replace $ref values
	if err := replaceRefsInData(&jsonData.Defs, definitionRefs); err != nil {
		return fmt.Errorf("error updating $ref values: %v", err)
	}

	// Marshal updated data
	updatedData, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling updated data: %v", err)
	}

	// Write updated data back to the file
	if err := os.WriteFile(jsonFile, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing updated file %s: %v", jsonFile, err)
	}

	fmt.Printf("Updated file: %s\n", jsonFile)
	return nil
}

// Does the actual replacing of references found within directories into the json file in
func replaceRefsInData(data *map[string]interface{}, replacements map[string][]string) error {

	for key, value := range *data {
		// finds the key value "$ref" in the schema file
		if key == "$ref" {
			if refValue, ok := value.(string); ok {
				name := extractName(refValue)
				if replacementPaths, found := replacements[name]; found && len(replacementPaths) > 0 {
					// Replace the $ref value with the path
					correctedPath, err := extractPath(replacementPaths[0], name)
					if err != nil {
						fmt.Println("Error: ", err)
						fmt.Printf("Error extracting path: %s\n", replacementPaths[0])
						return err
					}
					// Use the first matching path
					(*data)[key] = correctedPath
				}
			}
		} else {
			switch v := value.(type) {
			case map[string]interface{}:
				if err := replaceRefsInData(&v, replacements); err != nil {
					return err
				}
				(*data)[key] = v
			case []interface{}:
				for i, item := range v {
					// Replace refs within items of the array
					if err := replaceRefsInData(&map[string]interface{}{"item": item}, replacements); err != nil {
						return err
					}
					v[i] = item
				}
			}
		}
	}

	return nil
}

// Extracts a specific part of the schema path that fits the format of the references
func extractPath(path string, refValue string) (string, error) {
	// Check if the path is empty
	if path == "" {
		return "", errors.New("path is empty")
	}

	// Define the substring to search for
	const searchStr = "Schema/"

	// Find the index of the search string
	index := strings.Index(path, searchStr)
	if index == -1 {
		// If the search string is not found, return an error
		return "", errors.New("substring 'Schema/' not found in path")
	}

	// Calculate the starting index of the desired substring
	startIndex := index + len(searchStr)

	// Extract the substring from the start index to the end of the path
	extractedPath := path[startIndex:]

	// Extract the filename from the path
	filename := filepath.Base(extractedPath)

	return filename + "#/$defs/" + refValue, nil
}
