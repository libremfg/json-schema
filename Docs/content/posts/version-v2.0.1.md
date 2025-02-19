+++
title = 'Version 2.0.1'
date = 2025-02-19T12:50:23-05:00
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

Start out by importing the schema and using it in your JSON documents. As an example:

```
"$schema": "{{< siteurl >}}schemas/v2.0.1.equipment.schema.json"
```

Here is an example using the JSON schema in a `ProcessEquipment` message.

```json
{
    "$schema": "https://json.libremfg.ai/schemas/v2.0.1.equipment.schema.json",
    "ProcessEquipment": {
        "@releaseID": "1",
        "ApplicationArea": {
            "CreationDateTime": "2021-01-01T00:00:00Z",
            "Sender": {
                "LogicalID": "Rhize Manufacutring Data Hub"
            }
        },
        "DataArea": {
            "Process": {},
            "Equipment": [
                {
                    "ID": "Acme Inc.",
                    "Description": [
                        "Acme Inc. Manufacturing Line 1"
                    ],
                    "EquipmentLevel": "Enterprise"
                }
            ]
        }
    }
}
```

## Updates in 2.0.1

- Add concrete JSON schema definitions for `WorkflowSpecificationConnectionType`, `WorkflowSpecificationNodeType`, and `ResourceNetworkConnectionType`

Previously, these schema definitions mistakenly referenced the `xxxTypeType` entity instead of defining their own structures. This MR introduces the correct definitions for these types and updates references accordingly.

## Object Schema File Table

Use the search functionality of your browser to find an object of interet.

{{< table/propertiesV201 >}}
