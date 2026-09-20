package patient

import (
	"time"
	"uuid"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

type Patient struct {
	ID                  string                    `gorm:"primaryKey;size:36" json:"id"`
	MedicalRecordNo     string                    `gorm:"uniqueIndex;size:8;not null" json:"medical_record_no"`
	NIK                 string                    `gorm:"uniqueIndex:idx_patients_nik,where:nik != '' AND nik IS NOT NULL AND deleted_at IS NULL;size:16" json:"nik"`
	FamilyCardNo        string                    `gorm:"index;size:16" json:"family_card_no,omitempty"`
	IHSPatientID        string                    `gorm:"uniqueIndex:idx_patients_ihs,where:ihs_patient_id != '' AND ihs_patient_id IS NOT NULL AND deleted_at IS NULL;size:50" json:"ihs_patient_id,omitempty"`
	Title               string                    `gorm:"size:20" json:"title,omitempty"`
	ShortName           string                    `gorm:"not null;size:150" json:"short_name"`
	FullName            string                    `gorm:"index:idx_name_birth;not null;size:200" json:"full_name"`
	MotherName          string                    `gorm:"size:150" json:"mother_name,omitempty"`
	Gender              enums.Gender              `gorm:"size:10;not null" json:"gender"`
	BirthPlace          string                    `gorm:"size:50" json:"birth_place"`
	BirthDate           time.Time                 `gorm:"index:idx_name_birth;not null" json:"birth_date"`
	Phone               string                    `gorm:"index;not null;size:20" json:"phone"`
	Email               string                    `gorm:"size:100" json:"email"`
	MaritalStatus       string                    `gorm:"size:20" json:"marital_status"`
	Religion            string                    `gorm:"size:20" json:"religion"`
	Education           string                    `gorm:"size:20" json:"education"`
	Occupation          string                    `gorm:"size:20" json:"occupation"`
	Nationality         string                    `gorm:"size:20" json:"nationality"`
	BloodType           string                    `gorm:"size:20" json:"blood_type"`
	Rhesus              string                    `gorm:"size:10" json:"rhesus,omitempty"`
	SpecialNeeds        string                    `gorm:"size:150" json:"special_needs,omitempty"`
	IsUnknown           bool                      `gorm:"default:false" json:"is_unknown"`
	IsDeceased          bool                      `gorm:"default:false" json:"is_deceased"`
	DeceasedAt          *time.Time                `json:"deceased_at,omitempty"`
	InsuranceType       string                    `gorm:"size:20" json:"insurance_type"`
	InsuranceNumber     string                    `gorm:"index;size:50" json:"insurance_number"`
	InsuranceExpiryDate *time.Time                `gorm:"" json:"insurance_expiry_date"`
	EmergencyContacts   []PatientEmergencyContact `gorm:"foreignKey:PatientID;references:ID" json:"emergency_contacts,omitempty"`
	Relations           []PatientRelation         `gorm:"foreignKey:PatientID;references:ID" json:"relations,omitempty"`
	Allergies           []PatientAllergy          `gorm:"foreignKey:PatientID;references:ID" json:"allergies,omitempty"`
	Addresses           []PatientAddress          `gorm:"foreignKey:PatientID;references:ID" json:"addresses,omitempty"`
	DrugHistories       []PatientDrugHistory      `gorm:"foreignKey:PatientID;references:ID" json:"drug_histories,omitempty"`
	ChronicalDiseases   []PatientChronicalDisease `gorm:"foreignKey:PatientID;references:ID" json:"chronical_diseases,omitempty"`
	CreatedAt           time.Time                 `json:"created_at"`
	CreatedBy           string                    `json:"created_by"`
	UpdatedAt           time.Time                 `json:"updated_at"`
	UpdatedBy           string                    `json:"updated_by"`
	DeletedAt           gorm.DeletedAt            `gorm:"index" json:"-"`
	DeletedBy           string                    `json:"deleted_by"`
}

type PatientEmergencyContact struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	PatientID string         `gorm:"index:idx_pec_patient_active;size:36;not null" json:"patient_id"`
	Name      string         `gorm:"size:150;not null" json:"name"`
	Relation  string         `gorm:"size:20" json:"relation"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Address   string         `gorm:"size:255" json:"address"`
	IsActive  bool           `gorm:"index:idx_pec_patient_active;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by"`
}

type PatientRelation struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	PatientID string         `gorm:"index:idx_pr_patient_active;size:36;not null" json:"patient_id"`
	Name      string         `gorm:"size:150;not null" json:"name"`
	Relation  string         `gorm:"size:20" json:"relation"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Address   string         `gorm:"size:255" json:"address"`
	IsActive  bool           `gorm:"index:idx_pr_patient_active;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by"`
}

