+++
title = 'Version 2.0.0'
date = 2024-09-20T14:52:37-05:00
+++

> A JSON Schema implementation of the ANSI/ISA-95 and ANSI/ISA-88 standards.

An JSON implementation of the ANSI/ISA-95, Enterprise-Control System Integration, family of standards (ISA-95), known internationally as IEC/ISO 62264. B2MML consists of a set of JSON schemas written using the JSON Schema's Schema language (2020) that implement the data models in the ISA-95 standard.

This JSON schema can be used for design for syncrhonous and asynchronous APIs.

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

## Quick start

Start out by importing the schema and using it in your JSON documents. 

```
"$schema": "{{< siteurl >}}schemas"
```

Here is an example using the JSON schema in a `NotifyWorkDefined` message.

```json
{
    "$schema": "{{< siteurl >}}schemas/",
    "NotifyWorkDefined": {
        "ApplicationArea": {},
        "DataArea": {
            "Notify": {},
            "WorkDefined": {

            }
        }
    }
}
```

## Updates Introduced in Version 2.0.0

Refactored the v1.0.0.base.schema.json JSON file by separating it into multiple files that use the prefix "v2.0.0", aligning the structure with the existing B2MMl XML implementation https://github.com/MESAInternational/B2MML-BatchML/tree/master/Schema.

Implemented several automated scripts to assist in refactoring. While these scripts are tailored for specific use cases, they can generally be disregarded outside of those particular scenarios.

In Version 2.0.0, a JSON Schema validator is used for verification, replacing the previous Lint-based approach. This update compiles and validates schemas through the use of scripts, specifically compileSchemas.js, located in scripts, and validate.mjs.