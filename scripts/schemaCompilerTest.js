const Ajv2020 = require("ajv/dist/2020"); // Ensure you are using AJV with support for draft 2020-12
const fs = require("fs");

// Initialize AJV instance with support for draft-2020-12
const ajv = new Ajv2020({
    strict: true,
    allErrors: true,
    // Allow schema references to be relative paths
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

// Load the main schema file and reference schema
const mainSchemaPath = "v1.0.0.test.schema.json";
const refSchemaPath = "v1.0.0.testExtensions.schema.json";

async function compileSchema() {
    try {
        // Load the schemas
        const mainSchema = JSON.parse(fs.readFileSync(mainSchemaPath, "utf8"));
        const refSchema = JSON.parse(fs.readFileSync(refSchemaPath, "utf8"));

        // Add the reference schema to AJV
        ajv.addSchema(refSchema, refSchema.$id);

        // Compile the main schema
        const validate = ajv.compile(mainSchema);
        console.log("Schema compiled successfully.");

        // Optionally, validate some data against the compiled schema:
        // const valid = validate(someData);
        // if (!valid) console.log(validate.errors);
    } catch (err) {
        console.error("Schema is invalid:", err.errors || err.message);
    }
}

// Run the function
compileSchema();
