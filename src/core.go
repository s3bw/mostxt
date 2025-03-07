package src

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

// Read input from file or stdin
func ReadTemplate() (string, error) {
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

// fillTemplate replaces all placeholders with their values
func FillTemplate(template string, arguments []Argument, inputs map[string]string) string {
	// Replace each argument placeholder with its input
	for _, arg := range arguments {
		input := inputs[arg.Name]
		
		// Create pattern to match the full placeholder including properties
		pattern := fmt.Sprintf(`\{\{\s*%s(?::[\w]+)?(?:\s+'[^']+')?(?:\s+\w+\([^)]+\))*\s*\}\}`, arg.Name)
		re := regexp.MustCompile(pattern)
		
		template = re.ReplaceAllString(template, input)
	}
	
	return template
}

