package practitioner_test

import (
	"context"
	"errors"
	"testing"

	"hosim-go/internal/practitioner"
	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

type mockPractitionerRepo struct {
	practitioner.Repository
	createFn                  func(ctx context.Context, p *practitioner.Practitioner) (*practitioner.Practitioner, error)
	updateFn                  func(ctx context.Context, p *practitioner.Practitioner) (*practitioner.Practitioner, error)
	deleteFn                  func(ctx context.Context, id string, deletedBy string) error
	findByIDFn                func(ctx context.Context, id string) (*practitioner.Practitioner, error)
	findByNIKFn               func(ctx context.Context, nik string) (*practitioner.Practitioner, error)
	findByNIPFn               func(ctx context.Context, nip string) (*practitioner.Practitioner, error)
	findByIhsPractitionerIDFn func(ctx context.Context, ihsID string) (*practitioner.Practitioner, error)
	findAllFn                 func(ctx context.Context, params practitioner.ListParams) ([]practitioner.Practitioner, int, error)
}

func (m *mockPractitionerRepo) Create(ctx context.Context, p *practitioner.Practitioner) (*practitioner.Practitioner, error) {
	if m.createFn != nil {
		return m.createFn(ctx, p)
	}
	p.ID = "test-uuid"
	return p, nil
}

func (m *mockPractitionerRepo) Update(ctx context.Context, p *practitioner.Practitioner) (*practitioner.Practitioner, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, p)
	}
	return p, nil
}

func (m *mockPractitionerRepo) Delete(ctx context.Context, id string, deletedBy string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id, deletedBy)
	}
	return nil
}

func (m *mockPractitionerRepo) FindByID(ctx context.Context, id string) (*practitioner.Practitioner, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockPractitionerRepo) FindByNIK(ctx context.Context, nik string) (*practitioner.Practitioner, error) {
	if m.findByNIKFn != nil {
		return m.findByNIKFn(ctx, nik)
	}
	return nil, nil
}

func (m *mockPractitionerRepo) FindByNIP(ctx context.Context, nip string) (*practitioner.Practitioner, error) {
	if m.findByNIPFn != nil {
		return m.findByNIPFn(ctx, nip)
	}
	return nil, nil
}

func (m *mockPractitionerRepo) FindByIhsPractitionerID(ctx context.Context, ihsID string) (*practitioner.Practitioner, error) {
	if m.findByIhsPractitionerIDFn != nil {
		return m.findByIhsPractitionerIDFn(ctx, ihsID)
	}
	return nil, nil
}

