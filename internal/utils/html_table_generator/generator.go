package html_table_generator

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
)

func GenerateHTML(data interface{}) (string, error) {
	var headers []string
	var rows [][]interface{}

	dataSlice := reflect.ValueOf(data)
	if dataSlice.Kind() != reflect.Slice || dataSlice.Len() <= 0 {
		return "", fmt.Errorf("html_table_generator.generator.GenerateHTML -> dataSlice.Kind/Len: incorrect input data")
	}

	// Iterate over a slice to fill headers and rows
	firstElem := dataSlice.Index(0)
	elemStruct := firstElem.Type()
	for i := 0; i < firstElem.NumField(); i++ {
		headers = append(headers, elemStruct.Field(i).Name)
	}

	for i := 0; i < dataSlice.Len(); i++ {
		var row []interface{}
		elem := dataSlice.Index(i)
		for j := 0; j < elem.NumField(); j++ {
			row = append(row, elem.Field(j).Interface())
		}
		rows = append(rows, row)
	}

	// HTML generation
	tmpl, err := template.New("table").Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("html_table_generator.generator.GenerateHTML -> template.New: %w", err)
	}
	var html bytes.Buffer
	err = tmpl.Execute(&html, map[string]interface{}{
		"Headers": headers,
		"Rows":    rows,
	})
	if err != nil {
		return "", fmt.Errorf("html_table_generator.generator.GenerateHTML -> tmpl.Execute: %w", err)
	}

	return html.String(), nil
}
