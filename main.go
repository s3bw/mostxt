package main

import (
	"fmt"
	"os"

	"github.com/s3bw/mostxt/src"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: mostxt <template_file> <output_file>")
		os.Exit(1)
	}

	template, err := src.ReadTemplate()
	if err != nil {
		fmt.Println("Error reading template:", err)
		os.Exit(1)
	}

	arguments, err := src.ParseTemplate(template)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	// Collect all inputs
	inputs := make(map[string]string)
	for _, argument := range arguments {
		inputs[argument.Name] = src.GetUserInput(argument)
	}

	// Fill the template with inputs
	result := src.FillTemplate(template, arguments, inputs)

	// Write to output file
	outputFile := os.Args[2]
	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", outputFile, err)
		os.Exit(1)
	}

	fmt.Printf("\nOutput written to %s\n", outputFile)
}