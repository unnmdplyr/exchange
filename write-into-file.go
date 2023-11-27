package main

import (
	"fmt"
	"os"
)

// Write writes string into file.
func Write(filePath string, text string) error {

	fileHandler, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Didn't manage open file with file name:", filePath, ".\nError:", err)
		return err
	}

	defer func(f *os.File) {
		err2 := f.Close()
		if err2 != nil {
			fmt.Println("Didn't manage close file with file name:", filePath, ".\nError:", err)
		}
	}(fileHandler)

	n, err := fileHandler.WriteString(text)
	if n < len(text) {
		fmt.Println("Didn't manage to write the whole text into file. Starting with:",
			text[:min(10, len(text))], "...")
	}
	if err != nil {
		fmt.Println("Didn't manage to write the string into file:", filePath, ".\nError:", err)
		return err
	}

	return nil
}
