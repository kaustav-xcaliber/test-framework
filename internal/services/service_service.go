package services

import (
	"api-test-framework/internal/models"

	"gorm.io/gorm"
)

// ServiceService handles service operations
type ServiceService struct {
	db *gorm.DB
}

// NewServiceService creates a new service service
func NewServiceService(db *gorm.DB) *ServiceService {
	return &ServiceService{db: db}
}

// CreateService creates a new service
func (s *ServiceService) CreateService(service *models.Service) error {
	return s.db.Create(service).Error
}

// GetService retrieves a service by ID
func (s *ServiceService) GetService(id string) (*models.Service, error) {
	var service models.Service
	if err := s.db.First(&service, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

// ListServices retrieves all services with pagination
func (s *ServiceService) ListServices(limit, offset int) ([]models.Service, int64, error) {
	var services []models.Service
	var total int64

	// Get total count
	if err := s.db.Model(&models.Service{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := s.db.Limit(limit).Offset(offset).Find(&services).Error; err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

// UpdateService updates an existing service and returns the updated service
func (s *ServiceService) UpdateService(id string, service *models.Service) (*models.Service, error) {
	// First, get the existing service to preserve the ID
	var existingService models.Service
	if err := s.db.First(&existingService, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Update the service with the new data
	if err := s.db.Model(&existingService).Updates(service).Error; err != nil {
		return nil, err
	}

	// Fetch the updated service to get the latest data including timestamps
	var updatedService models.Service
	if err := s.db.First(&updatedService, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &updatedService, nil
}

// DeleteService deletes a service
func (s *ServiceService) DeleteService(id string) error {
	return s.db.Delete(&models.Service{}, "id = ?", id).Error
}
