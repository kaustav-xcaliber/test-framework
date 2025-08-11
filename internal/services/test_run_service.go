package services

import (
	"encoding/json"
	"fmt"
	"time"

	"api-test-framework/internal/models"
	"api-test-framework/internal/testrunner"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// TestRunService handles test execution and result management
type TestRunService struct {
	db              *gorm.DB
	testRunner      *testrunner.HTTPExpectExecutor
	redisClient     *redis.Client
}

// NewTestRunService creates a new test run service
func NewTestRunService(db *gorm.DB, testRunner *testrunner.HTTPExpectExecutor, redisClient *redis.Client) *TestRunService {
	return &TestRunService{
		db:          db,
		testRunner:  testRunner,
		redisClient: redisClient,
	}
}

// StartTestRun starts a new test execution run
func (s *TestRunService) StartTestRun(serviceID string, testIDs []string, name string) (*models.TestRun, error) {
	// Create test run
	testRun := &models.TestRun{
		Name:       name,
		Status:     "running",
		StartedAt:  time.Now(),
	}

	if err := s.db.Create(testRun).Error; err != nil {
		return nil, fmt.Errorf("failed to create test run: %v", err)
	}

	// Get test cases
	var testCases []models.TestCase
	query := s.db.Preload("Service")
	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}
	if len(testIDs) > 0 {
		query = query.Where("id IN ?", testIDs)
	}
	
	if err := query.Find(&testCases).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve test cases: %v", err)
	}

	testRun.TotalTests = len(testCases)
	if err := s.db.Save(testRun).Error; err != nil {
		return nil, fmt.Errorf("failed to update test run: %v", err)
	}

	// Execute tests asynchronously with timeout
	go func() {
		// Set a timeout of 5 minutes for test execution
		done := make(chan bool, 1)
		go func() {
			s.executeTests(testRun.ID, testCases)
			done <- true
		}()
		
		select {
		case <-done:
			// Test execution completed normally
		case <-time.After(5 * time.Minute):
			// Test execution timed out
			fmt.Printf("Test run %s timed out after 5 minutes\n", testRun.ID)
			completedAt := time.Now()
			s.db.Model(&models.TestRun{}).Where("id = ?", testRun.ID).Updates(map[string]interface{}{
				"status":          "failed",
				"passed_tests":    0,
				"failed_tests":    0,
				"execution_time_ms":  0,
				"completed_at":    completedAt,
			})
		}
	}()

	return testRun, nil
}

// executeTests executes all tests for a test run
func (s *TestRunService) executeTests(testRunID string, testCases []models.TestCase) {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in executeTests for test run %s: %v\n", testRunID, r)
			completedAt := time.Now()
			s.db.Model(&models.TestRun{}).Where("id = ?", testRunID).Updates(map[string]interface{}{
				"status":          "failed",
				"passed_tests":    0,
				"failed_tests":    0,
				"execution_time_ms":  0,
				"completed_at":    completedAt,
			})
		}
	}()

	passedTests := 0
	failedTests := 0

	fmt.Printf("Starting test execution for test run %s with %d test cases\n", testRunID, len(testCases))

	// Handle case where no test cases are found
	if len(testCases) == 0 {
		fmt.Printf("No test cases found for test run %s\n", testRunID)
		completedAt := time.Now()
		s.db.Model(&models.TestRun{}).Where("id = ?", testRunID).Updates(map[string]interface{}{
			"status":          "failed",
			"passed_tests":    0,
			"failed_tests":    0,
			"execution_time_ms":  0,
			"completed_at":    completedAt,
		})
		return
	}

	for i, testCase := range testCases {
		fmt.Printf("Executing test case %d/%d: %s\n", i+1, len(testCases), testCase.ID)
		
		// Parse test spec
		var testSpec models.TestSpec
		if err := json.Unmarshal([]byte(testCase.TestSpec), &testSpec); err != nil {
			fmt.Printf("Failed to parse test spec for test case %s: %v\n", testCase.ID, err)
			s.recordTestResult(testRunID, testCase.ID, "failed", 0, err.Error(), "")
			failedTests++
			continue
		}

		// Create test executor for this service
		executor := testrunner.NewHTTPExpectExecutor(testCase.Service.BaseURL)
		
		// Execute test
		result := executor.ExecuteTest(&testSpec)
		
		// Record result
		status := "passed"
		if result.Status == "FAILED" {
			status = "failed"
			failedTests++
		} else {
			passedTests++
		}

		fmt.Printf("Test case %s result: %s\n", testCase.ID, status)
		s.recordTestResult(testRunID, testCase.ID, status, int(result.Duration.Milliseconds()), result.ErrorMessage, result.ResponseData)
	}

	// Update test run status
	completedAt := time.Now()
	
	// Calculate execution time only if there are test cases
	var executionTime int64
	if len(testCases) > 0 {
		executionTime = completedAt.Sub(testCases[0].CreatedAt).Milliseconds()
	}
	
	status := "completed"
	if failedTests > 0 {
		status = "failed"
	}

	fmt.Printf("Completing test run %s: %s (passed: %d, failed: %d)\n", testRunID, status, passedTests, failedTests)

	s.db.Model(&models.TestRun{}).Where("id = ?", testRunID).Updates(map[string]interface{}{
		"status":          status,
		"passed_tests":    passedTests,
		"failed_tests":    failedTests,
		"execution_time_ms":  executionTime,
		"completed_at":    completedAt,
	})
}

// recordTestResult records a single test result
func (s *TestRunService) recordTestResult(testRunID, testCaseID, status string, executionTime int, errorMessage, responseData string) {
	// Ensure responseData is valid JSON for JSONB column
	if responseData == "" {
		responseData = "{}"
	}
	
	testResult := &models.TestResult{
		TestRunID:     testRunID,
		TestCaseID:    testCaseID,
		Status:        status,
		ExecutionTimeMs: executionTime,
		ErrorMessage:  errorMessage,
		ResponseData:  responseData,
	}

	s.db.Create(testResult)
}

// GetTestRun retrieves a test run by ID
func (s *TestRunService) GetTestRun(id string) (*models.TestRun, error) {
	var testRun models.TestRun
	if err := s.db.Preload("TestResults.TestCase").First(&testRun, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &testRun, nil
}

// GetTestResults retrieves test results for a test run
func (s *TestRunService) GetTestResults(testRunID string) ([]models.TestResult, error) {
	var testResults []models.TestResult
	err := s.db.Preload("TestCase").Where("test_run_id = ?", testRunID).Find(&testResults).Error
	return testResults, err
}

// ListTestRuns retrieves all test runs with pagination
func (s *TestRunService) ListTestRuns(limit, offset int) ([]models.TestRun, int64, error) {
	var testRuns []models.TestRun
	var total int64

	// Get total count
	if err := s.db.Model(&models.TestRun{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := s.db.Order("started_at DESC").Limit(limit).Offset(offset).Find(&testRuns).Error; err != nil {
		return nil, 0, err
	}

	return testRuns, total, nil
}
