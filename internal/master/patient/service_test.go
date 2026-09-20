package patient_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"hosim-go/internal/master/patient"
)

// 1. Buat Mock Repository (tanpa database asli)
type mockPatientRepo struct {
	patient.Repository
	findByNIKFn            func(ctx context.Context, nik string) (*patient.Patient, error)
	findByIDFn             func(ctx context.Context, id string) (*patient.Patient, error)
	getLastMedicalRecordFn func(ctx context.Context) (string, error)
	createFn               func(ctx context.Context, p *patient.Patient) (*patient.Patient, error)
	updateFn               func(ctx context.Context, p *patient.Patient) (*patient.Patient, error)
}

func (m *mockPatientRepo) FindByNIK(ctx context.Context, nik string) (*patient.Patient, error) {
	if m.findByNIKFn != nil {
		return m.findByNIKFn(ctx, nik)
	}
	return nil, nil
}

func (m *mockPatientRepo) FindByID(ctx context.Context, id string) (*patient.Patient, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockPatientRepo) GetLastMedicalRecordNo(ctx context.Context) (string, error) {
	if m.getLastMedicalRecordFn != nil {
		return m.getLastMedicalRecordFn(ctx)
	}
	return "", nil
}

func (m *mockPatientRepo) Create(ctx context.Context, p *patient.Patient) (*patient.Patient, error) {
	if m.createFn != nil {
		return m.createFn(ctx, p)
	}
	return p, nil
}

func (m *mockPatientRepo) Update(ctx context.Context, p *patient.Patient) (*patient.Patient, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, p)
	}
	return p, nil
}

// 2. Unit Test Skenario
func TestRegisterPatient_NIKDuplikat(t *testing.T) {
	// Setup mock: pura-pura NIK sudah ada di database
	mockRepo := &mockPatientRepo{
		findByNIKFn: func(ctx context.Context, nik string) (*patient.Patient, error) {
			return &patient.Patient{NIK: "3201234567890001", FullName: "Pasien Lama"}, nil
		},
	}

	svc := patient.NewService(mockRepo)

	req := patient.CreatePatientRequest{
		NIK:       "3201234567890001",
		FullName:  "Pasien Baru",
		BirthDate: "2000-01-01",
	}

	// Eksekusi
	res, err := svc.RegisterPatient(context.Background(), req, "operator-1")

	// Verifikasi: harus error karena NIK sudah ada
	if err == nil {
		t.Errorf("diharapkan error karena NIK duplikat, tapi err bernilai nil")
	}
	if res != nil {
		t.Errorf("diharapkan hasil nil saat error, tapi mendapat objek pasien")
	}
	if err.Error() != "pasien dengan NIK tersebut sudah terdaftar" {
		t.Errorf("pesan error tidak sesuai: %v", err)
	}
}

func TestRegisterPatient_GenerateNoRMSukses(t *testing.T) {
	// Setup mock: no RM terakhir di DB adalah 00000005
	mockRepo := &mockPatientRepo{
		findByNIKFn: func(ctx context.Context, nik string) (*patient.Patient, error) {
			return nil, nil // NIK belum terdaftar
		},
		getLastMedicalRecordFn: func(ctx context.Context) (string, error) {
			return "00000005", nil
		},
		createFn: func(ctx context.Context, p *patient.Patient) (*patient.Patient, error) {
			return p, nil
		},
	}

	svc := patient.NewService(mockRepo)

	req := patient.CreatePatientRequest{
		MedicalRecordNo: "", // sengaja dikosongkan agar generate otomatis
		NIK:             "3201234567890002",
		FullName:        "Budi Santoso",
		Gender:          "L",
		BirthDate:       "1995-05-20",
	}

	res, err := svc.RegisterPatient(context.Background(), req, "operator-1")

	if err != nil {
		t.Fatalf("seharusnya tidak error: %v", err)
	}

	// Verifikasi: no RM otomatis harus naik jadi 00000006
	expectedRM := "00000006"
	if res.MedicalRecordNo != expectedRM {
		t.Errorf("diharapkan no RM %s, tapi didapat %s", expectedRM, res.MedicalRecordNo)
	}
}

