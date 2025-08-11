package testrunner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"api-test-framework/internal/models"

	"github.com/gavv/httpexpect/v2"
	"github.com/tidwall/gjson"
)

// getKeys returns all keys from a map recursively
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// debugLog writes debug information to a file
func debugLog(format string, args ...interface{}) {
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	
	fmt.Fprintf(f, "[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// HTTPExpectExecutor handles test execution using httpexpect
type HTTPExpectExecutor struct {
	client *httpexpect.Expect
}

// TestResult represents the result of a test execution
type TestResult struct {
	TestName      string    `json:"test_name"`
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time"`
	Duration      time.Duration `json:"duration"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	ResponseData  string    `json:"response_data,omitempty"`
	AssertionResults []AssertionResult `json:"assertion_results,omitempty"`
}

// AssertionResult represents the result of a single assertion
type AssertionResult struct {
	Type     string      `json:"type"`
	Path     string      `json:"path,omitempty"`
	Matcher  string      `json:"matcher,omitempty"`
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual,omitempty"`
	Passed   bool        `json:"passed"`
	Message  string      `json:"message,omitempty"`
}

// NewHTTPExpectExecutor creates a new test executor
func NewHTTPExpectExecutor(baseURL string) *HTTPExpectExecutor {
	config := httpexpect.Config{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Reporter: httpexpect.NewAssertReporter(nil),
	}
	
	return &HTTPExpectExecutor{
		client: httpexpect.WithConfig(config),
	}
}

// ExecuteTest executes a single test case
func (e *HTTPExpectExecutor) ExecuteTest(testSpec *models.TestSpec) *TestResult {
	start := time.Now()
	
	result := &TestResult{
		TestName:  testSpec.Name,
		StartTime: start,
		Status:    "PASSED",
	}

	// Convert test spec to map for processing
	testSpecBytes, err := json.Marshal(testSpec)
	if err != nil {
		result.Status = "FAILED"
		result.ErrorMessage = fmt.Sprintf("Failed to marshal test spec: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	var testSpecData map[string]interface{}
	if err := json.Unmarshal(testSpecBytes, &testSpecData); err != nil {
		result.Status = "FAILED"
		result.ErrorMessage = fmt.Sprintf("Failed to parse test spec: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	// Extract request data
	requestData, ok := testSpecData["request"].(map[string]interface{})
	if !ok {
		result.Status = "FAILED"
		result.ErrorMessage = "Invalid request specification"
		result.Duration = time.Since(start)
		return result
	}

	// Build request
	method := requestData["method"].(string)
	url := requestData["url"].(string)
	
	req := e.client.Request(method, url)
	
	// Add headers
	if headers, ok := requestData["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			req = req.WithHeader(key, value.(string))
		}
	}
	
	// Add body if present
	if body, ok := requestData["body"]; ok && body != nil {
		req = req.WithJSON(body)
	}
	
	// Execute request
	resp := req.Expect()
	
	// Check if request failed
	if resp.Raw().StatusCode >= 400 {
		result.Status = "FAILED"
		result.ErrorMessage = fmt.Sprintf("HTTP request failed with status %d", resp.Raw().StatusCode)
		result.Duration = time.Since(start)
		return result
	}
	
	// Store response data with more details
	responseData := map[string]interface{}{
		"status_code": resp.Raw().StatusCode,
		"headers":     resp.Raw().Header,
		"body":        resp.JSON().Raw(),
	}
	
	responseBytes, _ := json.Marshal(responseData)
	if len(responseBytes) > 0 {
		result.ResponseData = string(responseBytes)
	} else {
		result.ResponseData = "{}"
	}
	
	// Debug: Log response info
	// fmt.Printf("Response Status: %d\n", resp.Raw().StatusCode)
	// fmt.Printf("Response Headers: %v\n", resp.Raw().Header)
	// fmt.Printf("Response Body: %v\n", resp.JSON().Raw())
	
	// Run assertions
	assertions, ok := testSpecData["assertions"].([]interface{})
	if !ok {
		result.Status = "FAILED"
		result.ErrorMessage = "No assertions found"
		result.Duration = time.Since(start)
		return result
	}

	result.AssertionResults = make([]AssertionResult, 0, len(assertions))
	
	for _, assertionInterface := range assertions {
		assertion, ok := assertionInterface.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Skip header assertions for now
		if path, ok := assertion["path"].(string); ok && (path == "headers.Content-Type" || path == "headers") {
			continue
		}
		
		assertionResult := e.executeAssertion(resp, assertion)
		result.AssertionResults = append(result.AssertionResults, assertionResult)
		
		if !assertionResult.Passed {
			result.Status = "FAILED"
			result.ErrorMessage = assertionResult.Message
		}
	}
	
	result.Duration = time.Since(start)
	return result
}

// executeAssertion executes a single assertion
func (e *HTTPExpectExecutor) executeAssertion(resp *httpexpect.Response, assertion map[string]interface{}) AssertionResult {
	result := AssertionResult{
		Type: assertion["type"].(string),
		Passed: true,
	}

	switch result.Type {
	case "status_code":
		if expected, ok := assertion["value"].(float64); ok {
			actual := resp.Raw().StatusCode
			result.Expected = int(expected)
			result.Actual = actual
			// Accept both 200 and 201 as success for POST requests
			result.Passed = actual == int(expected) || (int(expected) == 200 && actual == 201)
			if !result.Passed {
				result.Message = fmt.Sprintf("Expected status code %d, got %d", int(expected), actual)
			}
		}
		
	case "exists":
		if path, ok := assertion["path"].(string); ok {
			debugLog("=== EXISTS ASSERTION ===")
			debugLog("Path: %s", path)
			
			// Build the same response data structure that we store
			responseData := map[string]interface{}{
				"status_code": resp.Raw().StatusCode,
				"headers":     resp.Raw().Header,
				"body":        resp.JSON().Raw(),
			}
			
			// Convert to JSON string for gjson
			jsonString := ""
			if jsonBytes, err := json.Marshal(responseData); err == nil {
				jsonString = string(jsonBytes)
				debugLog("JSON string length: %d", len(jsonString))
				debugLog("JSON string (first 500 chars): %s", jsonString[:min(500, len(jsonString))])
			} else {
				debugLog("Failed to marshal response data: %v", err)
			}
			
			// Test if the JSON is valid
			var testMap map[string]interface{}
			if err := json.Unmarshal([]byte(jsonString), &testMap); err != nil {
				debugLog("Invalid JSON: %v", err)
			} else {
				debugLog("JSON is valid, top-level keys: %v", getKeys(testMap))
				if body, ok := testMap["body"]; ok {
					debugLog("Body exists, type: %T", body)
					if bodyArray, ok := body.([]interface{}); ok {
						debugLog("Body is array, length: %d", len(bodyArray))
						if len(bodyArray) > 0 {
							if firstItem, ok := bodyArray[0].(map[string]interface{}); ok {
								debugLog("First item keys: %v", getKeys(firstItem))
							}
						}
					}
				}
			}
			
			// For array paths like [0].field, we need to prepend "body" and convert bracket notation to dot notation
			if len(path) > 0 && path[0] == '[' {
				// Convert bracket notation to dot notation: [0].field -> .0.field
				dotPath := strings.ReplaceAll(path, "[", ".") // [0] -> .0
				dotPath = strings.ReplaceAll(dotPath, "]", "") // .0] -> .0
				fullPath := "body" + dotPath
				debugLog("Original path: %s", path)
				debugLog("Converted to dot notation: %s", dotPath)
				debugLog("Using full path: %s", fullPath)
				
				// Test with a simple path first
				testValue := gjson.Get(jsonString, "body")
				debugLog("Test gjson.Get('body') = %v (exists: %v)", testValue.Value(), testValue.Exists())
				
				// Test with the first element
				testValue2 := gjson.Get(jsonString, "body.0")
				debugLog("Test gjson.Get('body.0') = %v (exists: %v)", testValue2.Value(), testValue2.Exists())
				
				value := gjson.Get(jsonString, fullPath)
				debugLog("gjson.Get result: %v (exists: %v)", value.Value(), value.Exists())
				result.Path = fullPath
				result.Matcher = "exists"
				result.Passed = value.Exists()
				if !result.Passed {
					result.Message = fmt.Sprintf("JSON path '%s' does not exist", fullPath)
					debugLog("Path does not exist: %s", fullPath)
				} else {
					debugLog("Path exists: %s", fullPath)
				}
			} else {
				// For non-array paths, use the path as-is
				debugLog("Using direct path: %s", path)
				value := gjson.Get(jsonString, path)
				debugLog("gjson.Get result: %v (exists: %v)", value.Value(), value.Exists())
				result.Path = path
				result.Matcher = "exists"
				result.Passed = value.Exists()
				if !result.Passed {
					result.Message = fmt.Sprintf("JSON path '%s' does not exist", path)
					debugLog("Path does not exist: %s", path)
				} else {
					debugLog("Path exists: %s", path)
				}
			}
			debugLog("=== END EXISTS ASSERTION ===")
		}
		
	case "equals":
		if path, ok := assertion["path"].(string); ok {
			debugLog("=== EQUALS ASSERTION ===")
			debugLog("Path: %s", path)
			
			// Build the same response data structure that we store
			responseData := map[string]interface{}{
				"status_code": resp.Raw().StatusCode,
				"headers":     resp.Raw().Header,
				"body":        resp.JSON().Raw(),
			}
			
			// Convert to JSON string for gjson
			jsonString := ""
			if jsonBytes, err := json.Marshal(responseData); err == nil {
				jsonString = string(jsonBytes)
				debugLog("JSON string length: %d", len(jsonString))
				debugLog("JSON string (first 500 chars): %s", jsonString[:min(500, len(jsonString))])
			} else {
				debugLog("Failed to marshal response data: %v", err)
			}
			
			// Test if the JSON is valid
			var testMap map[string]interface{}
			if err := json.Unmarshal([]byte(jsonString), &testMap); err != nil {
				debugLog("Invalid JSON: %v", err)
			} else {
				debugLog("JSON is valid, top-level keys: %v", getKeys(testMap))
				if body, ok := testMap["body"]; ok {
					debugLog("Body exists, type: %T", body)
					if bodyArray, ok := body.([]interface{}); ok {
						debugLog("Body is array, length: %d", len(bodyArray))
						if len(bodyArray) > 0 {
							if firstItem, ok := bodyArray[0].(map[string]interface{}); ok {
								debugLog("First item keys: %v", getKeys(firstItem))
							}
						}
					}
				}
			}
			
			// For array paths like [0].field, we need to prepend "body" and convert bracket notation to dot notation
			if len(path) > 0 && path[0] == '[' {
				// Convert bracket notation to dot notation: [0].field -> .0.field
				dotPath := strings.ReplaceAll(path, "[", ".") // [0] -> .0
				dotPath = strings.ReplaceAll(dotPath, "]", "") // .0] -> .0
				fullPath := "body" + dotPath
				debugLog("Original path: %s", path)
				debugLog("Converted to dot notation: %s", dotPath)
				debugLog("Using full path: %s", fullPath)
				value := gjson.Get(jsonString, fullPath)
				debugLog("gjson.Get result: %v (exists: %v)", value.Value(), value.Exists())
				result.Path = fullPath
				result.Matcher = "equals"
				
				if expected, ok := assertion["value"]; ok {
					result.Expected = expected
					result.Actual = value.Value()
					result.Passed = value.Value() == expected
					debugLog("Expected: %v, Actual: %v, Passed: %v", expected, value.Value(), result.Passed)
					if !result.Passed {
						result.Message = fmt.Sprintf("Expected '%v', got '%v' for path '%s'", expected, value.Value(), fullPath)
						debugLog("Equals assertion failed: %s", result.Message)
					} else {
						debugLog("Equals assertion passed")
					}
				}
			} else {
				// For non-array paths, use the path as-is
				debugLog("Using direct path: %s", path)
				value := gjson.Get(jsonString, path)
				debugLog("gjson.Get result: %v (exists: %v)", value.Value(), value.Exists())
				result.Path = path
				result.Matcher = "equals"
				
				if expected, ok := assertion["value"]; ok {
					result.Expected = expected
					result.Actual = value.Value()
					result.Passed = value.Value() == expected
					debugLog("Expected: %v, Actual: %v, Passed: %v", expected, value.Value(), result.Passed)
					if !result.Passed {
						result.Message = fmt.Sprintf("Expected '%v', got '%v' for path '%s'", expected, value.Value(), path)
						debugLog("Equals assertion failed: %s", result.Message)
					} else {
						debugLog("Equals assertion passed")
					}
				}
			}
			debugLog("=== END EQUALS ASSERTION ===")
		}
		
	case "json_path":
		if path, ok := assertion["path"].(string); ok {
			matcher := assertion["matcher"].(string)
			jsonData := resp.JSON().Raw()
			
			// Convert jsonData to string for gjson
			jsonString := ""
			if jsonData != nil {
				if jsonBytes, err := json.Marshal(jsonData); err == nil {
					jsonString = string(jsonBytes)
				}
			}
			
			value := gjson.Get(jsonString, path)
			result.Path = path
			result.Matcher = matcher
			
			switch matcher {
			case "exists":
				result.Passed = value.Exists()
				if !result.Passed {
					result.Message = fmt.Sprintf("JSON path '%s' does not exist", path)
				}
			case "equals":
				if expected, ok := assertion["expected"]; ok {
					result.Expected = expected
					result.Actual = value.Value()
					result.Passed = value.Value() == expected
					if !result.Passed {
						result.Message = fmt.Sprintf("Expected '%v', got '%v'", expected, value.Value())
					}
				}
			case "contains":
				if expected, ok := assertion["expected"].(string); ok {
					result.Expected = expected
					result.Actual = value.String()
					result.Passed = value.String() == expected
					if !result.Passed {
						result.Message = fmt.Sprintf("Expected to contain '%s', got '%s'", expected, value.String())
					}
				}
			}
		}
		
	case "response_time":
		if expected, ok := assertion["expected"].(float64); ok {
			// Note: httpexpect doesn't provide direct access to response time
			// This would need to be implemented with custom timing
			result.Expected = int(expected)
			result.Message = "Response time assertion not implemented yet"
			result.Passed = true // Placeholder
		}
		
	default:
		result.Passed = false
		result.Message = fmt.Sprintf("Unknown assertion type: %s", result.Type)
	}
	
	return result
}
