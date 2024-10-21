const Ajv2020 = require("ajv/dist/2020");
const fs = require("fs");
const path = require("path");

// Initialize AJV instance with support for draft-2020-12
const ajv = new Ajv2020({
    strict: true,
    allErrors: true,
    allowUnionTypes: true,
    loadSchema: (uri) => {
        return new Promise((resolve, reject) => {
            fs.readFile(uri, "utf8", (err, data) => {
                if (err) {
                    reject(err);
                } else {
                    resolve(JSON.parse(data));
                }
            });
        });
    }
});

const schemasDir = path.join(__dirname, "schemas");

async function loadAndCompileSchemas() {
    try {
        const files = fs.readdirSync(schemasDir);
        const schemas = [];

        // Load and parse each schema file
        for (const file of files) {
            const schemaPath = path.join(schemasDir, file);
            const schema = JSON.parse(fs.readFileSync(schemaPath, "utf8"));

            // Add the schema to the array with its filename
            schemas.push({ schema, name: file });
        }

        console.log("\nAdding schemas:\n")
        
        // Add all schemas to AJV at once
        try {
            ajv.addSchema(schemas.map(s => s.schema));

            // Verify if each schema was correctly added
            schemas.forEach(({ schema, name }) => {
                const addedSchema = ajv.getSchema(schema.$id);
                if (addedSchema) {
                    console.log(`Schema ${name} (${schema.$id}) added successfully.`);
                } else {
                    console.error(`Schema ${name} (${schema.$id}) was not added.`);
                }
            });
        } catch (addError) {
            console.error(`Error adding schemas:`, addError.message);
        }

        console.log("\nCompiling schemas:\n")

        // Compile schemas
        for (const { schema, name } of schemas) {
            try {
                ajv.compile(schema);
                console.log(`Schema ${name} (${schema.$id}) compiled successfully.`);
            } catch (compilationError) {
                console.error(`Schema ${name} (${schema.$id}) compilation error:`, compilationError.errors || compilationError.message);
            }
        }

    } catch (err) {
        console.error("Schema loading or processing error:", err.message);
    }
}

// Run the function
loadAndCompileSchemas();