func (m *mockPractitionerRepo) FindAll(ctx context.Context, params practitioner.ListParams) ([]practitioner.Practitioner, int, error) {
	if m.findAllFn != nil {
		return m.findAllFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *mockPractitionerRepo) FindAllProfessions(ctx context.Context) ([]practitioner.Profession, error) {
	return []practitioner.Profession{{Code: "DOKTER", Name: "Dokter"}}, nil
}

func (m *mockPractitionerRepo) FindAllSpecialties(ctx context.Context) ([]practitioner.Specialty, error) {
	return []practitioner.Specialty{{Code: "SP-PD", Name: "Spesialis Penyakit Dalam"}}, nil
}

func TestCreatePractitioner_Success(t *testing.T) {
	mockRepo := &mockPractitionerRepo{
		findByNIKFn: func(ctx context.Context, nik string) (*practitioner.Practitioner, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createFn: func(ctx context.Context, p *practitioner.Practitioner) (*practitioner.Practitioner, error) {
			p.ID = "prac-001"
			return p, nil
		},
		findByIDFn: func(ctx context.Context, id string) (*practitioner.Practitioner, error) {
			return &practitioner.Practitioner{
				ID:     id,
				Name:   "dr. Budi Santoso, Sp.A",
				NIK:    "3201000000000001",
				Gender: enums.GenderMale,
			}, nil
		},
	}

	svc := practitioner.NewService(mockRepo)
	sipExpiry := "2028-12-31"
	req := practitioner.CreatePractitionerRequest{
		Name:          "dr. Budi Santoso, Sp.A",
		NIK:           "3201000000000001",
		Gender:        "laki-laki",
		SIP:           "SIP-12345",
		SIPExpiryDate: &sipExpiry,
	}

	res, err := svc.CreatePractitioner(context.Background(), req, "admin")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if res.ID != "prac-001" {
		t.Errorf("expected ID 'prac-001', got: %s", res.ID)
	}
	if res.Gender != enums.GenderMale {
		t.Errorf("expected gender 'L', got: %s", res.Gender)
	}
}

func TestCreatePractitioner_DuplicateNIK(t *testing.T) {
	mockRepo := &mockPractitionerRepo{
		findByNIKFn: func(ctx context.Context, nik string) (*practitioner.Practitioner, error) {
			return &practitioner.Practitioner{ID: "existing-id", NIK: nik}, nil
		},
	}

	svc := practitioner.NewService(mockRepo)
	req := practitioner.CreatePractitionerRequest{
		Name:   "dr. Test",
		NIK:    "3201000000000001",
		Gender: "L",
	}

	_, err := svc.CreatePractitioner(context.Background(), req, "admin")
	if !errors.Is(err, practitioner.ErrNIKAlreadyExists) {
		t.Fatalf("expected ErrNIKAlreadyExists, got: %v", err)
	}
}

func TestCreatePractitioner_DuplicateNIP(t *testing.T) {
	mockRepo := &mockPractitionerRepo{
		findByNIPFn: func(ctx context.Context, nip string) (*practitioner.Practitioner, error) {
			return &practitioner.Practitioner{ID: "existing-id", NIP: nip}, nil
		},
	}

	svc := practitioner.NewService(mockRepo)
	req := practitioner.CreatePractitionerRequest{
		Name:   "dr. Test",
		NIP:    "198501012010011001",
		Gender: "L",
	}

	_, err := svc.CreatePractitioner(context.Background(), req, "admin")
	if !errors.Is(err, practitioner.ErrNIPAlreadyExists) {
		t.Fatalf("expected ErrNIPAlreadyExists, got: %v", err)
	}
}

func TestCreatePractitioner_InvalidNIKFormat(t *testing.T) {
	svc := practitioner.NewService(&mockPractitionerRepo{})
	req := practitioner.CreatePractitionerRequest{
		Name:   "dr. Test",
		NIK:    "12345", // Bukan 16 digit
		Gender: "L",
	}

	_, err := svc.CreatePractitioner(context.Background(), req, "admin")
	if !errors.Is(err, practitioner.ErrInvalidNIKFormat) {
		t.Fatalf("expected ErrInvalidNIKFormat, got: %v", err)
	}
}

func TestCreatePractitioner_InvalidGender(t *testing.T) {
	svc := practitioner.NewService(&mockPractitionerRepo{})
	req := practitioner.CreatePractitionerRequest{
		Name:   "dr. Test",
		Gender: "Alien",
	}

	_, err := svc.CreatePractitioner(context.Background(), req, "admin")
	if !errors.Is(err, practitioner.ErrInvalidGender) {
		t.Fatalf("expected ErrInvalidGender, got: %v", err)
	}
}

func TestCreatePractitioner_Sanitization(t *testing.T) {
	var captured *practitioner.Practitioner
	mockRepo := &mockPractitionerRepo{
		createFn: func(ctx context.Context, p *practitioner.Practitioner) (*practitioner.Practitioner, error) {
			captured = p
			p.ID = "prac-002"
			return p, nil
		},
		findByIDFn: func(ctx context.Context, id string) (*practitioner.Practitioner, error) {
			return captured, nil
		},
	}

	svc := practitioner.NewService(mockRepo)
	emptyProf := "   "
	req := practitioner.CreatePractitionerRequest{
		Name:         "   dr. John Doe   ",
		Gender:       "  male ",
		ProfessionID: &emptyProf,
	}

	_, err := svc.CreatePractitioner(context.Background(), req, "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.Name != "dr. John Doe" {
		t.Errorf("expected trimmed name 'dr. John Doe', got: '%s'", captured.Name)
	}
	if captured.Gender != enums.GenderMale {
		t.Errorf("expected normalized gender 'L', got: '%s'", captured.Gender)
	}
	if captured.ProfessionID != nil {
		t.Errorf("expected empty string ProfessionID to be sanitized to nil, got: %v", captured.ProfessionID)
	}
}

func TestUpdatePractitioner_NotFound(t *testing.T) {
	mockRepo := &mockPractitionerRepo{
		findByIDFn: func(ctx context.Context, id string) (*practitioner.Practitioner, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	svc := practitioner.NewService(mockRepo)
	req := practitioner.UpdatePractitionerRequest{
		CreatePractitionerRequest: practitioner.CreatePractitionerRequest{
			Name:   "dr. Not Found",
			Gender: "L",
		},
	}

	_, err := svc.UpdatePractitioner(context.Background(), "non-existent-id", req, "admin")
	if !errors.Is(err, practitioner.ErrPractitionerNotFound) {
		t.Fatalf("expected ErrPractitionerNotFound, got: %v", err)
	}
}

func TestDeletePractitioner_Success(t *testing.T) {
	deleted := false
	mockRepo := &mockPractitionerRepo{
		findByIDFn: func(ctx context.Context, id string) (*practitioner.Practitioner, error) {
			return &practitioner.Practitioner{ID: id, Name: "dr. Delete"}, nil
		},
		deleteFn: func(ctx context.Context, id string, deletedBy string) error {
			deleted = true
			return nil
		},
	}

	svc := practitioner.NewService(mockRepo)
	err := svc.DeletePractitioner(context.Background(), "prac-003", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Errorf("expected delete to be executed")
	}
}
