package patient_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"hosim-go/internal/master/patient"

	"github.com/gin-gonic/gin"
)

// Mock Service untuk Handler
type mockPatientService struct {
	patient.Service
}

func TestHandler_CreatePatient_ValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Inisialisasi handler dengan mock service kosong
	handler := patient.NewHandler(&mockPatientService{})

	router := gin.New()
	router.POST("/api/v1/patients", handler.Create)

	// Body JSON tidak lengkap (kurang full_name, birth_date, dll yang binding:required)
	invalidJSON := `{"nik": "3201234567890001"}`

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/patients", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	// Perekam response HTTP
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verifikasi: Harusnya 400 Bad Request karena validasi binding JSON gagal
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("diharapkan status %d, tapi didapat %d", http.StatusBadRequest, recorder.Code)
	}
}

type mockPatientServiceForUpdate struct {
	patient.Service
}

func (m *mockPatientServiceForUpdate) UpdatePatient(ctx context.Context, id string, req patient.UpdatePatientRequest, op string) (*patient.Patient, error) {
	return nil, patient.ErrPatientNotFound
}

func TestHandler_UpdatePatient_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Inisialisasi mock service yang return ErrPatientNotFound
	handler := patient.NewHandler(&mockPatientServiceForUpdate{})

	router := gin.New()
	router.PUT("/api/v1/patients/:id", handler.Update)

	validJSON := `{"short_name": "Budi", "full_name": "Budi Santoso", "gender": "L", "birth_place": "Jkt", "birth_date": "1990-01-01", "phone": "0812345"}`

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/patients/patient-999", bytes.NewBufferString(validJSON))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verifikasi: Harusnya 404 Not Found, bukan 500
	if recorder.Code != http.StatusNotFound {
		t.Errorf("diharapkan status %d, tapi didapat %d", http.StatusNotFound, recorder.Code)
	}
}

