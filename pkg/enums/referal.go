package enums

type ReferalType string

const (
	ReferalTypeDoctor    ReferalType = "DOCTOR"
	ReferalTypeHospital  ReferalType = "HOSPITAL"
	ReferalTypePharmacy  ReferalType = "PHARMACY"
	ReferalTypeLab       ReferalType = "LAB"
	ReferalTypeRadiology ReferalType = "RADIOLOGY"
	ReferalTypeTherapy   ReferalType = "THERAPY"
	ReferalTypeOther     ReferalType = "OTHER"
)
