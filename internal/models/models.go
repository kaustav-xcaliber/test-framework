package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthConfig represents authentication configuration for a service
type AuthConfig struct {
	Type         string            `json:"type"`         // "bearer", "api_key", "basic", "oauth2"
	Token        string            `json:"token,omitempty"`
	KeyName      string            `json:"key_name,omitempty"`
	KeyValue     string            `json:"key_value,omitempty"`
	Username     string            `json:"username,omitempty"`
	Password     string            `json:"password,omitempty"`
	ClientID     string            `json:"client_id,omitempty"`
	ClientSecret string            `json:"client_secret,omitempty"`
	TokenURL     string            `json:"token_url,omitempty"`
	Extra        map[string]string `json:"extra,omitempty"`
}

// Value implements driver.Valuer interface
func (a AuthConfig) Value() (driver.Value, error) {
	if a.Type == "" {
		return "{}", nil
	}
	return json.Marshal(a)
}

// Scan implements sql.Scanner interface
func (a *AuthConfig) Scan(value interface{}) error {
	if value == nil {
		*a = AuthConfig{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			*a = AuthConfig{}
			return nil
		}
		return json.Unmarshal(v, a)
	case string:
		if v == "" {
			*a = AuthConfig{}
			return nil
		}
		return json.Unmarshal([]byte(v), a)
	default:
		*a = AuthConfig{}
		return nil
	}
}

// Service represents a microservice that can be tested
type Service struct {
	ID          string     `json:"id" gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	Name        string     `json:"name" gorm:"uniqueIndex;not null"`
	Description string     `json:"description"`
	BaseURL     string     `json:"base_url" gorm:"not null"`
	AuthConfig  AuthConfig `json:"auth_config" gorm:"type:jsonb;default:'{}'"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
}

// TestCase represents a test case for a service
type TestCase struct {
	ID          string    `json:"id" gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	ServiceID   string    `json:"service_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	TestSpec    string    `json:"test_spec" gorm:"type:jsonb;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Service     Service   `json:"service" gorm:"foreignKey:ServiceID;references:ID"`
}

// TestRun represents a test execution run
type TestRun struct {
	ID             string        `json:"id" gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	Name           string        `json:"name"`
	Status         string        `json:"status" gorm:"default:'running';check:status IN ('running', 'completed', 'failed')"`
	TotalTests     int           `json:"total_tests" gorm:"default:0"`
	PassedTests    int           `json:"passed_tests" gorm:"default:0"`
	FailedTests    int           `json:"failed_tests" gorm:"default:0"`
	ExecutionTimeMs int64        `json:"execution_time_ms" gorm:"default:0"`
	StartedAt      time.Time     `json:"started_at" gorm:"autoCreateTime"`
	CompletedAt    *time.Time    `json:"completed_at"`
	TestResults    []TestResult  `json:"test_results" gorm:"foreignKey:TestRunID"`
}

// TestResult represents the result of a single test execution
type TestResult struct {
	ID             string    `json:"id" gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	TestRunID      string    `json:"test_run_id" gorm:"not null"`
	TestCaseID     string    `json:"test_case_id" gorm:"not null"`
	Status         string    `json:"status" gorm:"not null;check:status IN ('passed', 'failed', 'skipped')"`
	ExecutionTimeMs int      `json:"execution_time_ms" gorm:"default:0"`
	ErrorMessage   string    `json:"error_message"`
	ResponseData   string    `json:"response_data" gorm:"type:jsonb"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	TestCase       TestCase  `json:"test_case" gorm:"foreignKey:TestCaseID;references:ID"`
}

// TestSpec represents the specification for a test case
type TestSpec struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ServiceName string            `json:"service_name"`
	Request     RequestSpec       `json:"request"`
	Assertions  []AssertionSpec   `json:"assertions"`
}

// RequestSpec represents the HTTP request specification
type RequestSpec struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

// AssertionSpec represents a single assertion to validate
type AssertionSpec struct {
	Type     string      `json:"type"`
	Path     string      `json:"path,omitempty"`
	Matcher  string      `json:"matcher,omitempty"`
	Expected interface{} `json:"expected"`
}

// BeforeCreate hooks for GORM
func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

func (s *Service) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

func (tc *TestCase) BeforeCreate(tx *gorm.DB) error {
	if tc.ID == "" {
		tc.ID = uuid.New().String()
	}
	return nil
}

func (tc *TestCase) BeforeUpdate(tx *gorm.DB) error {
	tc.UpdatedAt = time.Now()
	return nil
}

func (tr *TestRun) BeforeCreate(tx *gorm.DB) error {
	if tr.ID == "" {
		tr.ID = uuid.New().String()
	}
	return nil
}

func (tr *TestResult) BeforeCreate(tx *gorm.DB) error {
	if tr.ID == "" {
		tr.ID = uuid.New().String()
	}
	return nil
}