type PatientAllergy struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	PatientID string         `gorm:"index:idx_pa_patient_active;size:36;not null" json:"patient_id"`
	Allergy   string         `gorm:"size:150;not null" json:"allergy"`
	Reaction  string         `gorm:"size:255" json:"reaction"`
	IsActive  bool           `gorm:"index:idx_pa_patient_active;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by"`
}

type PatientAddress struct {
	ID          string         `gorm:"primaryKey;size:36" json:"id"`
	PatientID   string         `gorm:"index:idx_paddr_patient_type_active;size:36;not null" json:"patient_id"`
	AddressType string         `gorm:"index:idx_paddr_patient_type_active;size:20;not null" json:"address_type"` // "ktp", "domicile", dll
	AddressLine string         `gorm:"not null;size:255" json:"address_line"`
	RT          string         `gorm:"size:3" json:"rt"`
	RW          string         `gorm:"size:3" json:"rw"`
	PostalCode  string         `gorm:"size:10" json:"postal_code"`
	ProvinsiID  string         `gorm:"size:20" json:"provinsi_id"`
	KabupatenID string         `gorm:"size:20" json:"kabupaten_id"`
	KecamatanID string         `gorm:"size:20" json:"kecamatan_id"`
	KelurahanID string         `gorm:"size:20" json:"kelurahan_id"`
	IsActive    bool           `gorm:"index:idx_paddr_patient_type_active;default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	CreatedBy   string         `json:"created_by"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UpdatedBy   string         `json:"updated_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy   string         `json:"deleted_by"`
}

type PatientDrugHistory struct {
	ID           string         `gorm:"primaryKey;size:36" json:"id"`
	PatientID    string         `gorm:"index:idx_pdh_patient_active;size:36;not null" json:"patient_id"`
	DrugName     string         `gorm:"size:150;not null" json:"drug_name"`
	Indication   string         `gorm:"size:255" json:"indication"`
	IsChronic    bool           `gorm:"default:false" json:"is_chronic"`
	DurationUse  string         `gorm:"size:20" json:"duration_use"`
	Manufacturer string         `gorm:"size:100" json:"manufacturer"`
	IsActive     bool           `gorm:"index:idx_pdh_patient_active;default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	CreatedBy    string         `json:"created_by"`
	UpdatedAt    time.Time      `json:"updated_at"`
	UpdatedBy    string         `json:"updated_by"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy    string         `json:"deleted_by"`
}

type PatientChronicalDisease struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	PatientID string         `gorm:"index:idx_pcd_patient_active;size:36;not null" json:"patient_id"`
	Disease   string         `gorm:"size:150;not null" json:"disease"`
	IsActive  bool           `gorm:"index:idx_pcd_patient_active;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by"`
}

func (p *Patient) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewV7().String()
	}
	return nil
}

func (c *PatientEmergencyContact) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewV7().String()
	}
	return nil
}

func (r *PatientRelation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewV7().String()
	}
	return nil
}

func (a *PatientAllergy) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.NewV7().String()
	}
	return nil
}

func (addr *PatientAddress) BeforeCreate(tx *gorm.DB) error {
	if addr.ID == "" {
		addr.ID = uuid.NewV7().String()
	}
	return nil
}

func (d *PatientDrugHistory) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewV7().String()
	}
	return nil
}

func (cd *PatientChronicalDisease) BeforeCreate(tx *gorm.DB) error {
	if cd.ID == "" {
		cd.ID = uuid.NewV7().String()
	}
	return nil
}
