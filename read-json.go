package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadData(filePath string, data interface{}) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Couldn't read file:", filePath, ".\nError:", err)
		return err
	}

	err = json.Unmarshal(fileBytes, data)
	if err != nil {
		fmt.Println("Unmarshalling error:", err)
		return err
	}

	return nil
}
