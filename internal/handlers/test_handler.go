package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-test-framework/internal/models"
	"api-test-framework/internal/services"
	"api-test-framework/internal/utils"

	"github.com/gin-gonic/gin"
)

// TestHandler handles test-related HTTP requests
type TestHandler struct {
	testService *services.TestService
}

// NewTestHandler creates a new test handler
func NewTestHandler(testService *services.TestService) *TestHandler {
	return &TestHandler{testService: testService}
}

// ListTests handles GET /api/v1/tests
func (h *TestHandler) ListTests(c *gin.Context) {
	serviceID := c.Query("service_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tests, total, err := h.testService.ListTests(serviceID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve tests",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tests,
		"meta": gin.H{
			"total": total,
			"limit": limit,
			"offset": offset,
		},
	})
}

// CreateTest handles POST /api/v1/tests
func (h *TestHandler) CreateTest(c *gin.Context) {
	var testCase models.TestCase
	if err := c.ShouldBindJSON(&testCase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.testService.CreateTest(&testCase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create test",
			"details": err.Error(),
		})
		return
	}

	// Fetch the created test case with service preloaded
	createdTestCase, err := h.testService.GetTest(testCase.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve created test",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": createdTestCase,
	})
}

// CreateTestFromCurl handles POST /api/v1/tests/from-curl
func (h *TestHandler) CreateTestFromCurl(c *gin.Context) {
	var request struct {
		ServiceID   string `json:"service_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		CurlCommand string `json:"curl_command" binding:"required"`
		Assertions  []map[string]interface{} `json:"assertions"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Parse the curl command
	curlRequest, err := utils.ParseCurlCommand(request.CurlCommand)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid curl command",
			"details": err.Error(),
		})
		return
	}

	// Convert to test spec
	testSpec := curlRequest.ToTestSpec(request.Name, request.Description)
	
	// Override service_name with the actual service
	testSpec["service_name"] = "service-" + request.ServiceID

	// Add custom assertions if provided
	if len(request.Assertions) > 0 {
		testSpec["assertions"] = request.Assertions
	}

	// Convert test spec to JSON string
	testSpecJSON, err := json.Marshal(testSpec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to marshal test spec",
			"details": err.Error(),
		})
		return
	}

	// Create test case
	testCase := &models.TestCase{
		ServiceID:   request.ServiceID,
		Name:        request.Name,
		Description: request.Description,
		TestSpec:    string(testSpecJSON),
	}

	if err := h.testService.CreateTest(testCase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create test",
			"details": err.Error(),
		})
		return
	}

	// Fetch the created test case with service preloaded
	createdTestCase, err := h.testService.GetTest(testCase.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve created test",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": createdTestCase,
		"parsed_curl": curlRequest,
	})
}

// GetTest handles GET /api/v1/tests/:id
func (h *TestHandler) GetTest(c *gin.Context) {
	id := c.Param("id")

	testCase, err := h.testService.GetTest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Test not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": testCase,
	})
}

// UpdateTest handles PUT /api/v1/tests/:id
func (h *TestHandler) UpdateTest(c *gin.Context) {
	id := c.Param("id")

	var testCase models.TestCase
	if err := c.ShouldBindJSON(&testCase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	updatedTestCase, err := h.testService.UpdateTest(id, &testCase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update test",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedTestCase,
	})
}

// DeleteTest handles DELETE /api/v1/tests/:id
func (h *TestHandler) DeleteTest(c *gin.Context) {
	id := c.Param("id")

	if err := h.testService.DeleteTest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete test",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test deleted successfully",
	})
}
