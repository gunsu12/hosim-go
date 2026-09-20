package enums

import "strings"

// Gender merepresentasikan jenis kelamin di sistem SIMRS & SatuSehat
type Gender string

const (
	GenderMale    Gender = "L"
	GenderFemale  Gender = "P"
	GenderUnknown Gender = "Other"
)

// IsValid memvalidasi apakah nilai gender sesuai pilihan enum
func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderUnknown:
		return true
	default:
		return false
	}
}

// String mengembalikan string representasi gender
func (g Gender) String() string {
	return string(g)
}

// NormalizeGender memetakan berbagai variasi string input ke tipe enum Gender standar
func NormalizeGender(input string) (Gender, bool) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "l", "laki-laki", "male", "m":
		return GenderMale, true
	case "p", "perempuan", "female", "f":
		return GenderFemale, true
	case "other", "o", "unknown", "u", "tidak tahu", "lainnya":
		return GenderUnknown, true
	default:
		return "", false
	}
}
