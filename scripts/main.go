package main

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
	jsonFile := "/Users/mattprincev/Documents/Rhize/JSON Schema/b2mml-batchml/schemas/v2.0.0.generalRecipe.schema.json"
	// findPropRefs(jsonFile)
	changeReferences(jsonFile)

	// DTchangeReferences()
	// parse()
	// schemaFile := "oldSchema/v2.0.0.base.schema.json"

	// // Load the schema
	// schemaLoader := gojsonschema.NewReferenceLoader(schemaFile)

	// // Try to compile the schema
	// _, err := gojsonschema.NewSchema(schemaLoader)
	// if err != nil {
	// 	log.Fatalf("Schema compilation error: %s", err)
	// } else {
	// 	log.Println("Schema compiled successfully.")
	// }
}
