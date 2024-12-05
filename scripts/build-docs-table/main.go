package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PropertyInfo struct {
	FileName     string
	PropertyName string
}

func main() {
	dir := "../../schemas"

	// Array to store file and property information
	var propertyInfos []PropertyInfo

	// Walk through the directory to find JSON schema files
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Process only JSON files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") && !strings.HasPrefix(info.Name(), "v1") {
			// Parse the JSON schema file
			properties, err := parseJSONSchema(path)
			if err != nil {
				fmt.Printf("Error parsing file %s: %v\n", path, err)
				return nil
			}

			// Append properties to the array
			for _, prop := range properties {
				propertyInfos = append(propertyInfos, PropertyInfo{
					FileName:     filepath.Base(path),
					PropertyName: prop,
				})
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
		return
	}

	propertyInfos = filter(propertyInfos)

	// Print Markdown table
	markdownTableBytes := printMarkdownTable(propertyInfos)
	fmt.Println(markdownTableBytes.String())
	os.WriteFile("properties.md", markdownTableBytes.Bytes(), 0644)

	// Print HTML table
	htmlTableBytes := printHTMLTable(propertyInfos)
	fmt.Println(htmlTableBytes.String())
	os.WriteFile("../../Docs/layouts/shortcodes/table/propertiesV200.html", htmlTableBytes.Bytes(), 0644)
}

// filter removes properties that are not required
func filter(propertyInfos []PropertyInfo) []PropertyInfo {
	var initialFilter []PropertyInfo

	for _, prop := range propertyInfos {
		if strings.HasPrefix(prop.PropertyName, "Get") {
			initialFilter = append(initialFilter, prop)
		}
	}

	// Get the primary models (typically ones that have a Get/Show/Process e.t.c. method)
	for i := range initialFilter {
		newName := strings.Replace(initialFilter[i].PropertyName, "Get", "", -1)
		initialFilter[i].PropertyName = newName
	}

	var filtered []PropertyInfo
	for _, prop := range propertyInfos {
		contains := false
		for _, match := range initialFilter {
			// If the property name contains the match and is not an ID
			if strings.Contains(prop.PropertyName, match.PropertyName) && !strings.HasSuffix(prop.PropertyName, "ID") {
				contains = true
				break
			}
		}
		if contains {
			filtered = append(filtered, prop)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		files := strings.Compare(filtered[i].PropertyName, filtered[j].PropertyName)
		if files == 0 {
			return strings.Compare(filtered[i].FileName, filtered[j].FileName) < 0
		}
		return files < 0
	})

	return filtered
}

// parseJSONSchema reads a JSON file and extracts "properties"
func parseJSONSchema(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	// Extract properties if they exist
	props, ok := jsonData["properties"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	var propertyNames []string
	for name := range props {
		propertyNames = append(propertyNames, name)
	}
	return propertyNames, nil
}

// printMarkdownTable prints the properties in a Markdown table format
func printMarkdownTable(propertyInfos []PropertyInfo) bytes.Buffer {
	var buffer bytes.Buffer

	maxFileNameLen := 0
	maxPropertyNameLen := 0
	for _, info := range propertyInfos {
		if len(info.FileName) > maxFileNameLen {
			maxFileNameLen = len(info.FileName)
		}
		if len(info.PropertyName) > maxPropertyNameLen {
			maxPropertyNameLen = len(info.PropertyName)
		}
	}

	buffer.WriteString(fmt.Sprintf("| File Name " + strings.Repeat(" ", maxFileNameLen-9) + "| Property Name " + strings.Repeat(" ", maxPropertyNameLen-13) + "|\n"))
	buffer.WriteString(fmt.Sprintf("|-----------" + strings.Repeat("-", maxFileNameLen-9) + "|---------------" + strings.Repeat("-", maxPropertyNameLen-13) + "|\n"))

	for _, info := range propertyInfos {
		paddedFileName := "[" + info.FileName + "](/schemas/" + info.FileName + ")" + strings.Repeat(" ", maxFileNameLen-len(info.FileName))
		paddedPropertyName := info.PropertyName + strings.Repeat(" ", maxPropertyNameLen-len(info.PropertyName))
		buffer.WriteString(fmt.Sprintf("| %s | %s |\n", paddedFileName, paddedPropertyName))
	}

	return buffer
}

// printHTMLTable prints the properties in an HTML table format
func printHTMLTable(propertyInfos []PropertyInfo) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString(`<style>
	.fluid-table {
		width: 100vw;
		margin-left: calc(50% - 50vw);
		table-layout: auto;
	}
</style>`)
	buffer.WriteString("<div style=\"width:100%; margin-left: -300px\">\n<table>\n")
	buffer.WriteString("<tr><th>Object</th><th style=\"width:150px\">Schema File</th><th>Snippet</th></tr>\n")

	for _, info := range propertyInfos {
		buffer.WriteString(fmt.Sprintf(`<tr>
  <td>%s</td>
  <td style="width:150px">%s</td>
  <td style="width:300px"><pre tabindex="0" style="color:#f8f8f2;background-color:#272822;-moz-tab-size:4;-o-tab-size:4;tab-size:4;"><code class="language-json" data-lang="json"><span style="display:flex;"><span>{</span></span><span style="display:flex;"><span>    <span style="color:#f92672">"$schema"</span>: <span style="color:#e6db74">"{{ .Page.Site.BaseURL }}schemas/`+info.FileName+`"</span>,</span></span><span style="display:flex;"><span>    <span style="color:#f92672">"`+info.PropertyName+`"</span>: { }</span></span><span style="display:flex;"><span>}</span></span></code></pre></td>
</tr>`, info.PropertyName, `<a href="/schemas/`+info.FileName+`">`+info.FileName+`</a>`))
	}

	buffer.WriteString("</table>\n</div>\n")

	return buffer
}
