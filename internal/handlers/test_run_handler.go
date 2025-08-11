package handlers

import (
	"net/http"
	"strconv"

	"api-test-framework/internal/services"

	"github.com/gin-gonic/gin"
)

// TestRunHandler handles test run-related HTTP requests
type TestRunHandler struct {
	testRunService *services.TestRunService
}

// NewTestRunHandler creates a new test run handler
func NewTestRunHandler(testRunService *services.TestRunService) *TestRunHandler {
	return &TestRunHandler{testRunService: testRunService}
}

// StartTestRun handles POST /api/v1/test-runs
func (h *TestRunHandler) StartTestRun(c *gin.Context) {
	var request struct {
		ServiceID string   `json:"service_id"`
		TestIDs   []string `json:"test_ids"`
		Name      string   `json:"name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	testRun, err := h.testRunService.StartTestRun(request.ServiceID, request.TestIDs, request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start test run",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": testRun,
	})
}

// GetTestRun handles GET /api/v1/test-runs/:id
func (h *TestRunHandler) GetTestRun(c *gin.Context) {
	id := c.Param("id")

	testRun, err := h.testRunService.GetTestRun(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Test run not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": testRun,
	})
}

// GetTestResults handles GET /api/v1/test-runs/:id/results
func (h *TestRunHandler) GetTestResults(c *gin.Context) {
	id := c.Param("id")

	testResults, err := h.testRunService.GetTestResults(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve test results",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": testResults,
	})
}

// ListTestRuns handles GET /api/v1/test-runs
func (h *TestRunHandler) ListTestRuns(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	testRuns, total, err := h.testRunService.ListTestRuns(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve test runs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": testRuns,
		"meta": gin.H{
			"total": total,
			"limit": limit,
			"offset": offset,
		},
	})
}
