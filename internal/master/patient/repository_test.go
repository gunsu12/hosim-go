package patient_test

import (
	"context"
	"testing"
	"time"

	"hosim-go/internal/config"
	"hosim-go/internal/database"
	"hosim-go/internal/master/patient"
	"hosim-go/internal/master/payer"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	cfg := config.LoadConfig()
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("failed to connect to postgres database: %v", err)
	}

	err = db.AutoMigrate(
		&payer.PayerType{},
		&payer.Payer{},
		&patient.Patient{},
		&patient.PatientEmergencyContact{},
		&patient.PatientRelation{},
		&patient.PatientAllergy{},
		&patient.PatientAddress{},
		&patient.PatientDrugHistory{},
		&patient.PatientChronicalDisease{},
	)
	if err != nil {
		t.Fatalf("failed to auto migrate: %v", err)
	}

	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}

func TestRepository_UpdateAssociations(t *testing.T) {
	db := setupTestDB(t)
	repo := patient.NewRepository(db)
	ctx := context.Background()

	p := &patient.Patient{
		MedicalRecordNo: "00000001",
		NIK:             "3201000000000001",
		ShortName:       "Test",
		FullName:        "Test Patient",
		Gender:          "L",
		BirthPlace:      "Jakarta",
		BirthDate:       time.Now(),
		Phone:           "0812345",
		EmergencyContacts: []patient.PatientEmergencyContact{
			{Name: "Contact 1", Relation: "Ibu", Phone: "08111", IsActive: true},
		},
		Relations: []patient.PatientRelation{
			{Name: "Keluarga 1", Relation: "Ayah", Phone: "08112", IsActive: true},
		},
		Addresses: []patient.PatientAddress{
			{AddressType: "ktp", AddressLine: "Jalan 1", IsActive: true},
		},
	}

	created, err := repo.Create(ctx, p)
	if err != nil {
		t.Fatalf("failed to create patient: %v", err)
	}

	// Now update associations
	created.EmergencyContacts = []patient.PatientEmergencyContact{
		{Name: "Contact 2", Relation: "Ayah", Phone: "08222", IsActive: true},
	}
	created.Relations = []patient.PatientRelation{
		{Name: "Keluarga 2", Relation: "Istri", Phone: "08223", IsActive: true},
	}
	created.Addresses = []patient.PatientAddress{
		{AddressType: "ktp", AddressLine: "Jalan 2 Baru", IsActive: true},
	}

	updated, err := repo.Update(ctx, created)
	if err != nil {
		t.Fatalf("failed to update patient: %v", err)
	}

	// Fetch back
	fetched, err := repo.FindByID(ctx, updated.ID)
	if err != nil {
		t.Fatalf("failed to find patient: %v", err)
	}

	t.Logf("EmergencyContacts len: %d", len(fetched.EmergencyContacts))
	t.Logf("Relations len: %d", len(fetched.Relations))
	t.Logf("Addresses len: %d", len(fetched.Addresses))

	if len(fetched.EmergencyContacts) != 1 {
		t.Errorf("expected 1 emergency contact, got %d", len(fetched.EmergencyContacts))
	}
	if len(fetched.Relations) != 1 {
		t.Errorf("expected 1 relation, got %d", len(fetched.Relations))
	}
	if len(fetched.Addresses) != 1 {
		t.Errorf("expected 1 address, got %d", len(fetched.Addresses))
	}
	if len(fetched.EmergencyContacts) > 0 && fetched.EmergencyContacts[0].Name != "Contact 2" {
		t.Errorf("expected 'Contact 2', got '%s'", fetched.EmergencyContacts[0].Name)
	}
	if len(fetched.Relations) > 0 && fetched.Relations[0].Name != "Keluarga 2" {
		t.Errorf("expected 'Keluarga 2', got '%s'", fetched.Relations[0].Name)
	}
}

func TestRepository_CascadeDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := patient.NewRepository(db)
	ctx := context.Background()

	p := &patient.Patient{
		MedicalRecordNo: "00000002",
		NIK:             "3201000000000002",
		ShortName:       "Budi",
		FullName:        "Budi Santoso",
		BirthDate:       time.Now(),
		Phone:           "0812345",
		EmergencyContacts: []patient.PatientEmergencyContact{
			{Name: "Contact A", Relation: "Ibu", Phone: "08111", IsActive: true},
		},
		Relations: []patient.PatientRelation{
			{Name: "Penanggung Jawab A", Relation: "Ayah", Phone: "08112", IsActive: true},
		},
		Addresses: []patient.PatientAddress{
			{AddressType: "ktp", AddressLine: "Jalan Melati", IsActive: true},
		},
	}

	created, err := repo.Create(ctx, p)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Delete
	err = repo.Delete(ctx, created.ID, "operator-audit")
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	// Check child tables are soft-deleted
	var contactCount, relationCount, addressCount int64
	db.Model(&patient.PatientEmergencyContact{}).Where("patient_id = ?", created.ID).Count(&contactCount)
	db.Model(&patient.PatientRelation{}).Where("patient_id = ?", created.ID).Count(&relationCount)
	db.Model(&patient.PatientAddress{}).Where("patient_id = ?", created.ID).Count(&addressCount)

	if contactCount != 0 {
		t.Errorf("expected 0 active contacts after cascade delete, got %d", contactCount)
	}
	if relationCount != 0 {
		t.Errorf("expected 0 active relations after cascade delete, got %d", relationCount)
	}
	if addressCount != 0 {
		t.Errorf("expected 0 active addresses after cascade delete, got %d", addressCount)
	}
}

