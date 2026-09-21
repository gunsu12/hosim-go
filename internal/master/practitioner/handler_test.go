package practitioner_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"hosim-go/internal/master/practitioner"

	"github.com/gin-gonic/gin"
)

type mockPractitionerService struct {
	practitioner.Service
	createFn  func(ctx context.Context, req practitioner.CreatePractitionerRequest, op string) (*practitioner.Practitioner, error)
	getByIDFn func(ctx context.Context, id string) (*practitioner.Practitioner, error)
}

func (m *mockPractitionerService) CreatePractitioner(ctx context.Context, req practitioner.CreatePractitionerRequest, op string) (*practitioner.Practitioner, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req, op)
	}
	return nil, nil
}

func (m *mockPractitionerService) GetPractitionerByID(ctx context.Context, id string) (*practitioner.Practitioner, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func TestHandler_CreatePractitioner_ValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := practitioner.NewHandler(&mockPractitionerService{})
	router := gin.New()
	router.POST("/api/v1/practitioners", handler.Create)

	// Kurang field Name dan Gender yang binding:required
	invalidJSON := `{"nik": "3201234567890001"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/practitioners", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandler_CreatePractitioner_DuplicateConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockPractitionerService{
		createFn: func(ctx context.Context, req practitioner.CreatePractitionerRequest, op string) (*practitioner.Practitioner, error) {
			return nil, practitioner.ErrNIKAlreadyExists
		},
	}
	handler := practitioner.NewHandler(mockSvc)
	router := gin.New()
	router.POST("/api/v1/practitioners", handler.Create)

	body := `{"name": "dr. Budi", "gender": "L"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/practitioners", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected status %d (Conflict), got %d", http.StatusConflict, rec.Code)
	}
}

func TestHandler_GetPractitionerByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockPractitionerService{
		getByIDFn: func(ctx context.Context, id string) (*practitioner.Practitioner, error) {
			return nil, practitioner.ErrPractitionerNotFound
		},
	}
	handler := practitioner.NewHandler(mockSvc)
	router := gin.New()
	router.GET("/api/v1/practitioners/:id", handler.GetByID)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/practitioners/not-exist", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d (NotFound), got %d", http.StatusNotFound, rec.Code)
	}
}
