package services

import (
	"encoding/json"
	"fmt"

	"api-test-framework/internal/models"

	"gorm.io/gorm"
)

// TestService handles test case operations
type TestService struct {
	db *gorm.DB
}

// NewTestService creates a new test service
func NewTestService(db *gorm.DB) *TestService {
	return &TestService{db: db}
}

// CreateTest creates a new test case
func (s *TestService) CreateTest(testCase *models.TestCase) error {
	// Validate test spec JSON
	var testSpec models.TestSpec
	if err := json.Unmarshal([]byte(testCase.TestSpec), &testSpec); err != nil {
		return fmt.Errorf("invalid test spec JSON: %v", err)
	}

	// Check if service exists
	var service models.Service
	if err := s.db.First(&service, "id = ?", testCase.ServiceID).Error; err != nil {
		return fmt.Errorf("service not found: %v", err)
	}

	return s.db.Create(testCase).Error
}

// GetTest retrieves a test case by ID
func (s *TestService) GetTest(id string) (*models.TestCase, error) {
	var testCase models.TestCase
	if err := s.db.Preload("Service").First(&testCase, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &testCase, nil
}

// ListTests retrieves all test cases with optional filtering
func (s *TestService) ListTests(serviceID string, limit, offset int) ([]models.TestCase, int64, error) {
	var testCases []models.TestCase
	var total int64

	query := s.db.Model(&models.TestCase{}).Preload("Service")
	
	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Limit(limit).Offset(offset).Find(&testCases).Error; err != nil {
		return nil, 0, err
	}

	return testCases, total, nil
}

// UpdateTest updates an existing test case and returns the updated test case
func (s *TestService) UpdateTest(id string, testCase *models.TestCase) (*models.TestCase, error) {
	// Validate test spec JSON
	var testSpec models.TestSpec
	if err := json.Unmarshal([]byte(testCase.TestSpec), &testSpec); err != nil {
		return nil, fmt.Errorf("invalid test spec JSON: %v", err)
	}

	// First, get the existing test case to preserve the ID
	var existingTestCase models.TestCase
	if err := s.db.First(&existingTestCase, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Update the test case with the new data
	if err := s.db.Model(&existingTestCase).Updates(testCase).Error; err != nil {
		return nil, err
	}

	// Fetch the updated test case to get the latest data including timestamps
	var updatedTestCase models.TestCase
	if err := s.db.Preload("Service").First(&updatedTestCase, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &updatedTestCase, nil
}

// DeleteTest deletes a test case
func (s *TestService) DeleteTest(id string) error {
	return s.db.Delete(&models.TestCase{}, "id = ?", id).Error
}

// GetTestsByService retrieves all test cases for a specific service
func (s *TestService) GetTestsByService(serviceID string) ([]models.TestCase, error) {
	var testCases []models.TestCase
	err := s.db.Where("service_id = ?", serviceID).Find(&testCases).Error
	return testCases, err
}

// CleanTestSpec removes header assertions from a test spec
func (s *TestService) CleanTestSpec(testSpecJSON string) (string, error) {
	var testSpec map[string]interface{}
	if err := json.Unmarshal([]byte(testSpecJSON), &testSpec); err != nil {
		return testSpecJSON, err
	}
	
	assertions, ok := testSpec["assertions"].([]interface{})
	if !ok {
		return testSpecJSON, nil
	}
	
	var cleanedAssertions []interface{}
	for _, assertion := range assertions {
		if assertionMap, ok := assertion.(map[string]interface{}); ok {
			if path, ok := assertionMap["path"].(string); ok {
				// Skip header assertions
				if path == "headers.Content-Type" || path == "headers" {
					continue
				}
			}
		}
		cleanedAssertions = append(cleanedAssertions, assertion)
	}
	
	testSpec["assertions"] = cleanedAssertions
	
	cleanedJSON, err := json.Marshal(testSpec)
	if err != nil {
		return testSpecJSON, err
	}
	
	return string(cleanedJSON), nil
}
