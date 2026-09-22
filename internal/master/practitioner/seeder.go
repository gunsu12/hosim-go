package practitioner

import (
	"log"
	"time"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DefaultProfessions memuat daftar standar profesi tenaga medis / kesehatan di RS
var DefaultProfessions = []Profession{
	{Code: "DOKTER", Name: "Dokter"},
	{Code: "PERAWAT", Name: "Perawat"},
	{Code: "ANALIS", Name: "Analis Kesehatan"},
	{Code: "BIDAN", Name: "Bidan"},
	{Code: "RADIOGRAFER", Name: "Radiografer"},
	{Code: "GIZI", Name: "Ahli Gizi"},
	{Code: "TERAPIS", Name: "Terapis"},
	{Code: "FARMASI", Name: "Apoteker / Tenaga Teknis Kefarmasian"},
}

// DefaultSpecialties memuat daftar spesialisasi medis mengacu standar SatuSehat Kemenkes / Kolegium Kedokteran
var DefaultSpecialties = []Specialty{
	{Code: "DOKTER-UMUM", Name: "Dokter Umum (General Practitioner)"},
	{Code: "SP-PD", Name: "Spesialis Penyakit Dalam (Internal Medicine)"},
	{Code: "SP-A", Name: "Spesialis Anak (Pediatrics)"},
	{Code: "SP-OG", Name: "Spesialis Obstetri dan Ginekologi (Obgyn)"},
	{Code: "SP-B", Name: "Spesialis Bedah (General Surgery)"},
	{Code: "SP-S", Name: "Spesialis Saraf / Neurologi (Neurology)"},
	{Code: "SP-JP", Name: "Spesialis Jantung dan Pembuluh Darah (Cardiology)"},
	{Code: "SP-M", Name: "Spesialis Mata (Ophthalmology)"},
	{Code: "SP-THT", Name: "Spesialis Telinga Hidung Tenggorok (Otolaryngology)"},
	{Code: "SP-KJ", Name: "Spesialis Kedokteran Jiwa / Psikiatri (Psychiatry)"},
	{Code: "SP-DVE", Name: "Spesialis Dermatologi, Venereologi, dan Estetika (Kulit & Kelamin)"},
	{Code: "SP-P", Name: "Spesialis Paru / Pulmonologi (Pulmonology)"},
	{Code: "SP-AN", Name: "Spesialis Anestesiologi dan Terapi Intensif (Anesthesiology)"},
	{Code: "SP-RAD", Name: "Spesialis Radiologi (Radiology)"},
	{Code: "SP-PK", Name: "Spesialis Patologi Klinik (Clinical Pathology)"},
	{Code: "SP-PA", Name: "Spesialis Patologi Anatomi (Anatomical Pathology)"},
	{Code: "SP-OT", Name: "Spesialis Orthopaedi dan Traumatologi (Orthopedic Surgery)"},
	{Code: "SP-U", Name: "Spesialis Urologi (Urology)"},
	{Code: "SP-BS", Name: "Spesialis Bedah Saraf (Neurosurgery)"},
	{Code: "SP-KFR", Name: "Spesialis Kedokteran Fisik dan Rehabilitasi (Rehab Medik)"},
	{Code: "DOKTER-GIGI", Name: "Dokter Gigi (Dentist)"},
	{Code: "SP-KG", Name: "Spesialis Konservasi Gigi (Endodontics)"},
	{Code: "SP-BM", Name: "Spesialis Bedah Mulut dan Maksilofasial (Oral & Maxillofacial Surgery)"},
}

// SeedProfessionsAndSpecialties mengisi data awal (master reference) profesi dan spesialisasi secara idempoten
func SeedProfessionsAndSpecialties(db *gorm.DB) error {
	log.Println("[SEEDER] Memeriksa & melakukan seeding master Profesi...")
	for _, prof := range DefaultProfessions {
		var count int64
		db.Model(&Profession{}).Where("code = ?", prof.Code).Count(&count)
		if count == 0 {
			prof.CreatedBy = "system-seeder"
			prof.UpdatedBy = "system-seeder"
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&prof).Error; err != nil {
				return err
			}
		}
	}

	log.Println("[SEEDER] Memeriksa & melakukan seeding master Spesialisasi (SatuSehat)...")
	for _, spec := range DefaultSpecialties {
		var count int64
		db.Model(&Specialty{}).Where("code = ?", spec.Code).Count(&count)
		if count == 0 {
			spec.CreatedBy = "system-seeder"
			spec.UpdatedBy = "system-seeder"
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&spec).Error; err != nil {
				return err
			}
		}
	}

	log.Println("[SEEDER] Seeding master Profesi dan Spesialisasi berhasil.")
	return nil
}

// SeedDefaultPractitioner mengisi data awal tenaga medis dokter secara idempoten
func SeedDefaultPractitioner(db *gorm.DB) error {
	var count int64
	db.Model(&Practitioner{}).Where("name LIKE ?", "%Hendra Wijaya%").Count(&count)
	if count > 0 {
		return nil
	}

	var prof Profession
	_ = db.Where("code = ?", "DOKTER").First(&prof).Error

	var spec Specialty
	_ = db.Where("code = ?", "SP-B").First(&spec).Error

	expiry := time.Now().AddDate(3, 0, 0)
	doc := Practitioner{
		NIK:                 "3171012345678901",
		NIP:                 "198501152010121001",
		Name:                "dr. Hendra Wijaya, Sp.B",
		Gender:              enums.GenderMale,
		SIP:                 "503/SIP.DS/0123/2023",
		SIPExpiryDate:       &expiry,
		STR:                 "31.1.1.100.2.19.123456",
		Phone:               "081234567890",
		Email:               "hendra.wijaya@hosim.local",
		IhsPractitionerID:   "P12345678901",
		IhsPractitionerName: "dr. Hendra Wijaya, Sp.B",
		IsActive:            true,
		CreatedBy:           "SEEDER",
		UpdatedBy:           "SEEDER",
	}

	if prof.ID != "" {
		doc.ProfessionID = &prof.ID
	}
	if spec.ID != "" {
		doc.SpecialtyID = &spec.ID
	}

	if err := db.Create(&doc).Error; err != nil {
		log.Printf("[WARN] Seeder default practitioner gagal: %v\n", err)
		return err
	}

	log.Println("[SEEDER] Data dokter simulasi (dr. Hendra Wijaya, Sp.B) berhasil dibuat.")
	return nil
}
