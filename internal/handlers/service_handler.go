package handlers

import (
	"net/http"
	"strconv"

	"api-test-framework/internal/models"
	"api-test-framework/internal/services"

	"github.com/gin-gonic/gin"
)

// ServiceHandler handles service-related HTTP requests
type ServiceHandler struct {
	serviceService *services.ServiceService
}

// NewServiceHandler creates a new service handler
func NewServiceHandler(serviceService *services.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceService: serviceService}
}

// ListServices handles GET /api/v1/services
func (h *ServiceHandler) ListServices(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	services, total, err := h.serviceService.ListServices(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve services",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": services,
		"meta": gin.H{
			"total": total,
			"limit": limit,
			"offset": offset,
		},
	})
}

// CreateService handles POST /api/v1/services
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.serviceService.CreateService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create service",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": service,
	})
}

// GetService handles GET /api/v1/services/:id
func (h *ServiceHandler) GetService(c *gin.Context) {
	id := c.Param("id")

	service, err := h.serviceService.GetService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Service not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": service,
	})
}

// UpdateService handles PUT /api/v1/services/:id
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id := c.Param("id")

	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	updatedService, err := h.serviceService.UpdateService(id, &service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update service",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedService,
	})
}

// DeleteService handles DELETE /api/v1/services/:id
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")

	if err := h.serviceService.DeleteService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete service",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service deleted successfully",
	})
}
