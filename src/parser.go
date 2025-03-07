package src

import (
	"fmt"
	"regexp"
	"strings" 
)

type PropertyType string

const (
	PropertyTypeExample  PropertyType = "example"
	PropertyTypeDescribe PropertyType = "describe"
	PropertyTypeDefault  PropertyType = "default"
	PropertyTypeFormat   PropertyType = "format"
)

// Property represents a property attached to an argument
type Property struct {
	Name  string
	Value string
	Type  PropertyType
}

type ArgumentType string

const (
	ArgumentTypeString   ArgumentType = "string"
	ArgumentTypeList     ArgumentType = "list"
	ArgumentTypeDatetime ArgumentType = "datetime"
)

// Argument represents an argument with its properties and type
//
// Arguments are used to pass data to the template.
// They are defined in the template file using the {{ variable_name:type 'default_value' }} syntax.
//
// Example:
// {{ name:string 'John Doe' }}
type Argument struct {
	Name       string
	Type       ArgumentType
	Properties []Property
}

// Parser handles parsing of template files
type Parser struct {
	arguments map[string]Argument
}

// NewTemplateParser creates a new instance of TemplateParser
func NewParser() *Parser {
	return &Parser{
		arguments: make(map[string]Argument),
	}
}

// ParseTemplate parses the template content and returns a map of variables
func (p *Parser) ParseTemplate(content string) (map[string]Argument, error) {
	// Pattern to match {{ variable_name:type 'default_value' }}
	// and multiple property_name('value') pairs
	pattern := regexp.MustCompile(`\{\{\s*([^:}\s]+)(?::(\w+)(?:\s+'([^']+)')?)?(?:\s+(\w+)\(([^)]+)\))*(?:\s+\w+\([^)]+\))*\s*\}\}`)
	propertyPattern := regexp.MustCompile(`(\w+)\(([^)]+)\)`)

	matches := pattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 { // We only need the variable name to be present
			continue
		}

		varName := match[1]
		varType := match[2]
		defaultValue := match[3]
		if varType == "" {
			varType = "string"
		}

		// Create variable object
		argument := Argument{
			Name:       varName,
			Type:       ArgumentType(varType),
			Properties: make([]Property, 0),
		}

		// If it's a datetime type with a format, add it as a property
		if varType == "datetime" {
			if defaultValue == "" {
				return nil, fmt.Errorf("datetime type requires a format")
			}
			argument.Properties = append(argument.Properties, Property{
				Name:  "format",
				Value: defaultValue,
				Type:  PropertyTypeFormat,
			})
		}

		// Find all properties in the entire match
		fullMatch := match[0]
		propertyMatches := propertyPattern.FindAllStringSubmatch(fullMatch, -1)
		
		// Handle all properties
		for _, propMatch := range propertyMatches {
			if len(propMatch) == 3 {
				propertyName := propMatch[1]
				// Remove surrounding quotes if present
				propertyValue := strings.Trim(propMatch[2], "'\"")
				
				argument.Properties = append(argument.Properties, Property{
					Name:  propertyName,
					Value: propertyValue,
					Type:  PropertyType(propertyName),
				})
			}
		}

		p.arguments[varName] = argument
	}

	return p.arguments, nil
}

// GetVariable returns a specific variable by name
func (p *Parser) GetArgument(name string) (Argument, bool) {
	argument, exists := p.arguments[name]
	return argument, exists
}

// GetAllVariables returns all parsed variables
func (p *Parser) GetAllArguments() map[string]Argument {
	return p.arguments
}

// ParseTemplateFile reads and parses a template file
func ParseTemplate(content string) ([]Argument, error) {
	parser := NewParser()
	args, err := parser.ParseTemplate(content)
	if err != nil {
		return nil, err
	}
	
	// Convert map to slice
	var arguments []Argument
	for _, arg := range args {
		arguments = append(arguments, arg)
	}
	return arguments, nil
}