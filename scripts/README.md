# Assertion Generation Script

This directory contains a utility script for generating assertions based on JSON output for testing purposes.

## Script: `generate_assertions.go`

A Go script that analyzes JSON responses and automatically generates appropriate assertions for API testing.

### Features

- **Automatic Assertion Generation**: Analyzes JSON structure and generates existence and value assertions
- **Configurable Depth**: Limits recursion depth to avoid excessive assertions
- **Array Handling**: Processes arrays with configurable size limits
- **Status Code Support**: Optional status code assertion generation
- **Header Assertions**: Generates common header assertions
- **Multiple Output Formats**: Provides both JSON array and individual assertion formats

### Usage

```bash
# Basic usage with JSON file
go run scripts/generate_assertions.go response.json

# With status code
go run scripts/generate_assertions.go response.json 200

# Using the example file
go run scripts/generate_assertions.go scripts/example_response.json 200
```

### Output

The script generates two types of output:

1. **JSON Array Format**: Ready-to-use assertions array for your test spec
2. **Individual Assertions**: Numbered list for easy copy-paste

### Generated Assertion Types

- **`exists`**: Checks if a field exists in the response
- **`equals`**: Checks if a field equals a specific value
- **`status_code`**: Checks HTTP status code (when provided)

**Note**: Header assertions are currently disabled to avoid issues with response header access.

### Example Output

```json
[
  { "type": "status_code", "value": 200 },
  { "type": "exists", "path": "status" },
  { "type": "equals", "path": "status", "value": "success" },
  { "type": "exists", "path": "message" },
  {
    "type": "equals",
    "path": "message",
    "value": "Data retrieved successfully"
  },
  { "type": "exists", "path": "data" },
  { "type": "exists", "path": "data.id" },
  { "type": "equals", "path": "data.id", "value": "12345" },
  { "type": "exists", "path": "data.name" },
  { "type": "equals", "path": "data.name", "value": "John Doe" }
]
```

### Configuration Options

The script has several configurable options in the `AssertionGenerator` struct:

- **MaxDepth**: Maximum recursion depth (default: 5)
- **MaxArraySize**: Maximum array elements to process (default: 3)
- **IncludeNulls**: Whether to include null value assertions (default: false)

### Integration with Test Framework

1. **Save JSON Response**: Save your API response to a JSON file
2. **Generate Assertions**: Run the script on your JSON file
3. **Copy Assertions**: Use the generated JSON array in your test spec
4. **Customize**: Modify assertions as needed for your specific requirements

### Example with Healthcare API

For the encounter API response you provided:

```bash
# Save the encounter response to a file
# Then run:
go run scripts/generate_assertions.go encounter_response.json 200
```

This will generate assertions for:

- Status code (200)
- Response structure validation
- Specific field existence checks
- Value assertions for important fields

### Tips

1. **Review Generated Assertions**: Always review and customize the generated assertions
2. **Remove Unnecessary Assertions**: Remove assertions for fields that may change
3. **Add Custom Assertions**: Add specific business logic assertions
4. **Test the Assertions**: Verify the generated assertions work with your test runner

### File Structure

```
scripts/
├── generate_assertions.go    # Main assertion generation script
├── example_response.json     # Example JSON response for testing
└── README.md                # This documentation
```
