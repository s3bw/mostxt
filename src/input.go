package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// formatDatetime formats the current time according to the given format
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

// getUserInput prompts for and processes user input based on the argument type
func GetUserInput(argument Argument) string {
	// Handle datetime type
	if argument.Type == ArgumentTypeDatetime {
		for _, prop := range argument.Properties {
			if prop.Type == PropertyTypeFormat {
				return formatDatetime(prop.Value)
			}
		}
		// If no format found, use default format
		return formatDatetime("YYYY-MM-DD")
	}

	// Get example and description for the prompt
	var example, description string
	for _, prop := range argument.Properties {
		switch prop.Type {
		case PropertyTypeExample:
			example = prop.Value
		case PropertyTypeDescribe:
			description = prop.Value
		}
	}

	reader := bufio.NewReader(os.Stdin)
	green := color.New(color.FgGreen).SprintFunc()
	grey := color.New(color.FgHiBlack).SprintFunc()

	// Build the prompt
	prompt := fmt.Sprintf("%s %s", green("Enter"), argument.Name)
	if example != "" {
		prompt = fmt.Sprintf("%s %s (e.g: %s)", green("Enter"), argument.Name, example)
	}

	if description != "" {
		prompt += fmt.Sprintf("\n%s", grey(description))
	}
	prompt += "\n$ "

	fmt.Print(prompt)

	// Handle list type with multiple entries
	if argument.Type == ArgumentTypeList {
		var items []string
		for {
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			
			// Empty line signals end of input
			if input == "" {
				break
			}
			
			items = append(items, input)
			fmt.Print("$ ")
		}
		if len(items) > 0 {
			return strings.Join(items, ", ")
		}
		return ""
	}

	// Handle single line input for other types
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}