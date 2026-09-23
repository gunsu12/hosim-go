-- +goose Up
-- ========================================================
-- Modul: Master Pasien (Patient & Sub-Entities)
-- ========================================================

CREATE TABLE IF NOT EXISTS patients (
    id VARCHAR(36) PRIMARY KEY,
    medical_record_no VARCHAR(8) NOT NULL,
    nik VARCHAR(16),
    family_card_no VARCHAR(16),
    ihs_patient_id VARCHAR(50),
    title VARCHAR(20),
    short_name VARCHAR(150) NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    mother_name VARCHAR(150),
    gender VARCHAR(10) NOT NULL,
    birth_place VARCHAR(50),
    birth_date TIMESTAMPTZ NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    marital_status VARCHAR(20),
    religion VARCHAR(20),
    education VARCHAR(20),
    occupation VARCHAR(20),
    nationality VARCHAR(20),
    languages VARCHAR(20),
    blood_type VARCHAR(20),
    rhesus VARCHAR(10),
    special_needs VARCHAR(150),
    is_unknown BOOLEAN DEFAULT FALSE,
    is_deceased BOOLEAN DEFAULT FALSE,
    deceased_at TIMESTAMPTZ,
    customer_id VARCHAR(36) REFERENCES customers(id) ON DELETE SET NULL,
    insurance_type VARCHAR(20),
    insurance_number VARCHAR(50),
    insurance_expiry_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_mrn ON patients(medical_record_no);

CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_nik ON patients(nik)
WHERE nik IS NOT NULL AND nik != '' AND deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_ihs ON patients(ihs_patient_id)
WHERE ihs_patient_id IS NOT NULL AND ihs_patient_id != '' AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_patients_family_card_no ON patients(family_card_no);
CREATE INDEX IF NOT EXISTS idx_patients_name_birth ON patients(full_name, birth_date);
CREATE INDEX IF NOT EXISTS idx_patients_phone ON patients(phone);
CREATE INDEX IF NOT EXISTS idx_patients_customer_id ON patients(customer_id);
CREATE INDEX IF NOT EXISTS idx_patients_insurance_number ON patients(insurance_number);
CREATE INDEX IF NOT EXISTS idx_patients_deleted_at ON patients(deleted_at);

-- 1. Emergency Contacts
CREATE TABLE IF NOT EXISTS patient_emergency_contacts (
    id VARCHAR(36) PRIMARY KEY,
    patient_id VARCHAR(36) NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    relation VARCHAR(20),
    phone VARCHAR(20),
    address VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_pec_patient_active ON patient_emergency_contacts(patient_id, is_active);
CREATE INDEX IF NOT EXISTS idx_pec_deleted_at ON patient_emergency_contacts(deleted_at);

-- 2. Relations
CREATE TABLE IF NOT EXISTS patient_relations (
    id VARCHAR(36) PRIMARY KEY,
    patient_id VARCHAR(36) NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    relation VARCHAR(20),
    phone VARCHAR(20),
    address VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_pr_patient_active ON patient_relations(patient_id, is_active);
CREATE INDEX IF NOT EXISTS idx_pr_deleted_at ON patient_relations(deleted_at);

-- 3. Allergies
CREATE TABLE IF NOT EXISTS patient_allergies (
    id VARCHAR(36) PRIMARY KEY,
    patient_id VARCHAR(36) NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    allergy VARCHAR(150) NOT NULL,
    reaction VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_pa_patient_active ON patient_allergies(patient_id, is_active);
CREATE INDEX IF NOT EXISTS idx_pa_deleted_at ON patient_allergies(deleted_at);

-- 4. Addresses
CREATE TABLE IF NOT EXISTS patient_addresses (
    id VARCHAR(36) PRIMARY KEY,
    patient_id VARCHAR(36) NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    address_type VARCHAR(20) NOT NULL,
    address_line VARCHAR(255) NOT NULL,
    rt VARCHAR(3),
    rw VARCHAR(3),
    postal_code VARCHAR(10),
    provinsi_id VARCHAR(20),
    kabupaten_id VARCHAR(20),
    kecamatan_id VARCHAR(20),
    kelurahan_id VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_paddr_patient_type_active ON patient_addresses(patient_id, address_type, is_active);
CREATE INDEX IF NOT EXISTS idx_paddr_deleted_at ON patient_addresses(deleted_at);

-- 5. Drug Histories
CREATE TABLE IF NOT EXISTS patient_drug_histories (
    id VARCHAR(36) PRIMARY KEY,
    patient_id VARCHAR(36) NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    drug_name VARCHAR(150) NOT NULL,
    indication VARCHAR(255),
    is_chronic BOOLEAN DEFAULT FALSE,
    duration_use VARCHAR(20),
    manufacturer VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_pdh_patient_active ON patient_drug_histories(patient_id, is_active);
CREATE INDEX IF NOT EXISTS idx_pdh_deleted_at ON patient_drug_histories(deleted_at);

-- 6. Chronical Diseases
CREATE TABLE IF NOT EXISTS patient_chronical_diseases (
    id VARCHAR(36) PRIMARY KEY,
    patient_id VARCHAR(36) NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    disease VARCHAR(150) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_pcd_patient_active ON patient_chronical_diseases(patient_id, is_active);
CREATE INDEX IF NOT EXISTS idx_pcd_deleted_at ON patient_chronical_diseases(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS patient_chronical_diseases;
DROP TABLE IF EXISTS patient_drug_histories;
DROP TABLE IF EXISTS patient_addresses;
DROP TABLE IF EXISTS patient_allergies;
DROP TABLE IF EXISTS patient_relations;
DROP TABLE IF EXISTS patient_emergency_contacts;
DROP TABLE IF EXISTS patients;
