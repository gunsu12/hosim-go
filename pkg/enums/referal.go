package enums

type ReferalType string

const (
	ReferalTypeDoctor   ReferalType = "DOCTOR"
	ReferalTypeHospital ReferalType = "HOSPITAL"
	ReferalTypeFKTP     ReferalType = "FKTP"
	ReferalTypeApotek   ReferalType = "APOTEK"
	ReferalTypeAgent    ReferalType = "AGENT"
	ReferalTypeOther    ReferalType = "OTHER"
)

func (rt ReferalType) IsValid() bool {
	switch rt {
	case ReferalTypeDoctor, ReferalTypeHospital, ReferalTypeFKTP, ReferalTypeApotek, ReferalTypeAgent, ReferalTypeOther:
		return true
	default:
		return false
	}
}

func GetReferalType(code string) ReferalType {
	switch code {
	case "DOCTOR":
		return ReferalTypeDoctor
	case "HOSPITAL":
		return ReferalTypeHospital
	case "FKTP":
		return ReferalTypeFKTP
	case "APOTEK":
		return ReferalTypeApotek
	case "AGENT":
		return ReferalTypeAgent
	default:
		return ReferalTypeOther
	}
}
