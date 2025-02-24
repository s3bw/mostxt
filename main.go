package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Extract placeholders like {{ key:type('format') }}
func parsePlaceholders(template string) []Placeholder {
	re := regexp.MustCompile(`\{\{\s*(\w+)(?::(\w+))?(?:\('([^']+)'\))?\s*\}\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	var placeholders []Placeholder
	// Match is made up (placeholder, key, type, format)
	// 					(     0     ,  1 ,   2 ,    3  )
	for _, match := range matches {
		typ := StringHold
		format := ""

		if match[2] != "" {
			typ = ListHold
		}
		if match[3] != "" {
			typ = DateTimeHold
			format = match[3]
		}

		placeholders = append(placeholders, Placeholder{Key: match[1], Type: typ, Format: format})
	}

	return placeholders
}

// Struct to store extracted placeholders
type Placeholder struct {
	Key    string
	Type   HolderType
	Format string
}

type HolderType string

func (h HolderType) Is(other HolderType) bool {
	return h == other
}

const (
	StringHold   HolderType = "string"
	DateTimeHold HolderType = "datetime"
	ListHold     HolderType = "list"
)

// Format datetime based on user-specified format
func formatDatetime(format string) string {
	replacements := map[string]string{
		"YYYY": "2006", "MM": "01", "DD": "02",
		"HH": "15", "mm": "04", "ss": "05",
	}
	for k, v := range replacements {
		format = strings.ReplaceAll(format, k, v)
	}
	return time.Now().Format(format)
}

// Prompt user for input based on type
func getUserInput(p Placeholder) string {
	if p.Type.Is(DateTimeHold) {
		return formatDatetime(p.Format)
	}

	reader := bufio.NewReader(os.Stdin)

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s %s (%s): ", green("Enter"), p.Key, p.Type)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if p.Type.Is(ListHold) {
		return strings.Join(strings.Split(input, ","), ", ")
	}
	return input
}

// Replace placeholders in the template with user input
func fillTemplate(template string) string {
	placeholders := parsePlaceholders(template)
	values := make(map[string]string)

	// Get user input for each placeholder
	for _, p := range placeholders {
		values[p.Key] = getUserInput(p)
	}

	// Replace placeholders with user input
	re := regexp.MustCompile(`\{\{\s*(\w+)(?::(\w+))?(?:\('([^']+)'\))?\s*\}\}`)
	result := re.ReplaceAllStringFunc(template, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		return values[key]
	})

	return result
}

// Read input from file or stdin
func readTemplate() (string, error) {
	if len(os.Args) > 1 {
		// Read from file if filename is provided as argument
		filename := os.Args[1]
		content, err := os.ReadFile(filename)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	return "", errors.New("No provided input")
}

// Write output to a file
func writeToFile(output string, filename string) error {
	return os.WriteFile(filename, []byte(output), 0644)
}

func main() {
	template, err := readTemplate()
	if err != nil {
		fmt.Println("Error reading template:", err)
		os.Exit(1)
	}

	result := fillTemplate(template)

	// Determine output filename
	// outputFilename := "{{ datetime('YYYY-MM-DD'}}-{{ title }}.md"
	outputFilename := "output.txt"
	if len(os.Args) > 2 {
		outputFilename = os.Args[2]
	}

	err = writeToFile(result, outputFilename)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}

	fmt.Printf("\nFinal output written to %s\n", outputFilename)
}
