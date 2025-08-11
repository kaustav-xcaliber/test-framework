package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Assertion represents a single assertion
type Assertion struct {
	Type  string      `json:"type"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// AssertionGenerator generates assertions from JSON data
type AssertionGenerator struct {
	MaxDepth     int
	MaxArraySize int
	IncludeNulls bool
}

// NewAssertionGenerator creates a new assertion generator
func NewAssertionGenerator() *AssertionGenerator {
	return &AssertionGenerator{
		MaxDepth:     5,
		MaxArraySize: 3,
		IncludeNulls: false,
	}
}

// GenerateAssertions generates assertions from JSON data
func (ag *AssertionGenerator) GenerateAssertions(data interface{}, path string) []Assertion {
	var assertions []Assertion
	
	switch v := data.(type) {
	case map[string]interface{}:
		// Add existence assertion for the object itself
		if path != "" {
			assertions = append(assertions, Assertion{
				Type: "exists",
				Path: path,
			})
		}
		
		// Process each field
		for key, value := range v {
			currentPath := key
			if path != "" {
				currentPath = path + "." + key
			}
			
			// Add existence assertion for each field
			assertions = append(assertions, Assertion{
				Type: "exists",
				Path: currentPath,
			})
			
			// Recursively process nested values
			if ag.MaxDepth > 0 {
				ag.MaxDepth--
				nestedAssertions := ag.GenerateAssertions(value, currentPath)
				assertions = append(assertions, nestedAssertions...)
				ag.MaxDepth++
			}
		}
		
	case []interface{}:
		// Add existence assertion for the array itself
		if path != "" {
			assertions = append(assertions, Assertion{
				Type: "exists",
				Path: path,
			})
		}
		
		// Process array elements (limited by MaxArraySize)
		for i, value := range v {
			if i >= ag.MaxArraySize {
				break
			}
			
			currentPath := fmt.Sprintf("%s[%d]", path, i)
			if path == "" {
				currentPath = fmt.Sprintf("[%d]", i)
			}
			
			// Add existence assertion for each array element
			assertions = append(assertions, Assertion{
				Type: "exists",
				Path: currentPath,
			})
			
			// Recursively process nested values
			if ag.MaxDepth > 0 {
				ag.MaxDepth--
				nestedAssertions := ag.GenerateAssertions(value, currentPath)
				assertions = append(assertions, nestedAssertions...)
				ag.MaxDepth++
			}
		}
		
	case string:
		// Add value assertion for strings
		if path != "" {
			assertions = append(assertions, Assertion{
				Type:  "equals",
				Path:  path,
				Value: v,
			})
		}
		
	case float64:
		// Add value assertion for numbers
		if path != "" {
			assertions = append(assertions, Assertion{
				Type:  "equals",
				Path:  path,
				Value: v,
			})
		}
		
	case bool:
		// Add value assertion for booleans
		if path != "" {
			assertions = append(assertions, Assertion{
				Type:  "equals",
				Path:  path,
				Value: v,
			})
		}
		
	case nil:
		// Add null assertion if IncludeNulls is true
		if ag.IncludeNulls && path != "" {
			assertions = append(assertions, Assertion{
				Type:  "equals",
				Path:  path,
				Value: nil,
			})
		}
	}
	
	return assertions
}

// GenerateStatusAssertion generates a status code assertion
func GenerateStatusAssertion(statusCode int) Assertion {
	return Assertion{
		Type:  "status_code",
		Path:  "",
		Value: statusCode,
	}
}

// GenerateHeaderAssertions generates basic header assertions
func GenerateHeaderAssertions() []Assertion {
	return []Assertion{
		{
			Type: "exists",
			Path: "headers.Content-Type",
		},
		{
			Type:  "equals",
			Path:  "headers.Content-Type",
			Value: "application/json",
		},
	}
}

// GenerateCommonAssertions generates common assertions for API responses
func GenerateCommonAssertions() []Assertion {
	return []Assertion{
		{
			Type: "exists",
			Path: "data",
		},
		{
			Type: "exists",
			Path: "message",
		},
		{
			Type: "exists",
			Path: "status",
		},
	}
}

// FormatAssertion formats an assertion for display
func FormatAssertion(assertion Assertion) string {
	switch assertion.Type {
	case "status_code":
		return fmt.Sprintf(`{"type": "status_code", "value": %v}`, assertion.Value)
	case "exists":
		return fmt.Sprintf(`{"type": "exists", "path": "%s"}`, assertion.Path)
	case "equals":
		if assertion.Value == nil {
			return fmt.Sprintf(`{"type": "equals", "path": "%s", "value": null}`, assertion.Path)
		}
		switch v := assertion.Value.(type) {
		case string:
			return fmt.Sprintf(`{"type": "equals", "path": "%s", "value": "%s"}`, assertion.Path, v)
		default:
			return fmt.Sprintf(`{"type": "equals", "path": "%s", "value": %v}`, assertion.Path, v)
		}
	default:
		return fmt.Sprintf(`{"type": "%s", "path": "%s"}`, assertion.Type, assertion.Path)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/generate_assertions.go <json_file> [status_code]")
		fmt.Println("Example: go run scripts/generate_assertions.go response.json 200")
		os.Exit(1)
	}
	
	// Read JSON file
	jsonFile := os.Args[1]
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	
	// Parse JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	
	// Create assertion generator
	generator := NewAssertionGenerator()
	
	// Generate assertions
	assertions := generator.GenerateAssertions(jsonData, "")
	
	// Add status code assertion if provided
	if len(os.Args) > 2 {
		statusCode, err := strconv.Atoi(os.Args[2])
		if err == nil {
			assertions = append([]Assertion{GenerateStatusAssertion(statusCode)}, assertions...)
		}
	}
	
	// Add common header assertions (commented out for now)
	// headerAssertions := GenerateHeaderAssertions()
	// assertions = append(assertions, headerAssertions...)
	
	// Add common API response assertions
	commonAssertions := GenerateCommonAssertions()
	assertions = append(assertions, commonAssertions...)
	
	// Output results
	fmt.Println("Generated Assertions:")
	fmt.Println("====================")
	
	// Format as JSON array
	var formattedAssertions []string
	for _, assertion := range assertions {
		formattedAssertions = append(formattedAssertions, FormatAssertion(assertion))
	}
	
	// Remove duplicates while preserving order
	seen := make(map[string]bool)
	var uniqueAssertions []string
	for _, assertion := range formattedAssertions {
		if !seen[assertion] {
			seen[assertion] = true
			uniqueAssertions = append(uniqueAssertions, assertion)
		}
	}
	
	fmt.Printf("[\n  %s\n]\n", strings.Join(uniqueAssertions, ",\n  "))
	
	// Also output as individual assertions for easy copy-paste
	fmt.Println("\nIndividual Assertions:")
	fmt.Println("======================")
	for i, assertion := range uniqueAssertions {
		fmt.Printf("%d. %s\n", i+1, assertion)
	}
	
	// Output usage instructions
	fmt.Println("\nUsage Instructions:")
	fmt.Println("==================")
	fmt.Println("1. Copy the JSON array above and use it as the 'assertions' field in your test spec")
	fmt.Println("2. For status code assertions, ensure your test runner supports 'status_code' type")
	fmt.Println("3. For header assertions, ensure your test runner supports header path access")
	fmt.Println("4. Modify assertions as needed for your specific test requirements")
}
