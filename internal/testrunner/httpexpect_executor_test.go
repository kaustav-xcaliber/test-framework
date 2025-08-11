package testrunner

import (
	"testing"

	"api-test-framework/internal/models"
)

func TestHTTPExpectExecutor_ExecuteTest(t *testing.T) {
	// Create a test executor
	executor := NewHTTPExpectExecutor("https://jsonplaceholder.typicode.com")

	// Create a test specification
	testSpec := &models.TestSpec{
		Name:        "Test Get User",
		Description: "Test to get user by ID",
		ServiceName: "user-service",
		Request: models.RequestSpec{
			Method: "GET",
			URL:    "/users/1",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		Assertions: []models.AssertionSpec{
			{
				Type:     "status_code",
				Expected: 200,
			},
			{
				Type:     "json_path",
				Path:     "$.id",
				Matcher:  "exists",
				Expected: nil,
			},
		},
	}

	// Execute the test
	result := executor.ExecuteTest(testSpec)

	// Verify the result
	if result.Status != "PASSED" {
		t.Errorf("Expected test to pass, got status: %s", result.Status)
	}

	if result.Duration <= 0 {
		t.Errorf("Expected execution time to be positive, got: %v", result.Duration)
	}

	if len(result.AssertionResults) != 2 {
		t.Errorf("Expected 2 assertion results, got: %d", len(result.AssertionResults))
	}
}

func TestHTTPExpectExecutor_NewHTTPExpectExecutor(t *testing.T) {
	baseURL := "https://api.example.com"
	executor := NewHTTPExpectExecutor(baseURL)

	if executor == nil {
		t.Error("Expected executor to be created, got nil")
	}

	if executor.client == nil {
		t.Error("Expected client to be initialized, got nil")
	}
}
