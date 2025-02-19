# Rhize ISA-95/ISA-88 JSON Schema

> A JSON Schema implementation of the ANSI/ISA-95 and ANSI/ISA-88 standards.

An JSON implementation of the ANSI/ISA-95, Enterprise-Control System Integration, family of standards (ISA-95), known internationally as IEC/ISO 62264. B2MML consists of a set of JSON schemas written using the JSON Schema's Schema language (2020) that implement the data models in the ISA-95 standard.

This JSON schema is intented be used for design for syncrhonous and asynchronous APIs.

Based on the works of https://github.com/MESAInternational/B2MML-BatchML.

## Abbreviations

| Acronynm    | Description                                                   |
|-------------|---------------------------------------------------------------|
| **ANSI**    | American National Standards Institute                         |
| **API**     | Application Programming Interface                             |
| **B2MML**   | Business to (2) Manufacturing Markup Language                 |
| **BatchML** | Batch Markup Language                                         |
| **BOD**     | Business Object Document                                      | 
| **ISA**     | International Society of Automation                           |
| **IEC**     | International Electrotechnical Commission                     | 
| **JSON**    | JavaScipt Object Notation                                     |
| **XML**     | eXensible Markup Language                                     |
| **XSD**     | eXtensible markup language Schema Definition                  |

## Usage

The JSON schema can be used to validate JSON file content and can be referenced multiple ways. 

### Using https://json.libremfg.ai

For convenience, [Libre Technologies Inc.](http://www.libremfg.com/) hosts the schema at https://json.libremfg.ai/.

```
{
    "$schema": "https://json.libremfg.ai/schemas/v2.0.1.equipment.schema.json",
    "Equipment": {
        "ID": "Hello World"
    }
}
```

### Using Locally

You can clone this repository and use file references to the schema.

1. Clone this respository to a directory for example `git clone https://github.com/libremfg/json-schema ~/json-schema`
2. Create a new file `touch ~/equipment.json`. Edit the file with:

```json
{
    "$schema": "~/json-schema/schemas/v2.0.1.equipment.schema.json",
    "Equipment": {
        "ID": "Hello World"
    }
}
```

### Using with external tools

You can use the JSON schemas for validating ISA-95 based data exchanges. Here's an example of how to use it with a validation tool like [Ajv](https://ajv.js.org/):

```js
const Ajv = require("ajv");
const schema = require("./path/to/json-schema/v2.0.1.equipment.schema.json");

const ajv = new Ajv();
const validate = ajv.compile(schema);

const valid = validate(data);
if (!valid) {
  console.log(validate.errors);
}
```

## Contributing 

Contributions are welcome! Here's how you can get invovled:

1. Fork the repository.
2. Create a new branch for your feature or bugfix: `git checkout -b feature-name`. Following [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) is appreciated but not required.
3. Commit your changes: `git commit -m "Add a new feature"`
4. Push your branch: `git push origin feature-name`
5. Create a Pull Request.

## Issues

If you encounter any bugs or want to suggest enhancements, please [open an issue](https://github.com/libremfg/json-schema/issues/new) and provide detailed steps to reproduce the problem or your enhancement proposal.

## License

Rhize ISA-95/ISA-88 JSON Schema is distributed under [AGPL-3.0-only](LICENSE).

## Acknowledgements

- Based on the works of [MESA International](https://github.com/MESAInternational/B2MML-BatchML).
- Contributions from the ISA-95 and ISA-88 community are appreciated.

Copyright 2025, Libre Technologies Inc.