func TestRepository_FindAll_CaseInsensitive(t *testing.T) {
	db := setupTestDB(t)
	repo := patient.NewRepository(db)
	ctx := context.Background()

	p := &patient.Patient{
		MedicalRecordNo: "00000003",
		NIK:             "3201000000000003",
		ShortName:       "Siti",
		FullName:        "SITI NURHALIZA",
		BirthDate:       time.Now(),
		Phone:           "0812345",
	}
	_, err := repo.Create(ctx, p)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Search with lowercase "siti"
	patients, count, err := repo.FindAll(ctx, patient.ListParams{Page: 1, Limit: 10, Search: "siti"})
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}

	if count != 1 || len(patients) != 1 {
		t.Errorf("expected 1 match for lowercase search 'siti', got count=%d, len=%d", count, len(patients))
	}
}

func TestRepository_MultiplePatientsWithoutNIK(t *testing.T) {
	db := setupTestDB(t)
	repo := patient.NewRepository(db)
	ctx := context.Background()

	bbl1 := &patient.Patient{
		MedicalRecordNo: "00000010",
		NIK:             "", // BBL belum punya NIK
		ShortName:       "By Ny Siti 1",
		FullName:        "Bayi Ny Siti 1",
		BirthDate:       time.Now(),
		Phone:           "0812345",
	}
	_, err := repo.Create(ctx, bbl1)
	if err != nil {
		t.Fatalf("failed to create BBL 1: %v", err)
	}

	bbl2 := &patient.Patient{
		MedicalRecordNo: "00000011",
		NIK:             "", // BBL kedua juga belum punya NIK
		ShortName:       "By Ny Siti 2",
		FullName:        "Bayi Ny Siti 2",
		BirthDate:       time.Now(),
		Phone:           "0812345",
	}
	_, err = repo.Create(ctx, bbl2)
	if err != nil {
		t.Fatalf("FAILED to create BBL 2: %v", err)
	}
}

func TestRepository_SoftDeleteAndReRegisterNIK(t *testing.T) {
	db := setupTestDB(t)
	repo := patient.NewRepository(db)
	ctx := context.Background()

	p1 := &patient.Patient{
		MedicalRecordNo: "00000020",
		NIK:             "3201999999990001",
		ShortName:       "Pasien Hapus",
		FullName:        "Pasien Mau Dihapus",
		BirthDate:       time.Now(),
		Phone:           "0812345",
	}
	created, err := repo.Create(ctx, p1)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Soft delete p1
	err = repo.Delete(ctx, created.ID, "operator-audit")
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	// Re-register with the same NIK
	p2 := &patient.Patient{
		MedicalRecordNo: "00000021",
		NIK:             "3201999999990001", // NIK sama dengan yang dihapus
		ShortName:       "Pasien Baru",
		FullName:        "Pasien Daftar Ulang NIK Sama",
		BirthDate:       time.Now(),
		Phone:           "0899999",
	}
	_, err = repo.Create(ctx, p2)
	if err != nil {
		t.Fatalf("failed to re-register with same NIK after soft delete: %v", err)
	}
}

func TestRepository_DuplicateActiveNIK_Blocked(t *testing.T) {
	db := setupTestDB(t)
	repo := patient.NewRepository(db)
	ctx := context.Background()

	p1 := &patient.Patient{
		MedicalRecordNo: "00000030",
		NIK:             "3201888888880001",
		ShortName:       "P1",
		FullName:        "Pasien 1",
		BirthDate:       time.Now(),
		Phone:           "081",
	}
	_, err := repo.Create(ctx, p1)
	if err != nil {
		t.Fatalf("failed to create p1: %v", err)
	}

	p2 := &patient.Patient{
		MedicalRecordNo: "00000031",
		NIK:             "3201888888880001", // duplicate active NIK
		ShortName:       "P2",
		FullName:        "Pasien 2",
		BirthDate:       time.Now(),
		Phone:           "082",
	}
	_, err = repo.Create(ctx, p2)
	if err == nil {
		t.Fatalf("expected error duplicate NIK on active patients, but got nil")
	}
}

func TestRepository_GetLastMedicalRecordNo_IgnoresNonNumeric(t *testing.T) {
	db := setupTestDB(t)
	repo := patient.NewRepository(db)
	ctx := context.Background()

	// 1. Simpan pasien dengan RM standar 8 digit
	p1 := &patient.Patient{
		MedicalRecordNo: "00000045",
		NIK:             "3201777777770001",
		ShortName:       "P1",
		FullName:        "Pasien 45",
		BirthDate:       time.Now(),
		Phone:           "081",
	}
	if _, err := repo.Create(ctx, p1); err != nil {
		t.Fatalf("failed to create p1: %v", err)
	}

	// 2. Simpan pasien migrasi dengan RM berawalan huruf atau non-standar (misal "RM-00001", "ZZZZZZZZ")
	p2 := &patient.Patient{
		MedicalRecordNo: "RM-00001",
		NIK:             "3201777777770002",
		ShortName:       "P2",
		FullName:        "Pasien Migrasi",
		BirthDate:       time.Now(),
		Phone:           "082",
	}
	if _, err := repo.Create(ctx, p2); err != nil {
		t.Fatalf("failed to create p2: %v", err)
	}

	// 3. GetLastMedicalRecordNo harus mengembalikan "00000045", bukan "RM-00001"
	lastRM, err := repo.GetLastMedicalRecordNo(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting last RM: %v", err)
	}
	if lastRM != "00000045" {
		t.Errorf("expected last RM to be '00000045', got '%s'", lastRM)
	}
}



