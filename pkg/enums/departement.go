package enums

type DepartementType string

const (
	DepartementTypeEmergency  DepartementType = "emergency"
	DepartementTypeOutpatient DepartementType = "outpatient"
	DepartementTypeInpatient  DepartementType = "inpatient"
	DepartementTypeDiagnostic DepartementType = "diagnostic"
	DepartementTypeMCU        DepartementType = "medical_checkup"
	DepartementTypeOther      DepartementType = "other"
)