func TestRegisterPatient_NewFieldsMapping(t *testing.T) {
	mockRepo := &mockPatientRepo{
		findByNIKFn: func(ctx context.Context, nik string) (*patient.Patient, error) {
			return nil, nil
		},
		getLastMedicalRecordFn: func(ctx context.Context) (string, error) {
			return "00000001", nil
		},
		createFn: func(ctx context.Context, p *patient.Patient) (*patient.Patient, error) {
			return p, nil
		},
	}

	svc := patient.NewService(mockRepo)

	req := patient.CreatePatientRequest{
		Title:        "Tn.",
		ShortName:    "Budi",
		FullName:     "Budi Hartono",
		MotherName:   "Siti Rahayu",
		FamilyCardNo: "3201234567890009",
		Gender:       "L",
		BirthPlace:   "Jakarta",
		BirthDate:    "1990-12-15",
		Phone:        "08123456789",
		BloodType:    "O",
		Rhesus:       "+",
		SpecialNeeds: "Kursi Roda",
		IsUnknown:    false,
	}

	res, err := svc.RegisterPatient(context.Background(), req, "operator-1")
	if err != nil {
		t.Fatalf("seharusnya tidak error: %v", err)
	}

	if res.Title != "Tn." {
		t.Errorf("expected Title 'Tn.', got %s", res.Title)
	}
	if res.MotherName != "Siti Rahayu" {
		t.Errorf("expected MotherName 'Siti Rahayu', got %s", res.MotherName)
	}
	if res.FamilyCardNo != "3201234567890009" {
		t.Errorf("expected FamilyCardNo '3201234567890009', got %s", res.FamilyCardNo)
	}
	if res.Rhesus != "+" {
		t.Errorf("expected Rhesus '+', got %s", res.Rhesus)
	}
	if res.SpecialNeeds != "Kursi Roda" {
		t.Errorf("expected SpecialNeeds 'Kursi Roda', got %s", res.SpecialNeeds)
	}
}

func TestUpdatePatient_NotFound(t *testing.T) {
	mockRepo := &mockPatientRepo{
		findByIDFn: func(ctx context.Context, id string) (*patient.Patient, error) {
			return nil, nil // Pasien tidak ada
		},
	}

	svc := patient.NewService(mockRepo)
	_, err := svc.UpdatePatient(context.Background(), "unknown-id", patient.UpdatePatientRequest{}, "operator-1")

	if !errors.Is(err, patient.ErrPatientNotFound) {
		t.Errorf("expected ErrPatientNotFound, got %v", err)
	}
}

func TestUpdatePatient_NIKDuplikat(t *testing.T) {
	existingPatient := &patient.Patient{
		ID:        "patient-1",
		NIK:       "3201000000000001",
		FullName:  "Pasien Satu",
		BirthDate: time.Now(),
	}

	mockRepo := &mockPatientRepo{
		findByIDFn: func(ctx context.Context, id string) (*patient.Patient, error) {
			return existingPatient, nil
		},
		findByNIKFn: func(ctx context.Context, nik string) (*patient.Patient, error) {
			// NIK ini sudah dipakai oleh pasien lain (patient-2)
			return &patient.Patient{ID: "patient-2", NIK: nik}, nil
		},
	}

	svc := patient.NewService(mockRepo)
	req := patient.UpdatePatientRequest{
		CreatePatientRequest: patient.CreatePatientRequest{
			NIK:       "3201000000000002", // ganti ke NIK milik patient-2
			FullName:  "Pasien Satu Diedit",
			BirthDate: "1990-01-01",
		},
	}

	_, err := svc.UpdatePatient(context.Background(), "patient-1", req, "operator-1")
	if !errors.Is(err, patient.ErrNIKAlreadyExists) {
		t.Errorf("expected ErrNIKAlreadyExists, got %v", err)
	}
}

func TestRegisterPatient_BirthDateInFuture(t *testing.T) {
	mockRepo := &mockPatientRepo{}
	svc := patient.NewService(mockRepo)

	req := patient.CreatePatientRequest{
		FullName:  "Pasien Masa Depan",
		BirthDate: "2099-01-01",
	}

	_, err := svc.RegisterPatient(context.Background(), req, "operator-1")
	if !errors.Is(err, patient.ErrBirthDateInFuture) {
		t.Errorf("expected ErrBirthDateInFuture, got %v", err)
	}
}

func TestRegisterPatient_Sanitization(t *testing.T) {
	var capturedNIK, capturedName string
	mockRepo := &mockPatientRepo{
		findByNIKFn: func(ctx context.Context, nik string) (*patient.Patient, error) {
			capturedNIK = nik
			return nil, nil
		},
		getLastMedicalRecordFn: func(ctx context.Context) (string, error) {
			return "00000001", nil
		},
		createFn: func(ctx context.Context, p *patient.Patient) (*patient.Patient, error) {
			capturedName = p.FullName
			return p, nil
		},
	}

	svc := patient.NewService(mockRepo)
	req := patient.CreatePatientRequest{
		NIK:       "  3201000000000001  ",
		FullName:  "  Budi Santoso  ",
		Gender:    "male",
		BirthDate: "1995-01-01",
	}

	_, err := svc.RegisterPatient(context.Background(), req, "operator-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedNIK != "3201000000000001" {
		t.Errorf("expected trimmed NIK, got '%s'", capturedNIK)
	}
	if capturedName != "Budi Santoso" {
		t.Errorf("expected trimmed Name, got '%s'", capturedName)
	}
}

