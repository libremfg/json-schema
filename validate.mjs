import Ajv from 'ajv';
import addFormats from 'ajv-formats';
import fetch from 'node-fetch';
import fs from 'fs';
import path from 'path';

// Initialize Ajv instance
const ajv = new Ajv({ allErrors: true, strict: false });
addFormats(ajv);

// Define the meta-schema URLs
const metaSchemas = [
  'https://json-schema.org/draft/2020-12/schema',
  'https://json-schema.org/draft/2020-12/meta/core',
  'https://json-schema.org/draft/2020-12/meta/applicator',
  'https://json-schema.org/draft/2020-12/meta/validation',
  'https://json-schema.org/draft/2020-12/meta/meta-data',
  'https://json-schema.org/draft/2020-12/meta/format-annotation',
  'https://json-schema.org/draft/2020-12/meta/content'
];

// Function to fetch and add schemas to AJV
async function addMetaSchemas(ajv, urls) {
  for (const url of urls) {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch schema from ${url}: ${response.statusText}`);
    }
    const schema = await response.json();
    ajv.addMetaSchema(schema, url);
  }
}

(async () => {
  try {
    // Add meta-schemas to AJV
    await addMetaSchemas(ajv, metaSchemas);

    // Define the path to your schema file
    const schemaPath = path.resolve(process.argv[2]);
    const schema = JSON.parse(await fs.promises.readFile(schemaPath, 'utf8'));

    // Compile the schema
    const validate = ajv.compile(schema);

    // Optional: Test data against the schema
    const testData = {}; // Replace with your test data
    const valid = validate(testData);

    if (!valid) {
      console.error('Schema validation errors:');
      console.error(validate.errors);
      process.exit(1);
    } else {
      console.log('Schema is valid.');
    }
  } catch (err) {
    console.error(`Error: ${err.message}`);
    process.exit(1);
  }
})();
