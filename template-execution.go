package main

import (
	"bytes"
	"fmt"
	"text/template"
)

// ExecuteTemplate replaces templates in body and in url.
// Returns with the string replaced the template in it and a possible error.
func ExecuteTemplate(s string, envData *EnvironmentData) (string, error) {

	urlTemplate, err := template.New("test").Parse(s)
	if err != nil {
		fmt.Println("Template couldn't be parsed. It was malformed:", s, ".\nError:", err)
		return "", err
	}

	var buf bytes.Buffer
	err = urlTemplate.Execute(&buf, envData)
	if err != nil {
		fmt.Println("Template execution failed.\nError:", err)
		return "", err
	}

	return buf.String(), nil
}