func TestRegisterPatient_WithRelations(t *testing.T) {
	var capturedPatient *patient.Patient
	mockRepo := &mockPatientRepo{
		findByNIKFn: func(ctx context.Context, nik string) (*patient.Patient, error) {
			return nil, nil
		},
		getLastMedicalRecordFn: func(ctx context.Context) (string, error) {
			return "00000001", nil
		},
		createFn: func(ctx context.Context, p *patient.Patient) (*patient.Patient, error) {
			capturedPatient = p
			return p, nil
		},
	}

	svc := patient.NewService(mockRepo)
	req := patient.CreatePatientRequest{
		FullName:  "Pasien Dengan Keluarga",
		ShortName: "Pasien",
		Gender:    "L",
		BirthDate: "1990-01-01",
		Phone:     "0812345678",
		Relations: []patient.RelationRequest{
			{Name: "Siti Rahma", Relation: "Istri", Phone: "08111111", Address: "Jl Mawar"},
		},
	}

	res, err := svc.RegisterPatient(context.Background(), req, "operator-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(res.Relations))
	}
	if res.Relations[0].Name != "Siti Rahma" || res.Relations[0].Relation != "Istri" {
		t.Errorf("unexpected relation data: %+v", res.Relations[0])
	}
	if capturedPatient == nil || len(capturedPatient.Relations) != 1 {
		t.Errorf("expected repo create to receive relations")
	}
}

func TestRegisterPatient_InvalidManualMedicalRecordFormat(t *testing.T) {
	mockRepo := &mockPatientRepo{}
	svc := patient.NewService(mockRepo)

	req := patient.CreatePatientRequest{
		MedicalRecordNo: "RM-123", // bukan 8 digit angka
		FullName:        "Pasien RM Invalid",
		ShortName:       "Pasien",
		Gender:          "L",
		BirthDate:       "1990-01-01",
		Phone:           "0812345678",
	}

	_, err := svc.RegisterPatient(context.Background(), req, "operator-1")
	if !errors.Is(err, patient.ErrInvalidMedicalRecordFormat) {
		t.Errorf("expected ErrInvalidMedicalRecordFormat, got %v", err)
	}
}

func TestRegisterPatient_InvalidInsuranceExpiryDateFormat(t *testing.T) {
	mockRepo := &mockPatientRepo{}
	svc := patient.NewService(mockRepo)

	invalidDate := "31-12-2025" // format salah
	req := patient.CreatePatientRequest{
		FullName:            "Pasien Asuransi Format Salah",
		ShortName:           "Pasien",
		Gender:              "L",
		BirthDate:           "1990-01-01",
		Phone:               "0812345678",
		InsuranceExpiryDate: &invalidDate,
	}

	_, err := svc.RegisterPatient(context.Background(), req, "operator-1")
	if !errors.Is(err, patient.ErrInvalidInsuranceExpiryDateFormat) {
		t.Errorf("expected ErrInvalidInsuranceExpiryDateFormat, got %v", err)
	}
}

func TestRegisterPatient_InvalidDeceasedDateFormat(t *testing.T) {
	mockRepo := &mockPatientRepo{}
	svc := patient.NewService(mockRepo)

	invalidDate := "2025/12/31 10:00" // format salah
	req := patient.CreatePatientRequest{
		FullName:   "Pasien Meninggal Format Salah",
		ShortName:  "Pasien",
		Gender:     "L",
		BirthDate:  "1990-01-01",
		Phone:      "0812345678",
		IsDeceased: true,
		DeceasedAt: &invalidDate,
	}

	_, err := svc.RegisterPatient(context.Background(), req, "operator-1")
	if !errors.Is(err, patient.ErrInvalidDeceasedDate) {
		t.Errorf("expected ErrInvalidDeceasedDate, got %v", err)
	}
}

func TestUpdatePatient_ResetDeceasedWhenNotDeceased(t *testing.T) {
	pastDeceased := time.Now().Add(-24 * time.Hour)
	existingPatient := &patient.Patient{
		ID:         "patient-100",
		FullName:   "Pasien Salah Input Meninggal",
		BirthDate:  time.Now().AddDate(-30, 0, 0),
		IsDeceased: true,
		DeceasedAt: &pastDeceased,
	}

	var savedPatient *patient.Patient
	mockRepo := &mockPatientRepo{
		findByIDFn: func(ctx context.Context, id string) (*patient.Patient, error) {
			return existingPatient, nil
		},
		updateFn: func(ctx context.Context, p *patient.Patient) (*patient.Patient, error) {
			savedPatient = p
			return p, nil
		},
	}

	svc := patient.NewService(mockRepo)
	req := patient.UpdatePatientRequest{
		CreatePatientRequest: patient.CreatePatientRequest{
			FullName:   "Pasien Salah Input Meninggal",
			Gender:     "L",
			BirthDate:  "1995-01-01",
			Phone:      "08123456",
			IsDeceased: false, // Dikoreksi menjadi hidup
		},
	}

	res, err := svc.UpdatePatient(context.Background(), "patient-100", req, "operator-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsDeceased {
		t.Errorf("expected IsDeceased to be false")
	}
	if res.DeceasedAt != nil {
		t.Errorf("expected DeceasedAt to be nil when IsDeceased is false, got %v", res.DeceasedAt)
	}
	if savedPatient.DeceasedAt != nil {
		t.Errorf("expected saved patient DeceasedAt to be nil")
	}
}



